package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	apperror "steiger/internal/app_error"
	"steiger/internal/codec"
	"steiger/internal/services/joke"
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
				return apperror.New(err.Error(), http.StatusNotFound)
			}
			return apperror.New(err.Error(), http.StatusInternalServerError)
		}

		_ = codec.Encode(w, http.StatusOK, j)
		return nil
	}
}

func CreateJokeHandler(log *slog.Logger, jokeSvc *joke.JokeService) func(w http.ResponseWriter, r *http.Request) error {
	type request struct {
		Joke string `json:"joke"`
		Nsfw bool   `json:"nsfw"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		req, err := codec.Decode[request](r)
		if err != nil {
			return apperror.New(err.Error(), http.StatusBadRequest)
		}

		err = jokeSvc.CreateJoke(r.Context(), req.Joke, req.Nsfw)
		if err != nil {
			return apperror.New(err.Error(), http.StatusInternalServerError)
		}

		return nil
	}
}

func ListJokesHandler(log *slog.Logger, jokeSvc *joke.JokeService) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		jokes, err := jokeSvc.ListJokes(r.Context())
		if err != nil {
			return apperror.New(err.Error(), http.StatusInternalServerError)
		}

		return codec.Encode(w, http.StatusOK, jokes)
	}
}
