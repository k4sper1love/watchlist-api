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
       INSERT INTO films (user_id, is_favorite, title, year, genre, description, rating, image_url, comment, is_viewed, user_rating, review, url)       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)       RETURNING id, rating, user_rating, created_at, updated_at    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, f.UserID, f.IsFavorite, f.Title, f.Year, f.Genre, f.Description, f.Rating, f.ImageURL, f.Comment, f.IsViewed, f.UserRating, f.Review, f.URL).Scan(&f.ID, &f.Rating, &f.UserRating, &f.CreatedAt, &f.UpdatedAt)
}

// GetFilm retrieves a film by its ID.
func GetFilm(id int) (*models.Film, error) {
	query := `SELECT * FROM films WHERE id = $1`

	var f models.Film
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := GetDB().QueryRowContext(ctx, query, id).Scan(&f.ID, &f.UserID, &f.IsFavorite, &f.Title, &f.Year, &f.Genre, &f.Description, &f.Rating, &f.ImageURL, &f.Comment, &f.IsViewed, &f.UserRating, &f.Review, &f.URL, &f.CreatedAt, &f.UpdatedAt); err != nil {
		return nil, err
	}

	return &f, nil
}

// GetFilms retrieves films for a specific user based on filters and pagination.
func GetFilms(userID int, input *models.FilmsQueryInput) ([]models.Film, filters.Metadata, error) {
	query, args := buildFilmsQuery(userID, input)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	var films []models.Film
	totalRecords := 0

	for rows.Next() {
		var film models.Film
		if err := rows.Scan(&totalRecords, &film.ID, &film.UserID, &film.IsFavorite, &film.Title, &film.Year, &film.Genre, &film.Description, &film.Rating, &film.ImageURL, &film.Comment, &film.IsViewed, &film.UserRating, &film.Review, &film.URL, &film.CreatedAt, &film.UpdatedAt); err != nil {
			return nil, filters.Metadata{}, err
		}
		films = append(films, film)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, input.Filters.Page, input.Filters.PageSize)
	return films, metadata, nil
}

// UpdateFilm updates the details of an existing film.
func UpdateFilm(film *models.Film) error {
	query := `  
       UPDATE films      
       SET title = $3, year = $4, genre = $5, description = $6, rating = $7, image_url = $8, comment = $9, 
           is_viewed = $10, user_rating = $11, review = $12,  url = $13, is_favorite = $14, updated_at = CURRENT_TIMESTAMP     
       WHERE id = $1 AND updated_at = $2     
       RETURNING user_id, updated_at    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, film.ID, film.UpdatedAt, film.Title, film.Year, film.Genre, film.Description, film.Rating, film.ImageURL, film.Comment, film.IsViewed, film.UserRating, film.Review, film.URL, film.IsFavorite).Scan(&film.UserID, &film.UpdatedAt)
}

// DeleteFilm removes a film by its ID.
func DeleteFilm(id int) error {
	query := `DELETE FROM films WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, id)
	return err
}

// buildFilmsQuery constructs the SQL query and arguments for retrieving films.
func buildFilmsQuery(userID int, input *models.FilmsQueryInput) (string, []interface{}) {
	query := `
        SELECT COUNT(*) OVER(), f.*
        FROM films f
        WHERE f.user_id = $1
          AND (LOWER(f.title) ILIKE '%%' || LOWER($2) || '%%' OR $2 = '') 
          AND f.id NOT IN (
              SELECT cf.film_id
              FROM collection_films cf
              WHERE cf.collection_id = $3
          )
    `

	args := []interface{}{userID, input.Title, input.ExcludeCollection}

	return addFilmsFiltersToQuery(query, args, input)
}

func addFilmsFiltersToQuery(query string, args []interface{}, input *models.FilmsQueryInput) (string, []interface{}) {
	minRating, maxRating, ratingFlag, err := parseRangeOrExact(input.Rating)
	if err != nil {
		return "", nil
	}

	minUserRating, maxUserRating, userRatingFlag, err := parseRangeOrExact(input.UserRating)
	if err != nil {
		return "", nil
	}

	minYear, maxYear, yearFlag, err := parseYearOrRange(input.Year)
	if err != nil {
		return "", nil
	}

	if ratingFlag == 1 {
		query += fmt.Sprintf(" AND f.rating BETWEEN $%d AND $%d", len(args)+1, len(args)+2)
		args = append(args, minRating, maxRating)
	} else if ratingFlag == 0 {
		query += fmt.Sprintf(" AND f.rating = $%d", len(args)+1)
		args = append(args, minRating)
	}

	if userRatingFlag == 1 {
		query += fmt.Sprintf(" AND f.user_rating BETWEEN $%d AND $%d", len(args)+1, len(args)+2)
		args = append(args, minUserRating, maxUserRating)
	} else if userRatingFlag == 0 {
		query += fmt.Sprintf(" AND f.user_rating = $%d", len(args)+1)
		args = append(args, minUserRating)
	}

	if yearFlag == 1 {
		query += fmt.Sprintf(" AND f.year BETWEEN $%d AND $%d", len(args)+1, len(args)+2)
		args = append(args, minYear, maxYear)
	} else if yearFlag == 0 {
		query += fmt.Sprintf(" AND f.year = $%d", len(args)+1)
		args = append(args, minYear)
	}

	if input.IsViewed != nil {
		query += " AND f.is_viewed = $" + fmt.Sprint(len(args)+1)
		args = append(args, *input.IsViewed)
	}

	if input.IsFavorite != nil {
		query += " AND f.is_favorite = $" + fmt.Sprint(len(args)+1)
		args = append(args, *input.IsFavorite)
	}

	if input.HasURL != nil {
		if *input.HasURL {
			query += " AND f.url IS NOT NULL AND f.url <> ''"
		} else {
			query += " AND (f.url IS NULL OR f.url = '')"
		}
	}

	query += `
        ORDER BY %s %s, f.id
        LIMIT $%d OFFSET $%d
    `

	args = append(args, input.Filters.Limit(), input.Filters.Offset())
	query = fmt.Sprintf(query, input.Filters.SortColumn(), input.Filters.SortDirection(), len(args)-1, len(args))

	return query, args
}
