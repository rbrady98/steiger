package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rbrady98/steiger/internal/storage/model"
)

type SqliteJokeRepo struct {
	db *sqlx.DB
}

var _ model.JokeRepo = &SqliteJokeRepo{}

func NewSqliteJokeRepo(db *sqlx.DB) *SqliteJokeRepo {
	return &SqliteJokeRepo{
		db: db,
	}
}

func (r *SqliteJokeRepo) Get(ctx context.Context, id int) (model.Joke, error) {
	row := r.db.QueryRowxContext(ctx, `SELECT * FROM jokes WHERE id = ? LIMIT 1`, id)

	var j model.Joke
	err := row.StructScan(&j)
	if err != nil {
		return model.Joke{}, err
	}

	return j, nil
}

type CreateJokeParams struct {
	Joke string
	Nsfw bool
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

func (r *SqliteJokeRepo) List(ctx context.Context) ([]model.Joke, error) {
	var j []model.Joke
	err := r.db.Select(&j, `SELECT * FROM jokes LIMIT 50`)

	return j, err
}
