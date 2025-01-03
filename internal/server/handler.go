package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/rbrady98/steiger/internal/codec"

	apperror "github.com/rbrady98/steiger/internal/app_error"
)

type handlerFn func(http.ResponseWriter, *http.Request) error

func (fn handlerFn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var e *apperror.AppError
		if errors.As(err, &e) {
			if encodeErr := codec.Encode(w, e.Code, e); encodeErr != nil {
				log.Println("error encoding response", encodeErr)
			}
			return
		}

		http.Error(w, err.Error(), 500)
	}
}

func handler(fn handlerFn) http.HandlerFunc {
	return fn.ServeHTTP
}
