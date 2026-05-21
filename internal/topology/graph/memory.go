package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	topologyv1 "github.com/shadow0vortex/cortexops/api/v1"
)

// MemoryGraphStore implements a thread-safe, in-memory directed graph.
type MemoryGraphStore struct {
	mu    sync.RWMutex
	nodes map[string]*topologyv1.TopologyNode
	edges map[string][]*topologyv1.TopologyEdge // key is source_id
}

// NewMemoryGraphStore initializes the graph.
func NewMemoryGraphStore() *MemoryGraphStore {
	return &MemoryGraphStore{
		nodes: make(map[string]*topologyv1.TopologyNode),
		edges: make(map[string][]*topologyv1.TopologyEdge),
	}
}

// UpsertNode idempotently inserts or updates a node.
func (g *MemoryGraphStore) UpsertNode(ctx context.Context, node *topologyv1.TopologyNode) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if node == nil || node.Id == "" {
		return fmt.Errorf("invalid node: id is required")
	}
	
	g.nodes[node.Id] = node
	return nil
}

// DeleteNode removes a node and all its outward edges. (Inbound edges will be orphaned until cleanup).
func (g *MemoryGraphStore) DeleteNode(ctx context.Context, nodeID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.nodes, nodeID)
	delete(g.edges, nodeID)
	
	// A robust implementation requires scanning to remove inbound edges to `nodeID` as well,
	// handled asynchronously by the orphaned-edge cleanup job.
	return nil
}

// UpsertEdge idempotently adds a directional relationship.
func (g *MemoryGraphStore) UpsertEdge(ctx context.Context, sourceID, targetID string, relationship topologyv1.EdgeType) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check if edge exists to avoid duplication
	for _, edge := range g.edges[sourceID] {
		if edge.TargetId == targetID && edge.Relationship == relationship {
			return nil
		}
	}

	edge := &topologyv1.TopologyEdge{
		SourceId:     sourceID,
		TargetId:     targetID,
		Relationship: relationship,
	}
	
	g.edges[sourceID] = append(g.edges[sourceID], edge)
	return nil
}

// GetDependencies returns all node IDs that the given node has outward edges to.
func (g *MemoryGraphStore) GetDependencies(ctx context.Context, nodeID string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var deps []string
	edges := g.edges[nodeID]

	for _, edge := range edges {
		deps = append(deps, edge.TargetId)
	}

	return deps, nil
}

// CalculateBlastRadius returns all downstream impacted node IDs via BFS traversal.
func (g *MemoryGraphStore) CalculateBlastRadius(ctx context.Context, nodeID string) ([]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	var impacted []string
	queue := []string{nodeID}
	visited[nodeID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range g.edges[current] {
			if !visited[edge.TargetId] {
				visited[edge.TargetId] = true
				impacted = append(impacted, edge.TargetId)
				queue = append(queue, edge.TargetId)
			}
		}
	}

	return impacted, nil
}

// Snapshot serializes the current graph state (useful for async Postgres persistence).
func (g *MemoryGraphStore) Snapshot(ctx context.Context) ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	snapshot := struct {
		Nodes map[string]*topologyv1.TopologyNode `json:"nodes"`
		Edges map[string][]*topologyv1.TopologyEdge `json:"edges"`
	}{
		Nodes: g.nodes,
		Edges: g.edges,
	}

	return json.Marshal(snapshot)
}
