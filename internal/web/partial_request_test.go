package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWantsPartialResponseFromHeader(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/personnel", nil)
	request.Header.Set(partialRequestHeader, "1")

	if !wantsPartialResponse(request) {
		t.Fatal("expected partial response")
	}
}

func TestWantsPartialResponseFromQuery(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/personnel?partial=1", nil)

	if !wantsPartialResponse(request) {
		t.Fatal("expected partial response")
	}
}

func TestWantsPartialResponseFalse(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/personnel", nil)

	if wantsPartialResponse(request) {
		t.Fatal("expected full response")
	}
}
