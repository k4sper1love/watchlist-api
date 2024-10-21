package postgres

import (
	"context"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/models"
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

	if err := GetDB().QueryRowContext(ctx, query, c.Collection.ID, c.Film.ID).Scan(&c.AddedAt, &c.UpdatedAt); err != nil {
		return err
	}

	collection, err := GetCollection(c.Collection.ID)
	if err != nil {
		return err
	}

	film, err := GetFilm(c.Film.ID)
	if err != nil {
		return err
	}

	c.Collection = *collection
	c.Film = *film

	return nil
}

// GetCollectionFilm retrieves the association of a film in a collection by collection ID and film ID.
func GetCollectionFilm(c *models.CollectionFilm) error {
	query := `
		SELECT added_at, updated_at 
		FROM collection_films
		WHERE collection_id = $1
			AND film_id = $2
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := GetDB().QueryRowContext(ctx, query, c.Collection.ID, c.Film.ID).Scan(&c.AddedAt, &c.UpdatedAt); err != nil {
		return err
	}

	collection, err := GetCollection(c.Collection.ID)
	if err != nil {
		return err
	}

	film, err := GetFilm(c.Film.ID)
	if err != nil {
		return err
	}

	c.Collection = *collection
	c.Film = *film

	return nil
}

// GetCollectionFilms retrieves all films in a collection with optional pagination and sorting.
func GetCollectionFilms(c *models.CollectionFilms, title string, min, max float64, f filters.Filters) (filters.Metadata, error) {
	collection, err := GetCollection(c.Collection.ID)
	if err != nil {
		return filters.Metadata{}, err
	}

	films, metadata, err := GetFilms(-1, c.Collection.ID, title, min, max, f)
	if err != nil {
		return filters.Metadata{}, err
	}

	c.Collection = *collection
	c.Films = films

	return metadata, nil
}

// UpdateCollectionFilm updates the association of a film in a collection.
func UpdateCollectionFilm(c *models.CollectionFilm) error {
	query := `  
       UPDATE collection_films      
       SET added_at = $4, updated_at = CURRENT_TIMESTAMP  
       WHERE collection_id = $1 
         AND film_id = $2 
         AND updated_at = $3       
       RETURNING added_at, updated_at   
       `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, c.Collection.ID, c.Film.ID, c.UpdatedAt, c.AddedAt).Scan(&c.AddedAt, &c.UpdatedAt)
}

// DeleteCollectionFilm removes a film from a collection by collection ID and film ID.
func DeleteCollectionFilm(c *models.CollectionFilm) error {
	query := `  
       DELETE 
       FROM collection_films       
       WHERE collection_id = $1 
         AND film_id = $2    
         `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, c.Collection.ID, c.Film.ID)
	return err
}
