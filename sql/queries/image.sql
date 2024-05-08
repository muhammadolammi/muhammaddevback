
-- name: PostImage :one
INSERT INTO images (

image_url
)
VALUES ($1 )
RETURNING *;


-- name: GetImages :many
SELECT * FROM images;