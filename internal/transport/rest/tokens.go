package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"time"
)

// Token expiration durations
const (
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 24 * time.Hour
)

// tokenClaims defines the structure of JWT claims for user authentication.
type tokenClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

// generateToken creates a JWT token for a user with a specified expiration duration.
func generateToken(userId int, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	// Create the claims including user ID and expiration time.
	claims := &tokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create and sign the token using the HS256 signing method.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(config.JwtSecret))
}

// generateAccessToken creates a JWT access token for a user with a short expiration time.
func generateAccessToken(userId int) (string, error) {
	return generateToken(userId, accessTokenExpiration)
}

// generateAndSaveRefreshToken creates a JWT refresh token for a user with a longer expiration time.
// It also saves the refresh token in the database.
func generateAndSaveRefreshToken(userId int) (string, error) {
	tokenString, err := generateToken(userId, refreshTokenExpiration)
	if err != nil {
		return "", err
	}

	// Save the refresh token in the database for later use.
	expirationTime := time.Now().Add(refreshTokenExpiration)

	if err := postgres.SaveRefreshToken(tokenString, userId, expirationTime); err != nil {
		return "", err
	}

	return tokenString, nil
}
