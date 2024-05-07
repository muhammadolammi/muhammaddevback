package main

import (
	"database/sql"

	"github.com/google/uuid"
)

type Image struct {
	ID       uuid.UUID `json:"id"`
	ImageUrl string    `json:"image_url"`
}

type Playlist struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type Post struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	PostUrl     string         `json:"post_url"`
	Content     string         `json:"content"`
	YoutubeLink sql.NullString `json:"youtube_link"`
}

type Tutorial struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	TutorialUrl string    `json:"tutorial_url"`
	Description string    `json:"description"`
	YoutubeLink string    `json:"youtube_link"`
	PlaylistID  uuid.UUID `json:"playlist_id"`
}
