// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: image.sql

package database

import (
	"context"
)

const getImages = `-- name: GetImages :many
SELECT id, image_url FROM images
`

func (q *Queries) GetImages(ctx context.Context) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, getImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(&i.ID, &i.ImageUrl); err != nil {
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

const postImage = `-- name: PostImage :one
INSERT INTO images (

image_url
)
VALUES ($1 )
RETURNING id, image_url
`

func (q *Queries) PostImage(ctx context.Context, imageUrl string) (Image, error) {
	row := q.db.QueryRowContext(ctx, postImage, imageUrl)
	var i Image
	err := row.Scan(&i.ID, &i.ImageUrl)
	return i, err
}
