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

func (repo *Repository) Create(ctx context.Context, u *User) error {
	if err := utils.Time(u, true); err != nil {
		return err
	}
	query := "INSERT INTO users (id, password, email, date_created, date_updated) VALUES ($1, $2, $3, $4, $5)"
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

// FindByEmail loads a user by email. The bool reports whether a row was found.
func (repo *Repository) FindByEmail(ctx context.Context, u *User) (bool, error) {
	query := "SELECT id, email, date_created, date_updated FROM users WHERE email = $1"
	err := repo.db.QueryRow(ctx, query, u.Email).Scan(&u.Id, &u.Email, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Authenticate loads a user that matches the given email and password. The bool
// reports whether a matching user was found.
func (repo *Repository) Authenticate(ctx context.Context, u *User) (bool, error) {
	query := "SELECT id, email, date_created, date_updated FROM users WHERE email = $1 and password = $2"
	err := repo.db.QueryRow(ctx, query, u.Email, u.Password).Scan(&u.Id, &u.Email, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
