package web

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"cordell/internal/domain"
)

const (
	loginRateLimitMaxFailures = 5
	loginRateLimitWindow      = 15 * time.Minute
)

type loginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string]loginAttemptState
}

type loginAttemptState struct {
	failures   int
	windowEnds time.Time
}

func newLoginRateLimiter() *loginRateLimiter {
	return &loginRateLimiter{
		attempts: map[string]loginAttemptState{},
	}
}

func (l *loginRateLimiter) allow(r *http.Request, registrationID string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := loginRateLimitKey(r, registrationID)
	state, ok := l.attempts[key]
	if !ok {
		return true
	}

	if !now.Before(state.windowEnds) {
		delete(l.attempts, key)
		return true
	}

	return state.failures < loginRateLimitMaxFailures
}

func (l *loginRateLimiter) recordFailure(r *http.Request, registrationID string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	key := loginRateLimitKey(r, registrationID)
	state, ok := l.attempts[key]
	if !ok || !now.Before(state.windowEnds) {
		l.attempts[key] = loginAttemptState{
			failures:   1,
			windowEnds: now.Add(loginRateLimitWindow),
		}
		return
	}

	state.failures++
	l.attempts[key] = state
}

func (l *loginRateLimiter) recordSuccess(r *http.Request, registrationID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.attempts, loginRateLimitKey(r, registrationID))
}

func loginRateLimitKey(r *http.Request, registrationID string) string {
	normalized := domain.NormalizeRegistrationID(registrationID)
	if normalized == "" {
		normalized = strings.ToLower(strings.TrimSpace(registrationID))
	}

	return clientIP(r) + "|" + normalized
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
