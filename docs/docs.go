// Package docs serves the OpenAPI spec and a Swagger UI page. The spec is
// embedded, so it ships inside the binary; the UI page pulls Swagger UI's assets
// from a CDN, so viewing /docs in a browser needs internet access.
package docs

import (
	"embed"
	"net/http"
)

//go:embed openapi.yaml swagger.html
var files embed.FS

// Spec serves the OpenAPI document.
func Spec(w http.ResponseWriter, r *http.Request) {
	serve(w, "openapi.yaml", "application/yaml")
}

// UI serves the Swagger UI page pointed at the spec.
func UI(w http.ResponseWriter, r *http.Request) {
	serve(w, "swagger.html", "text/html; charset=utf-8")
}

func serve(w http.ResponseWriter, name, contentType string) {
	data, err := files.ReadFile(name)
	if err != nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
