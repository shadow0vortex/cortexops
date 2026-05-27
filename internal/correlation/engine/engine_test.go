package engine

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/correlation/causal"
	"github.com/shadow0vortex/cortexops/internal/correlation/heuristics"
)

// MockPublisher implements core.Publisher for testing.
type MockPublisher struct {
	Published map[string][]byte
}
func (m *MockPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	if m.Published == nil { m.Published = make(map[string][]byte) }
	m.Published[topic] = payload
	return nil
}

// MockMetricsRecorder implements core.MetricsRecorder for testing.
type MockMetricsRecorder struct{}
func (m *MockMetricsRecorder) IncCounter(ctx context.Context, name string, labels map[string]string) {}
func (m *MockMetricsRecorder) SetGauge(ctx context.Context, name string, value float64, labels map[string]string) {}
func (m *MockMetricsRecorder) ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string) {}

func TestEngine_ProcessAndFlush(t *testing.T) {
	scorer := heuristics.NewScorer(nil)
	chainBuilder := causal.NewChainBuilder(nil)
	publisher := &MockPublisher{}
	engine := NewEngine(scorer, chainBuilder, publisher, &MockMetricsRecorder{}, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	
	ctx := context.Background()
	event := &correlationv1.TelemetryEnvelope{
		Source: "node-1",
		Payload: &correlationv1.TelemetryEnvelope_K8SEvent{
			K8SEvent: &correlationv1.K8SEventMetadata{Message: "CPU Spike"},
		},
	}

	// 1. Process Event
	err := engine.ProcessEvent(ctx, event)
	if err != nil {
		t.Fatalf("ProcessEvent failed: %v", err)
	}

	if len(engine.incidents) != 1 {
		t.Errorf("expected 1 active incident, got %d", len(engine.incidents))
	}

	// 2. Force Flush (manually set activity time back)
	for _, inc := range engine.incidents {
		inc.LastActivity = time.Now().Add(-5 * time.Minute)
	}

	engine.FlushWindows(ctx)

	if len(engine.incidents) != 0 {
		t.Errorf("expected 0 active incidents after flush, got %d", len(engine.incidents))
	}

	if _, ok := publisher.Published["cortex.incident.correlated"]; !ok {
		t.Errorf("expected incident to be published")
	}
}

func TestEngine_DeriveSeverity(t *testing.T) {
	engine := &Engine{}
	
	if s := engine.deriveSeverity(11, 0.96); s != "CRITICAL" {
		t.Errorf("expected CRITICAL, got %s", s)
	}
	if s := engine.deriveSeverity(6, 0.86); s != "HIGH" {
		t.Errorf("expected HIGH, got %s", s)
	}
	if s := engine.deriveSeverity(3, 0.76); s != "MEDIUM" {
		t.Errorf("expected MEDIUM, got %s", s)
	}
	if s := engine.deriveSeverity(1, 0.1); s != "LOW" {
		t.Errorf("expected LOW, got %s", s)
	}
}
