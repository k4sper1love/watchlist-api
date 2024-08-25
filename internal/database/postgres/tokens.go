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

func IsRefreshTokenValid(refreshToken string) (bool, error) {
	hashedToken := hashToken(refreshToken)

	query := `SELECT revoked, expires_at FROM refresh_tokens WHERE token = $1`

	var revoked bool
	var exp time.Time
	err := db.QueryRow(query, hashedToken).Scan(&revoked, &exp)
	if err != nil {
		return false, err
	}

	if revoked || time.Now().After(exp) {
		return false, nil
	}

	return true, nil
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
