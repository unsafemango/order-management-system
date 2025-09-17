package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

// Open creates and returns a new pgx connection pool using the given DSN
//
// It parses the DSN into a pgxpool.Config
// Initializes a pool
// Verifies the connection by pinging the database with a 5-second timer
func Open(dsn string) (*pgxpool.Pool, error) {
	// Parse the connection string into a pgxpool.Config
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Create a new connection pool using the parsed configuration
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timout for the initial ping check
	// This ensures the database is reachable before returning the pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify the connection by pinging the database within the timeout
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

// NewStore creates and returns a new Store instance that uses the provided pgx connection pool for database access
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{db: pool}
}
