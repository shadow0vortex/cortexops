package diagnostics

import (
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
	mux.HandleFunc("GET /debug/healthz", a.handleHealth)
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

	nodes, err := a.graphStore.ListNodes(ctx)
	if err != nil {
		a.log.Error("Failed to list nodes", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to list nodes",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodes": nodes,
		"count": len(nodes),
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

	ctx := r.Context()
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
