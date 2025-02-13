package client

import (
	"context"
	"database/sql"
)

// Service represents a service that interacts with a database.
type DatabaseService interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	HealthCheck() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// ExecContext executes a query without returning any rows.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Transactions
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	GetDB() *sql.DB
}
