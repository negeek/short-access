package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/negeek/short-access/api"
	"github.com/negeek/short-access/api/v1/apikeys"
	"github.com/negeek/short-access/api/v1/urls"
	"github.com/negeek/short-access/api/v1/users"
	"github.com/negeek/short-access/db"
	v1middlewares "github.com/negeek/short-access/middlewares/v1"
	apikeyrepo "github.com/negeek/short-access/repository/v1/apikey"
	numberrepo "github.com/negeek/short-access/repository/v1/number"
	urlrepo "github.com/negeek/short-access/repository/v1/url"
	userrepo "github.com/negeek/short-access/repository/v1/user"
	routes "github.com/negeek/short-access/routes/v1"
	apikeyservice "github.com/negeek/short-access/service/v1/apikey"
	urlservice "github.com/negeek/short-access/service/v1/url"
	userservice "github.com/negeek/short-access/service/v1/user"
)

func main() {
	setupLogger()

	if os.Getenv("APP_ENV") == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			slog.Error("could not load .env file", "error", err)
			os.Exit(1)
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
	slog.Info("connecting to database")
	pool, err := db.Connect(dbURL)
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Wire the layers together: repositories -> services -> handlers.
	urlHandler := urls.NewHandler(urlservice.NewService(
		urlrepo.NewRepository(pool),
		numberrepo.NewRepository(pool),
	))
	userService := userservice.NewService(userrepo.NewRepository(pool))
	userHandler := users.NewHandler(userService)
	apiKeyService := apikeyservice.NewService(apikeyrepo.NewRepository(pool))
	apiKeyHandler := apikeys.NewHandler(apiKeyService)
	auth := v1middlewares.NewAuthenticator(userService, apiKeyService)

	// Routing.
	router := mux.NewRouter()
	router.Use(v1middlewares.CORS)
	router.HandleFunc("/", api.Home).Methods("GET")
	router.HandleFunc("/{slug}", urlHandler.UrlRedirect).Methods("GET")
	routes.V1routes(router.StrictSlash(true), urlHandler, userHandler, apiKeyHandler, auth)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in the background so we can wait for a shutdown signal.
	go func() {
		slog.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server stopped unexpectedly", "error", err)
		}
	}()

	// Wait for an interrupt, then shut down gracefully.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	os.Exit(0)
}

// setupLogger installs a structured logger as the default. The level comes from
// the LOG_LEVEL env var (debug, info, warn, error) and defaults to info.
func setupLogger() {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}
