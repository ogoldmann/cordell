package web

import (
	"crypto/subtle"
	"net/http"
)

const csrfFormFieldName = "csrf_token"

func (s *Server) requireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isSafeMethod(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		session, ok := currentSessionFromContext(r.Context())
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		submittedToken := r.FormValue(csrfFormFieldName)
		if !sameSecret(submittedToken, session.CSRFToken()) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isSafeMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions
}

func sameSecret(first string, second string) bool {
	if first == "" || second == "" {
		return false
	}

	if len(first) != len(second) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(first), []byte(second)) == 1
}

func csrfTokenFromContext(r *http.Request) string {
	session, ok := currentSessionFromContext(r.Context())
	if !ok {
		return ""
	}

	return session.CSRFToken()
}
