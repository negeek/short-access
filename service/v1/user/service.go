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

	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return "", apperr.Internal(err)
	}
	in.Password = hash

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
	plaintext := in.Password

	stored := &User{Email: in.Email}
	found, err := s.users.FindByEmail(ctx, stored)
	if err != nil {
		return "", apperr.Internal(err)
	}
	// Same message whether the email is unknown or the password is wrong, so we
	// don't reveal which emails are registered.
	if !found || !utils.CheckPassword(stored.Password, plaintext) {
		return "", apperr.BadRequest("Something Went Wrong. Check your email and password.")
	}

	token, err := utils.CreateJwtToken(stored.Id, stored.Email)
	if err != nil {
		return "", apperr.Internal(err)
	}
	return token, nil
}

// Authenticate verifies a JWT and returns the id of the user it names, after
// confirming that user still exists. It mirrors the api-key service's
// Authenticate so the middleware treats both credentials the same way.
func (s *Service) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, apperr.Unauthorized("Provide Auth Token")
	}

	claim, err := utils.VerifyJwt(token)
	if err != nil {
		return uuid.Nil, apperr.Unauthorized("Invalid Token")
	}

	exists, err := s.users.EmailExists(ctx, claim.Email)
	if err != nil {
		return uuid.Nil, apperr.Internal(err)
	}
	if !exists {
		return uuid.Nil, apperr.Unauthorized("Invalid User")
	}
	return claim.ID, nil
}
