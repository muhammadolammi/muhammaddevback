-- name: GetPlaylist :many
SELECT * FROM playlists;


-- name: GetPlaylistIdByName :one
SELECT name FROM playlists
WHERE $1=id;