package main

import "net/http"

func helloReady(w http.ResponseWriter, r *http.Request) { respondWithJson(w, 200, "hello") }
func errorReady(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 200, "this is an error test")
}
