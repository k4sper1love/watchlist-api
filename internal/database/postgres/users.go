package postgres

import (
	"errors"
	"github.com/k4sper1love/wishlist-api/internal/models"
)

func AddUser(user *models.User) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, created_at`

	return db.QueryRow(query, user.Username).Scan(&user.Id, &user.CreatedAt)
}

func GetUserById(id int) (*models.User, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM users WHERE id = $1`

	var user models.User
	err := db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.Username, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(user *models.User) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `UPDATE users SET username=$1 WHERE id=$2 RETURNING username, created_at`

	err := db.QueryRow(query, user.Username, user.Id).Scan(&user.Username, &user.CreatedAt)
	return err
}

func DeleteUser(id int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM users WHERE id=$1`

	_, err := db.Exec(query, id)
	return err
}
