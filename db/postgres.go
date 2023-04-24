package db

import (
    "context"
    "log"
    "os"

    "github.com/jackc/pgx/v4/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
    // Get the database URL from environment variable
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatalf("DATABASE_URL not set")
    }

    // Create a new database connection pool
    dbPool, err := pgxpool.Connect(context.Background(), dbURL)
    if err != nil {
        return nil, err
    }

    return dbPool, nil
}