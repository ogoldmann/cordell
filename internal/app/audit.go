package app

import (
	"context"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

const defaultListAuditEventsLimit = 100
const maxListAuditEventsLimit = 500

type RecordAuditEventCommand struct {
	ActorOperatorID domain.OperatorID
	EventType       domain.AuditEventType
	EntityType      domain.AuditEntityType
	EntityID        string
	Outcome         domain.AuditOutcome
	Metadata        map[string]string
}

type RecordAuditEventService struct {
	auditLogRepository ports.AuditLogRepository
	idGenerator        ports.IDGenerator
}

func NewRecordAuditEventService(
	auditLogRepository ports.AuditLogRepository,
	idGenerator ports.IDGenerator,
) *RecordAuditEventService {
	return &RecordAuditEventService{
		auditLogRepository: auditLogRepository,
		idGenerator:        idGenerator,
	}
}

func (s *RecordAuditEventService) Execute(ctx context.Context, cmd RecordAuditEventCommand) error {
	id, err := s.idGenerator.NewID()
	if err != nil {
		return err
	}

	event, err := domain.NewAuditEvent(
		domain.AuditEventID(id),
		cmd.ActorOperatorID,
		cmd.EventType,
		cmd.EntityType,
		cmd.EntityID,
		cmd.Outcome,
		cmd.Metadata,
	)
	if err != nil {
		return err
	}

	return s.auditLogRepository.Save(ctx, event)
}

type ListAuditEventsCommand struct {
	Limit int
}

type ListAuditEventsService struct {
	auditLogRepository ports.AuditLogRepository
}

func NewListAuditEventsService(auditLogRepository ports.AuditLogRepository) *ListAuditEventsService {
	return &ListAuditEventsService{
		auditLogRepository: auditLogRepository,
	}
}

func (s *ListAuditEventsService) Execute(ctx context.Context, cmd ListAuditEventsCommand) ([]ports.AuditEventEntry, error) {
	limit := normalizeListAuditEventsLimit(cmd.Limit)

	return s.auditLogRepository.List(ctx, limit)
}

func normalizeListAuditEventsLimit(limit int) int {
	if limit <= 0 {
		return defaultListAuditEventsLimit
	}

	if limit > maxListAuditEventsLimit {
		return maxListAuditEventsLimit
	}

	return limit
}
