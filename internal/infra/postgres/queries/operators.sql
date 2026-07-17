-- name: CreateOperator :exec
INSERT INTO operators (
    id,
    username,
    role,
    password_hash
) VALUES (
    @id,
    @username,
    @role,
    @password_hash
);

-- name: GetOperatorByID :one
SELECT
    id,
    username,
    role,
    password_hash,
    active,
    created_at,
    updated_at
FROM operators
WHERE id = @id;

-- name: GetOperatorByUsername :one
SELECT
    id,
    username,
    role,
    password_hash,
    active,
    created_at,
    updated_at
FROM operators
WHERE username = @username;

-- name: ListOperators :many
SELECT
    id,
    username,
    role,
    active,
    created_at
FROM operators
ORDER BY created_at DESC, id DESC
LIMIT @limit_count;