-- name: ListCurrentCustodySummaryByPersonnel :many
WITH latest_corrections AS (
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
        ct.transaction_type,
        COALESCE(lc.corrected_personnel_id, ct.personnel_id) AS effective_personnel_id,
        lc.id AS latest_correction_id
    FROM custody_transactions ct
    LEFT JOIN latest_corrections lc
        ON lc.corrected_transaction_id = ct.id
),
effective_lines AS (
    SELECT
        et.id AS custody_transaction_id,
        et.transaction_type,
        et.effective_personnel_id AS personnel_id,
        cl.asset_id,
        cl.quantity
    FROM effective_transactions et
    JOIN custody_lines cl
        ON cl.custody_transaction_id = et.id
    WHERE et.latest_correction_id IS NULL

    UNION ALL

    SELECT
        et.id AS custody_transaction_id,
        et.transaction_type,
        et.effective_personnel_id AS personnel_id,
        ccl.asset_id,
        ccl.quantity
    FROM effective_transactions et
    JOIN custody_correction_lines ccl
        ON ccl.custody_correction_id = et.latest_correction_id
    WHERE et.latest_correction_id IS NOT NULL
),
current_balances AS (
    SELECT
        effective_lines.personnel_id,
        effective_lines.asset_id,
        SUM(
            CASE effective_lines.transaction_type
                WHEN 'checkout' THEN effective_lines.quantity
                WHEN 'return' THEN -effective_lines.quantity
                ELSE 0
            END
        ) AS current_quantity
    FROM effective_lines
    GROUP BY
        effective_lines.personnel_id,
        effective_lines.asset_id
)
SELECT
    current_balances.personnel_id,
    COALESCE(SUM(current_balances.current_quantity), 0)::bigint AS current_custody_quantity
FROM current_balances
WHERE current_balances.current_quantity > 0
GROUP BY current_balances.personnel_id;

-- name: ListCurrentCustodySummaryByAsset :many
WITH latest_corrections AS (
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
        ct.transaction_type,
        COALESCE(lc.corrected_personnel_id, ct.personnel_id) AS effective_personnel_id,
        lc.id AS latest_correction_id
    FROM custody_transactions ct
    LEFT JOIN latest_corrections lc
        ON lc.corrected_transaction_id = ct.id
),
effective_lines AS (
    SELECT
        et.id AS custody_transaction_id,
        et.transaction_type,
        et.effective_personnel_id AS personnel_id,
        cl.asset_id,
        cl.quantity
    FROM effective_transactions et
    JOIN custody_lines cl
        ON cl.custody_transaction_id = et.id
    WHERE et.latest_correction_id IS NULL

    UNION ALL

    SELECT
        et.id AS custody_transaction_id,
        et.transaction_type,
        et.effective_personnel_id AS personnel_id,
        ccl.asset_id,
        ccl.quantity
    FROM effective_transactions et
    JOIN custody_correction_lines ccl
        ON ccl.custody_correction_id = et.latest_correction_id
    WHERE et.latest_correction_id IS NOT NULL
),
current_balances AS (
    SELECT
        effective_lines.personnel_id,
        effective_lines.asset_id,
        SUM(
            CASE effective_lines.transaction_type
                WHEN 'checkout' THEN effective_lines.quantity
                WHEN 'return' THEN -effective_lines.quantity
                ELSE 0
            END
        ) AS current_quantity
    FROM effective_lines
    GROUP BY
        effective_lines.personnel_id,
        effective_lines.asset_id
)
SELECT
    current_balances.asset_id,
    COUNT(DISTINCT current_balances.personnel_id)::bigint AS current_custodian_count
FROM current_balances
WHERE current_balances.current_quantity > 0
GROUP BY current_balances.asset_id;
