package api
import (
	"net/http"
	"github.com/negeek/short-access/utils"
)
func Home(w http.ResponseWriter, r *http.Request) {

	utils.JsonResponse(w, true, http.StatusOK , "Welcome to Short-Access!", map[string]interface{}{"name":"Short-Access","version":"1.0",
	"description":"URL shortener", "documentation":"https://github.com/negeek/short-access"})
	return	
}