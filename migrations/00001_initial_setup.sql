-- +goose Up
-- +goose StatementBegin
CREATE TABLE snippets(
    id INTEGER PRIMARY KEY NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL
);

CREATE TABLE tags(
    id INTEGER PRIMARY KEY NOT NULL,
    tag TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_tags_tag ON tags(tag);

CREATE TABLE snippets_tags(
    id INTEGER PRIMARY KEY NOT NULL,
    snippet_id INTEGER NOT NULL REFERENCES snippets(id),
    tag_id INTEGER NOT NULL REFERENCES tags(id)
);

CREATE TABLE urls(
    id INTEGER PRIMARY KEY NOT NULL,
    url TEXT NOT NULL,
    snippet_id INTEGER NOT NULL REFERENCES snippets(id)
);
CREATE UNIQUE INDEX idx_urls_url ON urls(url);

-- +goose StatementEnd
