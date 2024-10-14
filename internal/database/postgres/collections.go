package postgres

import (
	"context"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"log/slog"
	"time"
)

// AddCollection inserts a new collection into the collections table.
func AddCollection(c *models.Collection) error {
	query := `
		INSERT INTO collections (user_id, name, description) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, c.UserID, c.Name, c.Description).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// GetCollection retrieves a collection by its ID.
func GetCollection(collectionID int) (*models.Collection, error) {
	query := `SELECT * FROM collections WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var c models.Collection
	if err := GetDB().QueryRowContext(ctx, query, collectionID).Scan(&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCollections retrieves collections for a user with optional filtering and pagination.
func GetCollections(userID int, name string, f filters.Filters) ([]*models.Collection, filters.Metadata, error) {
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

	rows, err := GetDB().QueryContext(ctx, query, userID, name, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			sl.Log.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var collections []*models.Collection
	totalRecords := 0

	for rows.Next() {
		var c models.Collection
		if err := rows.Scan(&totalRecords, &c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
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
func UpdateCollection(c *models.Collection) error {
	query := `
		UPDATE collections 
		SET name = $3, description = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND updated_at = $2
		RETURNING user_id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return GetDB().QueryRowContext(ctx, query, c.ID, c.UpdatedAt, c.Name, c.Description).Scan(&c.UserID, &c.CreatedAt, &c.UpdatedAt)
}

// DeleteCollection removes a collection by its ID.
func DeleteCollection(id int) error {
	query := `DELETE FROM collections WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, id)
	return err
}
