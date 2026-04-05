-- +goose Up
-- +goose StatementBegin

DELETE FROM snippets_tags WHERE id NOT IN (
    SELECT MIN(id)
    FROM snippets_tags
    GROUP BY snippet_id, tag_id
);
CREATE UNIQUE INDEX idx_snippet_tag ON snippets_tags(snippet_id, tag_id);

-- +goose StatementEnd