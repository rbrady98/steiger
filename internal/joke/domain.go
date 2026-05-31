// Package joke contains domain data for jokes
package joke

import (
	"context"
	"errors"
	"time"
)

type Joke struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}

type Repository interface {
	Get(ctx context.Context, id int) (Joke, error)
	Create(ctx context.Context, content string, nsfw bool) error
	List(ctx context.Context) ([]Joke, error)
}

var ErrNotFound = errors.New("not found")
