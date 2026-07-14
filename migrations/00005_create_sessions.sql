-- +goose Up
CREATE TABLE operator_sessions (
    id TEXT PRIMARY KEY,
    operator_id TEXT NOT NULL REFERENCES operators(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT operator_sessions_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT operator_sessions_token_hash_not_empty CHECK (length(trim(token_hash)) > 0),
    CONSTRAINT operator_sessions_token_hash_unique UNIQUE (token_hash)
);

CREATE INDEX operator_sessions_operator_id_idx ON operator_sessions (operator_id);
CREATE INDEX operator_sessions_expires_at_idx ON operator_sessions (expires_at);

-- +goose Down
DROP TABLE operator_sessions;