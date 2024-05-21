package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"muhammaddev/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (config *Config) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	dbPosts, err := config.DB.GetPosts(r.Context())
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting posts err :%v", err))
		return
	}
	posts := dbPostsToPosts(dbPosts)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: posts}
	respondWithJson(w, 200, resp)
}

func (config *Config) postPosttHandler(w http.ResponseWriter, r *http.Request) {
	body := Post{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}


	dbPost, err := config.DB.PostPost(r.Context(), database.PostPostParams{

		Title:   body.Title,
		PostUrl: body.PostUrl,
		Content: body.Content,
		Thumbnail: sql.NullString{Valid: true, String: body.Thumbnail},
	})
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting playlist to db. err: %v", err))
		return
	}
	post := dbPostToPost(dbPost)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: post}
	respondWithJson(w, 200, resp)
}


func (config *Config) updatePosttHandler(w http.ResponseWriter, r *http.Request) {
	body := Post{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	id := chi.URLParam(r, "postID")
	postId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error parsing post id to uuid. err: %v", err))
		return
	}

	dbPost, err := config.DB.UpdatePost(r.Context(), database.UpdatePostParams{

		Title:   body.Title,
		PostUrl: body.PostUrl,
		Content: body.Content,
		Thumbnail: sql.NullString{Valid: true, String: body.Thumbnail},
		ID: postId,
	})
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting playlist to db. err: %v", err))
		return
	}
	post := dbPostToPost(dbPost)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: post}
	respondWithJson(w, 200, resp)
}


func (config *Config) deletePosttHandler(w http.ResponseWriter, r *http.Request) {
	
	id := chi.URLParam(r, "postID")
	postId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error parsing post id to uuid. err: %v", err))
		return
	}
	err = config.DB.DeletePost(context.Background(), postId)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error deleting post. err: %v", err))
		return
	}

	respondWithJson(w, 200, "post deleted successfully")
}