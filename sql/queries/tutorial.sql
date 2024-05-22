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



-- name: GetTutorialWithId :one
SELECT * FROM tutorials WHERE id = $1;

-- name: DeleteTutorial :exec
DELETE  FROM tutorials
 WHERE id = $1;

-- name: UpdateTutorial :one
UPDATE tutorials
SET title = $1, tutorial_url=$2, description=$3, youtube_link=$4,  thumbnail = $5,  playlist_id=$6
WHERE id = $7
RETURNING *;

