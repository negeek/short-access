package url

import (
	"time"
	"github.com/google/uuid"
	"github.com/negeek/short-access/repository/v1/user"
)

type Url struct {
	Id        	int   		`json:"-"`	
	Url  		string      `json:"url"`
	ShortUrl    string      `json:"short_url"`
	UserId    	uuid.UUID	`json:"-"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
	User        user.User   `json:"user"`
}
