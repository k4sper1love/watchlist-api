package postgres

import (
	"context"
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/models"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"log"
	"time"
)

// AddCollection inserts a new collection into the collections table.
//
// Returns an error if the insertion fails.
func AddCollection(c *models.Collection) error {
	query := `INSERT INTO collections (user_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, c.UserId, c.Name, c.Description).Scan(&c.Id, &c.CreatedAt, &c.UpdatedAt)
}

// GetCollection retrieves a collection by its ID.
//
// Returns the collection and an error if retrieval fails.
func GetCollection(collectionId int) (*models.Collection, error) {
	query := `SELECT * FROM collections WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var c models.Collection
	err := db.QueryRowContext(ctx, query, collectionId).Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCollections retrieves collections for a user with optional filtering and pagination.
//
// Returns a slice of collections, metadata, and an error if retrieval fails.
func GetCollections(userId int, name string, f filters.Filters) ([]*models.Collection, filters.Metadata, error) {
	query := fmt.Sprintf(
		`	
			SELECT count(*) OVER(), * 
			FROM collections 
			WHERE user_id = $1 
			  AND (LOWER(name) = LOWER($2) OR $2 = '')
			ORDER BY %s %s, id
			LIMIT $3 OFFSET $4
			`,
		f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userId, name, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	totalRecords := 0

	var collections []*models.Collection
	for rows.Next() {
		var c models.Collection
		err = rows.Scan(&totalRecords, &c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, filters.Metadata{}, err
		}

		collections = append(collections, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metadata := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)

	return collections, metadata, nil
}

// UpdateCollection updates an existing collection's details.
//
// Returns an error if the update fails.
func UpdateCollection(c *models.Collection) error {
	query := `
			UPDATE collections 
			SET name = $3, description = $4, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND updated_at = $2
        	RETURNING user_id, created_at, updated_at
        	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.QueryRowContext(ctx, query, c.Id, c.UpdatedAt, c.Name, c.Description).Scan(&c.UserId, &c.CreatedAt, &c.UpdatedAt)
}

// DeleteCollection removes a collection by its ID.
//
// Returns an error if the deletion fails.
func DeleteCollection(id int) error {
	query := `DELETE FROM collections WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, id)
	return err
}
