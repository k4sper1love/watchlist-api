package postgres

import (
	"context"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/lib/pq"
	"log/slog"
	"time"
)

// Permissions represents a slice of permission codes.
type Permissions []string

// Include checks if a given permission code is present in the Permissions slice.
func (p Permissions) Include(code string) bool {
	for _, v := range p {
		if code == v {
			return true
		}
	}
	return false
}

// AddPermission inserts a new permission into the permissions table.
func AddPermission(code string) error {
	query := `
		INSERT INTO permissions (code)
		VALUES ($1)
		ON CONFLICT (code) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, code)
	return err
}

// AddUserPermissions adds multiple permissions for a specific user.
func AddUserPermissions(userId int, codes ...string) error {
	query := `
		INSERT INTO user_permissions (user_id, permissions_id)
		SELECT $1, permissions.id
		FROM permissions
		WHERE permissions.code = ANY($2)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, userId, pq.Array(codes))
	return err
}

// GetUserPermissions retrieves all permission codes for a specific user.
func GetUserPermissions(userId int) (Permissions, error) {
	query := `
		SELECT permissions.code 
		FROM permissions
		JOIN user_permissions ON user_permissions.permissions_id = permissions.id
		WHERE user_permissions.user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := GetDB().QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			sl.Log.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var permissions Permissions
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// DeletePermissions deletes permission codes.
func DeletePermissions(codes ...string) error {
	query := `DELETE FROM permissions WHERE code = ANY($1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetDB().ExecContext(ctx, query, pq.Array(codes))
	return err
}
