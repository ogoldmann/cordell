package web

import (
	"net/http"

	"cordell/internal/domain"
)

func (s *Server) requireRole(role domain.OperatorRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			operator, ok := currentOperatorFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if operator.Role() != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return s.requireRole(domain.OperatorRoleAdmin)(next)
}
