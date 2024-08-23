package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddUserFilm(userFilm *models.UserFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			INSERT INTO user_films (user_id, film_id, comment)
			VALUES ($1, $2, $3)
			RETURNING user_id, film_id, comment, added_at
			`

	scanArgs := []interface{}{&userFilm.UserId, &userFilm.FilmId, &userFilm.Comment, &userFilm.AddedAt}

	return db.QueryRow(query, userFilm.UserId, userFilm.FilmId, &userFilm.Comment).Scan(scanArgs...)
}

func GetUserFilm(userId, filmId int) (*models.UserFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM user_films WHERE user_id = $1 AND film_id = $2`

	var userFilm models.UserFilm
	args := []interface{}{&userFilm.UserId, &userFilm.FilmId, &userFilm.Comment, &userFilm.AddedAt}
	err := db.QueryRow(query, userId, filmId).Scan(args...)
	if err != nil {
		return nil, err
	}

	return &userFilm, nil
}

func GetUserFilms(userId int) ([]*models.UserFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM user_films WHERE user_id = $1`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userFilms []*models.UserFilm
	for rows.Next() {
		var userFilm models.UserFilm
		err = rows.Scan(&userFilm.UserId, &userFilm.FilmId, &userFilm.Comment, &userFilm.AddedAt)
		if err != nil {
			return nil, err
		}
		userFilms = append(userFilms, &userFilm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userFilms, nil
}

func UpdateUserFilm(userFilm *models.UserFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `UPDATE user_films SET comment = $2 WHERE user_id = $1 RETURNING user_id, film_id, comment, added_at`

	args := []interface{}{&userFilm.UserId, &userFilm.FilmId, &userFilm.Comment, &userFilm.AddedAt}
	return db.QueryRow(query, userFilm.UserId, userFilm.Comment).Scan(args...)
}

func DeleteUserFilm(userId, filmId int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM user_films WHERE user_id = $1 AND film_id = $2`

	_, err := db.Exec(query, userId, filmId)
	return err
}
