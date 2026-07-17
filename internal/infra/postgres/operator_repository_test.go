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

func TestPostgresOperatorRepositoryList(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := domain.NewOperator(
		"operator-1",
		"admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := domain.NewOperator(
		"operator-2",
		"clerk",
		domain.OperatorRoleOperator,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	operators, err := operatorRepository.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error listing operators, got %v", err)
	}

	if len(operators) != 2 {
		t.Fatalf("expected 2 operators, got %d", len(operators))
	}

	for _, operator := range operators {
		if operator.Username == "" {
			t.Fatal("expected username not to be empty")
		}

		if operator.Role == "" {
			t.Fatal("expected role not to be empty")
		}

		if operator.CreatedAt.IsZero() {
			t.Fatal("expected created_at not to be zero")
		}
	}
}
