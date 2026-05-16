package temporal

import (
	"context"

	correlationv1 "github.com/cortexops/cortexops/api/v1"
	remediationv1 "github.com/cortexops/cortexops/api/v1"
	"github.com/cortexops/cortexops/pkg/core"
)

// Activities encapsulates the deterministic boundary implementations for Temporal Workflows.
type Activities struct {
	policyEngine core.PolicyEngine
	approval     core.ApprovalWorkflow
	executor     core.RemediationExecutor
}

func NewActivities(policy core.PolicyEngine, approval core.ApprovalWorkflow, executor core.RemediationExecutor) *Activities {
	return &Activities{
		policyEngine: policy,
		approval:     approval,
		executor:     executor,
	}
}

// EvaluatePolicy is a Temporal Activity that runs OPA governance checks.
func (a *Activities) EvaluatePolicy(ctx context.Context, incident *correlationv1.CorrelatedIncident, action *remediationv1.RemediationAction) (*remediationv1.PolicyDecision, error) {
	return a.policyEngine.Evaluate(ctx, incident, action)
}

// RequestApproval is a Temporal Activity that blocks until human Slack approval or timeout.
func (a *Activities) RequestApproval(ctx context.Context, action *remediationv1.RemediationAction) (*remediationv1.ApprovalDecision, error) {
	return a.approval.RequestApproval(ctx, action)
}

// DryRun executes the simulated K8s patch.
func (a *Activities) DryRun(ctx context.Context, action *remediationv1.RemediationAction) (bool, error) {
	return a.executor.DryRun(ctx, action)
}

// Execute performs the actual infrastructure mutation.
func (a *Activities) Execute(ctx context.Context, action *remediationv1.RemediationAction) error {
	return a.executor.Execute(ctx, action)
}

// Rollback triggers the reversion of the state.
func (a *Activities) Rollback(ctx context.Context, action *remediationv1.RemediationAction) error {
	return a.executor.Rollback(ctx, action)
}

// Verify loops telemetry bounds.
func (a *Activities) Verify(ctx context.Context, action *remediationv1.RemediationAction) (bool, error) {
	return a.executor.Verify(ctx, action)
}
