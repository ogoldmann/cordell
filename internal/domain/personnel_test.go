package domain

import "testing"

func TestNewPersonnel(t *testing.T) {
	personnel, err := NewPersonnel("personnel-1", "  John Doe  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if personnel.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", personnel.ID())
	}

	if personnel.FullName() != "John Doe" {
		t.Fatalf("expected trimmed full name John Doe, got %s", personnel.FullName())
	}

	if !personnel.Active() {
		t.Fatal("expected personnel to be active")
	}
}

func TestNewPersonnelRejectsEmptyID(t *testing.T) {
	_, err := NewPersonnel("", "John Doe")
	if err != ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestNewPersonnelRejectsEmptyName(t *testing.T) {
	_, err := NewPersonnel("personnel-1", "   ")
	if err != ErrEmptyPersonnelName {
		t.Fatalf("expected ErrEmptyPersonnelName, got %v", err)
	}
}
