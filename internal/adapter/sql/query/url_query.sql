-- name: DeleteShortUrlBeforeCreatedAt :many
DELETE FROM url
WHERE url.created_at < $1
  AND hash IN (SELECT hash FROM url LIMIT $2)
    RETURNING hash;

-- name: CreateShortUrl :one
INSERT INTO url (hash, url) VALUES ($1, $2)
    RETURNING *;

-- name: GetUniqueNumbers :many
SELECT nextval('unique_number_seq') FROM generate_series(1, $1);

-- name: SaveAllHashes :copyfrom
INSERT INTO hash (hash) VALUES ($1);

-- name: GetHashBatch :many
DELETE FROM hash
WHERE hash IN (SELECT hash FROM hash LIMIT $1)
    RETURNING hash;

-- name: FindUrlByHash :one
SELECT hash FROM url
WHERE hash = $1 LIMIT 1;