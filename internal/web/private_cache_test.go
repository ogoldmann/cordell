package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoStore(t *testing.T) {
	server := &Server{}

	handler := server.noStore(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %s", recorder.Header().Get("Cache-Control"))
	}

	if recorder.Header().Get("Pragma") != "no-cache" {
		t.Fatalf("expected Pragma no-cache, got %s", recorder.Header().Get("Pragma"))
	}

	if recorder.Header().Get("Expires") != "0" {
		t.Fatalf("expected Expires 0, got %s", recorder.Header().Get("Expires"))
	}
}
