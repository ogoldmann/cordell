package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
)

func TestPostgresOperatorRepositorySaveAndFindByUsername(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	operator, err := domain.NewOperator(
		"operator-1",
		"Admin.User",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	found, err := operatorRepository.FindByUsername(context.Background(), "ADMIN.USER")
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", found.Role())
	}

	if found.ID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", found.ID())
	}

	if found.Username() != "admin.user" {
		t.Fatalf("expected admin.user, got %s", found.Username())
	}

}

func TestPostgresOperatorRepositoryRejectsDuplicateUsername(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	firstOperator, err := domain.NewOperator(
		"operator-1",
		"admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid first operator, got %v", err)
	}

	secondOperator, err := domain.NewOperator(
		"operator-2",
		"ADMIN",
		domain.OperatorRoleOperator,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid second operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), firstOperator); err != nil {
		t.Fatalf("expected no error saving first operator, got %v", err)
	}

	err = operatorRepository.Save(context.Background(), secondOperator)
	if err != domain.ErrDuplicateOperatorUsername {
		t.Fatalf("expected ErrDuplicateOperatorUsername, got %v", err)
	}
}
