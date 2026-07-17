package urls

// ShortenRequest is the body of a shorten or custom request. Expiry is optional:
// leave time_unit empty for a link that never expires. short_url is only used by
// the custom endpoint.
type ShortenRequest struct {
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
	TimeUnit    string `json:"time_unit"`
	TimeValue   int    `json:"time_value"`
}

// DateTimeExpiryDetail is the body of a set-expiry request: which url, and how
// far in the future it should expire.
type DateTimeExpiryDetail struct {
	TimeUnit  string `json:"time_unit"`
	TimeValue int    `json:"time_value"`
	UrlId     int    `json:"url_id"`
}
