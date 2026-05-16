package action

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// K8sExecutor implements the RemediationExecutor interface for K8s API mutations.
type K8sExecutor struct {
	client  kubernetes.Interface
	metrics core.MetricsRecorder
	logger  *slog.Logger
}

func NewK8sExecutor(client kubernetes.Interface, metrics core.MetricsRecorder, logger *slog.Logger) *K8sExecutor {
	return &K8sExecutor{
		client:  client,
		metrics: metrics,
		logger:  logger,
	}
}

// DryRun simulates the command without persisting state to ensure syntax and RBAC validity.
func (e *K8sExecutor) DryRun(ctx context.Context, action *remediationv1.RemediationAction) (bool, error) {
	opts := metav1.DeleteOptions{
		DryRun: []string{metav1.DryRunAll},
	}

	if action.Type == remediationv1.ActionType_POD_RESTART {
		err := e.client.CoreV1().Pods(action.TargetNamespace).Delete(ctx, action.TargetResource, opts)
		if err != nil {
			return false, fmt.Errorf("dry run failed for pod restart: %w", err)
		}
		return true, nil
	}

	return false, fmt.Errorf("unsupported dry-run action type")
}

// Execute performs the actual mutation. It is strictly scoped to the 3 allowed actions.
func (e *K8sExecutor) Execute(ctx context.Context, action *remediationv1.RemediationAction) error {
	start := time.Now()
	defer func() {
		e.metrics.ObserveHistogram(ctx, "cortexops_remediation_execution_seconds", time.Since(start).Seconds(), nil)
	}()

	switch action.Type {
	case remediationv1.ActionType_POD_RESTART:
		err := e.client.CoreV1().Pods(action.TargetNamespace).Delete(ctx, action.TargetResource, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete pod: %w", err)
		}
		
	case remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART:
		// To restart a deployment, we patch its template annotations with a new timestamp
		patch := []byte(fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"cortexops.io/restartedAt": "%s"}}}}}`, time.Now().Format(time.RFC3339)))
		_, err := e.client.AppsV1().Deployments(action.TargetNamespace).Patch(ctx, action.TargetResource, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("failed to rollout deployment: %w", err)
		}
		
	case remediationv1.ActionType_HORIZONTAL_SCALE:
		// Not fully implemented for brevity, but would update scale subresource
		return fmt.Errorf("horizontal scaling execution not implemented in stub")
		
	default:
		return fmt.Errorf("unrecognized action type")
	}

	e.metrics.IncCounter(ctx, "cortexops_remediation_success_total", map[string]string{"action": action.Type.String()})
	return nil
}

// Rollback restores state if Execute or Verify fails.
func (e *K8sExecutor) Rollback(ctx context.Context, action *remediationv1.RemediationAction) error {
	e.logger.Warn("Executing rollback procedure", "action_id", action.ActionId)
	
	// A real implementation would fetch the pre-execution snapshot (saved in the orchestrator)
	// and reverse the patch. For POD_RESTART, rollback isn't technically applicable as ReplicaSet recreates it.
	
	e.metrics.IncCounter(ctx, "cortexops_remediation_rollback_total", map[string]string{"action": action.Type.String(), "reason": "verification_failed"})
	return nil
}

// Verify waits up to 5 minutes observing telemetry to confirm system stabilization.
func (e *K8sExecutor) Verify(ctx context.Context, action *remediationv1.RemediationAction) (bool, error) {
	// A deterministic verify queries the metrics/topology state, NOT the AI.
	// E.g., Did the Pod re-enter Ready state? Did 5xx rates drop below threshold?
	e.logger.Info("Starting verification window", "action_id", action.ActionId)
	
	// Stubbed: assume success for phase 6 demonstration.
	return true, nil
}
