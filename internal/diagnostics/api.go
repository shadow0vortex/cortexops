package diagnostics

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/shadow0vortex/cortexops/internal/topology/graph"
)

type API struct {
	store    *graph.MemoryGraphStore
	apiToken string
}

func NewAPI(store *graph.MemoryGraphStore) *API {
	return &API{
		store:    store,
		apiToken: os.Getenv("DIAG_API_TOKEN"),
	}
}

func (api *API) RegisterRoutes(mux *http.ServeMux) {
	// Health endpoint — unauthenticated (for K8s probes)
	mux.HandleFunc("GET /health", api.handleHealth)

	// Versioned API endpoints — protected by bearer token
	mux.HandleFunc("GET /v1/topology/nodes", api.requireAuth(api.handleListNodes))
	mux.HandleFunc("GET /v1/topology/blast-radius/{id}", api.requireAuth(api.handleBlastRadius))
}

// requireAuth wraps a handler with bearer token authentication.
// If DIAG_API_TOKEN is not set, authentication is disabled (dev mode).
func (api *API) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If no token configured, skip auth (development mode)
		if api.apiToken == "" {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] != api.apiToken {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func (api *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (api *API) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := api.store.ListNodes(r.Context())
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

	impacted, err := api.store.CalculateBlastRadius(r.Context(), id)
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
