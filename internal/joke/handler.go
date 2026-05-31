package joke

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	apperror "github.com/rbrady98/steiger/internal/app_error"
	"github.com/rbrady98/steiger/internal/codec"
)

func HandleGet(logger *slog.Logger, jokeSvc *Service) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger.InfoContext(r.Context(), "getting jokes")
		IDstr := chi.URLParam(r, "id")
		ID, err := strconv.Atoi(IDstr)
		if err != nil {
			return apperror.New("invalid joke id", http.StatusBadRequest)
		}

		j, err := jokeSvc.GetJoke(r.Context(), ID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return apperror.NewFromError(err, http.StatusNotFound)
			}

			return apperror.NewFromError(err, http.StatusInternalServerError)
		}

		_ = codec.Encode(w, http.StatusOK, j)
		return nil
	}
}

func HandleCreate(_ *slog.Logger, jokeSvc *Service) func(w http.ResponseWriter, r *http.Request) error {
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

		w.WriteHeader(http.StatusCreated)
		return nil
	}
}

func HandleList(logger *slog.Logger, jokeSvc *Service) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger.InfoContext(r.Context(), "getting jokes")

		jokes, err := jokeSvc.ListJokes(r.Context())
		if err != nil {
			return apperror.NewFromError(err, http.StatusInternalServerError)
		}

		return codec.Encode(w, http.StatusOK, jokes)
	}
}
