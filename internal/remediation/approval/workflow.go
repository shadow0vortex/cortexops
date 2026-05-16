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

// SlackWorkflow implements the ApprovalWorkflow interface requiring Human-in-the-Loop.
type SlackWorkflow struct {
	audit   core.AuditStore
	metrics core.MetricsRecorder
	logger  *slog.Logger
}

func NewSlackWorkflow(audit core.AuditStore, metrics core.MetricsRecorder, logger *slog.Logger) *SlackWorkflow {
	return &SlackWorkflow{
		audit:   audit,
		metrics: metrics,
		logger:  logger,
	}
}

// RequestApproval simulates sending a Slack block kit message and waiting for a callback.
func (w *SlackWorkflow) RequestApproval(ctx context.Context, action *remediationv1.RemediationAction) (*remediationv1.ApprovalDecision, error) {
	w.logger.Info("Requesting human approval via Slack", "action_id", action.ActionId, "risk_score", action.RiskScore)

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
	// For this phase, we mock a deterministic timeout/approval.
	
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
			ApproverId: "U123456 (SRE On-Call)",
			Comments:   "Approved via Slack",
			DecidedAt:  timestamppb.Now(),
		}, nil
	}
}
