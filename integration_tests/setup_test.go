// Package integration_tests holds end-to-end tests that drive the assembled
// HTTP stack against a real database. Unit tests live next to the code they
// cover; these cross-cutting flows live here.
package integration_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
