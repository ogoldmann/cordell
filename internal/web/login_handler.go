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
	ReturnTo string
}

func (s *Server) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	if _, ok := currentOperatorFromContext(r.Context()); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	returnTo := sanitizeReturnTo(r.URL.Query().Get("return_to"))

	if err := s.renderer.Render(w, http.StatusOK, "login.html", loginPageData{
		Title:    "Login",
		ReturnTo: returnTo,
	}); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := currentOperatorFromContext(r.Context()); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	returnTo := sanitizeReturnTo(r.FormValue("return_to"))
	now := time.Now().UTC()

	if !s.loginRateLimiter.allow(r, username, now) {
		s.renderLoginError(w, http.StatusTooManyRequests, username, returnTo, "Too many login attempts. Try again later.")
		return
	}

	operator, err := s.services.AuthenticateOperator.Execute(r.Context(), app.AuthenticateOperatorCommand{
		Username: username,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			s.loginRateLimiter.recordFailure(r, username, now)
			s.renderLoginError(w, http.StatusUnauthorized, username, returnTo, "Invalid username or password.")
			return
		}

		s.logger.Error("failed to authenticate operator", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.loginRateLimiter.recordSuccess(r, username)

	if err := s.services.DeleteExpiredOperatorSessions.Execute(r.Context(), app.DeleteExpiredOperatorSessionsCommand{
		Now: now,
	}); err != nil {
		s.logger.Error("failed to delete expired operator sessions", "error", err)
	}

	sessionResult, err := s.services.CreateOperatorSession.Execute(r.Context(), app.CreateOperatorSessionCommand{
		OperatorID: operator.ID(),
		Now:        now,
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

	http.Redirect(w, r, returnTo, http.StatusSeeOther)
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

func (s *Server) renderLoginError(
	w http.ResponseWriter,
	status int,
	username string,
	returnTo string,
	message string,
) {
	if err := s.renderer.Render(w, status, "login.html", loginPageData{
		Title:    "Login",
		Error:    message,
		Username: username,
		ReturnTo: returnTo,
	}); err != nil {
		s.handleRenderError(w, err)
	}
}
