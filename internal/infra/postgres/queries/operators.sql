-- name: CreateOperator :exec
INSERT INTO operators (
    id,
    username,
    password_hash
) VALUES (
    @id,
    @username,
    @password_hash
);

-- name: GetOperatorByID :one
SELECT
    id,
    username,
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
    password_hash,
    active,
    created_at,
    updated_at
FROM operators
WHERE username = @username;