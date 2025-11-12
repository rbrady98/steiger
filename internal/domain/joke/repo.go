package joke

import (
	"context"
	"errors"
)

type Repository interface {
	Get(ctx context.Context, id int) (Joke, error)
	Create(ctx context.Context, content string, nsfw bool) error
	List(ctx context.Context) ([]Joke, error)
}

var ErrNotFound = errors.New("not found")
