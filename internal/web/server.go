package web

import (
	"log/slog"
	"net/http"

	"cordell/internal/app"

	"github.com/go-chi/chi/v5"
)

// Server owns the HTTP dependencies and route definitions for the application.
type Server struct {
	logger              *slog.Logger
	services            app.Services
	renderer            *Renderer
	sessionCookieConfig sessionCookieConfig
}

// NewServer creates a Server with its required dependencies.
func NewServer(
	logger *slog.Logger,
	services app.Services,
	sessionCookieConfig sessionCookieConfig,
) (*Server, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return nil, err
	}

	return &Server{
		logger:              logger,
		services:            services,
		renderer:            renderer,
		sessionCookieConfig: sessionCookieConfig,
	}, nil
}

// Routes builds and returns the application's HTTP handler tree.
func (s *Server) Routes() http.Handler {
	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.Get("/", s.handleDashboard)
	router.Get("/search", s.handleGlobalSearch)
	router.Get("/health", s.handleHealthCheck)

	router.Get("/login", s.handleLoginForm)
	router.Post("/login", s.handleLogin)
	router.Post("/logout", s.handleLogout)

	router.Get("/personnel", s.handleListPersonnel)
	router.Get("/personnel/new", s.handleNewPersonnelForm)
	router.Post("/personnel", s.handleCreatePersonnel)
	router.Get("/personnel/{id}", s.handleShowPersonnel)

	router.Get("/assets", s.handleListAssets)
	router.Get("/assets/new", s.handleNewAssetForm)
	router.Post("/assets", s.handleCreateAsset)
	router.Get("/assets/{id}", s.handleShowAsset)

	router.Get("/custody/checkouts/new", s.handleNewCheckoutForm)
	router.Post("/custody/checkouts", s.handleCreateCheckout)

	router.Get("/custody/returns/new", s.handleNewReturnForm)
	router.Post("/custody/returns", s.handleCreateReturn)

	return router
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok","service":"cordell"}`))
}
