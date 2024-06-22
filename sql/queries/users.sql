
-- name: GetUsers :many
SELECT * FROM users;


-- name: CreateUser :one
INSERT INTO users (
first_name, last_name,
email, password, access_token  )
VALUES ( $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id=$1;

-- name: GetUserWithEmail :one
SELECT * FROM users
WHERE email=$1;


-- name: UpdatePassword :exec
UPDATE users
SET password = $1
WHERE email = $2
RETURNING *;

-- name: UpdateAccessToken :exec
UPDATE users
SET access_token = $1
WHERE email = $2
RETURNING *;




-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
);


-- name: AccessTokenExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE access_token = $1
);
