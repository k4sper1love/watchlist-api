package postgres

import (
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddFilm(f *models.Film) error {
	query := `
			INSERT INTO films (user_id, title, year, genre, description, rating, photo_url, comment, is_viewed, user_rating, review)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, created_at
			`

	queryArgs := []interface{}{f.UserId, f.Title, f.Year, f.Genre, f.Description, f.Rating, f.PhotoUrl, f.Comment, f.IsViewed, f.UserRating, f.Review}

	scanArgs := []interface{}{&f.Id, &f.CreatedAt}

	return db.QueryRow(query, queryArgs...).Scan(scanArgs...)
}

func GetFilm(id int) (*models.Film, error) {
	query := `SELECT * FROM films WHERE id = $1`

	var f models.Film
	args := []interface{}{&f.Id, &f.UserId, &f.Title, &f.Year, &f.Genre, &f.Description, &f.Rating, &f.PhotoUrl, &f.Comment, &f.IsViewed, &f.UserRating, &f.Review, &f.CreatedAt}

	err := db.QueryRow(query, id).Scan(args...)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetFilmsByUser(userId int) ([]*models.Film, error) {
	query := `SELECT * FROM films where user_id = $1`

	rows, err := db.Query(query, userId)
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
			&film.Comment,
			&film.IsViewed,
			&film.UserRating,
			&film.Review,
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
	query := `
			UPDATE films
			SET title = $2, year = $3, genre = $4, description = $5, rating = $6, photo_url = $7, comment = $8, is_viewed = $9, user_rating = $10, review = $11
			WHERE id = $1
			RETURNING user_id
			`

	queryArgs := []interface{}{
		film.Id,
		film.Title,
		film.Year,
		film.Genre,
		film.Description,
		film.Rating,
		film.PhotoUrl,
		film.Comment,
		film.IsViewed,
		film.UserRating,
		film.Review,
	}

	return db.QueryRow(query, queryArgs...).Scan(&film.UserId)
}

func DeleteFilm(id int) error {
	query := `DELETE FROM films WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}
