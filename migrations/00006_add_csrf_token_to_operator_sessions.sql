-- +goose Up
ALTER TABLE operator_sessions
ADD COLUMN csrf_token TEXT NOT NULL DEFAULT '__legacy_invalid_csrf_token__';

ALTER TABLE operator_sessions
ALTER COLUMN csrf_token DROP DEFAULT;

ALTER TABLE operator_sessions
ADD CONSTRAINT operator_sessions_csrf_token_not_empty CHECK (length(trim(csrf_token)) > 0);

-- +goose Down
ALTER TABLE operator_sessions
DROP CONSTRAINT operator_sessions_csrf_token_not_empty;

ALTER TABLE operator_sessions
DROP COLUMN csrf_token;