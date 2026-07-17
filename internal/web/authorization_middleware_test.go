package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cordell/internal/domain"
)

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	server := &Server{}

	operator, err := domain.NewOperator(
		"operator-1",
		"admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

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

	operator, err := domain.NewOperator(
		"operator-1",
		"clerk",
		domain.OperatorRoleOperator,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

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
