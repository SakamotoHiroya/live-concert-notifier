// Package repository provides PostgreSQL-backed persistence for the domain
// model, built on top of the sqlc-generated queries in ./sqlcgen.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect opens a connection pool to the database identified by databaseURL
// (e.g. "postgres://user:pass@host:5432/dbname?sslmode=disable").
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("repository: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("repository: ping: %w", err)
	}
	return pool, nil
}
