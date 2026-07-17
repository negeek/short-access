package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

// Authenticator guards routes by checking the request's JWT and confirming the
// user it names still exists.
type Authenticator struct {
	users *usersvc.Service
}

func NewAuthenticator(users *usersvc.Service) *Authenticator {
	return &Authenticator{users: users}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
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
