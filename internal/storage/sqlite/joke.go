package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Joke struct {
	ID        int
	Joke      string
	Nsfw      bool
	CreatedAt time.Time
}

type JokeRepo interface {
	Get(ctx context.Context, id int) (Joke, error)
	Create(ctx context.Context, params CreateJokeParams) error
	List(ctx context.Context) ([]Joke, error)
}

type SqliteJokeRepo struct {
	db *sqlx.DB
}

var _ JokeRepo = &SqliteJokeRepo{}

func NewSqliteJokeRepo(db *sqlx.DB) *SqliteJokeRepo {
	return &SqliteJokeRepo{
		db: db,
	}
}

func (r *SqliteJokeRepo) Get(ctx context.Context, id int) (Joke, error) {
	row := r.db.QueryRowxContext(ctx, `SELECT * FROM jokes WHERE id = ? LIMIT 1`, id)

	var j Joke
	err := row.StructScan(&j)
	if err != nil {
		return Joke{}, err
	}

	return j, nil
}

type CreateJokeParams struct {
	Joke string
	Nsfw bool
}

func (r *SqliteJokeRepo) Create(ctx context.Context, params CreateJokeParams) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO jokes (joke, nsfw, createdat) VALUES ( $1, $2, datetime('now'))`, params.Joke, params.Nsfw)
	return err
}

func (r *SqliteJokeRepo) List(ctx context.Context) ([]Joke, error) {
	var j []Joke
	err := r.db.Select(&j, `SELECT * FROM jokes LIMIT 50`)

	return j, err
}
