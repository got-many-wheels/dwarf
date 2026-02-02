-- name: GetURLByCode :one
SELECT * FROM urls WHERE code = $1;

-- name: CreateURL :one
INSERT INTO urls (code, long) VALUES ($1, $2) RETURNING *;

-- name: DeleteURLByCode :exec
DELETE FROM urls WHERE code = $1;

-- name: UpdateURL :exec
UPDATE urls set code = $2, long = $3 WHERE id = $1;
