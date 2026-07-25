package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrivateRoutesRequireAuthentication(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name         string
		path         string
		wantLocation string
	}{
		{
			name:         "dashboard",
			path:         "/",
			wantLocation: "/login?return_to=%2F",
		},
		{
			name:         "personnel",
			path:         "/personnel",
			wantLocation: "/login?return_to=%2Fpersonnel",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			recorder := httptest.NewRecorder()

			server.Routes().ServeHTTP(recorder, request)

			if recorder.Code != http.StatusSeeOther {
				t.Fatalf("expected status 303, got %d", recorder.Code)
			}

			if location := recorder.Header().Get("Location"); location != test.wantLocation {
				t.Fatalf("expected redirect to %q, got %q", test.wantLocation, location)
			}
		})
	}
}
