package main

import "muhammaddev/internal/database"

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}
type Config struct {
	PORT                          string
	DB                            *database.Queries
	API_KEY                       string
	AccessTokenExpirationMinutes  int
	RefreshTokenExpirationMinutes int
}

type Playlist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Post struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	PostUrl   string `json:"post_url"`
	Content   string `json:"content"`
	Thumbnail string `json:"thumbnail"`
}

type Tutorial struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	TutorialUrl string `json:"tutorial_url"`
	Description string `json:"description"`
	YoutubeLink string `json:"youtube_link"`
	PlaylistID  string `json:"playlist_id"`
	Thumbnail   string `json:"thumbnail"`
}
