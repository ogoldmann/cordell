-- +goose Up
CREATE TABLE personnel (
    id TEXT PRIMARY KEY,
    full_name TEXT NOT NULL,
    alias TEXT NOT NULL,
    rank TEXT NOT NULL,
    registration_id TEXT NOT NULL,
    section TEXT NOT NULL,
    organization_unit TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT personnel_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT personnel_full_name_not_empty CHECK (length(trim(full_name)) > 0),
    CONSTRAINT personnel_alias_not_empty CHECK (length(trim(alias)) > 0),
    CONSTRAINT personnel_rank_not_empty CHECK (length(trim(rank)) > 0),
    CONSTRAINT personnel_registration_id_not_empty CHECK (length(trim(registration_id)) > 0),
    CONSTRAINT personnel_section_not_empty CHECK (length(trim(section)) > 0),
    CONSTRAINT personnel_organization_unit_not_empty CHECK (length(trim(organization_unit)) > 0),
    CONSTRAINT personnel_registration_id_unique UNIQUE (registration_id)
);

CREATE INDEX personnel_full_name_idx ON personnel (full_name);
CREATE INDEX personnel_alias_idx ON personnel (alias);
CREATE INDEX personnel_rank_idx ON personnel (rank);
CREATE INDEX personnel_section_idx ON personnel (section);
CREATE INDEX personnel_organization_unit_idx ON personnel (organization_unit);

-- +goose Down
DROP TABLE personnel;