package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cordell/internal/domain"
)

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	server := &Server{}

	operator := mustBuildOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	nextCalled := false
	handler := server.requireRole(domain.OperatorRoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	request = request.WithContext(withCurrentOperator(request.Context(), operator))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestRequireRoleRejectsDifferentRole(t *testing.T) {
	server := &Server{}

	operator := mustBuildOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleOperator)

	handler := server.requireRole(domain.OperatorRoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	request = request.WithContext(withCurrentOperator(request.Context(), operator))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", recorder.Code)
	}
}

func TestRequireRoleRejectsMissingOperator(t *testing.T) {
	server := &Server{}

	handler := server.requireRole(domain.OperatorRoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	}))

	request := httptest.NewRequest(http.MethodGet, "/admin", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
