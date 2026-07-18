-- name: CreatePersonnel :exec
INSERT INTO personnel (
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active
) VALUES (
    @id,
    @full_name,
    @alias,
    @rank,
    @registration_id,
    @section,
    @organization_unit,
    @active
);

-- name: GetPersonnel :one
SELECT
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active,
    created_at,
    updated_at
FROM personnel
WHERE id = @id;

-- name: ListPersonnel :many
SELECT
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active,
    created_at,
    updated_at
FROM personnel
WHERE
    (
        @status_filter::text = 'all'
        OR (@status_filter::text = 'active' AND active = true)
        OR (@status_filter::text = 'inactive' AND active = false)
    )
ORDER BY created_at DESC, id DESC
LIMIT @limit_count;

-- name: SearchPersonnel :many
WITH search_terms AS (
    SELECT
        text_terms.search_pattern,
        digit_terms.registration_pattern
    FROM unnest(sqlc.arg(search_patterns)::text[]) WITH ORDINALITY AS text_terms(search_pattern, term_position)
    JOIN unnest(sqlc.arg(registration_patterns)::text[]) WITH ORDINALITY AS digit_terms(registration_pattern, term_position)
        USING (term_position)
)
SELECT
    id,
    full_name,
    alias,
    rank,
    registration_id,
    section,
    organization_unit,
    active,
    created_at,
    updated_at
FROM personnel
WHERE NOT EXISTS (
    SELECT 1
    FROM search_terms
    WHERE NOT (
        full_name ILIKE search_terms.search_pattern ESCAPE '\'
        OR alias ILIKE search_terms.search_pattern ESCAPE '\'
        OR registration_id ILIKE search_terms.search_pattern ESCAPE '\'
        OR registration_id ILIKE search_terms.registration_pattern ESCAPE '\'
        OR rank ILIKE search_terms.search_pattern ESCAPE '\'
        OR section ILIKE search_terms.search_pattern ESCAPE '\'
        OR organization_unit ILIKE search_terms.search_pattern ESCAPE '\'
    )
)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count);

-- name: DeactivatePersonnel :one
WITH updated AS (
    UPDATE personnel
    SET
        active = false,
        updated_at = now()
    WHERE id = @id
      AND active = true
    RETURNING id
)
SELECT count(*)::int
FROM updated;

-- name: ReactivatePersonnel :one
WITH updated AS (
    UPDATE personnel
    SET
        active = true,
        updated_at = now()
    WHERE id = @id
      AND active = false
    RETURNING id
)
SELECT count(*)::int
FROM updated;
