// Package server package contains all logic needed to create the http server
// this includes http server creation and http route registration
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/rbrady98/steiger/internal/config"
	"github.com/rbrady98/steiger/internal/services"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port string

	log *slog.Logger

	jokeSvc *services.JokeService
}

func NewServer(ctx context.Context, cfg config.Config, logger *slog.Logger, jokeSvc *services.JokeService) *http.Server {
	s := &Server{
		port:    cfg.Port,
		log:     logger,
		jokeSvc: jokeSvc,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		// Ensure that in-flight requests are not cancelled during graceful shutdown
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	return server
}
