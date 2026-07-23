package web

import (
	"testing"

	"cordell/internal/domain"
)

func TestActiveStatusLabel(t *testing.T) {
	if got := activeStatusLabel(true); got != "Ativo" {
		t.Fatalf("expected Ativo, got %q", got)
	}

	if got := activeStatusLabel(false); got != "Inativo" {
		t.Fatalf("expected Inativo, got %q", got)
	}
}

func TestCustodyTransactionTypeLabel(t *testing.T) {
	if got := custodyTransactionTypeLabel(domain.CustodyTransactionTypeCheckout); got != "Cautela" {
		t.Fatalf("expected Cautela, got %q", got)
	}

	if got := custodyTransactionTypeLabel(domain.CustodyTransactionTypeReturn); got != "Descautela" {
		t.Fatalf("expected Descautela, got %q", got)
	}
}

func TestOperatorRoleLabel(t *testing.T) {
	if got := operatorRoleLabel(domain.OperatorRoleAdmin); got != "Administrador" {
		t.Fatalf("expected Administrador, got %q", got)
	}

	if got := operatorRoleLabel(domain.OperatorRoleOperator); got != "Operador" {
		t.Fatalf("expected Operador, got %q", got)
	}
}
