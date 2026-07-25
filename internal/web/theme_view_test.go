package web

import "testing"

func TestNewThemeSelectorView(t *testing.T) {
	selector := newThemeSelectorView()

	if selector.Compact {
		t.Fatal("expected default theme selector not to be compact")
	}

	if len(selector.Options) != 3 {
		t.Fatalf("expected 3 theme options, got %d", len(selector.Options))
	}

	expectedValues := []string{"light", "dark", "sepia"}

	for index, expectedValue := range expectedValues {
		if selector.Options[index].Value != expectedValue {
			t.Fatalf("expected theme option %q, got %q", expectedValue, selector.Options[index].Value)
		}
	}

	for _, option := range selector.Options {
		if option.Label == "" {
			t.Fatalf("expected label for theme %q", option.Value)
		}

		if option.Icon.Name == "" {
			t.Fatalf("expected icon for theme %q", option.Value)
		}
	}
}

func TestNewCompactThemeSelectorView(t *testing.T) {
	selector := newCompactThemeSelectorView()

	if !selector.Compact {
		t.Fatal("expected compact theme selector")
	}

	if len(selector.Options) != 3 {
		t.Fatalf("expected 3 theme options, got %d", len(selector.Options))
	}
}
