package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rbrady98/steiger/internal/config"
	"github.com/rbrady98/steiger/internal/database"
	"github.com/rbrady98/steiger/internal/server"
	"github.com/rbrady98/steiger/internal/services/joke"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.NewConfig()

	db, err := database.New(cfg.DbURL)
	if err != nil {
		return err
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	jokeSvc := joke.NewJokeService(logger, db)

	srv := server.NewServer(cfg, logger, jokeSvc)

	go func() {
		log.Printf("listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// set up context to listen to ctrl+c
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (thanks Alessandro Rosetti)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}
