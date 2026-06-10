package context

import (
	"fmt"
	"strings"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
)

// Builder structures the deterministic context fed to the LLM.
type Builder struct {
	MaxTokens int
}

func NewBuilder(maxTokens int) *Builder {
	return &Builder{MaxTokens: maxTokens}
}

// Build creates a deterministic, grounded prompt payload from a CorrelatedIncident.
func (b *Builder) Build(incident *correlationv1.CorrelatedIncident) (string, []string, bool) {
	var sb strings.Builder
	var eventIDs []string
	truncated := false

	sb.WriteString(fmt.Sprintf("Incident ID: %s\n", incident.IncidentId))
	sb.WriteString(fmt.Sprintf("Title: %s\n", incident.Title))
	sb.WriteString(fmt.Sprintf("Severity: %s\n", incident.Severity))
	
	if incident.Confidence != nil {
		sb.WriteString(fmt.Sprintf("Correlation Confidence: %.2f\n", incident.Confidence.Value))
	}

	sb.WriteString("\n--- CAUSAL CHAIN ---\n")
	if incident.CausalChain != nil {
		for i, link := range incident.CausalChain.Links {
			sb.WriteString(fmt.Sprintf("%d. [%s] -> [%s] (Reason: %s)\n", i+1, link.CauseEventId, link.EffectEventId, link.Reasoning))
		}
	}

	sb.WriteString("\n--- TELEMETRY EVIDENCE ---\n")
	for _, ev := range incident.Evidence {
		eventIDs = append(eventIDs, ev.EventId)
		
		// Very naive token truncation: ~4 chars per token.
		// In production, use a library like tiktoken.
		if sb.Len() > (b.MaxTokens * 4) {
			sb.WriteString("\n... [EVIDENCE TRUNCATED TO PREVENT TOKEN OVERFLOW] ...\n")
			truncated = true
			break
		}

		sb.WriteString(fmt.Sprintf("Time: %s | Source: %s\n", ev.Timestamp.AsTime().Format("15:04:05"), ev.Source))
		if k8s := ev.GetK8SEvent(); k8s != nil {
			sb.WriteString(fmt.Sprintf("  Action: %s | Kind: %s | Name: %s | Message: %s\n", k8s.Action, k8s.ResourceKind, k8s.ResourceName, k8s.Message))
		}
	}

	return sb.String(), eventIDs, truncated
}
