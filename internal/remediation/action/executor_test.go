package action

import (
	"context"
	"log/slog"
	"os"
	"testing"

	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// MockMetricsRecorder implements core.MetricsRecorder for testing.
type MockMetricsRecorder struct{}
func (m *MockMetricsRecorder) IncCounter(ctx context.Context, name string, labels map[string]string) {}
func (m *MockMetricsRecorder) SetGauge(ctx context.Context, name string, value float64, labels map[string]string) {}
func (m *MockMetricsRecorder) ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string) {}

func TestK8sExecutor_DryRun(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"},
	})
	executor := NewK8sExecutor(client, &MockMetricsRecorder{}, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ctx := context.Background()

	action := &remediationv1.RemediationAction{
		Type:            remediationv1.ActionType_POD_RESTART,
		TargetNamespace: "default",
		TargetResource:  "test-pod",
	}

	success, err := executor.DryRun(ctx, action)
	if err != nil {
		t.Fatalf("DryRun failed: %v", err)
	}
	if !success {
		t.Errorf("expected DryRun success")
	}
}

func TestK8sExecutor_Rollback(t *testing.T) {
	client := fake.NewSimpleClientset()
	deploy := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "test-deploy", Namespace: "default"},
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"cortexops.io/restartedAt": "old-time"},
				},
			},
		},
	}
	client.AppsV1().Deployments("default").Create(context.Background(), deploy, metav1.CreateOptions{})
	
	executor := NewK8sExecutor(client, &MockMetricsRecorder{}, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	ctx := context.Background()

	action := &remediationv1.RemediationAction{
		Type:            remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART,
		TargetNamespace: "default",
		TargetResource:  "test-deploy",
		Parameters:      map[string]string{"previous_restartedAt": "old-time"},
	}

	err := executor.Rollback(ctx, action)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	updated, _ := client.AppsV1().Deployments("default").Get(ctx, "test-deploy", metav1.GetOptions{})
	if updated.Spec.Template.Annotations["cortexops.io/restartedAt"] != "old-time" {
		t.Errorf("expected 'old-time' annotation after rollback, got %s", updated.Spec.Template.Annotations["cortexops.io/restartedAt"])
	}
}

func int32Ptr(i int32) *int32 { return &i }
