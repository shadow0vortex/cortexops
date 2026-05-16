package temporal

import (
	"time"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// RemediationWorkflow orchestrates the safe execution lifecycle using Temporal durability.
func RemediationWorkflow(ctx workflow.Context, incident *correlationv1.CorrelatedIncident, action *remediationv1.RemediationAction) error {
	// Standard bounded retries for idempotent K8s API calls
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    3,
	}

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	var acts *Activities // Used for typed execution calling

	// 1. Policy Evaluation
	var policyDecision remediationv1.PolicyDecision
	err := workflow.ExecuteActivity(ctx, acts.EvaluatePolicy, incident, action).Get(ctx, &policyDecision)
	if err != nil {
		return err // Activity failure
	}
	if !policyDecision.Allowed {
		// Log Audit, fail workflow safely without panic
		workflow.GetLogger(ctx).Info("Remediation denied by policy", "reason", policyDecision.Reasoning)
		return nil
	}

	// 2. Human Approval (Long-running activity)
	approvalOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1, // Give human 1 hour to approve
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1, // Do not retry slack prompts
		},
	}
	approvalCtx := workflow.WithActivityOptions(ctx, approvalOpts)
	
	var approvalDecision remediationv1.ApprovalDecision
	err = workflow.ExecuteActivity(approvalCtx, acts.RequestApproval, action).Get(ctx, &approvalDecision)
	if err != nil || !approvalDecision.Approved {
		workflow.GetLogger(ctx).Info("Remediation rejected or timed out by human")
		return nil
	}

	// 3. Dry Run
	var dryRunSuccess bool
	err = workflow.ExecuteActivity(ctx, acts.DryRun, action).Get(ctx, &dryRunSuccess)
	if err != nil || !dryRunSuccess {
		workflow.GetLogger(ctx).Error("Dry run failed, aborting remediation")
		return nil
	}

	// 4. Execution
	err = workflow.ExecuteActivity(ctx, acts.Execute, action).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("Execution failed, triggering rollback")
		_ = workflow.ExecuteActivity(ctx, acts.Rollback, action).Get(ctx, nil)
		return err
	}

	// 5. Verification
	var verified bool
	err = workflow.ExecuteActivity(ctx, acts.Verify, action).Get(ctx, &verified)
	if err != nil || !verified {
		workflow.GetLogger(ctx).Error("Verification failed, triggering rollback")
		_ = workflow.ExecuteActivity(ctx, acts.Rollback, action).Get(ctx, nil)
		return err // Bubble error up to record workflow failure
	}

	workflow.GetLogger(ctx).Info("Remediation workflow completed successfully")
	return nil
}
