package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddFilm(f *models.Film) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			INSERT INTO films (user_id, title, year, genre, description, rating, photo_url)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, user_id, title, year, genre, description, rating, photo_url, created_at
			`

	queryArgs := []interface{}{f.UserId, f.Title, f.Year, f.Genre, f.Description, f.Rating, f.PhotoUrl}

	scanArgs := []interface{}{&f.Id, &f.UserId, &f.Title, &f.Year, &f.Genre, &f.Description, &f.Rating, &f.PhotoUrl, &f.CreatedAt}

	return db.QueryRow(query, queryArgs...).Scan(scanArgs...)
}

func GetFilm(id int) (*models.Film, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM films WHERE id = $1`

	var f models.Film
	args := []interface{}{&f.Id, &f.UserId, &f.Title, &f.Year, &f.Genre, &f.Description, &f.Rating, &f.PhotoUrl, &f.CreatedAt}

	err := db.QueryRow(query, id).Scan(args...)
	if err != nil {
		return nil, err
	}

	return &f, nil
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
			&film.UserId,
			&film.Title,
			&film.Year,
			&film.Genre,
			&film.Description,
			&film.Rating,
			&film.PhotoUrl,
			&film.CreatedAt,
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
			RETURNING id, user_id, title, year, genre, description, rating, photo_url, created_at
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
		&film.UserId,
		&film.Title,
		&film.Year,
		&film.Genre,
		&film.Description,
		&film.Rating,
		&film.PhotoUrl,
		&film.CreatedAt,
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
