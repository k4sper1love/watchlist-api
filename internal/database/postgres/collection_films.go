package postgres

import (
	"context"
	"fmt"
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

	query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), f.*
        FROM films f
        WHERE f.id IN (
            SELECT cf.film_id
            FROM collection_films cf
            WHERE cf.collection_id = $1
        )
 		  AND (LOWER(f.title) ILIKE '%%' || LOWER($2) || '%%' OR $2 = '') 
          AND (f.rating >= $3 OR $3 = 0)
          AND (f.rating <= $4 OR $4 = 0)
        ORDER BY %s %s, f.id
        LIMIT $5 OFFSET $6
    `,
		f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := GetDB().QueryContext(ctx, query, c.Collection.ID, title, min, max, f.Limit(), f.Offset())
	if err != nil {
		return filters.Metadata{}, err
	}
	defer rows.Close()

	var films []models.Film
	totalRecords := 0

	for rows.Next() {
		var film models.Film
		if err := rows.Scan(&totalRecords, &film.ID, &film.UserID, &film.Title, &film.Year, &film.Genre, &film.Description, &film.Rating, &film.ImageURL, &film.Comment, &film.IsViewed, &film.UserRating, &film.Review, &film.URL, &film.CreatedAt, &film.UpdatedAt); err != nil {
			return filters.Metadata{}, err
		}
		films = append(films, film)
	}

	if err = rows.Err(); err != nil {
		return filters.Metadata{}, err
	}

	c.Collection = *collection
	c.Films = films

	metadata := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
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
