-- name: CreateAsset :one
INSERT INTO assets (
    id,
    name,
    active
) VALUES (
    @id,
    @name,
    @active
)
RETURNING id, name, active, created_at, updated_at;

-- name: GetAsset :one
SELECT id, name, active, created_at, updated_at
FROM assets
WHERE id = @id;

-- name: ListAssets :many
SELECT id, name, active, created_at, updated_at
FROM assets
WHERE
    (
        @status_filter::text = 'all'
        OR (@status_filter::text = 'active' AND active = true)
        OR (@status_filter::text = 'inactive' AND active = false)
    )
ORDER BY created_at DESC, id DESC
LIMIT @limit_count;

-- name: SearchAssets :many
WITH search_terms AS (
    SELECT unnest(sqlc.arg(search_patterns)::text[]) AS search_pattern
)
SELECT
    id,
    name,
    active,
    created_at,
    updated_at
FROM assets
WHERE NOT EXISTS (
    SELECT 1
    FROM search_terms
    WHERE NOT (
        name ILIKE search_terms.search_pattern ESCAPE '\'
    )
)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count);
