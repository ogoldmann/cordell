package domain

import "strings"

type AuditEventID string

type AuditEventType string

const (
	AuditEventOperatorCreated        AuditEventType = "operator.created"
	AuditEventOperatorDeactivated    AuditEventType = "operator.deactivated"
	AuditEventOperatorReactivated    AuditEventType = "operator.reactivated"
	AuditEventOperatorRoleChanged    AuditEventType = "operator.role_changed"
	AuditEventOperatorPasswordReset  AuditEventType = "operator.password_reset"
	AuditEventCustodyCheckoutCreated AuditEventType = "custody.checkout_created"
	AuditEventCustodyReturnCreated   AuditEventType = "custody.return_created"
)

type AuditEntityType string

const (
	AuditEntityOperator           AuditEntityType = "operator"
	AuditEntityCustodyTransaction AuditEntityType = "custody_transaction"
)

type AuditOutcome string

const (
	AuditOutcomeSuccess AuditOutcome = "success"
	AuditOutcomeFailure AuditOutcome = "failure"
)

type AuditEvent struct {
	id              AuditEventID
	actorOperatorID OperatorID
	eventType       AuditEventType
	entityType      AuditEntityType
	entityID        string
	outcome         AuditOutcome
	metadata        map[string]string
}

func NewAuditEvent(
	id AuditEventID,
	actorOperatorID OperatorID,
	eventType AuditEventType,
	entityType AuditEntityType,
	entityID string,
	outcome AuditOutcome,
	metadata map[string]string,
) (AuditEvent, error) {
	if strings.TrimSpace(string(id)) == "" {
		return AuditEvent{}, ErrEmptyAuditEventID
	}

	if strings.TrimSpace(string(eventType)) == "" {
		return AuditEvent{}, ErrEmptyAuditEventType
	}

	if strings.TrimSpace(string(entityType)) == "" {
		return AuditEvent{}, ErrEmptyAuditEntityType
	}

	entityID = strings.TrimSpace(entityID)
	if entityID == "" {
		return AuditEvent{}, ErrEmptyAuditEntityID
	}

	if strings.TrimSpace(string(outcome)) == "" {
		return AuditEvent{}, ErrEmptyAuditOutcome
	}

	copiedMetadata := make(map[string]string, len(metadata))
	for key, value := range metadata {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		copiedMetadata[key] = sanitizeAuditMetadataValue(value)
	}

	return AuditEvent{
		id:              id,
		actorOperatorID: actorOperatorID,
		eventType:       eventType,
		entityType:      entityType,
		entityID:        entityID,
		outcome:         outcome,
		metadata:        copiedMetadata,
	}, nil
}

func (e AuditEvent) ID() AuditEventID {
	return e.id
}

func (e AuditEvent) ActorOperatorID() OperatorID {
	return e.actorOperatorID
}

func (e AuditEvent) Type() AuditEventType {
	return e.eventType
}

func (e AuditEvent) EntityType() AuditEntityType {
	return e.entityType
}

func (e AuditEvent) EntityID() string {
	return e.entityID
}

func (e AuditEvent) Outcome() AuditOutcome {
	return e.outcome
}

func (e AuditEvent) Metadata() map[string]string {
	copiedMetadata := make(map[string]string, len(e.metadata))
	for key, value := range e.metadata {
		copiedMetadata[key] = value
	}

	return copiedMetadata
}

func sanitizeAuditMetadataValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")

	return value
}
