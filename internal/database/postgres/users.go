package postgres

import (
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddUser(u *models.User) error {
	query := `
			INSERT INTO users (username, email, password)
			VALUES ($1, $2, $3) 
			RETURNING id, created_at
			`

	return db.QueryRow(query, u.Username, u.Email, u.Password).Scan(&u.Id, &u.CreatedAt)
}

func GetUserById(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user models.User
	err := db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user models.User
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsers() ([]*models.User, error) {
	query := `SELECT * FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return db.QueryRow(query, user.Id, user.Username).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
}

func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}
