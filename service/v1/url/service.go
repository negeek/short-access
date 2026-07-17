// Package url holds the business logic for shortening and managing urls.
package url

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/negeek/short-access/apperr"
	numberrepo "github.com/negeek/short-access/repository/v1/number"
	urlrepo "github.com/negeek/short-access/repository/v1/url"
	"github.com/negeek/short-access/utils"
)

// slugLength is how many characters every generated short url has.
const slugLength = 9

// Url is re-exported so handlers can talk in service types without importing
// the repository package directly.
type Url = urlrepo.Url

// counter hands out sequential numbers used to build short slugs. It reserves a
// block of numbers in the database and gives them out one at a time, only
// touching the database again once a block runs out.
type counter struct {
	mu     sync.Mutex
	number int
	step   int
	end    int
}

// Service coordinates the url and number repositories.
type Service struct {
	urls    *urlrepo.Repository
	numbers *numberrepo.Repository
	counter *counter
}

func NewService(urls *urlrepo.Repository, numbers *numberrepo.Repository) *Service {
	return &Service{
		urls:    urls,
		numbers: numbers,
		counter: &counter{step: 100, end: 100},
	}
}

// Shorten returns a short url for the given original url. If the same user has
// already shortened it, the existing record is returned instead of a new one.
func (s *Service) Shorten(ctx context.Context, userID uuid.UUID, in *Url) (*Url, error) {
	in.UserId = userID

	found, err := s.urls.FindByOriginalURL(ctx, in)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if found {
		return in, nil
	}

	number, err := s.nextNumber(ctx)
	if err != nil {
		return nil, apperr.Internal(err)
	}

	in.ShortUrl = utils.ShortAccess(number, slugLength)
	in.FillShortAccess()
	if err := s.urls.Create(ctx, in); err != nil {
		return nil, apperr.Internal(err)
	}
	return in, nil
}

// nextNumber returns the next counter value, reserving a new block from the
// database whenever the current block is used up.
func (s *Service) nextNumber(ctx context.Context) (int, error) {
	s.counter.mu.Lock()
	defer s.counter.mu.Unlock()

	row := &numberrepo.Number{Id: 1, Step: s.counter.step}

	// number == 0 means the server just (re)started and needs to learn where
	// the shared counter is before handing out values.
	if s.counter.number == 0 {
		found, err := s.numbers.FindByID(ctx, row)
		if err != nil {
			return 0, err
		}
		if !found {
			if err := s.numbers.CreateOrUpdate(ctx, row); err != nil {
				return 0, err
			}
			s.counter.number = 1
		} else {
			s.counter.number = row.Number + 1
			if err := s.numbers.CreateOrUpdate(ctx, row); err != nil {
				return 0, err
			}
		}
		s.counter.end = row.Number
		return s.counter.number, nil
	}

	// Current block used up: reserve the next one.
	if s.counter.number >= s.counter.end {
		if err := s.numbers.CreateOrUpdate(ctx, row); err != nil {
			return 0, err
		}
		s.counter.end = row.Number
	}
	s.counter.number++
	return s.counter.number, nil
}

// CreateCustom stores a url with a slug the user chose. It fails if that slug is
// already taken.
func (s *Service) CreateCustom(ctx context.Context, userID uuid.UUID, in *Url) (*Url, error) {
	in.UserId = userID

	found, err := s.urls.FindByShortURL(ctx, in)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if found {
		return nil, apperr.BadRequest("Url provided exists already")
	}

	in.FillShortAccess()
	in.IsCustom = true
	if err := s.urls.Create(ctx, in); err != nil {
		return nil, apperr.Internal(err)
	}
	return in, nil
}

// SetExpiry sets when a url should stop redirecting.
func (s *Service) SetExpiry(ctx context.Context, userID uuid.UUID, urlID int, unit string, value int) (*Url, error) {
	expireAt, err := utils.ExpiryDateTime(unit, value)
	if err != nil {
		return nil, apperr.BadRequest(err.Error())
	}

	target := &Url{Id: urlID, UserId: userID}
	found, err := s.urls.FindByID(ctx, target)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if !found {
		return nil, apperr.BadRequest("Url does not exist")
	}

	target.ExpireAt = expireAt
	if err := s.urls.Update(ctx, target); err != nil {
		return nil, apperr.Internal(err)
	}
	return target, nil
}

// GetByID loads a url by id, or reports that it does not exist.
func (s *Service) GetByID(ctx context.Context, id int) (*Url, error) {
	target := &Url{Id: id}
	found, err := s.urls.FindByID(ctx, target)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if !found {
		return nil, apperr.BadRequest("Url does not exist")
	}
	return target, nil
}

// Save writes an already-loaded url back to the database.
func (s *Service) Save(ctx context.Context, u *Url) (*Url, error) {
	if err := s.urls.Update(ctx, u); err != nil {
		return nil, apperr.Internal(err)
	}
	return u, nil
}

// Delete removes a url by id after checking it exists.
func (s *Service) Delete(ctx context.Context, id int) error {
	target := &Url{Id: id}
	found, err := s.urls.FindByID(ctx, target)
	if err != nil {
		return apperr.Internal(err)
	}
	if !found {
		return apperr.BadRequest("Url does not exist")
	}
	if err := s.urls.Delete(ctx, target); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

// List returns the caller's urls, filtered by the given query parameters.
func (s *Service) List(ctx context.Context, userID uuid.UUID, queryParams map[string][]string) ([]Url, error) {
	filter := Url{UserId: userID}
	query, values, err := utils.Filter(queryParams, filter, filter.TableName())
	if err != nil {
		return nil, apperr.BadRequest(err.Error())
	}
	urls, err := s.urls.UserURLs(ctx, query, values)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return urls, nil
}

// Redirect looks up a slug, makes sure it is still valid, counts the visit and
// returns the url to redirect to.
func (s *Service) Redirect(ctx context.Context, slug string) (*Url, error) {
	target := &Url{ShortUrl: slug}
	found, err := s.urls.FindByShortURL(ctx, target)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	if !found {
		return nil, apperr.BadRequest("Something went wrong. Make sure url is valid.")
	}
	if target.Expired() {
		return nil, apperr.BadRequest("Url has expired")
	}

	target.AccessCount++
	if err := s.urls.Update(ctx, target); err != nil {
		return nil, apperr.Internal(err)
	}
	return target, nil
}
