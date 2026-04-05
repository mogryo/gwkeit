-- +goose Up
-- +goose StatementBegin
ALTER TABLE snippets ADD COLUMN language TEXT;

CREATE INDEX IF NOT EXISTS idx_snippets_language ON snippets(language);
-- +goose StatementEnd
