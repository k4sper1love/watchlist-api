package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddFilm(film *models.Film) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			INSERT INTO films (creator_id, title, year, genre, description, rating, photo_url)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, creator_id, title, year, genre, description, rating, photo_url
			`

	queryArgs := []interface{}{
		film.CreatorId,
		film.Title,
		film.Year,
		film.Genre,
		film.Description,
		film.Rating,
		film.PhotoUrl,
	}

	scanArgs := []interface{}{
		&film.Id,
		&film.CreatorId,
		&film.Title,
		&film.Year,
		&film.Genre,
		&film.Description,
		&film.Rating,
		&film.PhotoUrl,
	}

	return db.QueryRow(query, queryArgs...).Scan(scanArgs...)
}

func GetFilm(id int) (*models.Film, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM films WHERE id = $1`

	var film models.Film
	args := []interface{}{
		&film.Id,
		&film.CreatorId,
		&film.Title,
		&film.Year,
		&film.Genre,
		&film.Description,
		&film.Rating,
		&film.PhotoUrl,
	}

	err := db.QueryRow(query, id).Scan(args...)
	if err != nil {
		return nil, err
	}

	return &film, nil
}

func GetFilms() ([]*models.Film, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM films`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []*models.Film
	for rows.Next() {
		var film models.Film
		args := []interface{}{
			&film.Id,
			&film.CreatorId,
			&film.Title,
			&film.Year,
			&film.Genre,
			&film.Description,
			&film.Rating,
			&film.PhotoUrl,
		}
		err = rows.Scan(args...)
		if err != nil {
			return nil, err
		}
		films = append(films, &film)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return films, nil
}

func UpdateFilm(film *models.Film) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			UPDATE films
			SET title = $2, year = $3, genre = $4, description = $5, rating = $6, photo_url = $7
			WHERE id = $1
			RETURNING id, creator_id, title, year, genre, description, rating, photo_url
			`

	queryArgs := []interface{}{
		film.Id,
		film.Title,
		film.Year,
		film.Genre,
		film.Description,
		film.Rating,
		film.PhotoUrl,
	}

	scanArgs := []interface{}{
		&film.Id,
		&film.CreatorId,
		&film.Title,
		&film.Year,
		&film.Genre,
		&film.Description,
		&film.Rating,
		&film.PhotoUrl,
	}
	return db.QueryRow(query, queryArgs...).Scan(scanArgs...)
}

func DeleteFilm(id int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM films WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}
