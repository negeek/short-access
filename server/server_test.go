package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/negeek/short-access/server"
	"github.com/negeek/short-access/testutil"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	os.Setenv("AUTH_KEY", "test-secret-key")
	os.Setenv("BASE_URL", "http://example.test")

	pool, cleanup, ok := testutil.Setup()
	if !ok {
		// No test database configured, so skip the whole package.
		os.Exit(0)
	}
	testPool = pool
	code := m.Run()
	cleanup()
	os.Exit(code)
}

// envelope matches the shape every endpoint responds with.
type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	testutil.Truncate(t, testPool)
	return httptest.NewServer(server.NewRouter(testPool))
}

// noRedirectClient returns an HTTP client that surfaces redirects instead of
// following them, so tests can assert on the redirect itself.
func noRedirectClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

func request(t *testing.T, client *http.Client, srv *httptest.Server, method, path string, headers map[string]string, body any) (*http.Response, envelope) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, srv.URL+path, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var env envelope
	raw, _ := io.ReadAll(resp.Body)
	if len(bytes.TrimSpace(raw)) > 0 {
		_ = json.Unmarshal(raw, &env)
	}
	return resp, env
}

func signup(t *testing.T, client *http.Client, srv *httptest.Server, email, password string) string {
	t.Helper()
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/user_mgt/join/", nil,
		map[string]string{"email": email, "password": password})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("signup: got %d (%s)", resp.StatusCode, env.Message)
	}
	var d struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(env.Data, &d)
	if d.AccessToken == "" {
		t.Fatal("signup: no access token returned")
	}
	return d.AccessToken
}

func createKey(t *testing.T, client *http.Client, srv *httptest.Server, token string) string {
	t.Helper()
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/user_mgt/api_keys/",
		map[string]string{"Authorization": "Bearer " + token}, map[string]any{"name": "test"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create key: got %d (%s)", resp.StatusCode, env.Message)
	}
	var d struct {
		ApiKey string `json:"api_key"`
	}
	json.Unmarshal(env.Data, &d)
	if d.ApiKey == "" {
		t.Fatal("create key: no raw key returned")
	}
	return d.ApiKey
}

func shorten(t *testing.T, client *http.Client, srv *httptest.Server, key, original string) (int, string) {
	t.Helper()
	resp, env := request(t, client, srv, http.MethodPost, "/api/v1/url_mgt/shorten/",
		map[string]string{"X-API-Key": key}, map[string]string{"original_url": original})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("shorten: got %d (%s)", resp.StatusCode, env.Message)
	}
	var u struct {
		Id       int    `json:"id"`
		ShortUrl string `json:"short_url"`
	}
	json.Unmarshal(env.Data, &u)
	if u.ShortUrl == "" {
		t.Fatal("shorten: no short url returned")
	}
	return u.Id, u.ShortUrl
}

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
	var list []struct {
		AccessCount int `json:"access_count"`
	}
	json.Unmarshal(env.Data, &list)
	if len(list) != 1 || list[0].AccessCount != 1 {
		t.Fatalf("expected one url with access_count 1, got %+v", list)
	}
}

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
