// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID
	Name        string
	Description sql.NullString
}

type Post struct {
	ID        uuid.UUID
	Title     string
	PostUrl   string
	Content   string
	Thumbnail sql.NullString
}

type Tutorial struct {
	ID          uuid.UUID
	Title       string
	TutorialUrl string
	Description string
	YoutubeLink string
	PlaylistID  uuid.UUID
	Thumbnail   sql.NullString
}

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessToken sql.NullString
}
