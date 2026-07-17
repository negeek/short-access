package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/negeek/short-access/api"
	"github.com/negeek/short-access/api/v1/urls"
	"github.com/negeek/short-access/api/v1/users"
	"github.com/negeek/short-access/db"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	numberrepo "github.com/negeek/short-access/repository/v1/number"
	urlrepo "github.com/negeek/short-access/repository/v1/url"
	userrepo "github.com/negeek/short-access/repository/v1/user"
	routes "github.com/negeek/short-access/routes/v1"
	urlservice "github.com/negeek/short-access/service/v1/url"
	userservice "github.com/negeek/short-access/service/v1/user"
)

func main() {
	if os.Getenv("APP_ENV") == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect to the database.
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Println("connecting to db")
	pool, err := db.Connect(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Wire the layers together: repositories -> services -> handlers.
	urlHandler := urls.NewHandler(urlservice.NewService(
		urlrepo.NewRepository(pool),
		numberrepo.NewRepository(pool),
	))
	userService := userservice.NewService(userrepo.NewRepository(pool))
	userHandler := users.NewHandler(userService)
	auth := v1middlewares.NewAuthenticator(userService)

	// Routing.
	router := mux.NewRouter()
	router.Use(v1middlewares.CORS)
	router.HandleFunc("/", api.Home).Methods("GET")
	router.HandleFunc("/{slug}", urlHandler.UrlRedirect).Methods("GET")
	routes.V1routes(router.StrictSlash(true), urlHandler, userHandler, auth)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in the background so we can wait for a shutdown signal.
	go func() {
		fmt.Println("start server")
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for an interrupt, then shut down gracefully.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	fmt.Println("shutting down")
	os.Exit(0)
}
