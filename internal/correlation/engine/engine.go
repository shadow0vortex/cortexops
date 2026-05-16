package engine

import (
	"context"
	"log/slog"
	"sync"
	"time"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	eventsv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/correlation/causal"
	"github.com/shadow0vortex/cortexops/internal/correlation/heuristics"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ActiveIncident represents a mutable temporal window of correlated events.
type ActiveIncident struct {
	mu           sync.Mutex
	ID           string
	Evidence     []*eventsv1.TelemetryEnvelope
	LastActivity time.Time
	MaxScore     *correlationv1.CorrelationScore
}

// Engine ingests raw telemetry, applies deterministic scoring, and publishes CorrelatedIncidents.
type Engine struct {
	mu        sync.RWMutex
	incidents map[string]*ActiveIncident

	scorer       *heuristics.Scorer
	chainBuilder *causal.ChainBuilder
	publisher    core.Publisher
	metrics      core.MetricsRecorder
	logger       *slog.Logger
}

func NewEngine(scorer *heuristics.Scorer, chainBuilder *causal.ChainBuilder, pub core.Publisher, metrics core.MetricsRecorder, logger *slog.Logger) *Engine {
	return &Engine{
		incidents:    make(map[string]*ActiveIncident),
		scorer:       scorer,
		chainBuilder: chainBuilder,
		publisher:    pub,
		metrics:      metrics,
		logger:       logger,
	}
}

// ProcessEvent applies temporal and topological correlation logic to an incoming telemetry envelope.
func (e *Engine) ProcessEvent(ctx context.Context, incoming *eventsv1.TelemetryEnvelope) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var bestIncident *ActiveIncident
	var highestScore float64
	var bestReasoning []string

	// 1. Evaluate incoming event against all active temporal windows
	for _, inc := range e.incidents {
		inc.mu.Lock()
		// Compare against the most recent event in the incident bucket
		latestEvent := inc.Evidence[len(inc.Evidence)-1]
		
		score, reasoning := e.scorer.Score(ctx, latestEvent, incoming)
		if score > highestScore {
			highestScore = score
			bestReasoning = reasoning
			bestIncident = inc
		}
		inc.mu.Unlock()
	}

	// 2. Correlation Threshold logic
	if highestScore >= 0.8 && bestIncident != nil {
		bestIncident.mu.Lock()
		bestIncident.Evidence = append(bestIncident.Evidence, incoming)
		bestIncident.LastActivity = time.Now()
		
		if highestScore > bestIncident.MaxScore.Value {
			bestIncident.MaxScore.Value = float32(highestScore)
			bestIncident.MaxScore.Reasoning = bestReasoning
		}
		bestIncident.mu.Unlock()
		
		e.logger.Debug("Correlated event into existing incident", "incident_id", bestIncident.ID, "score", highestScore)
	} else {
		// 3. Create a new temporal window if it doesn't correlate strongly
		// In production, we filter out informational events here. Assuming non-informational for simplicity.
		newID := uuid.New().String()
		e.incidents[newID] = &ActiveIncident{
			ID:           newID,
			Evidence:     []*eventsv1.TelemetryEnvelope{incoming},
			LastActivity: time.Now(),
			MaxScore: &correlationv1.CorrelationScore{
				Value:     0.0,
				Reasoning: []string{"Initial Event"},
			},
		}
		e.metrics.SetGauge(ctx, "cortexops_correlation_window_size", float64(len(e.incidents)), nil)
	}

	return nil
}

// FlushWindows isolates and finalizes incidents that have been inactive for > 2 minutes.
// This should be called via a periodic ticker in the main application loop.
func (e *Engine) FlushWindows(ctx context.Context) {
	e.mu.Lock()
	defer e.mu.Unlock()

	threshold := time.Now().Add(-2 * time.Minute)

	for id, inc := range e.incidents {
		inc.mu.Lock()
		if inc.LastActivity.Before(threshold) {
			// Window is closed, build causal chain and publish
			chain, err := e.chainBuilder.Build(ctx, inc.Evidence)
			if err != nil {
				e.logger.Error("Failed to build causal chain", "incident_id", id, "error", err)
			}

			// In a real system, BlastRadius is calculated here via the core.BlastRadiusEvaluator
			incident := &correlationv1.CorrelatedIncident{
				IncidentId:  id,
				State:       eventsv1.IncidentState_STATE_DETECTED,
				CreatedAt:   timestamppb.Now(),
				UpdatedAt:   timestamppb.Now(),
				Severity:    "HIGH", // Derived logically in a full implementation
				Evidence:    inc.Evidence,
				CausalChain: chain,
				Confidence:  inc.MaxScore,
			}

			payload, _ := proto.Marshal(incident)
			e.publisher.Publish(ctx, "cortex.incident.correlated", id, payload)
			
			e.metrics.IncCounter(ctx, "cortexops_correlation_incidents_created_total", map[string]string{"severity": "HIGH"})
			e.logger.Info("Flushed correlated incident", "incident_id", id, "evidence_count", len(inc.Evidence))
			
			// Remove from active windows
			delete(e.incidents, id)
		}
		inc.mu.Unlock()
	}
	e.metrics.SetGauge(ctx, "cortexops_correlation_window_size", float64(len(e.incidents)), nil)
}
