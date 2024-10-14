package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"strconv"
	"time"
)

// Token expiration durations
const (
	accessTokenExpiration  = 1 * time.Hour
	refreshTokenExpiration = 48 * time.Hour
)

// generateToken creates a JWT token for a user with a specified expiration duration.
func generateToken(ID int, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	// Create the claims including user ID and expiration time.
	claims := &models.JWTClaims{
		Sub: strconv.Itoa(ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create and sign the token using the HS256 signing method.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// generateAccessToken creates a JWT access token for a user with a short expiration time.
func generateAccessToken(ID int) (string, error) {
	return generateToken(ID, accessTokenExpiration)
}

// generateAndSaveRefreshToken creates a JWT refresh token for a user with a longer expiration time.
// It also saves the refresh token in the database.
func generateAndSaveRefreshToken(ID int) (string, error) {
	tokenString, err := generateToken(ID, refreshTokenExpiration)
	if err != nil {
		return "", err
	}

	// Save the refresh token in the database for later use.
	expirationTime := time.Now().Add(refreshTokenExpiration)

	if err := postgres.SaveRefreshToken(tokenString, ID, expirationTime); err != nil {
		return "", err
	}

	return tokenString, nil
}
