package v1

import (
	"github.com/gorilla/mux"
	"github.com/negeek/short-access/api/v1/handlers/urls"
	"github.com/negeek/short-access/api/v1/handlers/users"

)

func V1routes(r *mux.Router) {
	router := r.PathPrefix("/api/v1").Subrouter()
	users.Routes(router)
	urls.Routes(router)
	
}