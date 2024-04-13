package core

import (
	"context"
	"database/sql"
)

// ============================================================================
// Queryable
// ============================================================================

// Services should depend on Queryable so they can on start transactions
type Queryable interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
