-- +goose Up
CREATE TABLE operators (
    id TEXT PRIMARY KEY,
    registration_id TEXT NOT NULL,
    alias TEXT NOT NULL,
    rank TEXT NOT NULL,
    role TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT operators_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT operators_registration_id_not_empty CHECK (length(trim(registration_id)) > 0),
    CONSTRAINT operators_alias_not_empty CHECK (length(trim(alias)) > 0),
    CONSTRAINT operators_rank_not_empty CHECK (length(trim(rank)) > 0),
    CONSTRAINT operators_role_not_empty CHECK (length(trim(role)) > 0),
    CONSTRAINT operators_role_valid CHECK (role IN ('admin', 'operator')),
    CONSTRAINT operators_password_hash_not_empty CHECK (length(trim(password_hash)) > 0),
    CONSTRAINT operators_registration_id_unique UNIQUE (registration_id)
);

CREATE INDEX operators_active_idx ON operators (active);
CREATE INDEX operators_rank_idx ON operators (rank);
CREATE INDEX operators_alias_idx ON operators (alias);
CREATE INDEX operators_role_idx ON operators (role);

-- +goose Down
DROP TABLE operators;