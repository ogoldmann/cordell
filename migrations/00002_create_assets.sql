-- +goose Up
CREATE TABLE assets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT assets_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT assets_name_not_empty CHECK (length(trim(name)) > 0)
);

CREATE INDEX idx_assets_name ON assets (name);

-- +goose Down
DROP TABLE assets