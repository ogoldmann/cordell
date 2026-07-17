package domain

import "testing"

func TestNewOperator(t *testing.T) {
	operator, err := NewOperator("operator-1", "Admin.User", OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if operator.Username() != "admin.user" {
		t.Fatalf("expected normalized username admin.user, got %s", operator.Username())
	}

	if operator.Role() != OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", operator.Role())
	}

	if !operator.Active() {
		t.Fatal("expected operator to be active")
	}
}

func TestNewOperatorRejectsEmptyID(t *testing.T) {
	_, err := NewOperator("", "admin", OperatorRoleAdmin, "$argon2id$hash")
	if err != ErrEmptyOperatorID {
		t.Fatalf("expected ErrEmptyOperatorID, got %v", err)
	}
}

func TestNewOperatorRejectsInvalidUsername(t *testing.T) {
	_, err := NewOperator("operator-1", "admin user", OperatorRoleAdmin, "$argon2id$hash")
	if err != ErrInvalidOperatorUsername {
		t.Fatalf("expected ErrInvalidOperatorUsername, got %v", err)
	}
}

func TestNewOperatorRejectsInvalidRole(t *testing.T) {
	_, err := NewOperator("operator-1", "admin", OperatorRole("root"), "$argon2id$hash")
	if err != ErrInvalidOperatorRole {
		t.Fatalf("expected ErrInvalidOperatorRole, got %v", err)
	}
}

func TestNewOperatorRejectsEmptyPasswordHash(t *testing.T) {
	_, err := NewOperator("operator-1", "admin", OperatorRoleAdmin, "")
	if err != ErrEmptyOperatorPasswordHash {
		t.Fatalf("expected ErrEmptyOperatorPasswordHash, got %v", err)
	}
}
