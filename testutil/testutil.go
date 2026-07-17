// Package testutil sets up a real database for tests: it resets the schema and
// runs migrations once, then hands out a fast per-test table reset.
package testutil

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/db"
)

// Setup opens the test database named by TEST_DATABASE_URL, drops and recreates
// the public schema, and runs migrations. It returns the pool, a cleanup func
// and whether tests should run. When the variable is unset the caller should
// skip its tests instead of failing, so `go test` works without a database.
func Setup() (pool *pgxpool.Pool, cleanup func(), ok bool) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		return nil, nil, false
	}

	pool, err := db.Connect(dsn)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if _, err := pool.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		pool.Close()
		panic(err)
	}
	if err := db.MigrateUp(ctx, pool); err != nil {
		pool.Close()
		panic(err)
	}

	return pool, pool.Close, true
}

// Truncate empties every table (except the migrations bookkeeping) so each test
// starts from a clean slate without paying for a full migrate.
func Truncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	rows, err := pool.Query(ctx,
		"SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename <> 'schema_migrations'")
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			rows.Close()
			t.Fatalf("scan table name: %v", err)
		}
		tables = append(tables, name)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		t.Fatalf("list tables: %v", err)
	}
	if len(tables) == 0 {
		return
	}

	if _, err := pool.Exec(ctx, "TRUNCATE TABLE "+strings.Join(tables, ", ")+" RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}
