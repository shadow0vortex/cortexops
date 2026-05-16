package policy

import (
	"context"
	"fmt"
	"log/slog"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
)

// OPAEngine implements deterministic policy checking and risk scoring.
// In production, this would invoke the OPA Rego API. Here we stub the deterministic logic.
type OPAEngine struct {
	topology core.TopologyProvider
	metrics  core.MetricsRecorder
	logger   *slog.Logger
}

func NewOPAEngine(topology core.TopologyProvider, metrics core.MetricsRecorder, logger *slog.Logger) *OPAEngine {
	return &OPAEngine{
		topology: topology,
		metrics:  metrics,
		logger:   logger,
	}
}

// Evaluate applies hard boundaries. E.g., blocks "kube-system" operations.
func (e *OPAEngine) Evaluate(ctx context.Context, incident *correlationv1.CorrelatedIncident, action *remediationv1.RemediationAction) (*remediationv1.PolicyDecision, error) {
	if action.TargetNamespace == "kube-system" {
		e.metrics.IncCounter(ctx, "cortexops_remediation_policy_denials_total", map[string]string{"rule": "protected_namespace"})
		return &remediationv1.PolicyDecision{
			Allowed:         false,
			Reasoning:       "kube-system is immutable to automation",
			ViolatingRules:  []string{"DENY_KUBE_SYSTEM_MUTATION"},
		}, nil
	}

	// Validate action against allowed list
	allowed := false
	switch action.Type {
	case remediationv1.ActionType_POD_RESTART, remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART, remediationv1.ActionType_HORIZONTAL_SCALE:
		allowed = true
	}

	if !allowed {
		e.metrics.IncCounter(ctx, "cortexops_remediation_policy_denials_total", map[string]string{"rule": "disallowed_action_type"})
		return &remediationv1.PolicyDecision{
			Allowed:         false,
			Reasoning:       fmt.Sprintf("Action type %s is strictly forbidden in Phase 6", action.Type.String()),
			ViolatingRules:  []string{"DENY_UNAUTHORIZED_ACTION"},
		}, nil
	}

	return &remediationv1.PolicyDecision{
		Allowed:   true,
		Reasoning: "Passed all deterministic policy checks",
	}, nil
}

// CalculateRiskScore returns a deterministic [0.0, 1.0] float dictating approval gates.
func (e *OPAEngine) CalculateRiskScore(ctx context.Context, incident *correlationv1.CorrelatedIncident, action *remediationv1.RemediationAction) (float32, error) {
	var score float32 = 0.0

	// 1. Base Severity Contribution
	if incident.Severity == "CRITICAL" {
		score += 0.4
	} else if incident.Severity == "HIGH" {
		score += 0.2
	}

	// 2. Blast Radius Contribution (From Phase 3)
	if incident.BlastRadius != nil {
		impactDepth := incident.BlastRadius.TraversalDepth
		if impactDepth > 3 {
			score += 0.3
		} else if impactDepth > 1 {
			score += 0.1
		}
	}

	// 3. Action Criticality
	if action.Type == remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART {
		score += 0.2 // higher risk than single pod restart
	}

	if score > 1.0 {
		score = 1.0
	}

	return score, nil
}
