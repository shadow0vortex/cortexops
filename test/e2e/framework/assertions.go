package framework

import (
	"context"
	"time"

	"github.com/shadow0vortex/cortexops/pkg/core"
)

// Assertions provides a suite of deterministic, replay-safe checks.
type Assertions struct {
	harness *Harness
	audit   core.AuditStore // Interface to check DB state
}

func NewAssertions(h *Harness, auditStore core.AuditStore) *Assertions {
	return &Assertions{
		harness: h,
		audit:   auditStore,
	}
}

// AssertTopologyDrift checks if a deleted pod is correctly purged from the graph within threshold.
func (a *Assertions) AssertTopologyDrift(nodeID string, timeout time.Duration, graph core.TopologyProvider) {
	ctx, cancel := context.WithTimeout(a.harness.Ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.harness.T.Fatalf("Invariant violated: Node %s remained in topology graph after %v", nodeID, timeout)
		case <-ticker.C:
			// Attempt to fetch dependencies. If it returns an error or empty, the node is pruned.
			deps, err := graph.GetDependencies(ctx, nodeID)
			if err != nil || len(deps) == 0 {
				return // Success: Invariant holds
			}
		}
	}
}

// AssertReplayIdempotency checks that re-ingesting events doesn't spawn duplicate incidents.
func (a *Assertions) AssertReplayIdempotency(incidentID string, expectedCount int, queryFunc func(string) int) {
	// A naive assertion expecting exactly `expectedCount` records in the DB
	count := queryFunc(incidentID)
	if count != expectedCount {
		a.harness.T.Fatalf("Replay invariant violated: expected %d incident records, found %d", expectedCount, count)
	}
}
