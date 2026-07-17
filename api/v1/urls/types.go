package urls

// DateTimeExpiryDetail is the body of a set-expiry request: which url, and how
// far in the future it should expire.
type DateTimeExpiryDetail struct {
	TimeUnit  string `json:"time_unit"`
	TimeValue int    `json:"time_value"`
	UrlId     int    `json:"url_id"`
}
