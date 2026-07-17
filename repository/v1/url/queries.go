package url

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/utils"
)

// Repository runs url queries against the database pool it is given.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// scanTargets lists the columns FindBy* and UserURLs read, in order.
func scanTargets(u *Url) []interface{} {
	return []interface{}{
		&u.Id, &u.OriginalUrl, &u.ShortUrl, &u.ShortAccess,
		&u.IsCustom, &u.AccessCount, &u.ExpireAt, &u.DateCreated, &u.DateUpdated,
	}
}

func (repo *Repository) Create(ctx context.Context, u *Url) error {
	query, values, err := utils.CRUDQueryBuild(u, u.TableName(), "create")
	if err != nil {
		return err
	}
	query += " RETURNING id"
	return repo.db.QueryRow(ctx, query, values...).Scan(&u.Id)
}

func (repo *Repository) Update(ctx context.Context, u *Url) error {
	query, values, err := utils.CRUDQueryBuild(u, u.TableName(), "update")
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(ctx, query, values...)
	return err
}

func (repo *Repository) Delete(ctx context.Context, u *Url) error {
	query, values, err := utils.CRUDQueryBuild(u, u.TableName(), "delete")
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(ctx, query, values...)
	return err
}

// FindByID loads a url by its id. The bool reports whether a row was found.
func (repo *Repository) FindByID(ctx context.Context, u *Url) (bool, error) {
	query, values, err := utils.CRUDQueryBuild(u, u.TableName(), "retrieve")
	if err != nil {
		return false, err
	}
	return repo.queryOne(ctx, u, query, values...)
}

// FindByIDForUser loads a url only if it belongs to the given user. This is how
// we stop one user from touching another user's urls.
func (repo *Repository) FindByIDForUser(ctx context.Context, u *Url, userID uuid.UUID) (bool, error) {
	query := "SELECT id,original_url,short_url,short_access,is_custom,access_count,expire_at,date_created,date_updated FROM urls WHERE id=$1 AND user_id=$2"
	return repo.queryOne(ctx, u, query, u.Id, userID)
}

func (repo *Repository) FindByOriginalURL(ctx context.Context, u *Url) (bool, error) {
	query := "SELECT id,original_url,short_url,short_access,is_custom,access_count,expire_at,date_created,date_updated FROM urls WHERE original_url=$1 and user_id=$2"
	return repo.queryOne(ctx, u, query, u.OriginalUrl, u.UserId)
}

func (repo *Repository) FindByShortURL(ctx context.Context, u *Url) (bool, error) {
	query := "SELECT id,original_url,short_url,short_access,is_custom,access_count,expire_at,date_created,date_updated FROM urls WHERE short_url=$1"
	return repo.queryOne(ctx, u, query, u.ShortUrl)
}

// queryOne runs a single-row select and scans it into u. A missing row is not
// treated as an error: it returns (false, nil).
func (repo *Repository) queryOne(ctx context.Context, u *Url, query string, args ...interface{}) (bool, error) {
	err := repo.db.QueryRow(ctx, query, args...).Scan(scanTargets(u)...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UserURLs runs a caller-built filter query and returns the matching urls.
func (repo *Repository) UserURLs(ctx context.Context, query string, values []interface{}) ([]Url, error) {
	rows, err := repo.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []Url
	for rows.Next() {
		var u Url
		if err := rows.Scan(scanTargets(&u)...); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, rows.Err()
}
