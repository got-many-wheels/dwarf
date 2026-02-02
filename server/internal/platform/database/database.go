package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Initiate database connection to MongoDB and returns Database instance
func Connect(uri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v\n", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping psql connection: %v\n", err)
	}

	return pool, nil
}
