package postgres

import (
	"context"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
	"time"
)

// AddCollectionFilm adds a film to a collection.
func AddCollectionFilm(c *models.CollectionFilm) error {
	query := `
		INSERT INTO collection_films (collection_id, film_id)
		VALUES ($1, $2)
		RETURNING added_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, c.CollectionId, c.FilmId).Scan(&c.AddedAt, &c.UpdatedAt)
}

// GetCollectionFilm retrieves the association of a film in a collection by collection ID and film ID.
func GetCollectionFilm(collectionId, filmId int) (*models.CollectionFilm, error) {
	query := `
		SELECT *
		FROM collection_films
		WHERE collection_id = $1 AND film_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var c models.CollectionFilm
	if err := GetDB().QueryRowContext(ctx, query, collectionId, filmId).Scan(&c.CollectionId, &c.FilmId, &c.AddedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCollectionFilms retrieves all films in a collection with optional pagination and sorting.
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

	rows, err := GetDB().QueryContext(ctx, query, collectionId, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			sl.Log.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var collectionFilms []*models.CollectionFilm
	totalRecords := 0

	for rows.Next() {
		var c models.CollectionFilm
		if err := rows.Scan(&totalRecords, &c.CollectionId, &c.FilmId, &c.AddedAt, &c.UpdatedAt); err != nil {
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

// UpdateCollectionFilm updates the association of a film in a collection.
func UpdateCollectionFilm(c *models.CollectionFilm) error {
	query := `
		UPDATE collection_films 
		SET added_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE collection_id = $1 AND film_id = $2 AND updated_at = $3
		RETURNING added_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, c.CollectionId, c.FilmId, c.UpdatedAt, c.AddedAt).Scan(&c.AddedAt, &c.UpdatedAt)
}

// DeleteCollectionFilm removes a film from a collection by collection ID and film ID.
func DeleteCollectionFilm(collectionId, filmId int) error {
	query := `
		DELETE FROM collection_films
		WHERE collection_id = $1 AND film_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, collectionId, filmId)
	return err
}
