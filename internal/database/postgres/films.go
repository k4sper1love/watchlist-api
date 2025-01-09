package postgres

import (
	"context"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"time"
)

// AddFilm inserts a new film into the database and returns its ID, creation, and update timestamps.
func AddFilm(f *models.Film) error {
	query := `  
       INSERT INTO films (user_id, title, year, genre, description, rating, image_url, comment, is_viewed, user_rating, review, url)       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)       RETURNING id, created_at, updated_at    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, f.UserID, f.Title, f.Year, f.Genre, f.Description, f.Rating, f.ImageURL, f.Comment, f.IsViewed, f.UserRating, f.Review, f.URL).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

// GetFilm retrieves a film by its ID.
func GetFilm(id int) (*models.Film, error) {
	query := `SELECT * FROM films WHERE id = $1`

	var f models.Film
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := GetDB().QueryRowContext(ctx, query, id).Scan(&f.ID, &f.UserID, &f.Title, &f.Year, &f.Genre, &f.Description, &f.Rating, &f.ImageURL, &f.Comment, &f.IsViewed, &f.UserRating, &f.Review, &f.URL, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, err
	}

	return &f, nil
}

// GetFilms retrieves films for a specific user based on filters and pagination.
func GetFilms(userID int, title string, min, max float64, excludeCollectionID int, f filters.Filters) ([]models.Film, filters.Metadata, error) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), f.*
        FROM films f
        WHERE f.user_id = $1
 		  AND (LOWER(f.title) ILIKE '%%' || LOWER($2) || '%%' OR $2 = '') 
          AND (f.rating >= $3 OR $3 = 0)
          AND (f.rating <= $4 OR $4 = 0)
          AND f.id NOT IN (
              SELECT cf.film_id
              FROM collection_films cf
              WHERE cf.collection_id = $5
          )
        ORDER BY %s %s, f.id
        LIMIT $6 OFFSET $7
    `,
		f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := GetDB().QueryContext(ctx, query, userID, title, min, max, excludeCollectionID, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	var films []models.Film
	totalRecords := 0

	for rows.Next() {
		var film models.Film
		if err := rows.Scan(&totalRecords, &film.ID, &film.UserID, &film.Title, &film.Year, &film.Genre, &film.Description, &film.Rating, &film.ImageURL, &film.Comment, &film.IsViewed, &film.UserRating, &film.Review, &film.URL, &film.CreatedAt, &film.UpdatedAt); err != nil {
			return nil, filters.Metadata{}, err
		}
		films = append(films, film)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
	return films, metadata, nil
}

// UpdateFilm updates the details of an existing film.
func UpdateFilm(film *models.Film) error {
	query := `  
       UPDATE films      
       SET title = $3, year = $4, genre = $5, description = $6, rating = $7, image_url = $8, comment = $9, 
           is_viewed = $10, user_rating = $11, review = $12,  url = $13, updated_at = CURRENT_TIMESTAMP     
       WHERE id = $1 AND updated_at = $2     
       RETURNING user_id, updated_at    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, film.ID, film.UpdatedAt, film.Title, film.Year, film.Genre, film.Description, film.Rating, film.ImageURL, film.Comment, film.IsViewed, film.UserRating, film.Review, film.URL).Scan(&film.UserID, &film.UpdatedAt)
}

// DeleteFilm removes a film by its ID.
func DeleteFilm(id int) error {
	query := `DELETE FROM films WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, id)
	return err
}
