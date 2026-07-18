package domain

import "testing"

func TestNewAuditEvent(t *testing.T) {
	event, err := NewAuditEvent(
		"audit-1",
		"operator-1",
		AuditEventOperatorCreated,
		AuditEntityOperator,
		"operator-2",
		AuditOutcomeSuccess,
		map[string]string{
			"role": "operator",
		},
	)
	if err != nil {
		t.Fatalf("expected valid audit event, got %v", err)
	}

	if event.ID() != "audit-1" {
		t.Fatalf("expected audit-1, got %s", event.ID())
	}

	if event.ActorOperatorID() != "operator-1" {
		t.Fatalf("expected actor operator-1, got %s", event.ActorOperatorID())
	}

	if event.EntityID() != "operator-2" {
		t.Fatalf("expected entity operator-2, got %s", event.EntityID())
	}
}

func TestNewAuditEventRejectsEmptyID(t *testing.T) {
	_, err := NewAuditEvent(
		"",
		"operator-1",
		AuditEventOperatorCreated,
		AuditEntityOperator,
		"operator-2",
		AuditOutcomeSuccess,
		nil,
	)
	if err != ErrEmptyAuditEventID {
		t.Fatalf("expected ErrEmptyAuditEventID, got %v", err)
	}
}

func TestAuditEventMetadataReturnsCopy(t *testing.T) {
	event, err := NewAuditEvent(
		"audit-1",
		"operator-1",
		AuditEventOperatorCreated,
		AuditEntityOperator,
		"operator-2",
		AuditOutcomeSuccess,
		map[string]string{
			"role": "operator",
		},
	)
	if err != nil {
		t.Fatalf("expected valid audit event, got %v", err)
	}

	metadata := event.Metadata()
	metadata["role"] = "admin"

	if event.Metadata()["role"] != "operator" {
		t.Fatal("expected audit event metadata to be immutable from caller")
	}
}

func TestNewAuditEventSanitizesMetadataValues(t *testing.T) {
	event, err := NewAuditEvent(
		"audit-1",
		"operator-1",
		AuditEventOperatorCreated,
		AuditEntityOperator,
		"operator-2",
		AuditOutcomeSuccess,
		map[string]string{
			"note": "line one\nline two\rline three",
		},
	)
	if err != nil {
		t.Fatalf("expected valid audit event, got %v", err)
	}

	if event.Metadata()["note"] != "line one line two line three" {
		t.Fatalf("expected sanitized metadata, got %q", event.Metadata()["note"])
	}
}
