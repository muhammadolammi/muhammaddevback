// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: playlist.sql

package database

import (
	"context"
	"database/sql"
)

const getPlaylists = `-- name: GetPlaylists :many
SELECT id, name, description FROM playlists
`

func (q *Queries) GetPlaylists(ctx context.Context) ([]Playlist, error) {
	rows, err := q.db.QueryContext(ctx, getPlaylists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Playlist
	for rows.Next() {
		var i Playlist
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const postPlaylist = `-- name: PostPlaylist :one
INSERT INTO playlists (
name,
description )
VALUES ( $1, $2)
RETURNING id, name, description
`

type PostPlaylistParams struct {
	Name        string
	Description sql.NullString
}

func (q *Queries) PostPlaylist(ctx context.Context, arg PostPlaylistParams) (Playlist, error) {
	row := q.db.QueryRowContext(ctx, postPlaylist, arg.Name, arg.Description)
	var i Playlist
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}
