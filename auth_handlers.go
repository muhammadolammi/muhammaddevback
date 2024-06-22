package main

import (
	"encoding/json"
	"fmt"
	"log"
	"muhammaddev/internal/auth"
	"muhammaddev/internal/database"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var tokenBlacklist = make(map[string]bool)

func addToTokenBlacklist(token string) {
	tokenBlacklist[token] = true
}

func isInTokenBlacklist(token string) bool {
	return tokenBlacklist[token]
}

func (apiConfig *Config) signupHandler(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
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
	if body.Password == "" {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a password. err: %v", err))
		return
	}

	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if userExist {
		respondWithError(w, http.StatusBadRequest, "User already exist. Login")
		return
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}

	_, err = apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
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

	respondWithJson(w, 200, "")
}

func (apiConfig *Config) signinHandler(w http.ResponseWriter, r *http.Request) {

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
	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if !userExist {
		respondWithError(w, http.StatusUnauthorized, "No User with this mail. Signup")
		return
	}

	user, err := apiConfig.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			respondWithError(w, http.StatusUnauthorized, "Wrong password.")
			return
		}
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}

	err = auth.UpdateAccessToken([]byte(apiConfig.API_KEY), user.ID.String(), user.Email, apiConfig.AccessTokenExpirationMinutes, apiConfig.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating user access token. err: %v", err))
		return
	}

	err = auth.UpdateRefreshToken([]byte(apiConfig.API_KEY), user.ID.String(), apiConfig.RefreshTokenExpirationMinutes, w)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating user refresh token. err: %v", err))
		return
	}

	respondWithJson(w, 200, "")

}

func (apiConfig *Config) passwordChangeHandler(w http.ResponseWriter, r *http.Request) {

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

	user, err := apiConfig.DB.GetUserWithEmail(r.Context(), body.Email)

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
	err = apiConfig.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Email:    body.Email,
		Password: string(newHashedPassword),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating password. err: %v", err))
		return
	}
	respondWithJson(w, 200, "Password Updated")
}

func (apiConfig *Config) getUserHandler(w http.ResponseWriter, r *http.Request) {

	refreshtoken, err := r.Cookie("refreshtoken")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting refreshToken, Try login again. err: %v", err))
		return
	}

	refreshclaims := &jwt.RegisteredClaims{}

	refreshJwt, err := jwt.ParseWithClaims(
		refreshtoken.Value,
		refreshclaims,
		func(token *jwt.Token) (interface{}, error) { return []byte(apiConfig.API_KEY), nil },
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing jwt claims, err: %v", err))
		return
	}

	if refreshclaims.ExpiresAt != nil && refreshclaims.ExpiresAt.Time.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}

	userId, err := refreshJwt.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting issuer from jwt claims, err: %v", err))
		return
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing id, err: %v", err))
		return
	}
	user, err := apiConfig.DB.GetUser(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user, err: %v", err))
		return
	}

	response := struct {
		Data interface{} `json:"data"`
	}{
		Data: struct {
			ID           string `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			RefreshToken string `json:"refresh_token"`
			AccessToken  string `json:"access_token"`
		}{
			ID:           user.ID.String(),
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			RefreshToken: refreshtoken.Value,
			AccessToken:  user.AccessToken.String,
		},
	}

	respondWithJson(w, 200, response)
}

func (apiConfig *Config) refreshTokens(w http.ResponseWriter, r *http.Request) {
	params := struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding r body, err: %v", err))
		return
	}

	// Check if the refresh token is in the blacklist
	if isInTokenBlacklist(params.RefreshToken) {
		respondWithError(w, http.StatusInternalServerError, "refresh token is invalid")
		return
	}

	refreshclaims := &jwt.RegisteredClaims{}

	_, err = jwt.ParseWithClaims(
		params.RefreshToken,
		refreshclaims,
		func(token *jwt.Token) (interface{}, error) { return []byte(apiConfig.API_KEY), nil },
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing jwt claims, err: %v", err))
		return
	}

	userIdString := refreshclaims.Issuer
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing user id, err: %v", err))
		return
	}

	user, err := apiConfig.DB.GetUser(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user with id, err: %v", err))
		return
	}

	refreshExpiration := refreshclaims.ExpiresAt.Time

	if refreshExpiration.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusInternalServerError, "refresh token expired")
		return
	}

	err = auth.UpdateAccessToken([]byte(apiConfig.API_KEY), user.ID.String(), user.Email, apiConfig.AccessTokenExpirationMinutes, apiConfig.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating access token. err: %v", err))
		return
	}
	err = auth.UpdateRefreshToken([]byte(apiConfig.API_KEY), user.ID.String(), apiConfig.RefreshTokenExpirationMinutes, w)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating refresh token. err: %v", err))
		return
	}

	// Add the old refresh token to the blacklist
	addToTokenBlacklist(params.RefreshToken)

	// Debug
	log.Println(refreshExpiration)
	log.Println(user)

}

func (apiConfig *Config) validate(w http.ResponseWriter, r *http.Request) {

	respondWithJson(w, 200, "logged in")
}
