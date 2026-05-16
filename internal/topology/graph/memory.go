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
	edges map[string][]topologyv1.TopologyEdge // key is source_id
}

// NewMemoryGraphStore initializes the graph.
func NewMemoryGraphStore() *MemoryGraphStore {
	return &MemoryGraphStore{
		nodes: make(map[string]*topologyv1.TopologyNode),
		edges: make(map[string][]topologyv1.TopologyEdge),
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

	edge := topologyv1.TopologyEdge{
		SourceId:     sourceID,
		TargetId:     targetID,
		Relationship: relationship,
	}
	
	g.edges[sourceID] = append(g.edges[sourceID], edge)
	return nil
}

// GetDependencies returns all nodes that the given node has outward edges to.
func (g *MemoryGraphStore) GetDependencies(ctx context.Context, nodeID string) ([]*topologyv1.TopologyNode, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var deps []*topologyv1.TopologyNode
	edges := g.edges[nodeID]
	
	for _, edge := range edges {
		if node, exists := g.nodes[edge.TargetId]; exists {
			deps = append(deps, node)
		}
	}
	
	return deps, nil
}

// Snapshot serializes the current graph state (useful for async Postgres persistence).
func (g *MemoryGraphStore) Snapshot(ctx context.Context) ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	snapshot := struct {
		Nodes map[string]*topologyv1.TopologyNode `json:"nodes"`
		Edges map[string][]topologyv1.TopologyEdge `json:"edges"`
	}{
		Nodes: g.nodes,
		Edges: g.edges,
	}

	return json.Marshal(snapshot)
}
