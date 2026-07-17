package user

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/utils"
)

// Repository runs user queries against the database pool it is given.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create stores a new user. u.Password is expected to already be hashed.
func (repo *Repository) Create(ctx context.Context, u *User) error {
	if err := utils.Time(u, true); err != nil {
		return err
	}
	query := "INSERT INTO users (id, password_hash, email, date_created, date_updated) VALUES ($1, $2, $3, $4, $5)"
	_, err := repo.db.Exec(ctx, query, u.Id, u.Password, u.Email, u.DateCreated, u.DateUpdated)
	return err
}

// EmailExists reports whether a user with this email is already registered.
func (repo *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := repo.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

// FindByEmail loads a user by email, including the stored password hash so the
// service can verify a login. The bool reports whether a row was found.
func (repo *Repository) FindByEmail(ctx context.Context, u *User) (bool, error) {
	query := "SELECT id, email, password_hash, date_created, date_updated FROM users WHERE email = $1"
	err := repo.db.QueryRow(ctx, query, u.Email).Scan(&u.Id, &u.Email, &u.Password, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
