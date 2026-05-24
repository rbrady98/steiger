package server

import (
	"fmt"
	"net/http"

	"github.com/rbrady98/steiger/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rCtx := chi.NewRouteContext()
			routePattern := ""
			if router.Match(rCtx, r.Method, r.URL.Path) {
				routePattern = rCtx.RoutePattern()
			}

			newHandler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(semconv.URLFull(r.URL.String()), semconv.HTTPRoute(routePattern))

				// add the route pattern for metrics
				l, ok := otelhttp.LabelerFromContext(r.Context())
				if ok {
					l.Add(semconv.HTTPRoute(routePattern))
				}

				h.ServeHTTP(w, r)
			}), fmt.Sprintf("%s %s", r.Method, routePattern))
			newHandler.ServeHTTP(w, r)
		})
	})

	router.Get("/", handler(handlers.ListJokesHandler(s.log, s.jokeSvc)))
	router.Get("/{id}", handler(handlers.GetJokeHandler(s.log, s.jokeSvc)))
	router.Post("/", handler(handlers.CreateJokeHandler(s.log, s.jokeSvc)))

	return router
}
