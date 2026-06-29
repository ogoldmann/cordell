package domain

import "testing"

func TestNewQuantity(t *testing.T) {
	quantity, err := NewQuantity(3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if quantity.Int() != 3 {
		t.Fatalf("expected quantity 3, got %d", quantity.Int())
	}
}

func TestNewQuantityRejectsZero(t *testing.T) {
	_, err := NewQuantity(0)
	if err != ErrInvalidQuantity {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestNewQuantityRejectsNegativeValue(t *testing.T) {
	_, err := NewQuantity(-1)
	if err != ErrInvalidQuantity {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}
