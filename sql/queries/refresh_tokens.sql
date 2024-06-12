-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens 
WHERE user_id=$1;


-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
token, created_at, expires_at, user_id  )
VALUES ( $1, $2, $3, $4)
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token = $1
WHERE user_id = $2
RETURNING *;

