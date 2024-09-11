package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/internal/config"
	"github.com/k4sper1love/watchlist-api/internal/database/postgres"
	"time"
)

// tokenClaims defines the structure of JWT claims for user authentication.
type tokenClaims struct {
	UserId             int `json:"user_id"`
	jwt.StandardClaims     // Standard JWT claims, such as expiration time.
}

// generateAccessToken creates a JWT access token for a user with a short expiration time (15 minutes).
// Returns the token as a string and any error encountered.
func generateAccessToken(userId int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	// Create the claims including user ID and expiration time.
	claims := &tokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new token using the HS256 signing method and the claims.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	// Sign the token with the secret key and return the token string.
	return token.SignedString([]byte(config.TokenPass))
}

// generateAndSaveRefreshToken creates a JWT refresh token for a user with a longer expiration time (24 hours).
// It also saves the refresh token in the database.
// Returns the token as a string and any error encountered.
func generateAndSaveRefreshToken(userId int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims including user ID and expiration time.
	claims := &tokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new token using the HS256 signing method and the claims.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	// Sign the token with the secret key.
	tokenString, err := token.SignedString([]byte(config.TokenPass))
	if err != nil {
		return "", err
	}

	// Save the refresh token in the database for later use.
	err = postgres.SaveRefreshToken(tokenString, userId, expirationTime)
	if err != nil {
		return "", err
	}

	// Return the signed token string.
	return tokenString, nil
}
