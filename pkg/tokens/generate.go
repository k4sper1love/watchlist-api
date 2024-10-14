package tokens

import (
	"github.com/golang-jwt/jwt"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"strconv"
	"time"
)

// GenerateToken creates a JWT token for a user with a specified expiration duration.
func GenerateToken(secret string, id int, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)

	// Create the claims including user ID and expiration time.
	claims := &models.JWTClaims{
		Sub: strconv.Itoa(id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create and sign the token using the HS256 signing method.
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	return token.SignedString([]byte(secret))
}
