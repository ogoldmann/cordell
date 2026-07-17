-- +goose Up
CREATE TABLE operators (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    role TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT operators_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT operators_username_not_empty CHECK (length(trim(username)) > 0),
    CONSTRAINT operators_role_not_empty CHECK (length(trim(role)) > 0),
    CONSTRAINT operators_role_valid CHECK (role IN ('admin', 'operator')),
    CONSTRAINT operators_password_hash_not_empty CHECK (length(trim(password_hash)) > 0),
    CONSTRAINT operators_username_unique UNIQUE (username)
);

CREATE INDEX operators_active_idx ON operators (active);
CREATE INDEX operators_role_idx ON operators (role);

-- +goose Down
DROP TABLE operators;