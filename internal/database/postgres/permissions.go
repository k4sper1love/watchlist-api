package postgres

import (
	"context"
	"github.com/lib/pq"
	"log"
	"time"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for _, v := range p {
		if code == v {
			return true
		}
	}
	return false
}

func AddPermission(code string) error {
	query := `
			INSERT INTO permissions (code)
			VALUES ($1)
			ON CONFLICT (code) DO NOTHING 
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, code)
	return err
}

func AddUserPermissions(userId int, codes ...string) error {
	query := `
			INSERT INTO user_permissions (user_id, permissions_id)
			SELECT $1, permissions.id
			FROM permissions
			WHERE permissions.code = ANY($2)
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, userId, pq.Array(codes))
	return err
}

func GetUserPermissions(userId int) (Permissions, error) {
	query := `
			SELECT permissions.code 
			FROM permissions
			JOIN user_permissions ON user_permissions.permissions_id = permissions.id
			WHERE user_permissions.user_id = $1
			`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	var permissions Permissions
	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
