package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Server provides basic HTTP handlers for Kubernetes readiness and liveness probes.
type Server struct {
	mu    sync.RWMutex
	ready bool
}

// NewServer creates a new health check server.
func NewServer() *Server {
	return &Server{
		ready: false, // Initially false, should be marked ready after dependencies (e.g., DB, NATS) are connected.
	}
}

// SetReady updates the readiness status of the application.
func (s *Server) SetReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

// LivenessHandler responds with 200 OK as long as the HTTP server is running.
func (s *Server) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

// ReadinessHandler responds with 200 OK if the app is marked ready, otherwise 503.
func (s *Server) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if ready {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "not_ready"})
	}
}
