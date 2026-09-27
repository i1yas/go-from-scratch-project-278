-- name: GetLinks :many
SELECT
    id, original_url, shortcode
FROM links;

-- name: GetLinkByID :one
SELECT
    id, original_url, shortcode
FROM links
WHERE id = $1;

-- name: GetLinkByCode :one
SELECT
    id, original_url, shortcode
FROM links
WHERE shortcode = $1;

-- name: CreateLink :one
INSERT INTO links (original_url, shortcode)
VALUES ($1, $2)
RETURNING
    id, original_url, shortcode;

-- name: UpdateLink :one
UPDATE links
SET
    original_url = $2,
    shortcode = $3
WHERE id = $1
RETURNING
    id, original_url, shortcode;

-- name: DeleteLink :execrows
DELETE FROM links
WHERE id = $1;