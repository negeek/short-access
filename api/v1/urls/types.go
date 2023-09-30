package urls

type NumberStore struct {
	Number int
	Step int
	End int
}

type DateTimeExpiryDetail struct {
	TimeUnit string `json:"time_unit"`
	TimeValue int `json:"time_value"`
	ShortUrl string `json:"short_url"`
}