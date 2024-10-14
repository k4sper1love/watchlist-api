package rest

import (
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"github.com/k4sper1love/watchlist-api/pkg/tokens"
	"time"
)

// Token expiration durations
const (
	accessTokenExpiration  = 1 * time.Hour
	refreshTokenExpiration = 48 * time.Hour
)

// generateAccessToken creates a JWT access token for a user with a short expiration time.
func generateAccessToken(id int) (string, error) {
	return tokens.GenerateToken(id, accessTokenExpiration)
}

// generateAndSaveRefreshToken creates a JWT refresh token for a user with a longer expiration time.
// It also saves the refresh token in the database.
func generateAndSaveRefreshToken(id int) (string, error) {
	tokenString, err := tokens.GenerateToken(id, refreshTokenExpiration)
	if err != nil {
		return "", err
	}

	// Save the refresh token in the database for later use.
	expirationTime := time.Now().Add(refreshTokenExpiration)

	if err := postgres.SaveRefreshToken(tokenString, id, expirationTime); err != nil {
		return "", err
	}

	return tokenString, nil
}
