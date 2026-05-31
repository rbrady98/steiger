// Package sqlite contains implementations of repositories backed by sqlite
package joke

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rbrady98/steiger/internal/joke/db"
)

type SQLiteRepo struct {
	q *db.Queries
}

var _ Repository = &SQLiteRepo{}

func NewSQLiteRepo(sqldb *sql.DB) *SQLiteRepo {
	return &SQLiteRepo{q: db.New(sqldb)}
}

func (r *SQLiteRepo) Get(ctx context.Context, id int) (Joke, error) {
	j, err := r.q.GetJoke(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Joke{}, ErrNotFound
		}
		return Joke{}, err
	}

	return fromDBJoke(j), nil
}

func (r *SQLiteRepo) Create(ctx context.Context, content string, nsfw bool) error {
	return r.q.CreateJoke(ctx, db.CreateJokeParams{
		Joke: content,
		Nsfw: nsfw,
	})
}

func (r *SQLiteRepo) List(ctx context.Context) ([]Joke, error) {
	rows, err := r.q.ListJokes(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Joke, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBJoke(row))
	}
	return out, nil
}

func fromDBJoke(j db.Joke) Joke {
	return Joke{
		ID:        int(j.ID),
		Joke:      j.Joke,
		Nsfw:      j.Nsfw,
		CreatedAt: j.CreatedAt,
	}
}
