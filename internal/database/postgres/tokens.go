package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// hashToken hashes the given token using SHA-256 and returns its hexadecimal representation.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// SaveRefreshToken stores a refresh token, its associated user ID, and its expiration time in the database.
//
// Returns an error if insertion fails.
func SaveRefreshToken(refreshToken string, userId int, expiresAt time.Time) error {
	hashedToken := hashToken(refreshToken)

	query := `INSERT INTO refresh_tokens(token, user_id, expires_at) values ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, hashedToken, userId, expiresAt)
	return err
}

// RevokeRefreshToken marks a refresh token as revoked in the database.
//
// Returns an error if the update fails.
func RevokeRefreshToken(refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, hashedToken)
	return err
}

// IsRefreshTokenRevoked checks if a refresh token has been revoked.
//
// Returns true if the token is revoked, false otherwise, and an error if the query fails.
func IsRefreshTokenRevoked(refreshToken string) (bool, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT revoked FROM refresh_tokens WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var revoked bool
	err := db.QueryRowContext(ctx, query, hashedToken).Scan(&revoked)
	if err != nil {
		return false, err
	}

	if revoked {
		return true, nil
	}

	return false, nil
}

// GetIdFromRefreshToken retrieves the user ID associated with a refresh token.
//
// Returns the user ID and an error if the query fails.
func GetIdFromRefreshToken(refreshToken string) (int, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT user_id FROM refresh_tokens WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userId int
	err := db.QueryRowContext(ctx, query, hashedToken).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
