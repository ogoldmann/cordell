package web

import (
	"context"
	"testing"

	"cordell/internal/domain"
)

func TestCurrentOperatorFromContext(t *testing.T) {
	operator, err := domain.NewOperator("operator-1", "admin", "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	ctx := withCurrentOperator(context.Background(), operator)

	currentOperator, ok := currentOperatorFromContext(ctx)
	if !ok {
		t.Fatal("expected current operator in context")
	}

	if currentOperator.ID() != operator.ID() {
		t.Fatalf("expected operator %s, got %s", operator.ID(), currentOperator.ID())
	}
}

func TestCurrentOperatorFromContextReturnsFalseWhenMissing(t *testing.T) {
	_, ok := currentOperatorFromContext(context.Background())
	if ok {
		t.Fatal("expected no current operator")
	}
}
