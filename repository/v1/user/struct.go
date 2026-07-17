package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id          uuid.UUID `json:"-"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}
