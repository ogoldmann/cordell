package web

import (
	"errors"
	"net/http"
	"time"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

func (s *Server) loadCurrentOperator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := readSessionToken(r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		operator, err := s.services.GetOperatorBySessionToken.Execute(r.Context(), app.GetOperatorBySessionTokenCommand{
			Token: token,
			Now:   time.Now().UTC(),
		})
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrExpiredOperatorSession) {
				clearSessionCookie(w, s.sessionCookieConfig)
				next.ServeHTTP(w, r)
				return
			}

			s.logger.Error("failed to load current operator", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r.WithContext(withCurrentOperator(r.Context(), operator)))
	})
}

func (s *Server) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := currentOperatorFromContext(r.Context()); ok {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
