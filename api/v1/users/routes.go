package users

import (
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router, h *Handler) {
	router := r.PathPrefix("/user_mgt").Subrouter()
	router.HandleFunc("/join/", h.SignUp).Methods("POST")
	router.HandleFunc("/new_token/", h.NewToken).Methods("POST")
}
