package app

import (
	"errors"
	"testing"

	"cordell/internal/domain"
)

func TestDuplicatePersonnelRegistrationIDErrorWrapsDomainError(t *testing.T) {
	err := DuplicatePersonnelRegistrationIDError{
		ExistingPersonnelID: "personnel-1",
	}

	if !errors.Is(err, domain.ErrDuplicateRegistrationID) {
		t.Fatal("expected error to wrap ErrDuplicateRegistrationID")
	}

	existingID, ok := ExistingPersonnelIDFromDuplicateError(err)
	if !ok {
		t.Fatal("expected existing personnel ID")
	}

	if existingID != "personnel-1" {
		t.Fatalf("expected personnel-1, got %s", existingID)
	}
}

func TestDuplicateAssetNameErrorWrapsDomainError(t *testing.T) {
	err := DuplicateAssetNameError{
		ExistingAssetID: "asset-1",
	}

	if !errors.Is(err, domain.ErrDuplicateAssetName) {
		t.Fatal("expected error to wrap ErrDuplicateAssetName")
	}

	existingID, ok := ExistingAssetIDFromDuplicateError(err)
	if !ok {
		t.Fatal("expected existing asset ID")
	}

	if existingID != "asset-1" {
		t.Fatalf("expected asset-1, got %s", existingID)
	}
}
