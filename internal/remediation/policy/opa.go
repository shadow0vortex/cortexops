package policy

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/open-policy-agent/opa/rego"
	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
)

//go:embed remediation.rego
var defaultPolicy string

type OPAEngine struct {
	topology      core.TopologyProvider
	metrics       core.MetricsRecorder
	logger        *slog.Logger
	preparedQuery rego.PreparedEvalQuery
}

func NewOPAEngine(topology core.TopologyProvider, metrics core.MetricsRecorder, logger *slog.Logger) *OPAEngine {
	// In a real system, you would parse the policy once at startup and handle errors.
	// For resilience, if preparation fails, we panic here to prevent starting with broken policies.
	query, err := rego.New(
		rego.Query("data.cortexops.remediation"),
		rego.Module("remediation.rego", defaultPolicy),
	).PrepareForEval(context.Background())

	if err != nil {
		panic(fmt.Sprintf("failed to prepare OPA rego query: %v", err))
	}

	return &OPAEngine{
		topology:      topology,
		metrics:       metrics,
		logger:        logger,
		preparedQuery: query,
	}
}

func (e *OPAEngine) Evaluate(ctx context.Context, incident *correlationv1.CorrelatedIncident, action *remediationv1.RemediationAction) (*remediationv1.PolicyDecision, error) {
	start := time.Now()
	defer func() {
		e.metrics.ObserveHistogram(ctx, "cortexops_remediation_policy_evaluation_seconds", time.Since(start).Seconds(), nil)
	}()
	// Fail-closed guard: if the calling context is already cancelled (e.g., upstream
	// timeout, orchestrator abort), refuse to evaluate and deny by default. OPA's
	// PreparedEvalQuery.Eval() may complete too fast to observe cancellation internally.
	if err := ctx.Err(); err != nil {
		e.logger.Error("OPA evaluation skipped: context already cancelled", "error", err)
		e.metrics.IncCounter(ctx, "cortexops_remediation_policy_errors_total", nil)
		return &remediationv1.PolicyDecision{
			Allowed:        false,
			Reasoning:      fmt.Sprintf("Policy engine context cancelled: %v", err),
			ViolatingRules: []string{"CONTEXT_CANCELLED"},
		}, nil
	}

	input := map[string]interface{}{
		"action": map[string]interface{}{
			"type":             action.Type.String(),
			"target_namespace": action.TargetNamespace,
			"target_resource":  action.TargetResource,
		},
		"incident": map[string]interface{}{
			"severity": incident.Severity,
		},
	}

	results, err := e.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		e.logger.Error("OPA evaluation failed", "error", err)
		e.metrics.IncCounter(ctx, "cortexops_remediation_policy_errors_total", nil)
		// Fail closed
		return &remediationv1.PolicyDecision{
			Allowed:         false,
			Reasoning:       fmt.Sprintf("Policy engine error: %v", err),
			ViolatingRules:  []string{"INTERNAL_ENGINE_ERROR"},
		}, nil
	}

	if len(results) == 0 {
		return &remediationv1.PolicyDecision{
			Allowed:   false,
			Reasoning: "No policy results returned",
		}, nil
	}

	resMap, ok := results[0].Expressions[0].Value.(map[string]interface{})
	if !ok {
		return &remediationv1.PolicyDecision{
			Allowed:   false,
			Reasoning: "Invalid policy result format",
		}, nil
	}

	allowed, _ := resMap["allowed"].(bool)
	deniesInterface, _ := resMap["deny"].([]interface{})
	var denyReasons []string
	for _, d := range deniesInterface {
		if ds, ok := d.(string); ok {
			denyReasons = append(denyReasons, ds)
		}
	}

	if !allowed {
		e.metrics.IncCounter(ctx, "cortexops_remediation_policy_denials_total", map[string]string{"rule": "opa_deny"})
		reasoning := "Policy denied"
		if len(denyReasons) > 0 {
			reasoning = denyReasons[0]
		}
		
		e.logger.Warn("OPA Denied action", "action_id", action.ActionId, "reasons", denyReasons)
		return &remediationv1.PolicyDecision{
			Allowed:        false,
			Reasoning:      reasoning,
			ViolatingRules: denyReasons,
		}, nil
	}

	e.logger.Info("OPA Approved action", "action_id", action.ActionId, "latency_ms", time.Since(start).Milliseconds())
	return &remediationv1.PolicyDecision{
		Allowed:   true,
		Reasoning: "Passed all deterministic OPA policy checks",
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
