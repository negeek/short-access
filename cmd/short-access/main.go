package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/negeek/short-access/api/v1/handlers/users"
	"github.com/negeek/short-access/middlewares/v1"
	"github.com/negeek/short-access/api/v1/handlers/urls"
		)


func main(){

	envErr := godotenv.Load("../../internal/env/.env")
	
    if envErr != nil {
        log.Fatalf("Error loading .env file: %s", envErr)
    }
	
	//custom servermutiplexer
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", urls.UrlRedirect)

	user_mgt := router.PathPrefix("/api/v1/user_mgt").Subrouter()
	user_mgt.HandleFunc("/join/", users.SignUp)

	url_mgt:=router.PathPrefix("/api/v1/url").Subrouter()
	url_mgt.Use(middlewares.AuthenticationMiddleware)
	url_mgt.HandleFunc("/shorten/", urls.Shorten)

	
	//custom server
	server:=&http.Server{
		Addr: ":8080",
		Handler: router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// start server
	fmt.Println("server start")
	serverErr:=server.ListenAndServe()
	if serverErr != nil {
		fmt.Printf("error listening for server: %s\n", serverErr)
	}

}