-- name: CreatePersonnel :exec
INSERT INTO personnel (
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit
) VALUES (
    @id,
    @full_name,
    @alias,
    @rank,
    @registration_id,
    @section,
    @organization_unit
);

-- name: GetPersonnel :one
SELECT
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active,
    created_at,
    updated_at
FROM personnel
WHERE id = @id;

-- name: ListPersonnel :many
SELECT
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active,
    created_at,
    updated_at
FROM personnel
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count);
