package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rbrady98/steiger/internal/services/joke"
	"github.com/rbrady98/steiger/internal/storage"
)

type SqliteJokeRepo struct {
	db *sqlx.DB
}

var _ joke.JokeRepo = &SqliteJokeRepo{}

type JokeModel struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}

func NewSqliteJokeRepo(db *sqlx.DB) *SqliteJokeRepo {
	return &SqliteJokeRepo{
		db: db,
	}
}

func (r *SqliteJokeRepo) Get(ctx context.Context, id int) (joke.Joke, error) {
	var j JokeModel
	if err := r.db.GetContext(ctx, &j, `SELECT * FROM jokes WHERE id = ? LIMIT 1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return joke.Joke{}, storage.ErrNoRows
		}
	}

	return fromJokeModel(j), nil
}

func (r *SqliteJokeRepo) Create(ctx context.Context, content string, nsfw bool) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO jokes (joke, nsfw, createdat) VALUES ( $1, $2, datetime('now'))`,
		content,
		nsfw,
	)
	return err
}

func (r *SqliteJokeRepo) List(ctx context.Context) ([]joke.Joke, error) {
	var j []JokeModel
	err := r.db.Select(&j, `SELECT * FROM jokes LIMIT 50`)

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
