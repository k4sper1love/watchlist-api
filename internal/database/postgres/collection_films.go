package postgres

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/filters"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddCollectionFilm(c *models.CollectionFilm) error {
	query := `INSERT INTO collection_films (collection_id, film_id) VALUES ($1, $2) RETURNING added_at`

	return db.QueryRow(query, c.CollectionId, c.FilmId).Scan(&c.AddedAt)
}

func GetCollectionFilm(collectionId, filmId int) (*models.CollectionFilm, error) {
	query := `SELECT * FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	var c models.CollectionFilm
	err := db.QueryRow(query, collectionId, filmId).Scan(&c.CollectionId, &c.FilmId, &c.AddedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func GetCollectionFilms(collectionId int, f filters.Filters) ([]*models.CollectionFilm, filters.Metadata, error) {
	query := fmt.Sprintf(
		`	
			SELECT count(*) OVER(), *
			FROM collection_films 
			WHERE collection_id = $1 
			ORDER BY %s %s, film_id
			LIMIT $2 OFFSET $3
			`,
		f.SortColumn(), f.SortDirection())

	rows, err := db.Query(query, collectionId, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0

	var collectionFilms []*models.CollectionFilm
	for rows.Next() {
		var c models.CollectionFilm
		err = rows.Scan(&totalRecords, &c.CollectionId, &c.FilmId, &c.AddedAt)
		if err != nil {
			return nil, filters.Metadata{}, err
		}
		collectionFilms = append(collectionFilms, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)

	return collectionFilms, metadata, nil
}

func UpdateCollectionFilm(c *models.CollectionFilm) error {
	query := `
			UPDATE collection_films 
			SET added_at = $3 
			WHERE collection_id = $1 AND film_id = $2
			RETURNING *
			`

	args := []interface{}{&c.CollectionId, &c.FilmId, &c.AddedAt}

	return db.QueryRow(query, c.CollectionId, c.FilmId, c.AddedAt).Scan(args...)
}

func DeleteCollectionFilm(collectionId, filmId int) error {
	query := `DELETE FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	_, err := db.Exec(query, collectionId, filmId)
	return err
}
