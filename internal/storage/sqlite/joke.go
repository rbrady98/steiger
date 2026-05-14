// Package sqlite contains implementations of repositories backed by sqlite
package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rbrady98/steiger/internal/domain/joke"
	"github.com/rbrady98/steiger/internal/sqlitedb"
)

type JokeRepo struct {
	q *sqlitedb.Queries
}

var _ joke.Repository = &JokeRepo{}

func NewSqliteJokeRepo(db *sql.DB) *JokeRepo {
	return &JokeRepo{q: sqlitedb.New(db)}
}

func (r *JokeRepo) Get(ctx context.Context, id int) (joke.Joke, error) {
	j, err := r.q.GetJoke(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return joke.Joke{}, joke.ErrNotFound
		}
		return joke.Joke{}, err
	}

	return fromDBJoke(j), nil
}

func (r *JokeRepo) Create(ctx context.Context, content string, nsfw bool) error {
	return r.q.CreateJoke(ctx, sqlitedb.CreateJokeParams{
		Joke: content,
		Nsfw: nsfw,
	})
}

func (r *JokeRepo) List(ctx context.Context) ([]joke.Joke, error) {
	rows, err := r.q.ListJokes(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]joke.Joke, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBJoke(row))
	}
	return out, nil
}

func fromDBJoke(j sqlitedb.Joke) joke.Joke {
	return joke.Joke{
		ID:        int(j.ID),
		Joke:      j.Joke,
		Nsfw:      j.Nsfw,
		CreatedAt: j.CreatedAt,
	}
}
