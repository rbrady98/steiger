// Package sqlite contains implementations of repositories backed by sqlite
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rbrady98/steiger/internal/domain/joke"
)

type JokeRepo struct {
	db *sqlx.DB
}

var _ joke.Repository = &JokeRepo{}

type JokeModel struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}

func NewSqliteJokeRepo(db *sqlx.DB) *JokeRepo {
	return &JokeRepo{
		db: db,
	}
}

func (r *JokeRepo) Get(ctx context.Context, id int) (joke.Joke, error) {
	var j JokeModel
	if err := r.db.GetContext(ctx, &j, `SELECT * FROM jokes WHERE id = ? LIMIT 1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return joke.Joke{}, joke.ErrNotFound
		}
	}

	return fromJokeModel(j), nil
}

func (r *JokeRepo) Create(ctx context.Context, content string, nsfw bool) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO jokes (joke, nsfw, createdat) VALUES ( $1, $2, datetime('now'))`,
		content,
		nsfw,
	)
	return err
}

func (r *JokeRepo) List(ctx context.Context) ([]joke.Joke, error) {
	var j []JokeModel
	err := r.db.SelectContext(ctx, &j, `SELECT * FROM jokes LIMIT 50`)

	return fromJokeModelSlice(j), err
}

func fromJokeModel(m JokeModel) joke.Joke {
	return joke.Joke{
		ID:        m.ID,
		Joke:      m.Joke,
		Nsfw:      m.Nsfw,
		CreatedAt: m.CreatedAt,
	}
}

func fromJokeModelSlice(m []JokeModel) []joke.Joke {
	s := make([]joke.Joke, 0, len(m))

	for _, v := range m {
		s = append(s, fromJokeModel(v))
	}

	return s
}
