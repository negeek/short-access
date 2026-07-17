// Package apikey holds the business logic for issuing and checking API keys.
package apikey

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/negeek/short-access/apperr"
	apikeyrepo "github.com/negeek/short-access/repository/v1/apikey"
	"github.com/negeek/short-access/utils"
)

// ApiKey is re-exported so handlers can talk in service types without importing
// the repository package directly.
type ApiKey = apikeyrepo.ApiKey

// Service coordinates the api-key repository and key generation.
type Service struct {
	keys *apikeyrepo.Repository
}

func NewService(keys *apikeyrepo.Repository) *Service {
	return &Service{keys: keys}
}

// Create issues a new key for the user. It returns the raw key, which is shown
// to the user once and never stored, alongside the saved record.
func (s *Service) Create(ctx context.Context, userID uuid.UUID, name string, expireAt *time.Time) (string, *ApiKey, error) {
	raw, err := utils.GenerateAPIKey()
	if err != nil {
		return "", nil, apperr.Internal(err)
	}

	record := &ApiKey{
		UserId:   userID,
		KeyHash:  utils.HashAPIKey(raw),
		Name:     name,
		ExpireAt: expireAt,
	}
	if err := s.keys.Create(ctx, record); err != nil {
		return "", nil, apperr.Internal(err)
	}
	return raw, record, nil
}

// List returns the user's keys.
func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]ApiKey, error) {
	keys, err := s.keys.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return keys, nil
}

// Revoke disables a key the user owns.
func (s *Service) Revoke(ctx context.Context, userID uuid.UUID, id int) error {
	found, err := s.keys.Revoke(ctx, userID, id)
	if err != nil {
		return apperr.Internal(err)
	}
	if !found {
		return apperr.NotFound("API key not found")
	}
	return nil
}

// Delete removes a key the user owns.
func (s *Service) Delete(ctx context.Context, userID uuid.UUID, id int) error {
	found, err := s.keys.Delete(ctx, userID, id)
	if err != nil {
		return apperr.Internal(err)
	}
	if !found {
		return apperr.NotFound("API key not found")
	}
	return nil
}

// Authenticate resolves a raw API key to the id of the user it belongs to. It is
// used by the middleware that guards the url endpoints.
func (s *Service) Authenticate(ctx context.Context, rawKey string) (uuid.UUID, error) {
	if rawKey == "" {
		return uuid.Nil, apperr.Unauthorized("Provide API key")
	}
	record, found, err := s.keys.FindActiveByHash(ctx, utils.HashAPIKey(rawKey))
	if err != nil {
		return uuid.Nil, apperr.Internal(err)
	}
	if !found {
		return uuid.Nil, apperr.Unauthorized("Invalid API key")
	}
	return record.UserId, nil
}
