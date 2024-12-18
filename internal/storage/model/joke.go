package model

import (
	"context"
	"time"
)

type Joke struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}

type JokeRepo interface {
	Get(ctx context.Context, id int) (Joke, error)
	Create(ctx context.Context, content string, nsfw bool) error
	List(ctx context.Context) ([]Joke, error)
}
