package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-api/internal/models"
)

func AddCollection(c *models.Collection) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `INSERT INTO collections (user_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at`

	return db.QueryRow(query, c.UserId, c.Name, c.Description).Scan(&c.Id, &c.CreatedAt)
}

func GetCollection(id int) (*models.Collection, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM collections WHERE id = $1`

	var c models.Collection
	err := db.QueryRow(query, id).Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func GetCollections(userId int) ([]*models.Collection, error) {
	db := connectPostgres()
	if db == nil {
		return nil, errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `SELECT * FROM collections WHERE user_id = $1`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*models.Collection
	for rows.Next() {
		var c models.Collection
		err = rows.Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		collections = append(collections, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}

func UpdateCollection(c *models.Collection) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `
			UPDATE collections 
			SET name = $2, description = $3 
			WHERE id = $1 
        	RETURNING id, user_id, name, description, created_at
        	`

	return db.QueryRow(query, c.Id, c.Name, c.Description).Scan(&c.Id, &c.UserId, &c.Name, &c.Description, &c.CreatedAt)
}

func DeleteCollection(id int) error {
	db := connectPostgres()
	if db == nil {
		return errors.New("cannot connect to PostgreSQL")
	}
	defer db.Close()

	query := `DELETE FROM collections WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}
