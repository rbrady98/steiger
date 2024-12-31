package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rbrady98/steiger/internal/config"
	"github.com/rbrady98/steiger/internal/metrics"
	"github.com/rbrady98/steiger/internal/services/joke"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port string

	log *slog.Logger

	metrics *metrics.Metrics

	jokeSvc *joke.JokeService
}

func NewServer(
	cfg config.Config,
	logger *slog.Logger,
	metrics *metrics.Metrics,
	jokeSvc *joke.JokeService,
) *http.Server {
	s := &Server{
		port:    cfg.Port,
		log:     logger,
		metrics: metrics,
		jokeSvc: jokeSvc,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
