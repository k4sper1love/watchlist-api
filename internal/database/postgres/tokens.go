package postgres

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func SaveRefreshToken(refreshToken string, userId int, expiresAt time.Time) error {
	hashedToken := hashToken(refreshToken)

	query := `INSERT INTO refresh_tokens(token, user_id, expires_at) values ($1, $2, $3)`

	_, err := db.Exec(query, hashedToken, userId, expiresAt)
	return err
}

func RevokeRefreshToken(refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`

	_, err := db.Exec(query, hashedToken)
	return err
}

func IsRefreshTokenRevoked(refreshToken string) (bool, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT revoked FROM refresh_tokens WHERE token = $1`

	var revoked bool
	err := db.QueryRow(query, hashedToken).Scan(&revoked)
	if err != nil {
		return false, err
	}

	if revoked {
		return true, nil
	}

	return false, nil
}

func GetIdFromRefreshToken(refreshToken string) (int, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT user_id FROM refresh_tokens WHERE token = $1`

	var userId int
	err := db.QueryRow(query, hashedToken).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
