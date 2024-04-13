package core

import (
	"database/sql"

	"go-rest-starter.jtbergman.me/internal/xerrors"
)

// Call RowsAffect on the result of an ExecContext to return the rows affect or an error
func RowsAffected(result sql.Result, op string) (int64, *xerrors.AppError) {
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return 0, xerrors.DatabaseError(err, op)
	}

	return rowsAffected, nil
}
