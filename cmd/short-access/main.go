package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	"github.com/negeek/short-access/api/v1/urls"
	routes "github.com/negeek/short-access/routes/v1"
	"github.com/negeek/short-access/db"
	"os"
    "os/signal"
	"context"
	"syscall"
		)


func main(){

	err := godotenv.Load("../../internal/env/.env")
	
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	
	//custom servermutiplexer
	router := mux.NewRouter()
	router.Use(v1middlewares.CORS)
	router.HandleFunc("/{slug}", urls.UrlRedirect).Methods("GET")
	routes.V1routes(router.StrictSlash(true))

	// DB connection
	dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL not set")
    }
	if err = db.Connect(dbURL); err != nil {
		log.Fatal(err)
	}

	
	//custom server
	server:=&http.Server{
		Addr: ":8080",
		Handler: router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 *  time.Second,
	}

	// Run server in a goroutine so that it doesn't block.
	go func() {
		fmt.Println("server start")
		if err = server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	fmt.Println("shutting down")
	os.Exit(0)

}