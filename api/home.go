package api

import (
	"github.com/negeek/short-access/utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	utils.JsonResponse(w, true, http.StatusOK, "Welcome to Short-Access!", map[string]interface{}{"name": "Short-Access", "version": "1.0",
		"description": "URL shortener", "documentation": "https://github.com/negeek/short-access"})
	return
}
