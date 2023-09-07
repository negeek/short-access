package urls

import (
	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"

)

func Routes(r *mux.Router) {
	router := r.PathPrefix("/url").Subrouter()
	router.Use(v1middlewares.AuthenticationMiddleware)
	router.HandleFunc("/shorten/", Shorten).Methods("POST")
	router.HandleFunc("/custom/", Shorten).Methods("POST")
}