package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	server := &Server{
		securityHeadersConfig: NewSecurityHeadersConfig(false),
	}

	handler := server.securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options nosniff, got %s", recorder.Header().Get("X-Content-Type-Options"))
	}

	if recorder.Header().Get("Referrer-Policy") == "" {
		t.Fatal("expected Referrer-Policy header")
	}

	if recorder.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("expected Content-Security-Policy header")
	}

	if recorder.Header().Get("Strict-Transport-Security") != "" {
		t.Fatal("expected HSTS header to be disabled")
	}
}

func TestSecurityHeadersWithHSTS(t *testing.T) {
	server := &Server{
		securityHeadersConfig: NewSecurityHeadersConfig(true),
	}

	handler := server.securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Header().Get("Strict-Transport-Security") == "" {
		t.Fatal("expected HSTS header")
	}
}
