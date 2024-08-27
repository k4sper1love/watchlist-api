package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"time"
)

type tokenClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

func generateAccessToken(userId int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &tokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	return token.SignedString([]byte(config.TokenPassword))
}

func generateAndSaveRefreshToken(userId int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &tokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tokenString, err := token.SignedString([]byte(config.TokenPassword))
	if err != nil {
		return "", err
	}

	err = postgres.SaveRefreshToken(tokenString, userId, expirationTime)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
