package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddViewedFilm(v *models.ViewedFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `INSERT INTO viewed_films (user_id, film_id, rating, review) VALUES ($1, $2, $3, $4) RETURNING viewed_at`

	return db.QueryRow(query, v.UserId, v.FilmId, v.Rating, v.Review).Scan(&v.ViewedAt)
}

func GetViewedFilm(userId, filmId int) (*models.ViewedFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM viewed_films WHERE user_id = $1 AND film_id = $2`

	var v models.ViewedFilm
	err := db.QueryRow(query, userId, filmId).Scan(&v.UserId, &v.FilmId, &v.Rating, &v.Review, &v.ViewedAt)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func GetViewedFilms(userId int) ([]*models.ViewedFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM viewed_films WHERE user_id = $1`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewedFilms []*models.ViewedFilm
	for rows.Next() {
		var v models.ViewedFilm
		err = rows.Scan(&v.UserId, &v.FilmId, &v.Rating, &v.Review, &v.ViewedAt)
		if err != nil {
			return nil, err
		}
		viewedFilms = append(viewedFilms, &v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return viewedFilms, nil
}

func UpdateViewedFilm(v *models.ViewedFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			UPDATE viewed_films 
			SET rating = $3, review = $4
			WHERE user_id = $1 AND film_id = $2
			RETURNING *
			`

	args := []interface{}{&v.UserId, &v.FilmId, &v.Rating, &v.Review, &v.ViewedAt}

	return db.QueryRow(query, v.UserId, v.FilmId, v.Rating, v.Review).Scan(args...)
}

func DeleteViewedFilm(userId, filmId int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM viewed_films WHERE user_id = $1 AND film_id = $2`

	_, err := db.Exec(query, userId, filmId)
	return err
}
