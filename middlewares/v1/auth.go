package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

// Authenticator guards routes. It offers two schemes: JWT for a signed-in user
// managing their account, and API keys for an application calling the url API.
type Authenticator struct {
	users *usersvc.Service
	keys  *apikeysvc.Service
}

func NewAuthenticator(users *usersvc.Service, keys *apikeysvc.Service) *Authenticator {
	return &Authenticator{users: users, keys: keys}
}

// JWT checks the request's bearer token and confirms the user it names still
// exists. Used for account and key management.
func (a *Authenticator) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.JsonResponse(w, false, http.StatusUnauthorized, "Provide Auth Token", nil)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JsonResponse(w, false, http.StatusUnauthorized, "Invalid Authorisation Header", nil)
			return
		}

		claim, err := utils.VerifyJwt(parts[1])
		if err != nil {
			utils.JsonResponse(w, false, http.StatusUnauthorized, "Invalid Token", nil)
			return
		}

		exists, err := a.users.Exists(r.Context(), claim.Email)
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		if !exists {
			utils.JsonResponse(w, false, http.StatusUnauthorized, "Invalid User", nil)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), claim.ID)))
	})
}

// APIKey checks the X-API-Key header and resolves it to its owner. Used for the
// url endpoints an application calls.
func (a *Authenticator) APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := a.keys.Authenticate(r.Context(), r.Header.Get("X-API-Key"))
		if err != nil {
			utils.RespondError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), userID)))
	})
}
