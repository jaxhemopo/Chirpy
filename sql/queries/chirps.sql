-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetChirps :many
SELECT id, created_at, updated_at, user_id, body FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT id, created_at, updated_at, user_id, body FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;