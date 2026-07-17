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

	"github.com/negeek/short-access/db"
	"github.com/negeek/short-access/server"
	"github.com/negeek/short-access/utils"
)

func main() {
	setupLogger()
	loadEnv()

	// Subcommands: "migrate up" / "migrate down". With no arguments the binary
	// runs the HTTP server. One image can therefore both migrate and serve.
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "migrate" {
		runMigrate(args[1:])
		return
	}
	runServer()
}

func runServer() {
	slog.Info("connecting to database")
	pool, err := db.Connect(databaseURL())
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      server.NewRouter(pool),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in the background so we can wait for a shutdown signal.
	go func() {
		slog.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	srv.Shutdown(ctx)
	os.Exit(0)
}

func runMigrate(args []string) {
	direction := "up"
	if len(args) > 0 {
		direction = args[0]
	}

	pool, err := db.Connect(databaseURL())
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	ctx := context.Background()
	switch direction {
	case "up":
		if err := db.MigrateUp(ctx, pool); err != nil {
			slog.Error("migrate up failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied")
	case "down":
		if err := db.MigrateDown(ctx, pool); err != nil {
			slog.Error("migrate down failed", "error", err)
			os.Exit(1)
		}
		slog.Info("rolled back last migration")
	default:
		slog.Error("unknown migrate command; use 'migrate up' or 'migrate down'", "got", direction)
		os.Exit(1)
	}
}

// loadEnv reads a local .env file if one is present (handy for development).
// In Docker there is no .env file and the settings come from the environment,
// which already takes precedence over anything a .env file would set.
func loadEnv() {
	_ = utils.LoadEnvFile(".env")
}

// databaseURL builds the Postgres connection string from the environment.
func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
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
