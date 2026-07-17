// Package server wires the layers together into an HTTP router. Keeping this in
// one place lets both the binary and the tests build the exact same stack.
package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/negeek/short-access/api"
	"github.com/negeek/short-access/api/v1/apikeys"
	"github.com/negeek/short-access/api/v1/urls"
	"github.com/negeek/short-access/api/v1/users"
	"github.com/negeek/short-access/docs"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	apikeyrepo "github.com/negeek/short-access/repository/v1/apikey"
	numberrepo "github.com/negeek/short-access/repository/v1/number"
	urlrepo "github.com/negeek/short-access/repository/v1/url"
	userrepo "github.com/negeek/short-access/repository/v1/user"
	routes "github.com/negeek/short-access/routes/v1"
	apikeyservice "github.com/negeek/short-access/service/v1/apikey"
	urlservice "github.com/negeek/short-access/service/v1/url"
	userservice "github.com/negeek/short-access/service/v1/user"
)

// NewRouter builds the HTTP router from a database pool: repositories feed
// services, services feed handlers, handlers are mounted on routes.
func NewRouter(pool *pgxpool.Pool) http.Handler {
	urlHandler := urls.NewHandler(urlservice.NewService(
		urlrepo.NewRepository(pool),
		numberrepo.NewRepository(pool),
	))
	userService := userservice.NewService(userrepo.NewRepository(pool))
	userHandler := users.NewHandler(userService)
	apiKeyService := apikeyservice.NewService(apikeyrepo.NewRepository(pool))
	apiKeyHandler := apikeys.NewHandler(apiKeyService)
	auth := v1middlewares.NewAuthenticator(userService, apiKeyService)

	router := mux.NewRouter()
	router.Use(v1middlewares.CORS)
	router.HandleFunc("/", api.Home).Methods("GET")
	router.HandleFunc("/healthz", api.Health(pool)).Methods("GET")
	router.HandleFunc("/docs", docs.UI).Methods("GET")
	router.HandleFunc("/openapi.yaml", docs.Spec).Methods("GET")
	router.HandleFunc("/{slug}", urlHandler.UrlRedirect).Methods("GET")
	routes.V1routes(router.StrictSlash(true), urlHandler, userHandler, apiKeyHandler, auth)
	return router
}
