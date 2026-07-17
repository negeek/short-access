package integration_tests

import (
	"net/http"
	"testing"
)

func TestDocsEndpoints(t *testing.T) {
	srv := newServer(t)
	defer srv.Close()
	client := noRedirectClient()

	for _, path := range []string{"/openapi.yaml", "/docs"} {
		resp, _ := request(t, client, srv, http.MethodGet, path, nil, nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s: got %d", path, resp.StatusCode)
		}
	}
}
