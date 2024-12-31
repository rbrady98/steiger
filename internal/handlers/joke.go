package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	apperror "github.com/rbrady98/steiger/internal/app_error"
	"github.com/rbrady98/steiger/internal/codec"
	"github.com/rbrady98/steiger/internal/metrics"
	"github.com/rbrady98/steiger/internal/services/joke"
)

func GetJokeHandler(log *slog.Logger, jokeSvc *joke.JokeService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		IDstr := chi.URLParam(r, "id")
		ID, err := strconv.Atoi(IDstr)
		if err != nil {
			return apperror.New("invalid joke id", http.StatusBadRequest)
		}

		j, err := jokeSvc.GetJoke(r.Context(), ID)
		if err != nil {
			if errors.Is(err, joke.ErrNotFound) {
				return apperror.NewFromError(err, http.StatusNotFound)
			}
			return apperror.NewFromError(err, http.StatusInternalServerError)
		}

		_ = codec.Encode(w, http.StatusOK, j)
		return nil
	}
}

func CreateJokeHandler(
	log *slog.Logger,
	metrics *metrics.Metrics,
	jokeSvc *joke.JokeService,
) func(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Joke string `json:"joke"`
		Nsfw bool   `json:"nsfw"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := codec.Decode[request](r)
		if err != nil {
			return apperror.NewFromError(err, http.StatusBadRequest)
		}

		err = jokeSvc.CreateJoke(r.Context(), req.Joke, req.Nsfw)
		if err != nil {
			return apperror.NewFromError(err, http.StatusInternalServerError)
		}

		metrics.JokesCreated.Inc()

		return nil
	}
}

func ListJokesHandler(log *slog.Logger, jokeSvc *joke.JokeService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		jokes, err := jokeSvc.ListJokes(r.Context())
		if err != nil {
			return apperror.NewFromError(err, http.StatusInternalServerError)
		}

		return codec.Encode(w, http.StatusOK, jokes)
	}
}
