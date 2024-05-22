-- name: GetPosts :many
SELECT * FROM posts;


-- name: GetPostWithId :one
SELECT * FROM posts WHERE id = $1;


-- name: PostPost :one
INSERT INTO posts (
title,
post_url, content , thumbnail)
VALUES ( $1, $2, $3, $4)
RETURNING *;

-- name: DeletePost :exec
DELETE  FROM posts WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts
SET title = $1, content = $2, thumbnail = $3, post_url = $4
WHERE id = $5
RETURNING *;
