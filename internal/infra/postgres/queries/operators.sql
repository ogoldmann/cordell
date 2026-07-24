-- name: CreateOperator :exec
INSERT INTO operators (
    id,
    registration_id,
    alias,
    rank,
    role,
    password_hash
) VALUES (
    @id,
    @registration_id,
    @alias,
    @rank,
    @role,
    @password_hash
);

-- name: GetOperatorByID :one
SELECT
    id,
    registration_id,
    alias,
    rank,
    role,
    password_hash,
    active,
    created_at,
    updated_at
FROM operators
WHERE id = @id;

-- name: GetOperatorByRegistrationID :one
SELECT
    id,
    registration_id,
    alias,
    rank,
    role,
    password_hash,
    active,
    created_at,
    updated_at
FROM operators
WHERE registration_id = @registration_id;

-- name: GetOperatorSummaryByID :one
SELECT
    id,
    registration_id,
    alias,
    rank,
    role,
    active,
    created_at
FROM operators
WHERE id = @id;

-- name: ListOperators :many
WITH search_terms AS (
    SELECT unnest(sqlc.arg(search_patterns)::text[]) AS search_pattern
)
SELECT
    id,
    registration_id,
    alias,
    rank,
    role,
    active,
    created_at
FROM operators
WHERE (
    sqlc.arg(status_filter)::text = 'all'
    OR (sqlc.arg(status_filter)::text = 'active' AND active = true)
    OR (sqlc.arg(status_filter)::text = 'inactive' AND active = false)
)
AND (
    cardinality(sqlc.arg(search_patterns)::text[]) = 0
    OR NOT EXISTS (
        SELECT 1
        FROM search_terms
        WHERE NOT (
            alias ILIKE search_terms.search_pattern ESCAPE '\'
            OR registration_id ILIKE search_terms.search_pattern ESCAPE '\'
            OR rank ILIKE search_terms.search_pattern ESCAPE '\'
            OR role ILIKE search_terms.search_pattern ESCAPE '\'
        )
    )
)
ORDER BY created_at DESC, id DESC
LIMIT @limit_count;

-- name: CountActiveAdminOperators :one
SELECT count(*)::int
FROM operators
WHERE role = 'admin'
  AND active = true;

-- name: DeactivateOperator :one
WITH locked_active_admins AS (
    SELECT id
    FROM operators
    WHERE role = 'admin'
      AND active = true
    FOR UPDATE
),
active_admin_count AS (
    SELECT count(*) AS value
    FROM locked_active_admins
),
updated AS (
    UPDATE operators
    SET
        active = false,
        updated_at = now()
    WHERE operators.id = @id
      AND operators.active = true
      AND NOT (
          operators.role = 'admin'
          AND (SELECT value FROM active_admin_count) <= 1
      )
    RETURNING operators.id
)
SELECT count(*)::int
FROM updated;

-- name: ChangeOperatorRole :one
WITH locked_active_admins AS (
    SELECT id
    FROM operators
    WHERE role = 'admin'
      AND active = true
    FOR UPDATE
),
active_admin_count AS (
    SELECT count(*) AS value
    FROM locked_active_admins
),
updated AS (
    UPDATE operators
    SET
        role = sqlc.arg(role)::text,
        updated_at = now()
    WHERE operators.id = sqlc.arg(id)
      AND NOT (
          operators.role = 'admin'
          AND operators.active = true
          AND sqlc.arg(role)::text <> 'admin'
          AND (SELECT value FROM active_admin_count) <= 1
      )
    RETURNING operators.id
)
SELECT count(*)::int
FROM updated;

-- name: UpdateOperatorPasswordHash :one
WITH updated AS (
    UPDATE operators
    SET
        password_hash = @password_hash,
        updated_at = now()
    WHERE id = @id
      AND active = true
    RETURNING id
)
SELECT count(*)::int
FROM updated;

-- name: ReactivateOperator :one
WITH updated AS (
    UPDATE operators
    SET
        active = true,
        updated_at = now()
    WHERE id = @id
      AND active = false
    RETURNING id
)
SELECT count(*)::int
FROM updated;
