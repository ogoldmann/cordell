package web

import "testing"

func TestNewPageActionReturnsNilForEmptyValues(t *testing.T) {
	if action := newPageAction("", "/personnel/new"); action != nil {
		t.Fatal("expected nil action for empty label")
	}

	if action := newPageAction("Cadastrar militar", ""); action != nil {
		t.Fatal("expected nil action for empty URL")
	}
}

func TestNewPageAction(t *testing.T) {
	action := newPageAction("Cadastrar militar", "/personnel/new")
	if action == nil {
		t.Fatal("expected action")
	}

	if action.Label != "Cadastrar militar" {
		t.Fatalf("expected label Cadastrar militar, got %q", action.Label)
	}

	if action.URL != "/personnel/new" {
		t.Fatalf("expected URL /personnel/new, got %q", action.URL)
	}
}

func TestNewFormActions(t *testing.T) {
	actions := newFormActions("Salvar", "Cancelar", "/personnel")

	if actions.PrimaryLabel != "Salvar" {
		t.Fatalf("expected primary label Salvar, got %q", actions.PrimaryLabel)
	}

	if actions.SecondaryLabel != "Cancelar" {
		t.Fatalf("expected secondary label Cancelar, got %q", actions.SecondaryLabel)
	}

	if actions.SecondaryURL != "/personnel" {
		t.Fatalf("expected secondary URL /personnel, got %q", actions.SecondaryURL)
	}

	if actions.ShowSaveAndCreateAnother {
		t.Fatal("expected save and create another to be disabled")
	}
}

func TestNewErrorFeedback(t *testing.T) {
	feedback := newErrorFeedback("Erro de validação.")
	if feedback == nil {
		t.Fatal("expected feedback")
	}

	if feedback.Kind != "error" {
		t.Fatalf("expected kind error, got %q", feedback.Kind)
	}

	if feedback.Title != "Verifique as informações" {
		t.Fatalf("expected default error title, got %q", feedback.Title)
	}

	if feedback.Message != "Erro de validação." {
		t.Fatalf("expected message, got %q", feedback.Message)
	}
}

func TestNewErrorFeedbackReturnsNilForEmptyMessage(t *testing.T) {
	if feedback := newErrorFeedback(""); feedback != nil {
		t.Fatal("expected nil feedback")
	}
}

func TestNewDetailField(t *testing.T) {
	field := newDetailField("Nome", "Rádio")

	if field.Label != "Nome" {
		t.Fatalf("expected label Nome, got %q", field.Label)
	}

	if field.Value != "Rádio" {
		t.Fatalf("expected value Rádio, got %q", field.Value)
	}
}
