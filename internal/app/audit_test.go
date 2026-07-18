package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

type fakeAuditLogRepository struct {
	events  []domain.AuditEvent
	entries []ports.AuditEventEntry
}

func (r *fakeAuditLogRepository) Save(_ context.Context, event domain.AuditEvent) error {
	r.events = append(r.events, event)
	return nil
}

func (r *fakeAuditLogRepository) List(_ context.Context, limit int) ([]ports.AuditEventEntry, error) {
	if limit > len(r.entries) {
		limit = len(r.entries)
	}

	return r.entries[:limit], nil
}

func TestRecordAuditEventServiceExecute(t *testing.T) {
	repository := &fakeAuditLogRepository{}
	service := NewRecordAuditEventService(repository, fixedIDGenerator{id: "audit-1"})

	err := service.Execute(context.Background(), RecordAuditEventCommand{
		ActorOperatorID: "operator-1",
		EventType:       domain.AuditEventOperatorCreated,
		EntityType:      domain.AuditEntityOperator,
		EntityID:        "operator-2",
		Outcome:         domain.AuditOutcomeSuccess,
		Metadata: map[string]string{
			"role": "operator",
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repository.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(repository.events))
	}

	if repository.events[0].ID() != "audit-1" {
		t.Fatalf("expected audit-1, got %s", repository.events[0].ID())
	}
}

func TestListAuditEventsServiceExecute(t *testing.T) {
	occurredAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeAuditLogRepository{
		entries: []ports.AuditEventEntry{
			{
				ID:         "audit-1",
				EventType:  domain.AuditEventOperatorCreated,
				EntityType: domain.AuditEntityOperator,
				EntityID:   "operator-1",
				Outcome:    domain.AuditOutcomeSuccess,
				OccurredAt: occurredAt,
			},
		},
	}

	service := NewListAuditEventsService(repository)

	entries, err := service.Execute(context.Background(), ListAuditEventsCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}
