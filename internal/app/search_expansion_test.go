package app

import "testing"

func TestCanonicalSearchQueryExpandsRankAbbreviation(t *testing.T) {
	got := CanonicalSearchQuery("sd john")

	if got != "private john" {
		t.Fatalf("expected private john, got %q", got)
	}
}

func TestCanonicalSearchQueryExpandsCorporalAbbreviation(t *testing.T) {
	got := CanonicalSearchQuery("cb silva")

	if got != "corporal silva" {
		t.Fatalf("expected corporal silva, got %q", got)
	}
}

func TestCanonicalSearchQueryExpandsThirdSergeantAbbreviation(t *testing.T) {
	got := CanonicalSearchQuery("3º sgt silva")

	if got != "sergeant silva" {
		t.Fatalf("expected sergeant silva, got %q", got)
	}
}

func TestCanonicalSearchQueryExpandsSectionAbbreviation(t *testing.T) {
	got := CanonicalSearchQuery("s1 john")

	if got != "personnel john" {
		t.Fatalf("expected personnel john, got %q", got)
	}
}

func TestCanonicalSearchQueryExpandsSupplyAbbreviation(t *testing.T) {
	got := CanonicalSearchQuery("almx radio")

	if got != "supply radio" {
		t.Fatalf("expected supply radio, got %q", got)
	}
}
