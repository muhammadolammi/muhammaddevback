package main

import "net/http"

func (config *Config) getPlaylistsHandler(w http.ResponseWriter, r *http.Request) {

	respondWithJson(w, 200, "hello")
}
