package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

// A helper to check if a given permission if included in the Permissions slice
func (p Permissions) Includes(code string) bool {
	for i := range p {
		if p[i] == code {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	qurey := `SELECT P.code FROM users_premissions UP
						INNER JOIN permissions P
						ON UP.permission_id = P.id
						WHERE UP.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var permissions Permissions

	rows, err := m.DB.QueryContext(ctx, qurey, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var role string
		err = rows.Scan(&role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
