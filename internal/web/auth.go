package web

import (
	"net/http"
	"time"
)

const sessionCookieName = "cordell_session"

type sessionCookieConfig struct {
	Secure bool
}

// NewSessionCookieConfig creates the session cookie configuration.
func NewSessionCookieConfig(secure bool) sessionCookieConfig {
	return sessionCookieConfig{
		Secure: secure,
	}
}

func setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time, config sessionCookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt.UTC(),
		HttpOnly: true,
		Secure:   config.Secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, config sessionCookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   config.Secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func readSessionToken(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return "", false
	}

	if cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}
