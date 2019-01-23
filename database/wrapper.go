package database

import (
	"context"
	"database/sql"
)

type QueryContextFunc func(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

type ExecContextFunc func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

// Wrapper defines database common operations
type Wrapper interface {
	WrapQueryContext(fn QueryContextFunc, sql string, args ...interface{}) QueryContextFunc
	WrapExecContext(fn ExecContextFunc, sql string, args ...interface{}) ExecContextFunc
}
