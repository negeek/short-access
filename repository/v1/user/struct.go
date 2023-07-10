package user
import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	Id        	uuid.UUID   `json:"-"`
	Password  	string      `json:"password"`
	Email     	string      `json:"email"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
}