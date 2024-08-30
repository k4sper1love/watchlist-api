package postgres

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/internal/filters"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddCollection(c *models.Collection) error {
	query := `INSERT INTO collections (user_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at`

	return db.QueryRow(query, c.UserId, c.Name, c.Description).Scan(&c.Id, &c.CreatedAt)
}

func GetCollection(collectionId int) (*models.Collection, error) {
	query := `SELECT * FROM collections WHERE id = $1`

	var c models.Collection
	err := db.QueryRow(query, collectionId).Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

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

	rows, err := db.Query(query, userId, name, f.Limit(), f.Offset())
	if err != nil {
		return nil, filters.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0

	var collections []*models.Collection
	for rows.Next() {
		var c models.Collection
		err = rows.Scan(&totalRecords, &c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
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

func UpdateCollection(c *models.Collection) error {
	query := `
			UPDATE collections 
			SET name = $2, description = $3 
			WHERE id = $1 
        	RETURNING id, user_id, name, description, created_at
        	`

	return db.QueryRow(query, c.Id, c.Name, c.Description).Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
}

func DeleteCollection(id int) error {
	query := `DELETE FROM collections WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}
