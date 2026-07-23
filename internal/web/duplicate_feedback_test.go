package web

import (
	"testing"

	"cordell/internal/app"
)

func TestPersonnelFeedbackFromDuplicateErrorIncludesAction(t *testing.T) {
	feedback := personnelFeedbackFromError(app.DuplicatePersonnelRegistrationIDError{
		ExistingPersonnelID: "personnel-1",
	})

	if feedback == nil {
		t.Fatal("expected feedback")
	}

	if feedback.ActionLabel != "Abrir militar existente" {
		t.Fatalf("expected action label, got %q", feedback.ActionLabel)
	}

	if feedback.ActionURL != "/personnel/personnel-1" {
		t.Fatalf("expected action URL, got %q", feedback.ActionURL)
	}
}

func TestAssetFeedbackFromDuplicateErrorIncludesAction(t *testing.T) {
	feedback := assetFeedbackFromError(app.DuplicateAssetNameError{
		ExistingAssetID: "asset-1",
	})

	if feedback == nil {
		t.Fatal("expected feedback")
	}

	if feedback.ActionLabel != "Abrir material existente" {
		t.Fatalf("expected action label, got %q", feedback.ActionLabel)
	}

	if feedback.ActionURL != "/assets/asset-1" {
		t.Fatalf("expected action URL, got %q", feedback.ActionURL)
	}
}
