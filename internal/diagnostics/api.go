package diagnostics

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shadow0vortex/cortexops/internal/topology/graph"
)

type API struct {
	store *graph.MemoryGraphStore
}

func NewAPI(store *graph.MemoryGraphStore) *API {
	return &API{store: store}
}

func (api *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", api.handleHealth)
	mux.HandleFunc("GET /topology/nodes", api.handleListNodes)
	mux.HandleFunc("GET /topology/blast-radius/{id}", api.handleBlastRadius)
}

func (api *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (api *API) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := api.store.ListNodes(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(nodes),
		"nodes": nodes,
	})
}

func (api *API) handleBlastRadius(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing node ID", http.StatusBadRequest)
		return
	}

	impacted, err := api.store.CalculateBlastRadius(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"node":     id,
		"count":    len(impacted),
		"impacted": impacted,
	})
}
