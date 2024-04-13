package permissions

import (
	"context"
	"time"

	"github.com/lib/pq"
	"go-rest-starter.jtbergman.me/internal/models/core"
	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// ===========================================================================
// Interface
// ===========================================================================

type PermissionsRepository interface {
	GetByID(userID int64) (Perms, *xerrors.AppError)
	Insert(userID int64, codes ...string) (int64, *xerrors.AppError)
}

func Repository(db core.Queryable) PermissionsRepository {
	return &Permissions{DB: db}
}

// ===========================================================================
// Implementation
// ===========================================================================

// The Permission DAL
type Permissions struct {
	DB core.Queryable
}

// Gets permissions for the given user
func (m Permissions) GetByID(userID int64) (Perms, *xerrors.AppError) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN user_permissions ON user_permissions.permission_id = permissions.id
		INNER JOIN users ON user_permissions.user_id = users.id
		WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, xerrors.DatabaseError(err, "models.GetAllforUser.QueryContext")
	}
	defer rows.Close()

	all := Perms{}

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, xerrors.DatabaseError(err, "models.GetAllforUser.Scan")
		}
		all = append(all, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, xerrors.DatabaseError(err, "models.GetAllforUser.Err")
	}

	return all, nil
}

// Adds a variadic number of permissions for a user
func (m Permissions) Insert(userID int64, codes ...string) (int64, *xerrors.AppError) {
	query := `
		INSERT INTO user_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	if err != nil {
		return 0, xerrors.DatabaseError(err, "permissions.AddForUser")
	}

	return core.RowsAffected(result, "permissions.AddForUser")
}
