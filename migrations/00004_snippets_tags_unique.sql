-- +goose Up
-- +goose StatementBegin

CREATE UNIQUE INDEX idx_snippet_tag ON snippets_tags(snippet_id, tag_id);

-- +goose StatementEnd