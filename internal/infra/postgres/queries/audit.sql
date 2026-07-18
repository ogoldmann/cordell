-- name: CreateAuditEvent :exec
INSERT INTO audit_events (
    id,
    actor_operator_id,
    event_type,
    entity_type,
    entity_id,
    outcome,
    metadata
) VALUES (
    @id,
    sqlc.narg(actor_operator_id),
    @event_type,
    @entity_type,
    @entity_id,
    @outcome,
    @metadata
);

-- name: ListAuditEvents :many
SELECT
    ae.id,
    ae.actor_operator_id,
    COALESCE(o.alias, '') AS actor_alias,
    COALESCE(o.rank, '') AS actor_rank,
    ae.event_type,
    ae.entity_type,
    ae.entity_id,
    ae.outcome,
    ae.metadata,
    ae.occurred_at
FROM audit_events ae
LEFT JOIN operators o ON o.id = ae.actor_operator_id
ORDER BY ae.occurred_at DESC, ae.id DESC
LIMIT @limit_count;