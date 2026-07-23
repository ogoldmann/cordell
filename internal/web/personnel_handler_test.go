package web

import (
	"testing"

	"cordell/internal/domain"
)

func TestNewPersonnelViewUsesOperationalDisplayName(t *testing.T) {
	personnel, err := domain.NewPersonnel(
		"personnel-1",
		"John Doe",
		"John",
		domain.PersonnelRankSoldier,
		domain.RegistrationID("52998224725"),
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	view := newPersonnelView(personnel)

	if view.DisplayName != "Sd John" {
		t.Fatalf("expected display name Sd John, got %q", view.DisplayName)
	}

	if view.RankLabel != "Soldado (Sd)" {
		t.Fatalf("expected rank label Soldado (Sd), got %q", view.RankLabel)
	}
}
