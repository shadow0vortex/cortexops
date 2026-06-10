package policy

import (
	"context"
	"log/slog"
	"testing"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
)

type mockMetrics struct{}
func (m *mockMetrics) IncCounter(ctx context.Context, name string, labels map[string]string) {}
func (m *mockMetrics) ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string) {}
func (m *mockMetrics) SetGauge(ctx context.Context, name string, value float64, labels map[string]string) {}


func TestOPAEngine_Evaluate(t *testing.T) {
	metrics := &mockMetrics{}
	logger := slog.Default()
	engine := NewOPAEngine(nil, metrics, logger)
	ctx := context.Background()

	incident := &correlationv1.CorrelatedIncident{
		Severity: "HIGH",
	}

	tests := []struct {
		name       string
		action     *remediationv1.RemediationAction
		wantAllow  bool
		wantReason string
	}{
		{
			name: "Allowed Pod Restart",
			action: &remediationv1.RemediationAction{
				Type:            remediationv1.ActionType_POD_RESTART,
				TargetNamespace: "default",
				TargetResource:  "api-123",
			},
			wantAllow: true,
		},
		{
			name: "Denied Kube-System",
			action: &remediationv1.RemediationAction{
				Type:            remediationv1.ActionType_POD_RESTART,
				TargetNamespace: "kube-system",
				TargetResource:  "coredns",
			},
			wantAllow: false,
		},
		{
			name: "Denied Invalid Action Type",
			action: &remediationv1.RemediationAction{
				Type:            remediationv1.ActionType_ACTION_TYPE_UNSPECIFIED,
				TargetNamespace: "default",
			},
			wantAllow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := engine.Evaluate(ctx, incident, tt.action)
			if err != nil {
				t.Fatalf("Evaluate returned error: %v", err)
			}
			if decision.Allowed != tt.wantAllow {
				t.Errorf("got Allowed=%v, want %v", decision.Allowed, tt.wantAllow)
			}
		})
	}
}
