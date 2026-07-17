package integration_tests

import (
	"net/http"
	"testing"
)

func TestHealthz(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()

	resp, env := request(t, noRedirectClient(), srv, http.MethodGet, "/healthz", nil, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("healthz: got %d", resp.StatusCode)
	}
	if !env.Success {
		t.Fatalf("healthz: expected success, got %q", env.Message)
	}
}
