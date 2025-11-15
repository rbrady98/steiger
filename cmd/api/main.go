package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rbrady98/steiger/internal/config"
	"github.com/rbrady98/steiger/internal/database"
	"github.com/rbrady98/steiger/internal/server"
	"github.com/rbrady98/steiger/internal/services"
	"github.com/rbrady98/steiger/internal/storage/sqlite"

	_ "github.com/joho/godotenv/autoload"
)

const (
	shutdownPeriod     = 15 * time.Second
	shutdownHardPeriod = 3 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.NewConfig()

	db, err := database.New(cfg.DbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	jokeSvc := services.NewJokeService(logger, sqlite.NewSqliteJokeRepo(db))

	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Ensure in-flight requests aren't cancelled immediately
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	srv := server.NewServer(ongoingCtx, cfg, logger, jokeSvc)

	go func() {
		log.Println("Server starting on:", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	log.Println("waiting for cancellation")
	<-rootCtx.Done()
	stop()
	log.Println("Received shutdown signal, shutting down.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully")

	return nil
}
