package approval

import (
	"context"
	"log/slog"
	"time"

	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeterministicApprovalWorkflow implements the ApprovalWorkflow interface with deterministic mock behavior.
type DeterministicApprovalWorkflow struct {
	audit   core.AuditStore
	metrics core.MetricsRecorder
	logger  *slog.Logger
}

func NewDeterministicApprovalWorkflow(audit core.AuditStore, metrics core.MetricsRecorder, logger *slog.Logger) *DeterministicApprovalWorkflow {
	return &DeterministicApprovalWorkflow{
		audit:   audit,
		metrics: metrics,
		logger:  logger,
	}
}

// RequestApproval simulates a deterministic timeout/approval process.
func (w *DeterministicApprovalWorkflow) RequestApproval(ctx context.Context, action *remediationv1.RemediationAction) (*remediationv1.ApprovalDecision, error) {
	w.logger.Info("Requesting deterministic mock approval", "action_id", action.ActionId, "risk_score", action.RiskScore)

	// Persist Audit Record mapping the state transition to APPROVAL_PENDING
	_ = w.audit.Log(ctx, &remediationv1.AuditRecord{
		AuditId:    uuid.New().String(),
		IncidentId: action.IncidentId,
		ActionId:   action.ActionId,
		FromState:  remediationv1.ActionState_POLICY_EVALUATING,
		ToState:    remediationv1.ActionState_APPROVAL_PENDING,
		Actor:      "ApprovalWorkflow",
		Reasoning:  "Risk score breached threshold, awaiting human approval",
		Timestamp:  timestamppb.Now(),
	})

	// In production, this would block/suspend the state machine or use a callback webhook.
	// We retain this deterministic mock to preserve replay safety and operational simplicity.
	select {
	case <-ctx.Done():
		return &remediationv1.ApprovalDecision{
			Approved:  false,
			Comments:  "Approval request timed out",
			DecidedAt: timestamppb.Now(),
		}, nil
	case <-time.After(2 * time.Second): // Mocking immediate approval for demo
		return &remediationv1.ApprovalDecision{
			Approved:   true,
			ApproverId: "Deterministic_System",
			Comments:   "Approved via Deterministic Mock Workflow",
			DecidedAt:  timestamppb.Now(),
		}, nil
	}
}
