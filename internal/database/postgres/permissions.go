package postgres

import "github.com/lib/pq"

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

	_, err := db.Exec(query, code)
	return err
}

func AddUserPermissions(userId int, codes ...string) error {
	query := `
			INSERT INTO user_permissions (user_id, permissions_id)
			SELECT $1, permissions.id
			FROM permissions
			WHERE permissions.code = ANY($2)
			`

	_, err := db.Exec(query, userId, pq.Array(codes))
	return err
}

func GetUserPermissions(userId int) (Permissions, error) {
	query := `
			SELECT permissions.code 
			FROM permissions
			JOIN user_permissions ON user_permissions.permissions_id = permissions.id
			WHERE user_permissions.user_id = $1
			`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
