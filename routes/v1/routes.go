package v1

import (
	"github.com/gorilla/mux"
	"github.com/negeek/short-access/api/v1/apikeys"
	"github.com/negeek/short-access/api/v1/urls"
	"github.com/negeek/short-access/api/v1/users"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
)

func V1routes(r *mux.Router, urlHandler *urls.Handler, userHandler *users.Handler, apiKeyHandler *apikeys.Handler, auth *v1middlewares.Authenticator) {
	router := r.PathPrefix("/api/v1").Subrouter()
	users.Routes(router, userHandler)
	apikeys.Routes(router, apiKeyHandler, auth)
	urls.Routes(router, urlHandler, auth)
}
