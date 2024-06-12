package main

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (apiconfig *Config) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				respondWithError(w, http.StatusUnauthorized, "No access token")
				return
			}
			respondWithError(w, http.StatusBadRequest, "Bad request")
			return
		}

		tokenString := cookie.Value
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) { return []byte(apiconfig.API_KEY), nil },
		)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				respondWithError(w, http.StatusUnauthorized, "Invalid Signature")
				return
			}
			respondWithError(w, http.StatusBadRequest, "Bad request")

			return
		}
		log.Println("token string", tokenString)
		log.Println("token ", *token)
		if !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}
		issuer, _ := token.Claims.GetIssuer()
		if issuer != "Muhammad Olamide" {
			respondWithError(w, http.StatusUnauthorized, "Wrong Issuer")

			return
		}

		next.ServeHTTP(w, r)
	})
}
