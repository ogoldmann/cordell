package domain

import (
	"errors"
	"testing"
)

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

func TestPersonnelUpdateDetails(t *testing.T) {
	personnel := mustNewPersonnel(t)

	err := personnel.UpdateDetails(
		"John Updated",
		"Updated",
		PersonnelRankSoldier,
		RegistrationID("93541134780"),
		PersonnelSectionOperations,
		OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if personnel.FullName() != "John Updated" {
		t.Fatalf("expected updated full name, got %q", personnel.FullName())
	}

	if personnel.Alias() != "Updated" {
		t.Fatalf("expected updated alias, got %q", personnel.Alias())
	}

	if personnel.Rank() != PersonnelRankSoldier {
		t.Fatalf("expected updated rank, got %s", personnel.Rank())
	}

	if personnel.RegistrationID() != RegistrationID("93541134780") {
		t.Fatalf("expected normalized registration id, got %s", personnel.RegistrationID())
	}

	if personnel.Section() != PersonnelSectionOperations {
		t.Fatalf("expected updated section, got %s", personnel.Section())
	}
}

func TestPersonnelUpdateDetailsRejectsEmptyName(t *testing.T) {
	personnel := mustNewPersonnel(t)

	err := personnel.UpdateDetails(
		"",
		"Updated",
		PersonnelRankSoldier,
		RegistrationID("93541134780"),
		PersonnelSectionOperations,
		OrganizationUnitDefault,
	)
	if !errors.Is(err, ErrEmptyPersonnelName) {
		t.Fatalf("expected ErrEmptyPersonnelName, got %v", err)
	}
}

func mustNewPersonnel(t *testing.T) Personnel {
	t.Helper()

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
		t.Fatalf("expected valid personnel, got %v", err)
	}

	return personnel
}
