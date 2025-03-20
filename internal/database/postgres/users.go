package postgres

import (
	"context"
	"database/sql"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"log/slog"
	"time"
)

// AddUserWithCredentials inserts a new user with a username and password into the database.
func AddUserWithCredentials(c *models.Credentials) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, telegram_id, username, email, created_at, version
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var rawTelegramID sql.NullInt64
	var rawEmail sql.NullString

	err := GetDB().QueryRowContext(ctx, query, c.Username, c.Email, c.Password).Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &u.CreatedAt, &u.Version)
	if err != nil {
		return nil, err
	}

	u.TelegramID = extractInt(rawTelegramID)
	u.Email = extractString(rawEmail)
	return &u, nil
}

// AddUserByTelegramID inserts a new user with telegram_id and username into the database.
func AddUserByTelegramID(c *models.Credentials) (*models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username)
		VALUES ($1, $2)
		RETURNING id, telegram_id, username, created_at, version
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var rawTelegramID sql.NullInt64

	err := GetDB().QueryRowContext(ctx, query, c.TelegramID, c.Username).Scan(&u.ID, &rawTelegramID, &u.Username, &u.CreatedAt, &u.Version)
	if err != nil {
		return nil, err
	}

	u.TelegramID = extractInt(rawTelegramID)
	return &u, nil
}

// GetUserById retrieves a user by their ID.
func GetUserById(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var rawTelegramID sql.NullInt64
	var rawEmail sql.NullString
	var rawPassword sql.NullString

	if err := GetDB().QueryRowContext(ctx, query, id).Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &rawPassword, &u.CreatedAt, &u.Version); err != nil {
		return nil, err
	}

	u.TelegramID = extractInt(rawTelegramID)
	u.Email = extractString(rawEmail)
	u.Password = extractString(rawPassword)
	return &u, nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT * FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var rawTelegramID sql.NullInt64
	var rawEmail sql.NullString
	var rawPassword sql.NullString

	if err := GetDB().QueryRowContext(ctx, query, username).Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &rawPassword, &u.CreatedAt, &u.Version); err != nil {
		return nil, err
	}

	u.TelegramID = extractInt(rawTelegramID)
	u.Email = extractString(rawEmail)
	u.Password = extractString(rawPassword)
	return &u, nil
}

// GetUserByTelegramID retrieves a user by their telegram ID.
func GetUserByTelegramID(telegramID int) (*models.User, error) {
	query := `SELECT * FROM users WHERE telegram_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	var rawTelegramID sql.NullInt64
	var rawEmail sql.NullString
	var rawPassword sql.NullString

	if err := GetDB().QueryRowContext(ctx, query, telegramID).Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &rawPassword, &u.CreatedAt, &u.Version); err != nil {
		return nil, err
	}

	u.TelegramID = extractInt(rawTelegramID)
	u.Email = extractString(rawEmail)
	u.Password = extractString(rawPassword)
	return &u, nil
}

// GetUsers retrieves all users from the database.
func GetUsers() ([]*models.User, error) {
	query := `SELECT * FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := GetDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var users []*models.User
	for rows.Next() {
		var u models.User
		var rawTelegramID sql.NullInt64
		var rawEmail sql.NullString

		if err := rows.Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &u.CreatedAt, &u.Version); err != nil {
			return nil, err
		}

		u.TelegramID = extractInt(rawTelegramID)
		u.Email = extractString(rawEmail)
		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user's details based on their ID and version.
func UpdateUser(u *models.User) error {
	query := `
		UPDATE users 
		SET username = $3, email = $4, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING id, telegram_id, username, email, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rawTelegramID sql.NullInt64
	var rawEmail sql.NullString

	err := GetDB().QueryRowContext(ctx, query, u.ID, u.Version, u.Username, u.Email).Scan(&u.ID, &rawTelegramID, &u.Username, &rawEmail, &u.CreatedAt)
	if err != nil {
		return err
	}

	u.TelegramID = extractInt(rawTelegramID)
	u.Email = extractString(rawEmail)
	return nil
}

// DeleteUser removes a user from the database by their ID.
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, id)
	return err
}

func IsUsernameExists(username string) bool {
	query := `SELECT count(*) FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := GetDB().QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// extractInt converts the value of sql.NullInt64 to int.
// If the value is valid, returns it as int, otherwise returns -1.
func extractInt(id sql.NullInt64) int {
	if id.Valid {
		return int(id.Int64)
	}
	return -1
}

// extractString converts the value of sql.NullString to string.
// If the value is valid, it returns it as string; otherwise, it returns an empty string.
func extractString(str sql.NullString) string {
	if str.Valid {
		return str.String
	}
	return ""
}
