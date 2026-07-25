package web

import "testing"

func TestNewDashboardOperationalDockView(t *testing.T) {
	dock := newDashboardOperationalDockView()

	if len(dock.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(dock.Groups))
	}

	if dock.Groups[0].Title != "CADASTRAR" {
		t.Fatalf("expected first group CADASTRAR, got %q", dock.Groups[0].Title)
	}

	if dock.Groups[1].Title != "REGISTRAR" {
		t.Fatalf("expected second group REGISTRAR, got %q", dock.Groups[1].Title)
	}

	if len(dock.Groups[0].Actions) != 2 {
		t.Fatalf("expected 2 register actions, got %d", len(dock.Groups[0].Actions))
	}

	if len(dock.Groups[1].Actions) != 2 {
		t.Fatalf("expected 2 custody actions, got %d", len(dock.Groups[1].Actions))
	}
}

func TestRandomDashboardWelcomePhraseReturnsValue(t *testing.T) {
	phrase := randomDashboardWelcomePhrase()

	if phrase == "" {
		t.Fatal("expected welcome phrase")
	}
}
