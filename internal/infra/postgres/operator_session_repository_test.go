package postgres_test

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
	"cordell/internal/ports"
)

func TestPostgresOperatorSessionRepositoryDeleteByOperatorID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)
	sessionRepository := postgres.NewOperatorSessionRepository(queries)

	operator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	session, err := domain.NewOperatorSession(
		"session-1",
		operator.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	if err := sessionRepository.Save(context.Background(), session); err != nil {
		t.Fatalf("expected no error saving session, got %v", err)
	}

	if err := sessionRepository.DeleteByOperatorID(context.Background(), operator.ID()); err != nil {
		t.Fatalf("expected no error deleting sessions, got %v", err)
	}

	_, err = sessionRepository.FindByTokenHash(context.Background(), session.TokenHash())
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
