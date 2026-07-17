package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestListOperatorsServiceExecute(t *testing.T) {
	createdAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeOperatorRepository{
		summaries: []ports.OperatorSummary{
			{
				ID:        "operator-1",
				Username:  "admin",
				Role:      domain.OperatorRoleAdmin,
				Active:    true,
				CreatedAt: createdAt,
			},
			{
				ID:        "operator-2",
				Username:  "clerk",
				Role:      domain.OperatorRoleOperator,
				Active:    true,
				CreatedAt: createdAt,
			},
		},
	}

	service := NewListOperatorsService(repository)

	operators, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operators) != 2 {
		t.Fatalf("expected 2 operators, got %d", len(operators))
	}

	if operators[0].Username != "admin" {
		t.Fatalf("expected admin, got %s", operators[0].Username)
	}
}

func TestListOperatorsServiceLimitsMaximum(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewListOperatorsService(repository)

	_, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: maxListOperatorsLimit + 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
