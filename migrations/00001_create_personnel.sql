-- +goose Up
CREATE TABLE personnel (
    id TEXT PRIMARY KEY,
    full_name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT personnel_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT personnel_full_name_not_empty CHECK (length(trim(full_name)) > 0)
);

CREATE INDEX idx_personnel_full_name ON personnel (full_name);

-- +goose Down
DROP TABLE personnel;