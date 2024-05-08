package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"muhammaddev/internal/database"
	"net/http"
)

func (config *Config) getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	dbPlaylists, err := config.DB.GetPlaylists(r.Context())
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting playlists err :%v", err))
		return
	}
	playlists := dbPlaylistsToPlaylists(dbPlaylists)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: playlists}
	respondWithJson(w, 200, resp)
}

func (config *Config) postPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	body := Playlist{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}

	dbplaylist, err := config.DB.PostPlaylist(r.Context(), database.PostPlaylistParams{

		Description: sql.NullString{String: body.Description, Valid: true},
		Name:        body.Name,
	})
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting playlist to db. err: %v", err))
		return
	}
	playlist := dbPlaylistToPlaylist(dbplaylist)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: playlist}
	respondWithJson(w, 200, resp)
}
