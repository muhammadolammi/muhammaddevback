-- name: GetPlaylists :many
SELECT * FROM playlists;


-- name: PostPlaylist :one
INSERT INTO playlists (
name,
description )
VALUES ( $1, $2)
RETURNING *;




-- -- name: GetPlaylistIdByName :one
-- SELECT name FROM playlists
-- WHERE $1=id;


