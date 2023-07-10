package number

import (
	"time"
)

type Number struct {
	Id        	int   		`json:"-"`	
	Number     	int     	`json:"original_url"`
	Step       	int   		`json:"step"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
}
