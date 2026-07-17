package urls

import (
	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
)

func Routes(r *mux.Router, h *Handler, auth *v1middlewares.Authenticator) {
	router := r.PathPrefix("/url_mgt").Subrouter()
	router.Use(auth.Either)
	router.HandleFunc("/shorten/", h.Shorten).Methods("POST")
	router.HandleFunc("/custom/", h.CustomUrl).Methods("POST")
	router.HandleFunc("/url_expiry/", h.UrlExpiry).Methods("POST")
	router.HandleFunc("/", h.UrlFilter).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", h.UpdateDeleteUrl).Methods("PUT", "PATCH", "DELETE")
}
