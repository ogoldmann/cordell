package web

import (
	"testing"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestNewOperatorSummaryViewUsesOperationalDisplayName(t *testing.T) {
	operator := ports.OperatorSummary{
		ID:             "operator-1",
		RegistrationID: domain.RegistrationID("52998224725"),
		Alias:          "silva",
		Rank:           domain.RankThirdSergeant,
		Role:           domain.OperatorRoleOperator,
		Active:         true,
	}

	view := newOperatorSummaryView(operator)

	if view.DisplayName != "3º Sgt silva" {
		t.Fatalf("expected display name 3º Sgt silva, got %q", view.DisplayName)
	}

	if view.RankLabel != "Terceiro-Sargento (3º Sgt)" {
		t.Fatalf("expected full rank label, got %q", view.RankLabel)
	}
}
