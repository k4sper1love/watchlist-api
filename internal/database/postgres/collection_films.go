package postgres

import (
	"context"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/filters"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"log"
	"time"
)

func AddCollectionFilm(c *models.CollectionFilm) error {
	query := `INSERT INTO collection_films (collection_id, film_id) VALUES ($1, $2) RETURNING created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, c.CollectionId, c.FilmId).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func GetCollectionFilm(collectionId, filmId int) (*models.CollectionFilm, error) {
	query := `SELECT * FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var c models.CollectionFilm
	err := db.QueryRowContext(ctx, query, collectionId, filmId).Scan(&c.CollectionId, &c.FilmId, &c.CreatedAt, &c.UpdatedAt)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, collectionId, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	totalRecords := 0

	var collectionFilms []*models.CollectionFilm
	for rows.Next() {
		var c models.CollectionFilm
		err = rows.Scan(&totalRecords, &c.CollectionId, &c.FilmId, &c.CreatedAt, &c.UpdatedAt)
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
			SET created_at = $4, updated_at = CURRENT_TIMESTAMP
			WHERE collection_id = $1 AND film_id = $2 AND updated_at = $3
			RETURNING created_at, updated_at
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, c.CollectionId, c.FilmId, c.UpdatedAt, c.CreatedAt).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func DeleteCollectionFilm(collectionId, filmId int) error {
	query := `DELETE FROM collection_films WHERE collection_id = $1 AND film_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, collectionId, filmId)
	return err
}
