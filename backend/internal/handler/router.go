package handler

import (
	"log"
	"net/http"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
)

// RegisterRoutes registers all HTTP routes on the given mux.
// API routes defined in docs/api/openapi.yaml are served under /api/v1/,
// backed by apiHandler.
func RegisterRoutes(mux *http.ServeMux, apiHandler *APIHandler) {
	mux.HandleFunc("GET /healthz", handleHealthz)

	srv, err := oas.NewServer(apiHandler)
	if err != nil {
		log.Fatalf("failed to build oas server: %v", err)
	}
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
