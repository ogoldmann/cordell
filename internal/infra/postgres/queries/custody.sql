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
