package users

import (
	//"fmt"
	"github.com/gorilla/mux"
)

func Routes(r *mux.Router) {
	router := r.PathPrefix("/user_mgt").Subrouter()
	router.HandleFunc("/join/", SignUp).Methods("POST")
	router.HandleFunc("/new_token/", NewToken).Methods("POST")
}