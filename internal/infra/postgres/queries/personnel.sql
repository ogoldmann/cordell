-- name: CreatePersonnel :one
INSERT INTO personnel (
    id,
    full_name,
    active
) VALUES (
    @id,
    @full_name,
    @active
)
RETURNING id, full_name, active, created_at, updated_at;

-- name: GetPersonnel :one
SELECT id, full_name, active, created_at, updated_at
FROM personnel
WHERE id = @id;

-- name: ListPersonnel :many
SELECT id, full_name, active, created_at, updated_at
FROM personnel
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count);
