-- name: GetSnippets :many
SELECT * FROM snippets;

-- name: GetSnippetCount :one
SELECT COUNT(*) FROM snippets;

-- name: InsertSnippet :one
INSERT INTO snippets (title, body, description, url) VALUES (?, ?, ?, ?) RETURNING id;

-- name: InsertTag :one
INSERT INTO tags (tag) VALUES (?) RETURNING id;

-- name: InsertUrl :one
INSERT INTO urls (url, snippet_id) VALUES (?, ?) RETURNING id;

-- name: InsertSnippetTag :exec
INSERT INTO snippets_tags (snippet_id, tag_id) VALUES (?, ?);

-- name: UpdateSnippet :exec
UPDATE snippets SET title = ?, body = ?, description = ?, url = ?, updated_at = current_timestamp WHERE id = ?;

-- name: FindSnippetsByTags :many
SELECT s.*
FROM snippets s
JOIN snippets_tags st on s.id = st.snippet_id
JOIN tags t on st.tag_id = t.id
WHERE t.tag IN (sqlc.slice('tags'));

-- name: FindSnippetsByLikeTags :many
SELECT s.*
FROM snippets s
JOIN snippets_tags st on s.id = st.snippet_id
JOIN tags t on st.tag_id = t.id
WHERE t.tag LIKE (sqlc.arg('tag'));

-- name: FindUrlsBySnippetId :many
SELECT u.*
FROM urls u
WHERE u.snippet_id = ?;

-- name: FindTagsBySnippetId :many
SELECT t.*
FROM tags t
JOIN snippets_tags st on t.id = st.tag_id
WHERE st.snippet_id = ?;

-- name: FindTagsByTag :many
SELECT * FROM tags WHERE tag IN (sqlc.slice('tags'));

-- name: FindSnippetDataById :one
SELECT s.*, group_concat(COALESCE(t.tag, ''), ' ') as tag_list, group_concat(COALESCE(u.url, ''), ' ') as url_list
FROM snippets s
LEFT JOIN snippets_tags st ON s.id = st.snippet_id
LEFT JOIN tags t ON st.tag_id = t.id
LEFT JOIN urls u on s.id = u.snippet_id
WHERE s.id = ?
GROUP BY s.id;

-- name: DeleteSnippetTags :exec
DELETE FROM snippets_tags
WHERE id IN (
    SELECT st.id
    FROM snippets_tags st
    LEFT JOIN tags ON snippets_tags.tag_id = tags.id
    WHERE snippets_tags.tag_id IN (sqlc.slice('tagIds')) AND snippets_tags.snippet_id = ?
);

-- name: DeleteSnippetUrls :exec
DELETE FROM urls WHERE id IN (sqlc.slice('ids'));

-- name: DeleteTagByTag :exec
DELETE FROM tags WHERE tag = ?;

-- name: SnippetTagExists :one
SELECT EXISTS (SELECT 1 FROM snippets_tags JOIN tags ON snippets_tags.tag_id = tags.id WHERE tags.tag = ?);