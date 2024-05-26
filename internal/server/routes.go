package server

import (
	"net/http"

	"github.com/rbrady98/steiger/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handler(handlers.ListJokesHandler(s.log, s.jokeSvc)))
	r.Get("/{id}", handler(handlers.GetJokeHandler(s.log, s.jokeSvc)))
	r.Post("/", handler(handlers.CreateJokeHandler(s.log, s.jokeSvc)))

	return r
}
