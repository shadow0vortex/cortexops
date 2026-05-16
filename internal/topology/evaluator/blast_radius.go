package evaluator

import (
	"context"
	"log/slog"
	"time"

	topologyv1 "github.com/cortexops/cortexops/api/v1"
	"github.com/cortexops/cortexops/pkg/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// BlastRadiusEvaluator calculates the propagation of failure based on graph dependencies.
type BlastRadiusEvaluator struct {
	store   core.TopologyProvider // Interface abstracting GraphStore
	metrics core.MetricsRecorder
	logger  *slog.Logger
}

func NewBlastRadiusEvaluator(store core.TopologyProvider, metrics core.MetricsRecorder, logger *slog.Logger) *BlastRadiusEvaluator {
	return &BlastRadiusEvaluator{
		store:   store,
		metrics: metrics,
		logger:  logger,
	}
}

// Evaluate traverses the graph up to maxDepth to find all impacted downstream nodes.
func (e *BlastRadiusEvaluator) Evaluate(ctx context.Context, sourceNodeID string, maxDepth int) (*topologyv1.BlastRadiusResult, error) {
	start := time.Now()
	defer func() {
		e.metrics.ObserveHistogram(ctx, "cortexops_topology_blast_radius_traversal_duration_seconds", time.Since(start).Seconds(), nil)
	}()

	visited := make(map[string]bool)
	var impacted []string

	type queueItem struct {
		id    string
		depth int
	}
	queue := []queueItem{{id: sourceNodeID, depth: 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.id] {
			continue
		}
		visited[current.id] = true

		if current.id != sourceNodeID {
			impacted = append(impacted, current.id)
		}

		if current.depth >= maxDepth {
			continue
		}

		deps, err := e.store.GetDependencies(ctx, current.id)
		if err != nil {
			e.logger.Error("Failed to fetch dependencies during traversal", "node", current.id, "error", err)
			continue
		}

		for _, dep := range deps {
			// e.store.GetDependencies from core.TopologyProvider returns []string
			if !visited[dep] {
				queue = append(queue, queueItem{id: dep, depth: current.depth + 1})
			}
		}
	}

	return &topologyv1.BlastRadiusResult{
		SourceNodeId:    sourceNodeID,
		ImpactedNodeIds: impacted,
		TraversalDepth:  int32(maxDepth),
		EvaluatedAt:     timestamppb.Now(),
	}, nil
}
