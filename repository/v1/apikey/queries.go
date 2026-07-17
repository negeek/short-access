package apikey

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Repository runs api-key queries against the database pool it is given.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create stores a new key and fills in the generated id and timestamps.
func (repo *Repository) Create(ctx context.Context, a *ApiKey) error {
	query := "INSERT INTO api_keys (user_id, key_hash, name, expire_at) VALUES ($1, $2, $3, $4) RETURNING id, revoked, date_created, date_updated"
	return repo.db.QueryRow(ctx, query, a.UserId, a.KeyHash, a.Name, a.ExpireAt).
		Scan(&a.Id, &a.Revoked, &a.DateCreated, &a.DateUpdated)
}

// ListByUser returns a user's keys, newest first.
func (repo *Repository) ListByUser(ctx context.Context, userID uuid.UUID) ([]ApiKey, error) {
	query := "SELECT id, name, revoked, expire_at, date_created, date_updated FROM api_keys WHERE user_id = $1 ORDER BY date_created DESC"
	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []ApiKey
	for rows.Next() {
		var k ApiKey
		if err := rows.Scan(&k.Id, &k.Name, &k.Revoked, &k.ExpireAt, &k.DateCreated, &k.DateUpdated); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, rows.Err()
}

// FindActiveByHash looks up a usable key by its hash: it must exist, not be
// revoked, and not be past its expiry. The bool reports whether one was found.
func (repo *Repository) FindActiveByHash(ctx context.Context, hash string) (*ApiKey, bool, error) {
	query := "SELECT id, user_id, name, revoked, expire_at, date_created, date_updated FROM api_keys WHERE key_hash = $1 AND revoked = false AND (expire_at IS NULL OR expire_at > now())"
	var k ApiKey
	err := repo.db.QueryRow(ctx, query, hash).
		Scan(&k.Id, &k.UserId, &k.Name, &k.Revoked, &k.ExpireAt, &k.DateCreated, &k.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &k, true, nil
}

// Revoke marks a key as revoked. The bool reports whether the user owned a key
// with that id.
func (repo *Repository) Revoke(ctx context.Context, userID uuid.UUID, id int) (bool, error) {
	query := "UPDATE api_keys SET revoked = true, date_updated = now() WHERE id = $1 AND user_id = $2"
	tag, err := repo.db.Exec(ctx, query, id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// Delete removes a key. The bool reports whether the user owned a key with that id.
func (repo *Repository) Delete(ctx context.Context, userID uuid.UUID, id int) (bool, error) {
	query := "DELETE FROM api_keys WHERE id = $1 AND user_id = $2"
	tag, err := repo.db.Exec(ctx, query, id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
