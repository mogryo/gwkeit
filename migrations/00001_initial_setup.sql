-- +goose Up
-- +goose StatementBegin
CREATE TABLE snippets(
    id INTEGER PRIMARY KEY NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- And an external content fts5 table to index snippets.
CREATE VIRTUAL TABLE snippets_fts USING fts5(title, body, description, url, content='snippets', content_rowid='id');

-- Triggers to keep the FTS index up to date.
CREATE TRIGGER snippets_ai AFTER INSERT ON snippets BEGIN
    INSERT INTO snippets_fts(rowid, title, body, description, url)
    VALUES (new.id, new.title, new.body, new.description, new.url);
END;
CREATE TRIGGER snippets_ad AFTER DELETE ON snippets BEGIN
    INSERT INTO snippets_fts(snippets_fts, rowid, title, body, description, url)
    VALUES('delete', old.id, old.title, old.body, old.description, old.url);
END;
CREATE TRIGGER snippets_au AFTER UPDATE ON snippets BEGIN
    INSERT INTO snippets_fts(snippets_fts, rowid, title, body, description, url)
    VALUES('delete', old.id, old.title, old.body, old.description, old.url);
    INSERT INTO snippets_fts(rowid, title, body, description, url)
    VALUES (old.id, old.title, old.body, old.description, old.url);
    UPDATE snippets SET updated_at=current_timestamp WHERE rowid = new.rowid;
END;

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
