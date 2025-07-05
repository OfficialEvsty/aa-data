package db

import (
	"context"
	"database/sql"
)

// ISqlExecutor interface for tx and db implementation
type ISqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// ITxExecutor interface for providing tx implementation
type ITxExecutor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}
