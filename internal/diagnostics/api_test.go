package diagnostics

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterRoutesHealth(t *testing.T) {
	for _, path := range []string{"/health", "/debug/healthz"} {
		t.Run(path, func(t *testing.T) {
			mux := http.NewServeMux()
			(&API{}).RegisterRoutes(mux)

			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))

			if response.Code != http.StatusOK {
				t.Fatalf("GET %s returned %d, want %d", path, response.Code, http.StatusOK)
			}
		})
	}
}