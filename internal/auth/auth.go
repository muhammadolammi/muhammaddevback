package auth

import (
	"context"
	"database/sql"
	"fmt"
	"muhammaddev/internal/database"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJwtTokenString(signgingKey []byte, userId, tokenName string, tokenExpiration int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("%v", userId),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(tokenExpiration) * time.Minute)),
		Subject:   tokenName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signgingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func UpdateRefreshToken(signgingKey []byte, userId string, expirationTime int, w http.ResponseWriter) error {

	jwtRefreshTokenString, err := MakeJwtTokenString(signgingKey, userId, "refreshtoken", expirationTime)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshtoken",
		Value:    jwtRefreshTokenString,
		Expires:  time.Now().UTC().Add(time.Duration(expirationTime) * time.Minute),
		HttpOnly: true,
	})
	return nil

}
func UpdateAccessToken(signgingKey []byte, userId string, userEmail string, expirationTime int, DB *database.Queries) error {
	jwtAccessTokenString, err := MakeJwtTokenString(signgingKey, userId, "accesstoken", expirationTime)
	if err != nil {
		return err
	}

	err = DB.UpdateAccessToken(context.Background(), database.UpdateAccessTokenParams{
		AccessToken: sql.NullString{
			Valid:  true,
			String: jwtAccessTokenString,
		},
		Email: userEmail,
	})
	if err != nil {
		return err
	}
	return nil
}
