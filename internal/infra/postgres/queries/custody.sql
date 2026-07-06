-- name: CreateCustodyTransaction :one
INSERT INTO custody_transactions (
    id,
    transaction_type,
    personnel_id,
    notes
) VALUES (
    @id,
    @transaction_type,
    @personnel_id,
    @notes
)
RETURNING id, transaction_type, personnel_id, notes, created_at;

-- name: CreateCustodyLine :one
INSERT INTO custody_lines (
    custody_transaction_id,
    asset_id,
    quantity
) VALUES (
    @custody_transaction_id,
    @asset_id,
    @quantity
)
RETURNING id, custody_transaction_id, asset_id, quantity, created_at;

-- name: GetCustodyBalanceQuantity :one
SELECT quantity
FROM custody_balances
WHERE personnel_id = @personnel_id
  AND asset_id = @asset_id;

-- name: IncreaseCustodyBalanceForCheckout :exec
INSERT INTO custody_balances (
    personnel_id,
    asset_id,
    quantity
) VALUES (
    @personnel_id,
    @asset_id,
    @quantity
)
ON CONFLICT (personnel_id, asset_id)
DO UPDATE SET
    quantity = custody_balances.quantity + EXCLUDED.quantity,
    updated_at = now();

-- name: DecreaseCustodyBalanceForReturn :execrows
UPDATE custody_balances
SET
    quantity = quantity - @quantity,
    updated_at = now()
WHERE personnel_id = @personnel_id
  AND asset_id = @asset_id
  AND quantity >= @quantity;

-- name: ListCurrentCustodyByPersonnel :many
SELECT
    cb.personnel_id,
    cb.asset_id,
    a.name AS asset_name,
    cb.quantity,
    cb.updated_at
FROM custody_balances cb
JOIN assets a ON a.id = cb.asset_id
WHERE cb.personnel_id = @personnel_id
  AND cb.quantity > 0
ORDER BY a.name ASC, cb.asset_id ASC;
