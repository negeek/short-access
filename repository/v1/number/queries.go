package number

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/utils"
)

// Repository runs number-counter queries against the database pool it is given.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateOrUpdate bumps the shared counter row by one step and returns the new
// value in n.Number. The first call inserts the row; later calls add to it.
func (repo *Repository) CreateOrUpdate(ctx context.Context, n *Number) error {
	if err := utils.Time(n, false); err != nil {
		return err
	}
	query := "INSERT INTO numbers (id, number, date_updated) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET number = numbers.number + $2 RETURNING number"
	return repo.db.QueryRow(ctx, query, n.Id, n.Step, n.DateUpdated).Scan(&n.Number)
}

// FindByID loads the counter row. The bool reports whether a row was found.
func (repo *Repository) FindByID(ctx context.Context, n *Number) (bool, error) {
	query := "SELECT number FROM numbers WHERE id = $1"
	err := repo.db.QueryRow(ctx, query, n.Id).Scan(&n.Number)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
