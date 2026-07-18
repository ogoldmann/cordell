package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
)

func TestPostgresAuditLogRepositorySaveAndList(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)
	auditLogRepository := postgres.NewAuditLogRepository(queries)

	actorOperator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)
	if err := operatorRepository.Save(context.Background(), actorOperator); err != nil {
		t.Fatalf("expected no error saving actor operator, got %v", err)
	}

	event, err := domain.NewAuditEvent(
		"audit-1",
		actorOperator.ID(),
		domain.AuditEventOperatorCreated,
		domain.AuditEntityOperator,
		"operator-2",
		domain.AuditOutcomeSuccess,
		map[string]string{
			"role": "operator",
		},
	)
	if err != nil {
		t.Fatalf("expected valid audit event, got %v", err)
	}

	if err := auditLogRepository.Save(context.Background(), event); err != nil {
		t.Fatalf("expected no error saving audit event, got %v", err)
	}

	events, err := auditLogRepository.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error listing audit events, got %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].ID != "audit-1" {
		t.Fatalf("expected audit-1, got %s", events[0].ID)
	}

	if events[0].ActorOperatorID != actorOperator.ID() {
		t.Fatalf("expected actor operator %s, got %s", actorOperator.ID(), events[0].ActorOperatorID)
	}

	if events[0].ActorAlias != "silva" {
		t.Fatalf("expected actor alias silva, got %s", events[0].ActorAlias)
	}

	if events[0].ActorRank != domain.RankSergeant {
		t.Fatalf("expected actor rank %s, got %s", domain.RankSergeant, events[0].ActorRank)
	}

	if events[0].EventType != domain.AuditEventOperatorCreated {
		t.Fatalf("expected event type %s, got %s", domain.AuditEventOperatorCreated, events[0].EventType)
	}

	if events[0].EntityType != domain.AuditEntityOperator {
		t.Fatalf("expected entity type %s, got %s", domain.AuditEntityOperator, events[0].EntityType)
	}

	if events[0].EntityID != "operator-2" {
		t.Fatalf("expected entity operator-2, got %s", events[0].EntityID)
	}

	if events[0].Outcome != domain.AuditOutcomeSuccess {
		t.Fatalf("expected outcome %s, got %s", domain.AuditOutcomeSuccess, events[0].Outcome)
	}

	if events[0].Metadata["role"] != "operator" {
		t.Fatalf("expected metadata role operator, got %s", events[0].Metadata["role"])
	}

	if events[0].OccurredAt.IsZero() {
		t.Fatal("expected occurred_at not to be zero")
	}
}
