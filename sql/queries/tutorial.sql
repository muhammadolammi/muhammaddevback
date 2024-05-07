-- name: GetTutorials :many
SELECT * FROM tutorials;

-- name: GetTutorial :one
SELECT * FROM tutorials
WHERE $1=title;



-- name: GetPlaylistTutorials :many
SELECT * FROM tutorials
WHERE $1=playlist_id;



