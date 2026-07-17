package utils

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/negeek/short-access/db"
)

// Setup loads env vars and opens a database pool for tests. The full test
// harness is reworked in a later step; this keeps the package building.
func Setup() (*pgxpool.Pool, error) {
	if os.Getenv("APP_ENV") == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			if err := godotenv.Load("../../internal/env/.env"); err != nil {
				return nil, err
			}
		}
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	return db.Connect(dbURL)
}
