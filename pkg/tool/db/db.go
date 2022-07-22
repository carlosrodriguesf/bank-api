//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type (
	Connection interface {
		PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
		NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
		NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
		QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
		QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
		NamedGetContext(ctx context.Context, query string, dest, arg interface{}) error
		NamedSelectContext(ctx context.Context, query string, dest, arg interface{}) error
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}
	ExtendedDB interface {
		Connection
		Close() error
	}
	ExtendedTx interface {
		Connection
		Commit() error
		Rollback() error
	}
)
