package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJwtTokenString(signgingKey []byte, userid string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "Muhammad Olamide",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
		Subject:   fmt.Sprintf("%v", userid),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signgingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func GenerateRefreshToken() (string, error) {

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
