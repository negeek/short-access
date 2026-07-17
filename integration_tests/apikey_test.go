package integration_tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestRevokedKeyIsRejected(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "b@example.com", "secret123")
	key := createKey(t, client, srv, token)

	// Find the key's id.
	_, env := request(t, client, srv, http.MethodGet, "/api/v1/user_mgt/api_keys/",
		map[string]string{"Authorization": "Bearer " + token}, nil)
	var keys []struct {
		Id int `json:"id"`
	}
	json.Unmarshal(env.Data, &keys)
	if len(keys) != 1 {
		t.Fatalf("expected one key, got %d", len(keys))
	}

	// Revoke it.
	resp, env := request(t, client, srv, http.MethodPost,
		"/api/v1/user_mgt/api_keys/"+strconv.Itoa(keys[0].Id)+"/revoke",
		map[string]string{"Authorization": "Bearer " + token}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("revoke: got %d (%s)", resp.StatusCode, env.Message)
	}

	// The revoked key can no longer shorten.
	resp, _ = request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key}, map[string]string{"original_url": "https://example.com"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("shorten with revoked key: expected 401, got %d", resp.StatusCode)
	}
}
