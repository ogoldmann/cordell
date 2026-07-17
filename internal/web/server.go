package web

import (
	"log/slog"
	"net/http"

	"cordell/internal/app"

	"github.com/go-chi/chi/v5"
)

// Server owns the HTTP dependencies and route definitions for the application.
type Server struct {
	logger                *slog.Logger
	services              app.Services
	renderer              *Renderer
	sessionCookieConfig   sessionCookieConfig
	securityHeadersConfig securityHeadersConfig
	loginRateLimiter      *loginRateLimiter
}

// NewServer creates a Server with its required dependencies.
func NewServer(
	logger *slog.Logger,
	services app.Services,
	sessionCookieConfig sessionCookieConfig,
	securityHeadersConfig securityHeadersConfig,
) (*Server, error) {
	renderer, err := NewRenderer()
	if err != nil {
		return nil, err
	}

	return &Server{
		logger:                logger,
		services:              services,
		renderer:              renderer,
		sessionCookieConfig:   sessionCookieConfig,
		securityHeadersConfig: securityHeadersConfig,
		loginRateLimiter:      newLoginRateLimiter(),
	}, nil
}

// Routes builds and returns the application's HTTP handler tree.
func (s *Server) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(s.securityHeaders)

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.Group(func(public chi.Router) {
		public.Use(s.loadCurrentOperator)

		public.Get("/login", s.handleLoginForm)
		public.Post("/login", s.handleLogin)
	})

	router.Group(func(private chi.Router) {
		private.Use(s.loadCurrentOperator)
		private.Use(s.requireAuthentication)
		private.Use(s.requireCSRF)

		private.Get("/", s.handleDashboard)
		private.Get("/search", s.handleGlobalSearch)

		private.Group(func(admin chi.Router) {
			admin.Use(s.requireAdmin)

			admin.Get("/admin", s.handleAdminIndex)
			admin.Get("/admin/operators", s.handleAdminOperatorsIndex)
			admin.Get("/admin/operators/new", s.handleNewAdminOperatorForm)
			admin.Post("/admin/operators", s.handleCreateAdminOperator)

		})

		private.Get("/personnel", s.handleListPersonnel)
		private.Get("/personnel/new", s.handleNewPersonnelForm)
		private.Post("/personnel", s.handleCreatePersonnel)
		private.Get("/personnel/{id}", s.handleShowPersonnel)

		private.Get("/assets", s.handleListAssets)
		private.Get("/assets/new", s.handleNewAssetForm)
		private.Post("/assets", s.handleCreateAsset)
		private.Get("/assets/{id}", s.handleShowAsset)

		private.Get("/custody/checkouts/new", s.handleNewCheckoutForm)
		private.Post("/custody/checkouts", s.handleCreateCheckout)

		private.Get("/custody/returns/new", s.handleNewReturnForm)
		private.Post("/custody/returns", s.handleCreateReturn)

		private.Post("/logout", s.handleLogout)
	})

	return router
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok","service":"cordell"}`))
}
