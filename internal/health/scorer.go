package health

import (
	"context"
	"sync"
	"time"

	"github.com/cortexops/cortexops/pkg/core"
)

// Subsystem represents a core component that provides health metrics.
type Subsystem string

const (
	Broker     Subsystem = "NATS_JETSTREAM"
	Topology   Subsystem = "TOPOLOGY_GRAPH"
	AI         Subsystem = "AI_INFERENCE"
	Temporal   Subsystem = "TEMPORAL_ORCHESTRATOR"
)

// PlatformHealth represents the aggregated degraded status of CortexOps.
type PlatformHealth struct {
	GlobalScore      float64 // 0.0 (Down) to 1.0 (Healthy)
	DegradedSystems  []Subsystem
	IsRemediationSafe bool
}

// Scorer calculates the deterministic health of the platform based on observability metrics.
type Scorer struct {
	mu            sync.RWMutex
	metrics       core.MetricsRecorder
	lastHeartbeat map[Subsystem]time.Time
}

func NewScorer(metrics core.MetricsRecorder) *Scorer {
	return &Scorer{
		metrics:       metrics,
		lastHeartbeat: make(map[Subsystem]time.Time),
	}
}

// ReportHeartbeat records an active ping from a subsystem.
func (s *Scorer) ReportHeartbeat(sys Subsystem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastHeartbeat[sys] = time.Now()
}

// Calculate interrogates heartbeats and bounded metrics to produce the health summary.
func (s *Scorer) Calculate(ctx context.Context) *PlatformHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()

	score := 1.0
	var degraded []Subsystem
	safeForRemediation := true

	threshold := time.Now().Add(-2 * time.Minute)

	// Evaluate Temporal / Broker (Critical for Remediation)
	if hb, ok := s.lastHeartbeat[Temporal]; !ok || hb.Before(threshold) {
		score -= 0.5
		degraded = append(degraded, Temporal)
		safeForRemediation = false
	}

	if hb, ok := s.lastHeartbeat[Broker]; !ok || hb.Before(threshold) {
		score -= 0.3
		degraded = append(degraded, Broker)
		safeForRemediation = false
	}

	// Evaluate AI (Non-critical, causes degraded RCA mode)
	if hb, ok := s.lastHeartbeat[AI]; !ok || hb.Before(threshold) {
		score -= 0.1
		degraded = append(degraded, AI)
	}

	if score < 0.0 {
		score = 0.0
	}

	// Emit to Prometheus
	s.metrics.SetGauge(ctx, "cortexops_platform_health_score", score, nil)

	return &PlatformHealth{
		GlobalScore:      score,
		DegradedSystems:  degraded,
		IsRemediationSafe: safeForRemediation,
	}
}
