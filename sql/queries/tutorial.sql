-- name: GetTutorials :many
SELECT * FROM tutorials;

-- name: GetTutorial :one
SELECT * FROM tutorials
WHERE $1=id;



-- name: GetPlaylistTutorials :many
SELECT * FROM tutorials
WHERE $1=playlist_id;


-- name: PostTutorial :one
INSERT INTO tutorials (

title,
tutorial_url,
description ,
youtube_link,
playlist_id
)
VALUES ($1, $2, $3 , $4, $5)
RETURNING *;
