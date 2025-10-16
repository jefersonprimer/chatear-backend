package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Adapter represents the PostgreSQL database adapter
type Adapter struct {
	Pool *pgxpool.Pool
}

// NewAdapter creates a new PostgreSQL database adapter
func NewAdapter(databaseURL string) (*Adapter, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Adapter{
		Pool: pool,
	}, nil
}

// Close closes the database connection pool
func (a *Adapter) Close() {
	a.Pool.Close()
}

// RunMigrations runs database migrations from the specified directory
func (a *Adapter) RunMigrations(migrationsPath string) error {
	// For now, skip migrations in the worker context
	// This should be handled by the main application
	log.Println("Migrations skipped in worker context")
	return nil
}

// Health checks the database connection health
func (a *Adapter) Health(ctx context.Context) error {
	return a.Pool.Ping(ctx)
}
