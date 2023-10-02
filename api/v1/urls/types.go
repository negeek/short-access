package urls

type NumberStore struct {
	Number int
	Step int
	End int
}

type DateTimeExpiryDetail struct {
	TimeUnit string `json:"time_unit"`
	TimeValue int `json:"time_value"`
	UrlId int `json:"url_id"`
}