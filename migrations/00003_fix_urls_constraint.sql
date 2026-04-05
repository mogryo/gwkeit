-- +goose Up
-- +goose StatementBegin

DROP INDEX idx_urls_url;
CREATE UNIQUE INDEX idx_urls_url ON urls(url, snippet_id);

-- +goose StatementEnd