package db

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// migration is one versioned pair of up/down SQL.
type migration struct {
	version string
	up      string
	down    string
}

const migrationsTable = `CREATE TABLE IF NOT EXISTS schema_migrations (
	version    text PRIMARY KEY,
	applied_at timestamptz NOT NULL DEFAULT now()
)`

// MigrateUp applies every migration that has not run yet, in version order.
func MigrateUp(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, migrationsTable); err != nil {
		return err
	}
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, m := range migrations {
		applied, err := isApplied(ctx, pool, m.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := runInTx(ctx, pool, m.up, func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", m.version)
			return err
		}); err != nil {
			return fmt.Errorf("migrate up %s: %w", m.version, err)
		}
	}
	return nil
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, migrationsTable); err != nil {
		return err
	}
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	var latest string
	err = pool.QueryRow(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&latest)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	for _, m := range migrations {
		if m.version != latest {
			continue
		}
		if err := runInTx(ctx, pool, m.down, func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", m.version)
			return err
		}); err != nil {
			return fmt.Errorf("migrate down %s: %w", m.version, err)
		}
		break
	}
	return nil
}

// loadMigrations reads the embedded .sql files and pairs each version's up and
// down halves together.
func loadMigrations() ([]migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	byVersion := map[string]*migration{}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		version, _, ok := strings.Cut(name, "_")
		if !ok {
			continue
		}
		content, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return nil, err
		}

		m := byVersion[version]
		if m == nil {
			m = &migration{version: version}
			byVersion[version] = m
		}
		if strings.HasSuffix(name, ".up.sql") {
			m.up = string(content)
		} else if strings.HasSuffix(name, ".down.sql") {
			m.down = string(content)
		}
	}

	migrations := make([]migration, 0, len(byVersion))
	for _, m := range byVersion {
		migrations = append(migrations, *m)
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	return migrations, nil
}

func isApplied(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	return exists, err
}

// runInTx runs the migration SQL and its bookkeeping in a single transaction, so
// a failure leaves nothing half-applied.
func runInTx(ctx context.Context, pool *pgxpool.Pool, sql string, record func(pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if strings.TrimSpace(sql) != "" {
		if _, err := tx.Exec(ctx, sql); err != nil {
			return err
		}
	}
	if err := record(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
