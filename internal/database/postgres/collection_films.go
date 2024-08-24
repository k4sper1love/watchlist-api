package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddCollectionFilm(c *models.CollectionFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `INSERT INTO collection_films (collection_id, film_id, comment) VALUES ($1, $2, $3) RETURNING added_at`

	return db.QueryRow(query, c.CollectionId, c.FilmId, c.Comment).Scan(&c.AddedAt)
}

func GetCollectionFilm(collectionId, filmId int) (*models.CollectionFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	var c models.CollectionFilm
	err := db.QueryRow(query, collectionId, filmId).Scan(&c.CollectionId, &c.FilmId, &c.Comment, &c.AddedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func GetCollectionFilms(collectionId int) ([]*models.CollectionFilm, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM collection_films WHERE collection_id = $1`

	rows, err := db.Query(query, collectionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collectionFilms []*models.CollectionFilm
	for rows.Next() {
		var c models.CollectionFilm
		err = rows.Scan(&c.CollectionId, &c.FilmId, &c.Comment, &c.AddedAt)
		if err != nil {
			return nil, err
		}
		collectionFilms = append(collectionFilms, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collectionFilms, nil
}

func UpdateCollectionFilm(c *models.CollectionFilm) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			UPDATE collection_films 
			SET comment = $3 
			WHERE collection_id = $1 AND film_id = $2
			RETURNING *
			`

	args := []interface{}{&c.CollectionId, &c.FilmId, &c.Comment, &c.AddedAt}

	return db.QueryRow(query, c.CollectionId, c.FilmId, c.Comment).Scan(args...)
}

func DeleteCollectionFilm(collectionId, filmId int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	_, err := db.Exec(query, collectionId, filmId)
	return err
}
