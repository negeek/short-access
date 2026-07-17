package apikey

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id          int        `json:"id"`
	UserId      uuid.UUID  `json:"-"`
	KeyHash     string     `json:"-"`
	Name        string     `json:"name"`
	Revoked     bool       `json:"revoked"`
	ExpireAt    *time.Time `json:"expire_at"` // nil means the key never expires
	DateCreated time.Time  `json:"date_created"`
	DateUpdated time.Time  `json:"date_updated"`
}

func (a *ApiKey) TableName() string {
	return "api_keys"
}
