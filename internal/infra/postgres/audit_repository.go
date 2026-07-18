package postgres

import (
	"context"
	"encoding/json"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuditLogRepository struct {
	queries *db.Queries
}

func NewAuditLogRepository(queries *db.Queries) *AuditLogRepository {
	return &AuditLogRepository{
		queries: queries,
	}
}

func (r *AuditLogRepository) Save(ctx context.Context, event domain.AuditEvent) error {
	metadata, err := json.Marshal(event.Metadata())
	if err != nil {
		return err
	}

	var actorOperatorID pgtype.Text
	if event.ActorOperatorID() != "" {
		actorOperatorID = pgtype.Text{
			String: string(event.ActorOperatorID()),
			Valid:  true,
		}
	}

	return r.queries.CreateAuditEvent(ctx, db.CreateAuditEventParams{
		ID:              string(event.ID()),
		ActorOperatorID: actorOperatorID,
		EventType:       string(event.Type()),
		EntityType:      string(event.EntityType()),
		EntityID:        event.EntityID(),
		Outcome:         string(event.Outcome()),
		Metadata:        metadata,
	})
}

func (r *AuditLogRepository) List(ctx context.Context, limit int) ([]ports.AuditEventEntry, error) {
	rows, err := r.queries.ListAuditEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	entries := make([]ports.AuditEventEntry, 0, len(rows))

	for _, row := range rows {
		metadata := map[string]string{}
		if len(row.Metadata) > 0 {
			if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
				return nil, err
			}
		}

		var actorOperatorID domain.OperatorID
		if row.ActorOperatorID.Valid {
			actorOperatorID = domain.OperatorID(row.ActorOperatorID.String)
		}

		occurredAt, err := timestamptzToTime(row.OccurredAt)
		if err != nil {
			return nil, err
		}

		entries = append(entries, ports.AuditEventEntry{
			ID:              domain.AuditEventID(row.ID),
			ActorOperatorID: actorOperatorID,
			ActorAlias:      row.ActorAlias,
			ActorRank:       domain.Rank(row.ActorRank),
			EventType:       domain.AuditEventType(row.EventType),
			EntityType:      domain.AuditEntityType(row.EntityType),
			EntityID:        row.EntityID,
			Outcome:         domain.AuditOutcome(row.Outcome),
			Metadata:        metadata,
			OccurredAt:      occurredAt,
		})
	}

	return entries, nil
}
