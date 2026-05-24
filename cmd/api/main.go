package main

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/rbrady98/steiger/internal/telemetry"
	"go.opentelemetry.io/contrib/processors/minsev"

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

	logger := telemetry.NewLogger(cfg.Env, minsev.SeverityDebug)

	jokeSvc := services.NewJokeService(logger, sqlite.NewSqliteJokeRepo(db))

	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Ensure in-flight requests aren't cancelled immediately
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	defer stopOngoingGracefully()

	telemetryShutdown, err := telemetry.Setup(ongoingCtx, cfg.Env, minsev.SeverityDebug)
	if err != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
		defer cancel()
		return errors.Join(err, telemetryShutdown(shutdownCtx))
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
		defer cancel()
		err = errors.Join(err, telemetryShutdown(shutdownCtx))
	}()

	srv := server.NewServer(ongoingCtx, cfg, logger, jokeSvc)

	go func() {
		logger.InfoContext(ongoingCtx, "starting server", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("waiting for cancellation")
	<-rootCtx.Done()
	stop()
	logger.Info("Received shutdown signal, shutting down.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		logger.Info("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(shutdownHardPeriod)
	}

	logger.Info("Server shut down gracefully")

	return nil
}
