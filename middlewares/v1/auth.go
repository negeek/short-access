package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/negeek/short-access/apperr"
	apikeysvc "github.com/negeek/short-access/service/v1/apikey"
	usersvc "github.com/negeek/short-access/service/v1/user"
	"github.com/negeek/short-access/utils"
)

// contextKey is unexported so no other package can collide with our context values.
type contextKey string

const userContextKey contextKey = "user"

// WithUser stores the authenticated user id on the context.
func WithUser(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userContextKey, id)
}

// UserID reads the authenticated user id back out of the context.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userContextKey).(uuid.UUID)
	return id, ok
}

// Authenticator guards routes. It supports two schemes: a JWT for a signed-in
// user managing their account, and an API key for an application calling the
// url API. The real work lives in the user and api-key services; the middleware
// only pulls the credential out of the right header.
type Authenticator struct {
	users *usersvc.Service
	keys  *apikeysvc.Service
}

func NewAuthenticator(users *usersvc.Service, keys *apikeysvc.Service) *Authenticator {
	return &Authenticator{users: users, keys: keys}
}

// JWT accepts only a bearer token. Used for account and key management.
func (a *Authenticator) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		userID, err := a.users.Authenticate(r.Context(), token)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), userID)))
	})
}

// Either accepts whichever credential the caller sends: an X-API-Key header
// (an application) or a bearer token (a signed-in user). The API key is checked
// first when both are present.
func (a *Authenticator) Either(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID uuid.UUID
		var err error

		switch {
		case r.Header.Get("X-API-Key") != "":
			userID, err = a.keys.Authenticate(r.Context(), r.Header.Get("X-API-Key"))
		case r.Header.Get("Authorization") != "":
			var token string
			if token, err = bearerToken(r); err == nil {
				userID, err = a.users.Authenticate(r.Context(), token)
			}
		default:
			err = apperr.Unauthorized("Provide an API key or bearer token")
		}

		if err != nil {
			utils.RespondError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), userID)))
	})
}

// bearerToken pulls the token out of an "Authorization: Bearer <token>" header.
// A missing header yields an empty token so the service can report it uniformly;
// a malformed header is rejected here.
func bearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", nil
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", apperr.Unauthorized("Invalid Authorisation Header")
	}
	return parts[1], nil
}
