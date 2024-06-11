// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const deletePost = `-- name: DeletePost :exec
DELETE  FROM posts WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostWithId = `-- name: GetPostWithId :one
SELECT id, title, post_url, content, thumbnail FROM posts WHERE id = $1
`

func (q *Queries) GetPostWithId(ctx context.Context, id uuid.UUID) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostWithId, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PostUrl,
		&i.Content,
		&i.Thumbnail,
	)
	return i, err
}

const getPostWithTitle = `-- name: GetPostWithTitle :one
SELECT id, title, post_url, content, thumbnail FROM posts WHERE title = $1
`

func (q *Queries) GetPostWithTitle(ctx context.Context, title string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostWithTitle, title)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PostUrl,
		&i.Content,
		&i.Thumbnail,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
SELECT id, title, post_url, content, thumbnail FROM posts
`

func (q *Queries) GetPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.PostUrl,
			&i.Content,
			&i.Thumbnail,
		); err != nil {
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

const postPost = `-- name: PostPost :one
INSERT INTO posts (
title,
post_url, content , thumbnail)
VALUES ( $1, $2, $3, $4)
RETURNING id, title, post_url, content, thumbnail
`

type PostPostParams struct {
	Title     string
	PostUrl   string
	Content   string
	Thumbnail sql.NullString
}

func (q *Queries) PostPost(ctx context.Context, arg PostPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, postPost,
		arg.Title,
		arg.PostUrl,
		arg.Content,
		arg.Thumbnail,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PostUrl,
		&i.Content,
		&i.Thumbnail,
	)
	return i, err
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET title = $1, content = $2, thumbnail = $3, post_url = $4
WHERE id = $5
RETURNING id, title, post_url, content, thumbnail
`

type UpdatePostParams struct {
	Title     string
	Content   string
	Thumbnail sql.NullString
	PostUrl   string
	ID        uuid.UUID
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost,
		arg.Title,
		arg.Content,
		arg.Thumbnail,
		arg.PostUrl,
		arg.ID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PostUrl,
		&i.Content,
		&i.Thumbnail,
	)
	return i, err
}
