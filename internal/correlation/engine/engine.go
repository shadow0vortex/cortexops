package engine

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/correlation/causal"
	"github.com/shadow0vortex/cortexops/internal/correlation/heuristics"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ActiveIncident represents a mutable temporal window of correlated events.
type ActiveIncident struct {
	mu           sync.Mutex
	ID           string
	Evidence     []*correlationv1.TelemetryEnvelope
	Fingerprints map[string]struct{}
	LastActivity time.Time
	MaxScore     *correlationv1.CorrelationScore
}

const maxEvidenceCount = 100

func generateFingerprint(env *correlationv1.TelemetryEnvelope) string {
	if env == nil {
		return ""
	}
	traceID := ""
	if env.TraceContext != nil {
		traceID = env.TraceContext.TraceId
	}
	switch payload := env.Payload.(type) {
	case *correlationv1.TelemetryEnvelope_K8SEvent:
		k := payload.K8SEvent
		return fmt.Sprintf("k8s:%s:%s:%s:%s", k.Namespace, k.ResourceKind, k.Uid, traceID)
	case *correlationv1.TelemetryEnvelope_MetricEvent:
		m := payload.MetricEvent
		return fmt.Sprintf("metric:%s:%.2f:%s", m.MetricName, m.Value, traceID)
	default:
		return fmt.Sprintf("raw:%s:%s", env.EventId, traceID)
	}
}

// Engine ingests raw telemetry, applies deterministic scoring, and publishes CorrelatedIncidents.
type Engine struct {
	mu        sync.RWMutex
	incidents map[string]*ActiveIncident

	scorer       *heuristics.Scorer
	chainBuilder *causal.ChainBuilder
	publisher    core.Publisher
	metrics      core.MetricsRecorder
	logger       *slog.Logger
}

func NewEngine(scorer *heuristics.Scorer, chainBuilder *causal.ChainBuilder, pub core.Publisher, metrics core.MetricsRecorder, logger *slog.Logger) *Engine {
	return &Engine{
		incidents:    make(map[string]*ActiveIncident),
		scorer:       scorer,
		chainBuilder: chainBuilder,
		publisher:    pub,
		metrics:      metrics,
		logger:       logger,
	}
}

// ProcessEvent applies temporal and topological correlation logic to an incoming telemetry envelope.
func (e *Engine) ProcessEvent(ctx context.Context, incoming *correlationv1.TelemetryEnvelope) error {
	fingerprint := generateFingerprint(incoming)

	e.mu.Lock()
	defer e.mu.Unlock()

	var bestIncident *ActiveIncident
	var highestScore float64
	var bestReasoning []string

	// 1. Evaluate incoming event against all active temporal windows
	for _, inc := range e.incidents {
		inc.mu.Lock()
		
		// Semantic Deduplication Check
		if _, exists := inc.Fingerprints[fingerprint]; exists {
			inc.mu.Unlock()
			continue // Skip processing, already seen this semantic fingerprint
		}

		// Compare against the most recent event in the incident bucket
		latestEvent := inc.Evidence[len(inc.Evidence)-1]
		
		score, reasoning := e.scorer.Score(ctx, latestEvent, incoming)
		if score > highestScore {
			highestScore = score
			bestReasoning = reasoning
			bestIncident = inc
		}
		inc.mu.Unlock()
	}

	// 2. Correlation Threshold logic
	if highestScore >= 0.8 && bestIncident != nil {
		bestIncident.mu.Lock()
		
		if len(bestIncident.Evidence) >= maxEvidenceCount {
			e.metrics.IncCounter(ctx, "cortexops_correlation_evidence_overflow_total", map[string]string{"incident_id": bestIncident.ID})
			// Evict oldest event (index 0) to make room
			bestIncident.Evidence = bestIncident.Evidence[1:]
		}
		
		bestIncident.Evidence = append(bestIncident.Evidence, incoming)
		bestIncident.Fingerprints[fingerprint] = struct{}{}
		bestIncident.LastActivity = time.Now()
		
		if highestScore > float64(bestIncident.MaxScore.Value) {
			bestIncident.MaxScore.Value = float32(highestScore)
			bestIncident.MaxScore.Reasoning = bestReasoning
		}
		bestIncident.mu.Unlock()
		
		e.logger.Debug("Correlated event into existing incident", "incident_id", bestIncident.ID, "score", highestScore)
	} else {
		// Enforce global cap
		if len(e.incidents) >= 1000 { // maxGlobalIncidents
			e.metrics.IncCounter(ctx, "cortexops_correlation_global_incident_overflow_total", nil)
			e.logger.Warn("Global incident cap reached, dropping event")
			return nil
		}

		// 3. Create a new temporal window if it doesn't correlate strongly
		// In production, we filter out informational events here. Assuming non-informational for simplicity.
		newID := uuid.New().String()
		e.incidents[newID] = &ActiveIncident{
			ID:           newID,
			Evidence:     []*correlationv1.TelemetryEnvelope{incoming},
			Fingerprints: map[string]struct{}{fingerprint: {}},
			LastActivity: time.Now(),
			MaxScore: &correlationv1.CorrelationScore{
				Value:     0.0,
				Reasoning: []string{"Initial Event"},
			},
		}
		e.metrics.SetGauge(ctx, "cortexops_correlation_window_size", float64(len(e.incidents)), nil)
	}

	return nil
}

// FlushWindows isolates and finalizes incidents that have been inactive for > 2 minutes.
// This should be called via a periodic ticker in the main application loop.
func (e *Engine) FlushWindows(ctx context.Context) {
	e.mu.Lock()
	threshold := time.Now().Add(-2 * time.Minute)
	var closedIncidents []*ActiveIncident

	for id, inc := range e.incidents {
		inc.mu.Lock()
		if inc.LastActivity.Before(threshold) {
			closedIncidents = append(closedIncidents, inc)
			delete(e.incidents, id)
		}
		inc.mu.Unlock()
	}
	e.metrics.SetGauge(ctx, "cortexops_correlation_window_size", float64(len(e.incidents)), nil)
	e.mu.Unlock()

	// Process closed incidents outside the global lock
	for _, inc := range closedIncidents {
		inc.mu.Lock()
		id := inc.ID
		evidence := inc.Evidence
		maxScore := inc.MaxScore
		inc.mu.Unlock()

		chain, err := e.chainBuilder.Build(ctx, evidence)
		if err != nil {
			e.logger.Error("Failed to build causal chain", "incident_id", id, "error", err)
		}

		severity := e.deriveSeverity(len(evidence), maxScore.Value)

		incident := &correlationv1.CorrelatedIncident{
			IncidentId:  id,
			State:       correlationv1.IncidentState_DETECTED,
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
			Severity:    severity,
			Evidence:    evidence,
			CausalChain: chain,
			Confidence:  maxScore,
		}

		payload, err := proto.Marshal(incident)
		if err != nil {
			e.logger.Error("Failed to marshal correlated incident", "incident_id", id, "error", err)
			continue
		}

		if err := e.publisher.Publish(ctx, "cortex.incident.correlated", id, payload); err != nil {
			e.logger.Error("Failed to publish correlated incident", "incident_id", id, "error", err)
		}
		
		e.metrics.IncCounter(ctx, "cortexops_correlation_incidents_created_total", map[string]string{"severity": severity})
		e.logger.Info("Flushed correlated incident", "incident_id", id, "evidence_count", len(evidence), "severity", severity)
	}
}

func (e *Engine) deriveSeverity(evidenceCount int, maxScore float32) string {
	if evidenceCount > 10 || maxScore > 0.95 {
		return "CRITICAL"
	}
	if evidenceCount > 5 || maxScore > 0.85 {
		return "HIGH"
	}
	if evidenceCount > 2 || maxScore > 0.75 {
		return "MEDIUM"
	}
	return "LOW"
}
