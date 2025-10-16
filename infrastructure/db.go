package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jefersonprimer/chatear-backend/infrastructure/db/postgres"
)

// DB holds the database connection pool.
// Deprecated: Use postgres.Adapter instead
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB creates a new database connection.
// Deprecated: Use postgres.NewAdapter instead
func NewDB(databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection.
func (db *DB) Close() {
	db.Pool.Close()
}

// RunMigrations runs the database migrations.
// Deprecated: Use postgres.Adapter.RunMigrations instead
func (db *DB) RunMigrations(migrationsPath string) error {
	// For now, we'll just log that we're running migrations.
	// In a real application, you would use a library like golang-migrate/migrate.
	fmt.Printf("Running migrations from %s\n", migrationsPath)
	return nil
}

// NewPostgresAdapter creates a new PostgreSQL database adapter
func NewPostgresAdapter(databaseURL string) (*postgres.Adapter, error) {
	return postgres.NewAdapter(databaseURL)
}

// NewRedisClient creates a new Redis client
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}
	client := redis.NewClient(opt)
	// Ping to check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return client, nil
}

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	Postgres *postgres.Adapter
	Redis    *redis.Client
}

// NewInfrastructure creates and initializes all infrastructure components
func NewInfrastructure(databaseURL, redisURL string) (*Infrastructure, error) {
	pgAdapter, err := NewPostgresAdapter(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres adapter: %w", err)
	}

	redisClient, err := NewRedisClient(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	return &Infrastructure{
		Postgres: pgAdapter,
		Redis:    redisClient,
	}, nil
}

// Close closes all infrastructure connections
func (i *Infrastructure) Close() {
	if i.Postgres != nil {
		i.Postgres.Close()
	}
	if i.Redis != nil {
		i.Redis.Close()
	}
}
