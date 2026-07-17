// Package user holds the business logic for signing up and issuing tokens.
package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/negeek/short-access/apperr"
	userrepo "github.com/negeek/short-access/repository/v1/user"
	"github.com/negeek/short-access/utils"
)

// User is re-exported so handlers can talk in service types without importing
// the repository package directly.
type User = userrepo.User

// Service coordinates the user repository and token creation.
type Service struct {
	users *userrepo.Repository
}

func NewService(users *userrepo.Repository) *Service {
	return &Service{users: users}
}

// SignUp registers a new user and returns a fresh auth token. It fails if the
// email is already taken.
func (s *Service) SignUp(ctx context.Context, in *User) (string, error) {
	in.Id = uuid.New()

	exists, err := s.users.EmailExists(ctx, in.Email)
	if err != nil {
		return "", apperr.Internal(err)
	}
	if exists {
		return "", apperr.BadRequest("Email already exist")
	}

	if err := s.users.Create(ctx, in); err != nil {
		return "", apperr.Internal(err)
	}

	token, err := utils.CreateJwtToken(in.Id, in.Email)
	if err != nil {
		return "", apperr.Internal(err)
	}
	return token, nil
}

// NewToken checks the given credentials and returns a fresh auth token.
func (s *Service) NewToken(ctx context.Context, in *User) (string, error) {
	ok, err := s.users.Authenticate(ctx, in)
	if err != nil {
		return "", apperr.Internal(err)
	}
	if !ok {
		return "", apperr.BadRequest("Something Went Wrong. Check your email and password.")
	}

	token, err := utils.CreateJwtToken(in.Id, in.Email)
	if err != nil {
		return "", apperr.Internal(err)
	}
	return token, nil
}

// Exists reports whether a user with this email is registered. The auth
// middleware uses it to confirm a token still points at a real user.
func (s *Service) Exists(ctx context.Context, email string) (bool, error) {
	exists, err := s.users.EmailExists(ctx, email)
	if err != nil {
		return false, apperr.Internal(err)
	}
	return exists, nil
}
