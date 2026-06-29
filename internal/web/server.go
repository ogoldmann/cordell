package web

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server owns the HTTP dependencies and route definitions for the application.
type Server struct {
	logger *slog.Logger
}

// NewServer creates a Server with its required dependencies.
func NewServer(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

// Routes builds and returns the application's HTTP handler tree.
func (s *Server) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", s.handleHealthCheck)

	return router
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok","service":"cordell"}`))
}
