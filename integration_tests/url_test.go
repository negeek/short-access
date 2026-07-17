package integration_tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestShortenAndRedirect(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "a@example.com", "secret123")
	key := createKey(t, client, srv, token)
	_, shortURL := shorten(t, client, srv, key, "https://example.com/page")

	// Following the slug redirects to the original url.
	resp, _ := request(t, client, srv, http.MethodGet, "/"+shortURL, nil, nil)
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("redirect: got status %d", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != "https://example.com/page" {
		t.Fatalf("redirect: got location %q", loc)
	}

	// The visit is counted.
	resp, env := request(t, client, srv, http.MethodGet, "/api/v1/url_mgt/",
		map[string]string{"X-API-Key": key}, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: got status %d", resp.StatusCode)
	}
	var page struct {
		Items []struct {
			AccessCount int `json:"access_count"`
		} `json:"items"`
		Count int `json:"count"`
	}
	json.Unmarshal(env.Data, &page)
	if page.Count != 1 || len(page.Items) != 1 || page.Items[0].AccessCount != 1 {
		t.Fatalf("expected one url with access_count 1, got %+v", page)
	}
}

func TestListPagination(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "pager@example.com", "secret123")
	key := createKey(t, client, srv, token)

	// Create three urls.
	for _, u := range []string{"https://a.example.com", "https://b.example.com", "https://c.example.com"} {
		shorten(t, client, srv, key, u)
	}

	// First page of two: two items and has_more true.
	_, env := request(t, client, srv, http.MethodGet, "/api/v1/url_mgt/?limit=2",
		map[string]string{"X-API-Key": key}, nil)
	var first struct {
		Items   []json.RawMessage `json:"items"`
		HasMore bool              `json:"has_more"`
	}
	json.Unmarshal(env.Data, &first)
	if len(first.Items) != 2 || !first.HasMore {
		t.Fatalf("first page: expected 2 items and has_more, got %d items has_more=%v", len(first.Items), first.HasMore)
	}

	// Second page: the remaining one, has_more false.
	_, env = request(t, client, srv, http.MethodGet, "/api/v1/url_mgt/?limit=2&offset=2",
		map[string]string{"X-API-Key": key}, nil)
	var second struct {
		Items   []json.RawMessage `json:"items"`
		HasMore bool              `json:"has_more"`
	}
	json.Unmarshal(env.Data, &second)
	if len(second.Items) != 1 || second.HasMore {
		t.Fatalf("second page: expected 1 item and no has_more, got %d items has_more=%v", len(second.Items), second.HasMore)
	}
}

func TestShortenWithExpiryAtCreation(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "expiry@example.com", "secret123")
	key := createKey(t, client, srv, token)

	// Both time_unit and time_value set: the url gets a real expiry.
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key},
		map[string]any{"original_url": "https://example.com/a", "time_unit": "d", "time_value": 1})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("shorten with expiry: got %d (%s)", resp.StatusCode, env.Message)
	}
	var created struct {
		ExpireAt time.Time `json:"expire_at"`
	}
	json.Unmarshal(env.Data, &created)
	if !created.ExpireAt.After(time.Now()) {
		t.Fatalf("expected a future expiry, got %v", created.ExpireAt)
	}

	// Only one of the two set: rejected.
	resp, _ = request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key},
		map[string]any{"original_url": "https://example.com/b", "time_unit": "d"})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("shorten with only time_unit: expected 400, got %d", resp.StatusCode)
	}
}

func TestReshortenRevivesExpiredUrl(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "revive@example.com", "secret123")
	key := createKey(t, client, srv, token)
	target := "https://example.com/revive"

	// Shorten with a 1 second expiry.
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key},
		map[string]any{"original_url": target, "time_unit": "s", "time_value": 1})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("shorten: %d (%s)", resp.StatusCode, env.Message)
	}
	var first struct {
		ShortUrl string `json:"short_url"`
	}
	json.Unmarshal(env.Data, &first)

	// Let it expire, then confirm the link is dead.
	time.Sleep(1500 * time.Millisecond)
	resp, _ = request(t, client, srv, http.MethodGet, "/"+first.ShortUrl, nil, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected expired link to be rejected, got %d", resp.StatusCode)
	}

	// Re-shorten the same url with no expiry: it comes back on the same slug.
	resp, env = request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key},
		map[string]string{"original_url": target})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("re-shorten: %d (%s)", resp.StatusCode, env.Message)
	}
	var second struct {
		ShortUrl string `json:"short_url"`
	}
	json.Unmarshal(env.Data, &second)
	if second.ShortUrl != first.ShortUrl {
		t.Fatalf("expected the same slug, got %q vs %q", second.ShortUrl, first.ShortUrl)
	}

	// The revived link redirects again.
	resp, _ = request(t, client, srv, http.MethodGet, "/"+second.ShortUrl, nil, nil)
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected revived link to redirect, got %d", resp.StatusCode)
	}
}

func TestUrlEndpointAcceptsBearerToken(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	token := signup(t, client, srv, "bearer@example.com", "secret123")

	// The url endpoints accept a bearer token, not only an API key.
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"Authorization": "Bearer " + token},
		map[string]string{"original_url": "https://example.com/bearer"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("shorten with bearer token: got %d (%s)", resp.StatusCode, env.Message)
	}
}

func TestCannotTouchAnotherUsersUrl(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	ownerToken := signup(t, client, srv, "owner@example.com", "secret123")
	ownerKey := createKey(t, client, srv, ownerToken)
	id, _ := shorten(t, client, srv, ownerKey, "https://owner.example.com")

	otherToken := signup(t, client, srv, "other@example.com", "secret123")
	otherKey := createKey(t, client, srv, otherToken)

	// A different user cannot delete the owner's url.
	resp, _ := request(t, client, srv, http.MethodDelete, "/api/v1/url_mgt/"+strconv.Itoa(id),
		map[string]string{"X-API-Key": otherKey}, nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("cross-user delete: expected 400, got %d", resp.StatusCode)
	}

	// The owner still can.
	resp, _ = request(t, client, srv, http.MethodDelete, "/api/v1/url_mgt/"+strconv.Itoa(id),
		map[string]string{"X-API-Key": ownerKey}, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("owner delete: expected 204, got %d", resp.StatusCode)
	}
}
