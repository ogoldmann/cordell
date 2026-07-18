-- +goose Up
CREATE TABLE audit_events (
    id TEXT PRIMARY KEY,
    actor_operator_id TEXT,
    event_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    outcome TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT audit_events_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT audit_events_event_type_not_empty CHECK (length(trim(event_type)) > 0),
    CONSTRAINT audit_events_entity_type_not_empty CHECK (length(trim(entity_type)) > 0),
    CONSTRAINT audit_events_entity_id_not_empty CHECK (length(trim(entity_id)) > 0),
    CONSTRAINT audit_events_outcome_not_empty CHECK (length(trim(outcome)) > 0),
    CONSTRAINT audit_events_actor_operator_id_fk
        FOREIGN KEY (actor_operator_id)
        REFERENCES operators (id)
        ON UPDATE RESTRICT
        ON DELETE SET NULL
);

CREATE INDEX idx_audit_events_occurred_at
    ON audit_events (occurred_at DESC);

CREATE INDEX idx_audit_events_actor_operator_id
    ON audit_events (actor_operator_id);

CREATE INDEX idx_audit_events_event_type
    ON audit_events (event_type);

CREATE INDEX idx_audit_events_entity
    ON audit_events (entity_type, entity_id);

-- +goose Down
DROP TABLE audit_events;