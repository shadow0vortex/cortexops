package heuristics

import (
	"context"
	"math"

	eventsv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
)

// Scorer provides deterministic correlation logic between two telemetry events.
type Scorer struct {
	topology core.TopologyProvider
}

func NewScorer(topology core.TopologyProvider) *Scorer {
	return &Scorer{topology: topology}
}

// Score determines the causal confidence between an incoming event (e2) and an existing event (e1).
func (s *Scorer) Score(ctx context.Context, e1, e2 *eventsv1.TelemetryEnvelope) (float64, []string) {
	var score float64
	var reasoning []string

	// 1. Trace ID Match (Absolute Certainty)
	if e1.TraceContext != nil && e2.TraceContext != nil {
		if e1.TraceContext.TraceId != "" && e1.TraceContext.TraceId == e2.TraceContext.TraceId {
			return 1.0, []string{"TraceID match"}
		}
	}

	// 2. Temporal Proximity
	// Compare protobuf timestamps. If within 30 seconds, add 0.2
	if e1.Timestamp != nil && e2.Timestamp != nil {
		diff := math.Abs(float64(e1.Timestamp.Seconds - e2.Timestamp.Seconds))
		if diff <= 30 {
			score += 0.2
			reasoning = append(reasoning, "Temporal proximity (<=30s)")
		}
	}

	// 3. Namespace Affinity (If both are K8s events)
	k1 := e1.GetK8SEvent()
	k2 := e2.GetK8SEvent()

	if k1 != nil && k2 != nil {
		if k1.Namespace == k2.Namespace {
			score += 0.1
			reasoning = append(reasoning, "Namespace affinity")
		}

		// 4. Topology Affinity (Requires topology graph lookup)
		// Assuming NodeIDs are formatted as "kind/namespace/name"
		n1ID := k1.ResourceKind + "/" + k1.Namespace + "/" + k1.ResourceName
		n2ID := k2.ResourceKind + "/" + k2.Namespace + "/" + k2.ResourceName

		deps, _ := s.topology.GetDependencies(ctx, n1ID)
		isAdjacent := false
		for _, dep := range deps {
			if dep == n2ID {
				isAdjacent = true
				break
			}
		}

		// Also check the reverse direction
		if !isAdjacent {
			depsRev, _ := s.topology.GetDependencies(ctx, n2ID)
			for _, dep := range depsRev {
				if dep == n1ID {
					isAdjacent = true
					break
				}
			}
		}

		if isAdjacent {
			score += 0.6
			reasoning = append(reasoning, "Topology affinity (1st-degree edge)")
		}
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score, reasoning
}
