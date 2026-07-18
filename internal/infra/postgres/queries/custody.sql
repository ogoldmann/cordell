-- name: CreateCustodyTransaction :one
INSERT INTO custody_transactions (
    id,
    transaction_type,
    personnel_id,
    operator_id,
    notes
) VALUES (
    @id,
    @transaction_type,
    @personnel_id,
    @operator_id,
    @notes
)
RETURNING id, transaction_type, personnel_id, operator_id, notes, created_at;

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

-- name: ListCustodyHistoryByPersonnel :many
WITH recent_transactions AS (
    SELECT
        id,
        transaction_type,
        personnel_id,
        operator_id,
        notes,
        created_at
    FROM custody_transactions
    WHERE personnel_id = @personnel_id
    ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg(limit_count)
)
SELECT
    rt.id AS transaction_id,
    rt.transaction_type,
    rt.personnel_id,
    rt.operator_id,
    o.alias AS operator_alias,
    o.rank AS operator_rank,
    rt.notes,
    rt.created_at AS transaction_created_at,
    cl.asset_id,
    a.name AS asset_name,
    cl.quantity
FROM recent_transactions rt
JOIN custody_lines cl ON cl.custody_transaction_id = rt.id
JOIN assets a ON a.id = cl.asset_id
JOIN operators o ON o.id = rt.operator_id
ORDER BY rt.created_at DESC, rt.id DESC, cl.id ASC;

-- name: ListCurrentCustodyByAsset :many
SELECT
    cb.asset_id,
    cb.personnel_id,
    p.full_name AS personnel_full_name,
    cb.quantity,
    cb.updated_at
FROM custody_balances cb
JOIN personnel p ON p.id = cb.personnel_id
WHERE cb.asset_id = @asset_id
  AND cb.quantity > 0
ORDER BY p.full_name ASC, cb.personnel_id ASC;
