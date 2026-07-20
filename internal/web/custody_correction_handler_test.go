package web

import (
	"testing"

	"cordell/internal/domain"
)

func TestHumanizeCustodyCorrectionErrorExplainsInsufficientBalance(t *testing.T) {
	message := humanizeCustodyCorrectionError(domain.ErrInsufficientCustodyBalance)

	if message == "" {
		t.Fatal("expected message")
	}

	if message == "This edit cannot be applied because it would make a custody balance negative." {
		t.Fatal("expected expanded insufficient balance explanation")
	}
}
