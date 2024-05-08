package main

import (
	"encoding/json"
	"fmt"
	"log"
	"muhammaddev/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (config *Config) getPlaylistTutorials(w http.ResponseWriter, r *http.Request) {

	// params := struct {
	// 	PlaylistId string `json:"playlist_id"`
	// }{}
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&params)
	// if err != nil {
	// 	respondWithError(w, 500, fmt.Sprintf("error decoding params. err: %v", err))
	// 	return
	// }

	id := chi.URLParam(r, "playlistID")
	log.Println(id)
	playlistId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error parsing playlist id. err: %v", err))
		return
	}
	dbPlaylistTutorials, err := config.DB.GetPlaylistTutorials(r.Context(), playlistId)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting playlist tutorials err :%v", err))
		return
	}
	playlistTutorials := dbTutorialsToTutorials(dbPlaylistTutorials)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: playlistTutorials}
	respondWithJson(w, 200, resp)
}

func (config *Config) getTutorials(w http.ResponseWriter, r *http.Request) {

	dbTutorials, err := config.DB.GetTutorials(r.Context())
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting playlist tutorials err :%v", err))
		return
	}
	tutorials := dbTutorialsToTutorials(dbTutorials)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: tutorials}
	respondWithJson(w, 200, resp)
}

func (config *Config) postTutorialHandler(w http.ResponseWriter, r *http.Request) {
	body := Tutorial{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}

	playlistId, err := uuid.Parse(body.PlaylistID)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error parsing playlist id. err: %v", err))
		return
	}

	dbTutorial, err := config.DB.PostTutorial(r.Context(), database.PostTutorialParams{

		Title:       body.Title,
		Description: body.Description,
		TutorialUrl: body.TutorialUrl,
		YoutubeLink: body.YoutubeLink,
		PlaylistID:  playlistId,
	})
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting tutorial to db. err: %v", err))
		return
	}

	tutorial := dbTutorialToTutorial(dbTutorial)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: tutorial}
	respondWithJson(w, 200, resp)
}
