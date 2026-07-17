package apikeys

import (
	"github.com/gorilla/mux"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
)

func Routes(r *mux.Router, h *Handler, auth *v1middlewares.Authenticator) {
	router := r.PathPrefix("/user_mgt/api_keys").Subrouter()
	router.Use(auth.JWT)
	router.HandleFunc("/", h.Create).Methods("POST")
	router.HandleFunc("/", h.List).Methods("GET")
	router.HandleFunc("/{id:[0-9]+}", h.Delete).Methods("DELETE")
	router.HandleFunc("/{id:[0-9]+}/revoke", h.Revoke).Methods("POST")
}
