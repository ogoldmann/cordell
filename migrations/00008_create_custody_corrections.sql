-- +goose Up
CREATE TABLE custody_corrections (
    id TEXT PRIMARY KEY,
    corrected_transaction_id TEXT NOT NULL,
    operator_id TEXT NOT NULL,
    corrected_personnel_id TEXT NOT NULL,
    corrected_notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT custody_corrections_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT custody_corrections_corrected_transaction_id_fk
        FOREIGN KEY (corrected_transaction_id)
        REFERENCES custody_transactions (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT custody_corrections_operator_id_fk
        FOREIGN KEY (operator_id)
        REFERENCES operators (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT custody_corrections_corrected_personnel_id_fk
        FOREIGN KEY (corrected_personnel_id)
        REFERENCES personnel (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);

CREATE TABLE custody_correction_lines (
    id BIGSERIAL PRIMARY KEY,
    custody_correction_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT custody_correction_lines_quantity_positive CHECK (quantity > 0),
    CONSTRAINT custody_correction_lines_correction_id_fk
        FOREIGN KEY (custody_correction_id)
        REFERENCES custody_corrections (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT custody_correction_lines_asset_id_fk
        FOREIGN KEY (asset_id)
        REFERENCES assets (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);

CREATE INDEX idx_custody_corrections_transaction_created_at
    ON custody_corrections (corrected_transaction_id, created_at DESC, id DESC);

CREATE INDEX idx_custody_correction_lines_correction_id
    ON custody_correction_lines (custody_correction_id);

CREATE INDEX idx_custody_correction_lines_asset_id
    ON custody_correction_lines (asset_id);

-- +goose Down
DROP TABLE custody_correction_lines;
DROP TABLE custody_corrections;
