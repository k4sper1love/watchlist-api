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
func SaveRefreshToken(refreshToken string, userID int, expiresAt time.Time) error {
	hashedToken := hashToken(refreshToken)

	query := `INSERT INTO refresh_tokens(token, user_id, expires_at) values ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, hashedToken, userID, expiresAt)
	return err
}

// RevokeRefreshToken marks a refresh token as revoked in the database.
func RevokeRefreshToken(refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, hashedToken)
	return err
}

// IsRefreshTokenRevoked checks if a refresh token has been revoked.
func IsRefreshTokenRevoked(refreshToken string) (bool, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT revoked FROM refresh_tokens WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var revoked bool
	if err := GetDB().QueryRowContext(ctx, query, hashedToken).Scan(&revoked); err != nil {
		return false, err
	}

	return revoked, nil
}

// GetIdFromRefreshToken retrieves the user ID associated with a refresh token.
func GetIdFromRefreshToken(refreshToken string) (int, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT user_id FROM refresh_tokens WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userID int
	if err := GetDB().QueryRowContext(ctx, query, hashedToken).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}
