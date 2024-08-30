package postgres

import (
	"context"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"time"
)

func AddUser(u *models.User) error {
	query := `
			INSERT INTO users (username, email, password)
			VALUES ($1, $2, $3) 
			RETURNING id, created_at
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, u.Username, u.Email, u.Password).Scan(&u.Id, &u.CreatedAt)
}

func GetUserById(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	err := db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	err := db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsers() ([]*models.User, error) {
	query := `SELECT * FROM users`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.Username, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = $2 WHERE id = $1 RETURNING id, username, email, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, user.Id, user.Username).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
}

func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, id)
	return err
}
