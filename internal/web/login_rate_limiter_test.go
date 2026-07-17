package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginRateLimiterBlocksAfterFailures(t *testing.T) {
	limiter := newLoginRateLimiter()
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	request.RemoteAddr = "127.0.0.1:12345"

	for i := 0; i < loginRateLimitMaxFailures; i++ {
		if !limiter.allow(request, "admin", now) {
			t.Fatalf("expected attempt %d to be allowed", i+1)
		}

		limiter.recordFailure(request, "admin", now)
	}

	if limiter.allow(request, "admin", now) {
		t.Fatal("expected login to be rate limited")
	}
}

func TestLoginRateLimiterAllowsAfterWindow(t *testing.T) {
	limiter := newLoginRateLimiter()
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	request.RemoteAddr = "127.0.0.1:12345"

	for i := 0; i < loginRateLimitMaxFailures; i++ {
		limiter.recordFailure(request, "admin", now)
	}

	if !limiter.allow(request, "admin", now.Add(loginRateLimitWindow+time.Second)) {
		t.Fatal("expected login to be allowed after rate limit window")
	}
}

func TestLoginRateLimiterClearsOnSuccess(t *testing.T) {
	limiter := newLoginRateLimiter()
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	request.RemoteAddr = "127.0.0.1:12345"

	for i := 0; i < loginRateLimitMaxFailures; i++ {
		limiter.recordFailure(request, "admin", now)
	}

	limiter.recordSuccess(request, "admin")

	if !limiter.allow(request, "admin", now) {
		t.Fatal("expected login to be allowed after success")
	}
}
