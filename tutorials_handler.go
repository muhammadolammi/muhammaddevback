package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"muhammaddev/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiConfig *Config) getPlaylistTutorialsHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "playlistID")
	log.Println(id)
	playlistId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error parsing playlist id. err: %v", err))
		return
	}
	dbPlaylistTutorials, err := apiConfig.DB.GetPlaylistTutorials(r.Context(), playlistId)
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

func (apiConfig *Config) getTutorialsHandler(w http.ResponseWriter, r *http.Request) {

	dbTutorials, err := apiConfig.DB.GetTutorials(r.Context())
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

func (apiConfig *Config) postTutorialHandler(w http.ResponseWriter, r *http.Request) {
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

	dbTutorial, err := apiConfig.DB.PostTutorial(r.Context(), database.PostTutorialParams{

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

func (apiConfig *Config) getTutorialWithIdHandler(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "tutorialID")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error  parsing id to uuid. err :%v", err))
		return
	}
	dbTutorial, err := apiConfig.DB.GetTutorialWithId(r.Context(), id)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting tutorial. err :%v", err))
		return
	}
	tutorial := dbTutorialToTutorial(dbTutorial)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: tutorial}
	respondWithJson(w, 200, resp)
}

func (apiConfig *Config) getTutorialWithTitleHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "tutorialTitle")

	// rawQuery := r.URL.RawQuery
	// title, err := url.QueryUnescape(encodedTitle)
	// if err != nil {
	//     respondWithError(w, 400, fmt.Sprintf("error decoding title: %v", err))
	//     return
	// }
	log.Println(title)
	// log.Println(rawQuery)

	dbTutorial, err := apiConfig.DB.GetTutorialWithTitle(r.Context(), title)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting tutorial. err :%v", err))
		return
	}
	tutorial := dbTutorialToTutorial(dbTutorial)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: tutorial}
	respondWithJson(w, 200, resp)
}

func (apiConfig *Config) updateTutorialHandler(w http.ResponseWriter, r *http.Request) {
	body := Tutorial{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	id := chi.URLParam(r, "tutorialID")
	tutorialID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error parsing tutorial id to uuid. err: %v", err))
		return
	}
	var playlistId uuid.UUID
	if body.PlaylistID != "" {
		playlistId, err = uuid.Parse(body.PlaylistID)
		if err != nil {
			respondWithError(w, 501, fmt.Sprintf("error parsing playlist id to uuid. err: %v", err))
			return
		}
	}

	dbTutorial, err := apiConfig.DB.UpdateTutorial(r.Context(), database.UpdateTutorialParams{

		Title:       body.Title,
		TutorialUrl: body.TutorialUrl,
		Description: body.Description,
		Thumbnail:   sql.NullString{Valid: true, String: body.Thumbnail},
		YoutubeLink: body.YoutubeLink,
		PlaylistID:  playlistId,
		ID:          tutorialID,
	})
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting Tutorial to db. err: %v", err))
		return
	}
	tutorial := dbTutorialToTutorial(dbTutorial)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: tutorial}
	respondWithJson(w, 200, resp)
}

func (apiConfig *Config) deleteTutorialHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "tutorialID")
	tutorialID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error parsing tutorial id to uuid. err: %v", err))
		return
	}
	err = apiConfig.DB.DeleteTutorial(context.Background(), tutorialID)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error deleting tutorial. err: %v", err))
		return
	}

	respondWithJson(w, 200, "tutorial deleted successfully")
}
