package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect opens a connection pool to Postgres and returns it. Callers pass the
// pool into the repositories that need it, so there is no shared global state.
func Connect(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
