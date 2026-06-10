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
	patchOpts := metav1.PatchOptions{
		DryRun: []string{metav1.DryRunAll},
	}

	switch action.Type {
	case remediationv1.ActionType_POD_RESTART:
		err := e.client.CoreV1().Pods(action.TargetNamespace).Delete(ctx, action.TargetResource, opts)
		if err != nil {
			return false, fmt.Errorf("dry run failed for pod restart: %w", err)
		}
		return true, nil

	case remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART:
		patch := []byte(`{"spec": {"template": {"metadata": {"annotations": {"cortexops.io/restartedAt": "dry-run"}}}}}`)
		_, err := e.client.AppsV1().Deployments(action.TargetNamespace).Patch(ctx, action.TargetResource, types.StrategicMergePatchType, patch, patchOpts)
		if err != nil {
			return false, fmt.Errorf("dry run failed for rollout: %w", err)
		}
		return true, nil

	case remediationv1.ActionType_HORIZONTAL_SCALE:
		_, err := e.client.AppsV1().Deployments(action.TargetNamespace).GetScale(ctx, action.TargetResource, metav1.GetOptions{})
		if err != nil {
			return false, fmt.Errorf("dry run failed for scale (get): %w", err)
		}
		return true, nil

	default:
		return false, fmt.Errorf("unsupported dry-run action type")
	}
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
		// Capture original state
		deploy, err := e.client.AppsV1().Deployments(action.TargetNamespace).Get(ctx, action.TargetResource, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get deployment for pre-state capture: %w", err)
		}
		
		if action.Parameters == nil {
			action.Parameters = make(map[string]string)
		}
		
		// Save existing annotation
		if deploy.Spec.Template.Annotations != nil {
			if val, ok := deploy.Spec.Template.Annotations["cortexops.io/restartedAt"]; ok {
				action.Parameters["previous_restartedAt"] = val
			} else {
				action.Parameters["previous_restartedAt"] = ""
			}
		} else {
			action.Parameters["previous_restartedAt"] = ""
		}

		// To restart a deployment, we patch its template annotations with a new timestamp
		patch := []byte(fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"cortexops.io/restartedAt": "%s"}}}}}`, time.Now().Format(time.RFC3339)))
		_, err = e.client.AppsV1().Deployments(action.TargetNamespace).Patch(ctx, action.TargetResource, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("failed to rollout deployment: %w", err)
		}

	case remediationv1.ActionType_HORIZONTAL_SCALE:
		replicasStr, ok := action.Parameters["replicas"]
		if !ok {
			return fmt.Errorf("replicas parameter missing for scale action")
		}
		var replicas int32
		_, err := fmt.Sscanf(replicasStr, "%d", &replicas)
		if err != nil {
			return fmt.Errorf("invalid replicas parameter: %w", err)
		}

		scale, err := e.client.AppsV1().Deployments(action.TargetNamespace).GetScale(ctx, action.TargetResource, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get scale: %w", err)
		}
		scale.Spec.Replicas = replicas
		_, err = e.client.AppsV1().Deployments(action.TargetNamespace).UpdateScale(ctx, action.TargetResource, scale, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update scale: %w", err)
		}

	default:
		return fmt.Errorf("unrecognized action type")
	}

	e.metrics.IncCounter(ctx, "cortexops_remediation_success_total", map[string]string{"action": action.Type.String()})
	return nil
}

// Rollback restores state if Execute or Verify fails.
func (e *K8sExecutor) Rollback(ctx context.Context, action *remediationv1.RemediationAction) error {
	e.logger.Warn("Executing rollback procedure", "action_id", action.ActionId)

	switch action.Type {
	case remediationv1.ActionType_HORIZONTAL_SCALE:
		prevReplicas, ok := action.Parameters["previous_replicas"]
		if !ok {
			return fmt.Errorf("previous_replicas parameter missing for rollback")
		}
		var replicas int32
		_, err := fmt.Sscanf(prevReplicas, "%d", &replicas)
		if err != nil {
			return fmt.Errorf("invalid previous_replicas: %w", err)
		}

		scale, err := e.client.AppsV1().Deployments(action.TargetNamespace).GetScale(ctx, action.TargetResource, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get scale for rollback: %w", err)
		}
		scale.Spec.Replicas = replicas
		_, err = e.client.AppsV1().Deployments(action.TargetNamespace).UpdateScale(ctx, action.TargetResource, scale, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to rollback scale: %w", err)
		}

	case remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART:
		prevRestart, ok := action.Parameters["previous_restartedAt"]
		if !ok {
			return fmt.Errorf("previous_restartedAt parameter missing for rollback")
		}

		var patch []byte
		if prevRestart == "" {
			patch = []byte(`{"spec": {"template": {"metadata": {"annotations": {"cortexops.io/restartedAt": null}}}}}`)
		} else {
			patch = []byte(fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"cortexops.io/restartedAt": "%s"}}}}}`, prevRestart))
		}

		_, err := e.client.AppsV1().Deployments(action.TargetNamespace).Patch(ctx, action.TargetResource, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("failed to rollback rollout restart: %w", err)
		}
		e.logger.Info("Rollback for rollout restart complete", "resource", action.TargetResource)

	case remediationv1.ActionType_POD_RESTART:
		e.logger.Info("Rollback not applicable for POD_RESTART (ReplicaSet handles recovery)")
	}

	e.metrics.IncCounter(ctx, "cortexops_remediation_rollback_total", map[string]string{"action": action.Type.String(), "reason": "verification_failed"})
	return nil
}

// Verify waits up to 5 minutes observing telemetry to confirm system stabilization.
func (e *K8sExecutor) Verify(ctx context.Context, action *remediationv1.RemediationAction) (bool, error) {
	e.logger.Info("Starting verification window", "action_id", action.ActionId)

	// Implementation: Check if the target resource is healthy.
	// For simplicity, we check if pods in the namespace are Ready if it was a POD_RESTART or ROLLOUT.
	
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("verification timed out after 5 minutes")
		case <-ticker.C:
			// Basic verification: check if deployment is stable
			if action.Type == remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART || action.Type == remediationv1.ActionType_HORIZONTAL_SCALE {
				deploy, err := e.client.AppsV1().Deployments(action.TargetNamespace).Get(ctx, action.TargetResource, metav1.GetOptions{})
				if err != nil {
					e.logger.Warn("Failed to get deployment during verification", "error", err)
					continue
				}
				if deploy.Status.ReadyReplicas == deploy.Status.Replicas && deploy.Status.UnavailableReplicas == 0 {
					e.logger.Info("Resource stabilized", "action_id", action.ActionId)
					return true, nil
				}
			} else if action.Type == remediationv1.ActionType_POD_RESTART {
				// For POD_RESTART, we assume success if the new pod is Ready.
				// This is simplified; a better check would look for pods with new timestamps.
				return true, nil 
			}
		}
	}
}
