-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, password, is_chirpy_red FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT id, email
FROM users u
JOIN refresh_tokens rt ON rt.user_id = u.id
WHERE rt.token = $1 
AND rt.expires_at > NOW() 
AND rt.revoked_at IS NULL;

-- name: SetIsChirpyRed :exec
UPDATE users
SET is_chirpy_red = $1,
    updated_at = NOW()
WHERE id = $2;

-- name: UpdateUserCredentials :one
UPDATE users
SET email = $2,
    password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, password, is_chirpy_red;