package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (config *Config) postImageHandler(w http.ResponseWriter, r *http.Request) {
	body := Image{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}

	dbImage, err := config.DB.PostImage(r.Context(), body.ImageUrl)
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error posting image to db. err: %v", err))
		return
	}

	image := dbImageToImage(dbImage)

	resp := struct {
		Data interface{} `json:"data"`
	}{Data: image}
	respondWithJson(w, 200, resp)

}

func (config *Config) getImagesHandler(w http.ResponseWriter, r *http.Request) {
	dbImages, err := config.DB.GetImages(r.Context())
	if err != nil {
		respondWithError(w, 501, fmt.Sprintf("error getting playlists err :%v", err))
		return
	}
	images := dbImagesToImages(dbImages)
	resp := struct {
		Data interface{} `json:"data"`
	}{Data: images}
	respondWithJson(w, 200, resp)
}
