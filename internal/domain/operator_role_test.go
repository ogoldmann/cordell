package domain

import "testing"

func TestNewOperatorRole(t *testing.T) {
	role, err := NewOperatorRole(" Admin ")
	if err != nil {
		t.Fatalf("expected valid role, got %v", err)
	}

	if role != OperatorRoleAdmin {
		t.Fatalf("expected admin, got %s", role)
	}
}

func TestNewOperatorRoleRejectsEmptyRole(t *testing.T) {
	_, err := NewOperatorRole(" ")
	if err != ErrEmptyOperatorRole {
		t.Fatalf("expected ErrEmptyOperatorRole, got %v", err)
	}
}

func TestNewOperatorRoleRejectsInvalidRole(t *testing.T) {
	_, err := NewOperatorRole("root")
	if err != ErrInvalidOperatorRole {
		t.Fatalf("expected ErrInvalidOperatorRole, got %v", err)
	}
}

func TestOperatorRoleCanManageOperators(t *testing.T) {
	if !OperatorRoleAdmin.CanManageOperators() {
		t.Fatal("expected admin to manage operators")
	}

	if OperatorRoleOperator.CanManageOperators() {
		t.Fatal("expected operator not to manage operators")
	}
}
