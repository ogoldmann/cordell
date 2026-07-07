package domain

import "testing"

func TestNewPersonnel(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	personnel, err := NewPersonnel(
		"personnel-1",
		"John Doe",
		"Doe",
		PersonnelRankSergeant,
		registrationID,
		PersonnelSectionOperations,
		OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if personnel.ID() != "personnel-1" {
		t.Fatalf("expected id personnel-1, got %s", personnel.ID())
	}

	if personnel.FullName() != "John Doe" {
		t.Fatalf("expected full name John Doe, got %s", personnel.FullName())
	}

	if personnel.Alias() != "Doe" {
		t.Fatalf("expected alias Doe, got %s", personnel.Alias())
	}

	if personnel.Rank() != PersonnelRankSergeant {
		t.Fatalf("expected rank sergeant, got %s", personnel.Rank())
	}

	if personnel.RegistrationID() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", personnel.RegistrationID())
	}

	if personnel.Section() != PersonnelSectionOperations {
		t.Fatalf("expected section operations, got %s", personnel.Section())
	}

	if personnel.OrganizationUnit() != OrganizationUnitDefault {
		t.Fatalf("expected organization unit default_unit, got %s", personnel.OrganizationUnit())
	}

	if !personnel.Active() {
		t.Fatal("expected personnel to be active")
	}
}

func TestNewPersonnelRejectsEmptyID(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewPersonnel(
		"",
		"John Doe",
		"Doe",
		PersonnelRankSergeant,
		registrationID,
		PersonnelSectionOperations,
		OrganizationUnitDefault,
	)
	if err != ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestNewPersonnelRejectsEmptyName(t *testing.T) {
	registrationID, err := NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	_, err = NewPersonnel(
		"personnel-1",
		"   ",
		"Doe",
		PersonnelRankSergeant,
		registrationID,
		PersonnelSectionOperations,
		OrganizationUnitDefault,
	)
	if err != ErrEmptyPersonnelName {
		t.Fatalf("expected ErrEmptyPersonnelName, got %v", err)
	}
}
