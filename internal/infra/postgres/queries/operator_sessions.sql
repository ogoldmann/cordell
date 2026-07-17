-- name: CreateOperatorSession :exec
INSERT INTO operator_sessions (
    id,
    operator_id,
    token_hash,
    csrf_token,
    expires_at,
    created_at
) VALUES (
    @id,
    @operator_id,
    @token_hash,
    @csrf_token,
    @expires_at,
    @created_at
);

-- name: GetOperatorSessionByTokenHash :one
SELECT
    id,
    operator_id,
    token_hash,
    csrf_token,
    expires_at,
    created_at
FROM operator_sessions
WHERE token_hash = @token_hash;

-- name: DeleteOperatorSessionByTokenHash :exec
DELETE FROM operator_sessions
WHERE token_hash = @token_hash;

-- name: DeleteExpiredOperatorSessions :exec
DELETE FROM operator_sessions
WHERE expires_at <= @now;

-- name: DeleteOperatorSessionsByOperatorID :exec
DELETE FROM operator_sessions
WHERE operator_id = @operator_id;