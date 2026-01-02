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
    snippet_id INTEGER NOT NULL REFERENCES snippets(id),
    tag_id INTEGER NOT NULL REFERENCES tags(id)
);

CREATE TABLE urls(
    id INTEGER PRIMARY KEY NOT NULL,
    url TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_urls_url ON urls(url);

CREATE TABLE snippets_urls(
    snippet_id INTEGER NOT NULL REFERENCES snippets(id),
    url_id INTEGER NOT NULL REFERENCES urls(id)
);
-- +goose StatementEnd
