package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"time"
)

func GenerateAccessToken(userId int) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.TokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.TokenPassword))
}

func GenerateAndSaveRefreshToken(userId int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.TokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
