package url

import (
	"time"
	"github.com/google/uuid"
	"github.com/negeek/short-access/repository/v1/user"
)

type Url struct {
	Id        	int   		`json:"id"`	
	OriginalUrl  		string      `json:"original_url"`
	ShortUrl    string      `json:"short_url"`
	IsCustom    bool      `json:"is_custom"`
	AccessCount int `json:"access_count"`
	UserId    	uuid.UUID	`json:"-"`
	ExpireAt    time.Time    `json:"expire_at"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
	User        user.User   `json:"-"`
}

func (u *Url)TableName()string{
	return "urls"
}