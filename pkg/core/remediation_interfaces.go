package core

import (
	"context"

	apiv1 "github.com/shadow0vortex/cortexops/api/v1"
)

// RemediationExecutor provides the contract for dry-running, executing, and rolling back actions.
type RemediationExecutor interface {
	DryRun(ctx context.Context, action *apiv1.RemediationAction) (bool, error)
	Execute(ctx context.Context, action *apiv1.RemediationAction) error
	Rollback(ctx context.Context, action *apiv1.RemediationAction) error
	Verify(ctx context.Context, action *apiv1.RemediationAction) (bool, error)
}

// PolicyEngine evaluates deterministic governance rules.
type PolicyEngine interface {
	Evaluate(ctx context.Context, incident *apiv1.CorrelatedIncident, action *apiv1.RemediationAction) (*apiv1.PolicyDecision, error)
	CalculateRiskScore(ctx context.Context, incident *apiv1.CorrelatedIncident, action *apiv1.RemediationAction) (float32, error)
}

// ApprovalWorkflow handles asynchronous human-in-the-loop authorization.
type ApprovalWorkflow interface {
	RequestApproval(ctx context.Context, action *apiv1.RemediationAction) (*apiv1.ApprovalDecision, error)
}

// AuditStore abstracts the immutable storage of remediation state changes.
type AuditStore interface {
	Log(ctx context.Context, record *apiv1.AuditRecord) error
}
