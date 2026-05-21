package diagnostics

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shadow0vortex/cortexops/internal/topology/graph"
)

// API provides HTTP endpoints for topology diagnostics.
type API struct {
	graphStore *graph.MemoryGraphStore
	log        *slog.Logger
}

// NewAPI initializes the diagnostics API with a graph store.
func NewAPI(graphStore *graph.MemoryGraphStore) *API {
	return &API{
		graphStore: graphStore,
		log:        slog.Default(),
	}
}

// RegisterRoutes registers all diagnostics API routes.
func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", a.handleHealth)
	mux.HandleFunc("GET /topology/nodes", a.handleListNodes)
	mux.HandleFunc("GET /topology/blast-radius/{nodeID}", a.handleBlastRadius)
}

// handleHealth is a simple health check endpoint.
func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handleListNodes returns all topology nodes.
func (a *API) handleListNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// In a full implementation, this would query the graph store
	// For now, return a placeholder response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodes": []interface{}{},
		"count": 0,
	})
}

// handleBlastRadius calculates the blast radius for a given node.
func (a *API) handleBlastRadius(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("nodeID")
	if nodeID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "nodeID is required",
		})
		return
	}

	ctx := context.Background()
	impacted, err := a.graphStore.CalculateBlastRadius(ctx, nodeID)
	if err != nil {
		a.log.Error("Failed to calculate blast radius", "nodeID", nodeID, "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to calculate blast radius",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodeID":  nodeID,
		"impacted": impacted,
		"count":   len(impacted),
	})
}
