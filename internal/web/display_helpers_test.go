package web

import (
	"testing"

	"cordell/internal/domain"
)

func TestMilitaryDisplayNameUsesRankAbbreviation(t *testing.T) {
	got := militaryDisplayName(domain.RankSoldier, "John")

	if got != "Sd John" {
		t.Fatalf("expected Sd John, got %q", got)
	}
}

func TestMilitaryDisplayNameUsesThirdSergeantAbbreviation(t *testing.T) {
	got := militaryDisplayName(domain.RankThirdSergeant, "Silva")

	if got != "3º Sgt Silva" {
		t.Fatalf("expected 3º Sgt Silva, got %q", got)
	}
}

func TestMilitaryDisplayNameWithoutAliasUsesRankAbbreviation(t *testing.T) {
	got := militaryDisplayName(domain.RankSoldier, "")

	if got != "Sd" {
		t.Fatalf("expected Sd, got %q", got)
	}
}
