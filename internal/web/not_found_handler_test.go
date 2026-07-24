package web

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cordell/internal/app"
)

func TestRoutesRendersCustomNotFoundPage(t *testing.T) {
	server, err := NewServer(
		slog.Default(),
		app.Services{},
		NewSessionCookieConfig(false),
		NewSecurityHeadersConfig(false),
	)
	if err != nil {
		t.Fatalf("expected server, got %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	response := httptest.NewRecorder()

	server.Routes().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}

	body := response.Body.String()

	if !strings.Contains(body, "Página não encontrada") {
		t.Fatalf("expected custom not found page, got body %q", body)
	}

	if !strings.Contains(body, "Home") {
		t.Fatalf("expected home button, got body %q", body)
	}
}
