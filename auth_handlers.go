package main

import (
	"encoding/json"
	"fmt"
	"muhammaddev/internal/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (config *Config) signupHandler(w http.ResponseWriter, r *http.Request) {
	body := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}
	_, err = config.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hashedPassword),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user. err: %v", err))
		return
	}
	respondWithJson(w, 200, "user created successfully")
}

func (config *Config) loginHandler(w http.ResponseWriter, r *http.Request) {

	body := User{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}

	user, err := config.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Wrong password. err: %v", err))
		return
	}
	respondWithJson(w, 200, "Login successfully")
}
