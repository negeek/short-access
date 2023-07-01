package urls

import (
	"github.com/gorilla/mux"
	"github.com/negeek/short-access/middlewares/v1"

)

func UrlRoutes(r *mux.Router) {
	router := r.PathPrefix("/url").Subrouter()
	router.Use(middlewares.AuthenticationMiddleware)
	router.HandleFunc("/shorten/", Shorten).Methods("POST")
}