package graph

import (
	"context"
	"testing"

	topologyv1 "github.com/shadow0vortex/cortexops/api/v1"
)

func TestMemoryGraphStore_UpsertAndDeleteNode(t *testing.T) {
	store := NewMemoryGraphStore()
	ctx := context.Background()
	node := &topologyv1.TopologyNode{Id: "node-1", Name: "Node 1"}

	// Test Upsert
	err := store.UpsertNode(ctx, node)
	if err != nil {
		t.Fatalf("failed to upsert node: %v", err)
	}

	nodes, _ := store.ListNodes(ctx)
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}

	// Test Delete
	err = store.DeleteNode(ctx, "node-1")
	if err != nil {
		t.Fatalf("failed to delete node: %v", err)
	}

	nodes, _ = store.ListNodes(ctx)
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes after delete, got %d", len(nodes))
	}
}

func TestMemoryGraphStore_OrphanedEdgeCleanup(t *testing.T) {
	store := NewMemoryGraphStore()
	ctx := context.Background()

	store.UpsertNode(ctx, &topologyv1.TopologyNode{Id: "node-1"})
	store.UpsertNode(ctx, &topologyv1.TopologyNode{Id: "node-2"})
	store.UpsertEdge(ctx, "node-2", "node-1", topologyv1.EdgeType_ROUTES_TO)

	// Verify edge exists
	deps, _ := store.GetDependencies(ctx, "node-2")
	if len(deps) != 1 || deps[0] != "node-1" {
		t.Errorf("expected edge node-2 -> node-1")
	}

	// Delete target node
	store.DeleteNode(ctx, "node-1")

	// Verify inbound edge was cleaned up
	deps, _ = store.GetDependencies(ctx, "node-2")
	if len(deps) != 0 {
		t.Errorf("expected inbound edge to be cleaned up, still got %v", deps)
	}
}

func TestMemoryGraphStore_CalculateBlastRadius(t *testing.T) {
	store := NewMemoryGraphStore()
	ctx := context.Background()

	// Setup graph: 1 -> 2 -> 3
	store.UpsertNode(ctx, &topologyv1.TopologyNode{Id: "node-1"})
	store.UpsertNode(ctx, &topologyv1.TopologyNode{Id: "node-2"})
	store.UpsertNode(ctx, &topologyv1.TopologyNode{Id: "node-3"})
	store.UpsertEdge(ctx, "node-1", "node-2", topologyv1.EdgeType_ROUTES_TO)
	store.UpsertEdge(ctx, "node-2", "node-3", topologyv1.EdgeType_ROUTES_TO)

	impacted, err := store.CalculateBlastRadius(ctx, "node-1")
	if err != nil {
		t.Fatalf("failed to calculate blast radius: %v", err)
	}

	if len(impacted) != 2 {
		t.Errorf("expected blast radius size 2, got %d", len(impacted))
	}
}
