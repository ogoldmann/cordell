package web

import (
	"errors"
	"net/http"
	"time"

	"cordell/internal/app"
	"cordell/internal/domain"
)

type loginPageData struct {
	Title    string
	Error    string
	Username string
}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	if err := s.renderer.Render(w, http.StatusOK, "login.html", loginPageData{
		Title: "Login",
	}); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	operator, err := s.services.AuthenticateOperator.Execute(r.Context(), app.AuthenticateOperatorCommand{
		Username: username,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			s.renderLoginError(w, username)
			return
		}

		s.logger.Error("failed to authenticate operator", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sessionResult, err := s.services.CreateOperatorSession.Execute(r.Context(), app.CreateOperatorSessionCommand{
		OperatorID: operator.ID(),
		Now:        time.Now().UTC(),
	})
	if err != nil {
		s.logger.Error("failed to create operator session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(
		w,
		sessionResult.Token,
		sessionResult.Session.ExpiresAt(),
		s.sessionCookieConfig,
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if token, ok := readSessionToken(r); ok {
		if err := s.services.DeleteOperatorSession.Execute(r.Context(), app.DeleteOperatorSessionCommand{
			Token: token,
		}); err != nil {
			s.logger.Error("failed to delete operator session", "error", err)
		}
	}

	clearSessionCookie(w, s.sessionCookieConfig)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) renderLoginError(w http.ResponseWriter, username string) {
	if err := s.renderer.Render(w, http.StatusUnauthorized, "login.html", loginPageData{
		Title:    "Login",
		Error:    "Invalid username or password.",
		Username: username,
	}); err != nil {
		s.handleRenderError(w, err)
	}
}
