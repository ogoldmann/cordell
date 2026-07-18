package web

import (
	"encoding/json"
	"time"

	"cordell/internal/ports"
)

type auditEventView struct {
	ID           string
	ActorDisplay string
	EventType    string
	EntityType   string
	EntityID     string
	Outcome      string
	Metadata     string
	OccurredAt   string
}

func newAuditEventView(entry ports.AuditEventEntry) auditEventView {
	metadata := "{}"
	if len(entry.Metadata) > 0 {
		encoded, err := json.MarshalIndent(entry.Metadata, "", "  ")
		if err == nil {
			metadata = string(encoded)
		}
	}

	actorDisplay := "System"
	if entry.ActorOperatorID != "" {
		actorDisplay = militaryDisplayName(entry.ActorRank, entry.ActorAlias)
	}

	return auditEventView{
		ID:           string(entry.ID),
		ActorDisplay: actorDisplay,
		EventType:    string(entry.EventType),
		EntityType:   string(entry.EntityType),
		EntityID:     entry.EntityID,
		Outcome:      string(entry.Outcome),
		Metadata:     metadata,
		OccurredAt:   formatAuditTimestamp(entry.OccurredAt),
	}
}

func formatAuditTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("2006-01-02 15:04")
}
