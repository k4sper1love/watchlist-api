package postgres

import (
	"context"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"time"
)

// AddUser inserts a new user into the database and returns the user with ID, creation timestamp, and version.
func AddUser(u *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, u.Username, u.Email, u.Password).Scan(&u.Id, &u.CreatedAt, &u.Version)
}

// GetUserById retrieves a user by their ID.
func GetUserById(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	if err := GetDB().QueryRowContext(ctx, query, id).Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.Version); err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserByEmail retrieves a user by their email address.
func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	if err := GetDB().QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.Version); err != nil {
		return nil, err
	}

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
			sl.Log.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.Version); err != nil {
			return nil, err
		}
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
		SET username = $3, version = version + 1
		WHERE id = $1 AND version = $2
		RETURNING id, username, email, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, u.Id, u.Version, u.Username).Scan(&u.Id, &u.Username, &u.Email, &u.CreatedAt)
}

// DeleteUser removes a user from the database by their ID.
func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, id)
	return err
}
