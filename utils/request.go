package utils

import (
	"encoding/json"
	"net/http"
)

// DecodeBody reads the JSON request body into a fresh value of type T.
// It saves every handler from repeating the same read-and-unmarshal steps.
func DecodeBody[T any](r *http.Request) (T, error) {
	var body T
	err := DecodeBodyInto(r, &body)
	return body, err
}

// DecodeBodyInto reads the JSON request body into dst. Use this when you want
// to merge the incoming fields onto a struct you already have, e.g. a record
// loaded from the database before a partial update.
func DecodeBodyInto(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
