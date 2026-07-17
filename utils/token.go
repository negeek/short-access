package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// defaultTokenTTLHours is how long an auth token stays valid when TOKEN_TTL_HOURS
// is not set.
const defaultTokenTTLHours = 24

type UserClaim struct {
	jwt.RegisteredClaims
	ID    uuid.UUID
	Email string
}

func CreateJwtToken(id uuid.UUID, email string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL())),
		},
		ID:    id,
		Email: email,
	})

	return token.SignedString([]byte(os.Getenv("AUTH_KEY")))
}

func VerifyJwt(jwtToken string) (*UserClaim, error) {
	// ParseWithClaims checks the signature and, because the claims carry an
	// expiry, rejects tokens that have expired.
	token, err := jwt.ParseWithClaims(jwtToken, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*UserClaim), nil
}

// tokenTTL reads how long tokens should last from TOKEN_TTL_HOURS, falling back
// to the default when the value is missing or invalid.
func tokenTTL() time.Duration {
	hours := defaultTokenTTLHours
	if v := os.Getenv("TOKEN_TTL_HOURS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			hours = parsed
		}
	}
	return time.Duration(hours) * time.Hour
}
