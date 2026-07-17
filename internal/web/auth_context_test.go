package web

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
)

func TestCurrentOperatorFromContext(t *testing.T) {
	operator := mustBuildOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	ctx := withCurrentOperator(context.Background(), operator)

	currentOperator, ok := currentOperatorFromContext(ctx)
	if !ok {
		t.Fatal("expected current operator in context")
	}

	if currentOperator.ID() != operator.ID() {
		t.Fatalf("expected operator %s, got %s", operator.ID(), currentOperator.ID())
	}
}

func mustBuildOperator(
	t *testing.T,
	id domain.OperatorID,
	registrationIDValue string,
	alias string,
	rank domain.Rank,
	role domain.OperatorRole,
) domain.Operator {
	t.Helper()

	registrationID, err := domain.NewRegistrationID(registrationIDValue)
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	operator, err := domain.NewOperator(
		id,
		registrationID,
		alias,
		rank,
		role,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	return operator
}

func TestCurrentOperatorFromContextReturnsFalseWhenMissing(t *testing.T) {
	_, ok := currentOperatorFromContext(context.Background())
	if ok {
		t.Fatal("expected no current operator")
	}
}

func TestCurrentSessionFromContext(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	session, err := domain.NewOperatorSession(
		"session-1",
		"operator-1",
		"token-hash",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	ctx := withCurrentSession(context.Background(), session)

	currentSession, ok := currentSessionFromContext(ctx)
	if !ok {
		t.Fatal("expected current session in context")
	}

	if currentSession.ID() != session.ID() {
		t.Fatalf("expected session %s, got %s", session.ID(), currentSession.ID())
	}
}
