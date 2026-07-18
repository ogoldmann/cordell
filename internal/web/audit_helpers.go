package web

import (
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
)

func (s *Server) recordAuditEvent(
	r *http.Request,
	eventType domain.AuditEventType,
	entityType domain.AuditEntityType,
	entityID string,
	metadata map[string]string,
) error {
	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		return domain.ErrEmptyOperatorID
	}

	return s.services.RecordAuditEvent.Execute(r.Context(), app.RecordAuditEventCommand{
		ActorOperatorID: currentOperator.ID(),
		EventType:       eventType,
		EntityType:      entityType,
		EntityID:        entityID,
		Outcome:         domain.AuditOutcomeSuccess,
		Metadata:        metadata,
	})
}

func (s *Server) recordAuditEventOrLog(
	r *http.Request,
	eventType domain.AuditEventType,
	entityType domain.AuditEntityType,
	entityID string,
	metadata map[string]string,
) {
	if err := s.recordAuditEvent(r, eventType, entityType, entityID, metadata); err != nil {
		s.logger.Error(
			"failed to record audit event",
			"error", err,
			"event_type", eventType,
			"entity_type", entityType,
			"entity_id", entityID,
		)
	}
}
