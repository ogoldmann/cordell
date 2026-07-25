package web

import "testing"

func TestNewSearchBar(t *testing.T) {
	view := newSearchBar("search-input", "radio", "Pesquisar...")

	if view.ID != "search-input" {
		t.Fatalf("expected id search-input, got %q", view.ID)
	}

	if view.Name != "q" {
		t.Fatalf("expected name q, got %q", view.Name)
	}

	if view.Value != "radio" {
		t.Fatalf("expected value radio, got %q", view.Value)
	}

	if view.Placeholder != "Pesquisar..." {
		t.Fatalf("expected placeholder, got %q", view.Placeholder)
	}

	if view.Variant != searchBarVariantDefault {
		t.Fatalf("expected default variant, got %q", view.Variant)
	}

	if view.IconClass != "size-5" {
		t.Fatalf("expected size-5 icon, got %q", view.IconClass)
	}
}

func TestNewHeroSearchBar(t *testing.T) {
	view := newHeroSearchBar("hero-search", "", "Pesquisar...")

	if view.Variant != searchBarVariantHero {
		t.Fatalf("expected hero variant, got %q", view.Variant)
	}

	if !view.Autofocus {
		t.Fatal("expected hero search bar to autofocus")
	}

	if view.IconClass != "size-6" {
		t.Fatalf("expected size-6 icon, got %q", view.IconClass)
	}
}

func TestNewCompactSearchBar(t *testing.T) {
	view := newCompactSearchBar("compact-search", "", "Pesquisar...")

	if view.Variant != searchBarVariantCompact {
		t.Fatalf("expected compact variant, got %q", view.Variant)
	}

	if view.IconClass != "size-4" {
		t.Fatalf("expected size-4 icon, got %q", view.IconClass)
	}
}
