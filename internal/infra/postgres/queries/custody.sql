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

-- name: IncreaseCustodyBalance :exec
INSERT INTO custody_balances (
    personnel_id,
    asset_id,
    quantity,
    updated_at
) VALUES (
    @personnel_id,
    @asset_id,
    @quantity,
    now()
)
ON CONFLICT (personnel_id, asset_id)
DO UPDATE SET
    quantity = custody_balances.quantity + EXCLUDED.quantity,
    updated_at = now();

-- name: DecreaseCustodyBalanceIfAvailable :one
UPDATE custody_balances
SET
    quantity = quantity - @quantity,
    updated_at = now()
WHERE personnel_id = @personnel_id
  AND asset_id = @asset_id
  AND quantity >= @quantity
RETURNING quantity;

-- name: ListCurrentCustodyByPersonnel :many
SELECT
    cb.personnel_id,
    cb.asset_id,
    a.name AS asset_name,
    a.active AS asset_active,
    cb.quantity,
    cb.updated_at
FROM custody_balances cb
JOIN assets a ON a.id = cb.asset_id
WHERE cb.personnel_id = @personnel_id
  AND cb.quantity > 0
ORDER BY a.name ASC, cb.asset_id ASC;

-- name: ListPersonnelWithCurrentCustody :many
SELECT
    p.id,
    p.full_name,
    p.alias,
    p.rank,
    p.registration_id,
    p.active,
    sum(cb.quantity)::int AS total_quantity
FROM custody_balances cb
JOIN personnel p ON p.id = cb.personnel_id
WHERE cb.quantity > 0
GROUP BY
    p.id,
    p.full_name,
    p.alias,
    p.rank,
    p.registration_id,
    p.active
ORDER BY p.active DESC, p.rank ASC, p.alias ASC, p.full_name ASC;

-- name: ListCustodyHistoryByPersonnel :many
WITH correction_counts AS (
    SELECT
        corrected_transaction_id,
        count(*)::int AS edit_count
    FROM custody_corrections
    GROUP BY corrected_transaction_id
),
latest_corrections AS (
    SELECT DISTINCT ON (corrected_transaction_id)
        id,
        corrected_transaction_id,
        corrected_personnel_id,
        corrected_notes
    FROM custody_corrections
    ORDER BY corrected_transaction_id, created_at DESC, id DESC
),
effective_transactions AS (
    SELECT
        ct.id,
        ct.transaction_type,
        COALESCE(lc.corrected_personnel_id, ct.personnel_id) AS effective_personnel_id,
        ct.personnel_id AS original_personnel_id,
        ct.operator_id,
        COALESCE(lc.corrected_notes, ct.notes) AS effective_notes,
        ct.created_at,
        lc.id AS latest_correction_id,
        COALESCE(cc.edit_count, 0)::int AS edit_count
    FROM custody_transactions ct
    LEFT JOIN latest_corrections lc
        ON lc.corrected_transaction_id = ct.id
    LEFT JOIN correction_counts cc
        ON cc.corrected_transaction_id = ct.id
)
SELECT
    history.transaction_id,
    history.transaction_type,
    history.personnel_id,
    history.operator_id,
    history.operator_alias,
    history.operator_rank,
    history.notes,
    history.transaction_created_at,
    history.asset_id,
    history.asset_name,
    history.asset_active,
    history.quantity,
    history.has_correction,
    history.edit_count
FROM (
    SELECT
        et.id AS transaction_id,
        et.transaction_type,
        et.effective_personnel_id AS personnel_id,
        et.operator_id,
        o.alias AS operator_alias,
        o.rank AS operator_rank,
        et.effective_notes AS notes,
        et.created_at AS transaction_created_at,
        COALESCE(ccl.asset_id, cl.asset_id) AS asset_id,
        a.name AS asset_name,
        a.active AS asset_active,
        COALESCE(ccl.quantity, cl.quantity) AS quantity,
        (et.edit_count > 0) AS has_correction,
        et.edit_count
    FROM effective_transactions et
    LEFT JOIN custody_lines cl
        ON cl.custody_transaction_id = et.id
        AND et.latest_correction_id IS NULL
    LEFT JOIN custody_correction_lines ccl
        ON ccl.custody_correction_id = et.latest_correction_id
        AND et.latest_correction_id IS NOT NULL
    JOIN assets a ON a.id = COALESCE(ccl.asset_id, cl.asset_id)
    JOIN operators o ON o.id = et.operator_id
) history
WHERE personnel_id = sqlc.arg(personnel_id)
ORDER BY transaction_created_at DESC, transaction_id DESC, asset_name ASC
LIMIT sqlc.arg(limit_count);

-- name: ListCustodyTransactionSummaries :many
WITH correction_counts AS (
    SELECT
        corrected_transaction_id,
        count(*)::int AS edit_count
    FROM custody_corrections
    GROUP BY corrected_transaction_id
),
latest_corrections AS (
    SELECT DISTINCT ON (corrected_transaction_id)
        id,
        corrected_transaction_id,
        corrected_personnel_id
    FROM custody_corrections
    ORDER BY corrected_transaction_id, created_at DESC, id DESC
),
effective_transactions AS (
    SELECT
        ct.id,
        row_number() OVER (ORDER BY ct.created_at ASC, ct.id ASC)::int AS sequence_number,
        ct.transaction_type,
        ct.personnel_id AS original_personnel_id,
        COALESCE(lc.corrected_personnel_id, ct.personnel_id) AS effective_personnel_id,
        ct.operator_id,
        ct.created_at,
        lc.id AS latest_correction_id,
        COALESCE(cc.edit_count, 0)::int AS edit_count
    FROM custody_transactions ct
    LEFT JOIN latest_corrections lc
        ON lc.corrected_transaction_id = ct.id
    LEFT JOIN correction_counts cc
        ON cc.corrected_transaction_id = ct.id
),
effective_lines AS (
    SELECT
        et.id AS transaction_id,
        cl.id AS line_order,
        cl.asset_id,
        a.name AS asset_name,
        a.active AS asset_active,
        cl.quantity
    FROM effective_transactions et
    JOIN custody_lines cl
        ON cl.custody_transaction_id = et.id
    JOIN assets a
        ON a.id = cl.asset_id
    WHERE et.latest_correction_id IS NULL

    UNION ALL

    SELECT
        et.id AS transaction_id,
        ccl.id AS line_order,
        ccl.asset_id,
        a.name AS asset_name,
        a.active AS asset_active,
        ccl.quantity
    FROM effective_transactions et
    JOIN custody_correction_lines ccl
        ON ccl.custody_correction_id = et.latest_correction_id
    JOIN assets a
        ON a.id = ccl.asset_id
    WHERE et.latest_correction_id IS NOT NULL
),
line_totals AS (
    SELECT
        transaction_id,
        sum(quantity)::int AS total_quantity
    FROM effective_lines
    GROUP BY transaction_id
),
filtered_transactions AS (
    SELECT et.*
    FROM effective_transactions et
    JOIN personnel op
        ON op.id = et.original_personnel_id
    JOIN personnel ep
        ON ep.id = et.effective_personnel_id
    JOIN operators o
        ON o.id = et.operator_id
    WHERE (
        sqlc.arg(transaction_type_filter)::text = 'all'
        OR et.transaction_type = sqlc.arg(transaction_type_filter)::text
    )
    AND (
        sqlc.arg(edit_status_filter)::text = 'all'
        OR (
            sqlc.arg(edit_status_filter)::text = 'edited'
            AND et.edit_count > 0
        )
        OR (
            sqlc.arg(edit_status_filter)::text = 'unedited'
            AND et.edit_count = 0
        )
    )
    AND (
        sqlc.arg(search_pattern)::text = ''
        OR et.id ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR ep.full_name ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR ep.alias ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR ep.registration_id ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR op.full_name ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR op.alias ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR op.registration_id ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR o.alias ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        OR EXISTS (
            SELECT 1
            FROM effective_lines el
            WHERE el.transaction_id = et.id
              AND el.asset_name ILIKE sqlc.arg(search_pattern)::text ESCAPE '\'
        )
    )
),
selected_transactions AS (
    SELECT *
    FROM filtered_transactions
    ORDER BY created_at DESC, id DESC
    LIMIT @limit_count
)
SELECT
    st.id,
    st.sequence_number,
    st.transaction_type,

    op.id AS original_personnel_id,
    op.full_name AS original_personnel_full_name,
    op.alias AS original_personnel_alias,
    op.rank AS original_personnel_rank,
    op.active AS original_personnel_active,

    ep.id AS effective_personnel_id,
    ep.full_name AS effective_personnel_full_name,
    ep.alias AS effective_personnel_alias,
    ep.rank AS effective_personnel_rank,
    ep.active AS effective_personnel_active,

    st.operator_id,
    o.alias AS operator_alias,
    o.rank AS operator_rank,
    o.role AS operator_role,
    o.active AS operator_active,

    COALESCE(lt.total_quantity, 0)::int AS total_quantity,
    st.created_at,
    (st.edit_count > 0) AS has_correction,
    st.edit_count,

    el.asset_id,
    el.asset_name,
    el.asset_active,
    el.quantity
FROM selected_transactions st
JOIN personnel op ON op.id = st.original_personnel_id
JOIN personnel ep ON ep.id = st.effective_personnel_id
JOIN operators o ON o.id = st.operator_id
LEFT JOIN line_totals lt ON lt.transaction_id = st.id
LEFT JOIN effective_lines el ON el.transaction_id = st.id
ORDER BY st.created_at DESC, st.id DESC, el.asset_name ASC, el.line_order ASC;

-- name: ListCurrentCustodyByAsset :many
SELECT
    cb.asset_id,
    cb.personnel_id,
    p.full_name AS personnel_full_name,
    p.alias AS personnel_alias,
    p.rank AS personnel_rank,
    p.active AS personnel_active,
    cb.quantity,
    cb.updated_at
FROM custody_balances cb
JOIN personnel p ON p.id = cb.personnel_id
WHERE cb.asset_id = @asset_id
  AND cb.quantity > 0
ORDER BY p.full_name ASC, cb.personnel_id ASC;

-- name: GetCustodyTransactionReceiptByID :many
SELECT
    ct.id,
    ct.transaction_type,
    ct.personnel_id,
    p.full_name AS personnel_full_name,
    p.alias AS personnel_alias,
    p.rank AS personnel_rank,
    p.registration_id AS personnel_registration_id,
    p.active AS personnel_active,
    ct.operator_id,
    o.alias AS operator_alias,
    o.rank AS operator_rank,
    o.role AS operator_role,
    o.active AS operator_active,
    ct.notes,
    ct.created_at,
    cl.asset_id,
    a.name AS asset_name,
    a.active AS asset_active,
    cl.quantity
FROM custody_transactions ct
JOIN personnel p ON p.id = ct.personnel_id
JOIN operators o ON o.id = ct.operator_id
JOIN custody_lines cl ON cl.custody_transaction_id = ct.id
JOIN assets a ON a.id = cl.asset_id
WHERE ct.id = @id
ORDER BY a.name ASC, cl.id ASC;

-- name: GetCustodyCorrectionContextsByTransactionID :many
SELECT
    cc.id,
    cc.corrected_transaction_id,
    cc.operator_id,
    o.alias AS operator_alias,
    o.rank AS operator_rank,
    o.role AS operator_role,
    o.active AS operator_active,
    cc.corrected_personnel_id,
    p.full_name AS corrected_personnel_full_name,
    p.alias AS corrected_personnel_alias,
    p.rank AS corrected_personnel_rank,
    p.registration_id AS corrected_personnel_registration_id,
    p.active AS corrected_personnel_active,
    cc.corrected_notes,
    cc.created_at,
    ccl.asset_id,
    a.name AS asset_name,
    a.active AS asset_active,
    ccl.quantity
FROM custody_corrections cc
JOIN operators o ON o.id = cc.operator_id
JOIN personnel p ON p.id = cc.corrected_personnel_id
JOIN custody_correction_lines ccl ON ccl.custody_correction_id = cc.id
JOIN assets a ON a.id = ccl.asset_id
WHERE cc.corrected_transaction_id = @corrected_transaction_id
ORDER BY cc.created_at ASC, cc.id ASC, a.name ASC, ccl.id ASC;
