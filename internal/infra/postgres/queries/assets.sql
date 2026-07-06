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
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count);