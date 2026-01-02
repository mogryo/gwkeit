-- name: GetSnippets :many
SELECT * FROM snippets;

-- name: GetSnippetCount :one
SELECT COUNT(*) FROM snippets;

-- name: InsertSnippet :one
INSERT INTO snippets (title, body) VALUES (?, ?) RETURNING id;

-- name: InsertTag :one
INSERT INTO tags (tag) VALUES (?) RETURNING id;

-- name: InsertUrl :one
INSERT INTO urls (url) VALUES (?) RETURNING id;

-- name: InsertSnippetTag :exec
INSERT INTO snippets_tags (snippet_id, tag_id) VALUES (?, ?);

-- name: InsertUrlSnippet :exec
INSERT INTO snippets_urls (url_id, snippet_id) VALUES (?, ?);

-- name: FindSnippetsByTags :many
SELECT s.*
FROM snippets s
JOIN snippets_tags st on s.id = st.snippet_id
JOIN tags t on st.tag_id = t.id
WHERE t.tag IN (sqlc.slice('tags'));

-- name: FindUrlsBySnippetId :many
SELECT u.*
FROM urls u
JOIN snippets_urls su on u.id = su.url_id
WHERE su.snippet_id = ?;

-- name: FindTagsBySnippetId :many
SELECT t.*
FROM tags t
JOIN snippets_tags st on t.id = st.tag_id
WHERE st.snippet_id = ?;

-- name: FindTagsByTag :many
SELECT * FROM tags WHERE tag IN (sqlc.slice('tags'));