package urls

import (
	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"

)

func Routes(r *mux.Router) {
	router := r.PathPrefix("/url_mgt").Subrouter()
	router.Use(v1middlewares.AuthenticationMiddleware)
	router.HandleFunc("/shorten/", Shorten).Methods("POST")
	router.HandleFunc("/custom/", CustomUrl).Methods("POST")
	router.HandleFunc("/", UrlFilter).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", UrlFilter).Methods("PUT","PATCH","DELETE")
}