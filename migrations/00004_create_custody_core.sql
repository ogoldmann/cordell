-- +goose Up
CREATE TABLE custody_transactions (
    id TEXT PRIMARY KEY,
    transaction_type TEXT NOT NULL,
    personnel_id TEXT NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT custody_transactions_id_not_empty CHECK (length(trim(id)) > 0),
    CONSTRAINT custody_transactions_type_valid CHECK (transaction_type IN ('checkout', 'return')),
    CONSTRAINT custody_transactions_personnel_id_fk
        FOREIGN KEY (personnel_id)
        REFERENCES personnel (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);

CREATE TABLE custody_lines (
    id BIGSERIAL PRIMARY KEY,
    custody_transaction_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT custody_lines_quantity_positive CHECK (quantity > 0),
    CONSTRAINT custody_lines_transaction_id_fk
        FOREIGN KEY (custody_transaction_id)
        REFERENCES custody_transactions (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT custody_lines_asset_id_fk
        FOREIGN KEY (asset_id)
        REFERENCES assets (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);

CREATE TABLE custody_balances (
    personnel_id TEXT NOT NULL,
    asset_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (personnel_id, asset_id),

    CONSTRAINT custody_balances_quantity_non_negative CHECK (quantity >= 0),
    CONSTRAINT custody_balances_personnel_id_fk
        FOREIGN KEY (personnel_id)
        REFERENCES personnel (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT custody_balances_asset_id_fk
        FOREIGN KEY (asset_id)
        REFERENCES assets (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);

CREATE INDEX idx_custody_transactions_personnel_created_at
    ON custody_transactions (personnel_id, created_at DESC);

CREATE INDEX idx_custody_lines_asset_id
    ON custody_lines (asset_id);

CREATE INDEX idx_custody_lines_transaction_id
    ON custody_lines (custody_transaction_id);

CREATE INDEX idx_custody_balances_asset_id
    ON custody_balances (asset_id);

-- +goose Down
DROP TABLE custody_balances;
DROP TABLE custody_lines;
DROP TABLE custody_transactions;