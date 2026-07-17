package domain

import "testing"

func TestNewOperator(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	operator, err := NewOperator(
		"operator-1",
		registrationID,
		"silva",
		RankSergeant,
		OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if operator.RegistrationID() != registrationID {
		t.Fatalf("expected registration id %s, got %s", registrationID, operator.RegistrationID())
	}

	if operator.Alias() != "silva" {
		t.Fatalf("expected alias silva, got %s", operator.Alias())
	}

	if operator.Rank() != RankSergeant {
		t.Fatalf("expected rank sergeant, got %s", operator.Rank())
	}

	if operator.Role() != OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", operator.Role())
	}

	if !operator.Active() {
		t.Fatal("expected operator to be active")
	}
}

func TestNewOperatorRejectsEmptyID(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewOperator("", registrationID, "silva", RankSergeant, OperatorRoleAdmin, "$argon2id$hash")
	if err != ErrEmptyOperatorID {
		t.Fatalf("expected ErrEmptyOperatorID, got %v", err)
	}
}

func TestNewOperatorRejectsEmptyAlias(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewOperator("operator-1", registrationID, " ", RankSergeant, OperatorRoleAdmin, "$argon2id$hash")
	if err != ErrEmptyOperatorAlias {
		t.Fatalf("expected ErrEmptyOperatorAlias, got %v", err)
	}
}

func TestNewOperatorRejectsInvalidRank(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewOperator("operator-1", registrationID, "silva", Rank("general"), OperatorRoleAdmin, "$argon2id$hash")
	if err != ErrInvalidOperatorRank {
		t.Fatalf("expected ErrInvalidOperatorRank, got %v", err)
	}
}

func TestNewOperatorRejectsInvalidRole(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewOperator("operator-1", registrationID, "silva", RankSergeant, OperatorRole("root"), "$argon2id$hash")
	if err != ErrInvalidOperatorRole {
		t.Fatalf("expected ErrInvalidOperatorRole, got %v", err)
	}
}

func TestNewOperatorRejectsEmptyPasswordHash(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewOperator("operator-1", registrationID, "silva", RankSergeant, OperatorRoleAdmin, "")
	if err != ErrEmptyOperatorPasswordHash {
		t.Fatalf("expected ErrEmptyOperatorPasswordHash, got %v", err)
	}
}
