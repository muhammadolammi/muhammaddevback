package main

import (
	"encoding/json"
	"fmt"
	"muhammaddev/internal/auth"
	"muhammaddev/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	if body.Email == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a mail. err: %v", err))
		return
	}
	if body.Password == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a password. err: %v", err))
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}
	user, err := config.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hashedPassword),
	})
	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint \"users_email_key\`) {
			respondWithError(w, http.StatusInternalServerError, "User already created")
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user. err: %v", err))
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating refresh token. err: %v", err))
		return
	}

	config.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		UserID:    user.ID,
	})

	respondWithJson(w, 200, "user created successfully")
}

func (config *Config) loginHandler(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	if body.Email == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter your registered email. err: %v", err))
		return
	}
	if body.Password == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a password. err: %v", err))
		return
	}

	user, err := config.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			respondWithError(w, http.StatusInternalServerError, "Wrong password.")
			return
		}
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}

	jwtTokenString, err := auth.MakeJwtTokenString([]byte(config.API_KEY), user.ID.String())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating jwt string. err: %v", err))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwtTokenString,
		Expires: time.Now().UTC().Add(2 * time.Minute),
		Secure:  false,
	})

	userRefreshToken, err := config.DB.GetRefreshToken(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting refresh token. err: %v", err))
		return
	}
	//   Check if reresh token expired
	if userRefreshToken.ExpiresAt.Before(time.Now()) {
		refreshToken, err := auth.GenerateRefreshToken()
		if err != nil {

			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error generating refresh token. err: %v", err))
			return
		}

		config.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Minute * 5),
			UserID:    user.ID,
		})
	}

	respondWithJson(w, 200, "Login successfully")

}

func (config *Config) passwordChangeHandler(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	if body.Email == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a mail. err: %v", err))
		return
	}
	if body.OldPassword == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a password. err: %v", err))
		return
	}
	if body.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a new password. err: %v", err))
		return
	}

	user, err := config.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	// AUTHENTICATE THE USER
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			respondWithError(w, http.StatusInternalServerError, "Wrong password.")
			return
		}
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}
	// UPDATE THE PASSWORD
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}
	err = config.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Email:    body.Email,
		Password: string(newHashedPassword),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating password. err: %v", err))
		return
	}
	respondWithJson(w, 200, "Password Updated")
}

func (config *Config) validateHandler(w http.ResponseWriter, r *http.Request) {

	respondWithJson(w, 200, "true")

}

func (apiconfig *Config) refresh(w http.ResponseWriter, r *http.Request) {

	params := struct {
		RefreshToken string `json:"refresh_token"`
		UserId       string `json:"user_id"`
	}{}

	// Parse the refresh token from the body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body || or bad request")
		return
	}
	userId, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing user id")
		return
	}
	userRefreshToken, err := apiconfig.DB.GetRefreshToken(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting user user refresh token")
		return
	}
	if userRefreshToken.Token != params.RefreshToken {
		respondWithError(w, http.StatusInternalServerError, "wrong refresh token")
		return
	}

	if userRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusInternalServerError, "Expired refresh token Login again")

		return
	}

	tokenString, err := auth.MakeJwtTokenString([]byte(apiconfig.API_KEY), params.UserId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new jwt  || or bad request")
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().UTC().Add(2 * time.Minute),
		Secure:  false,
	})
}
