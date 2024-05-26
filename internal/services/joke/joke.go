package joke

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rbrady98/steiger/internal/storage/sqlite"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type JokeService struct {
	log  *slog.Logger
	repo sqlite.JokeRepo
}

type Joke struct {
	ID        int       `json:"id"`
	Joke      string    `json:"joke"`
	Nsfw      bool      `json:"nsfw"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewJokeService(log *slog.Logger, db *sqlx.DB) *JokeService {
	return &JokeService{
		log:  log,
		repo: sqlite.NewSqliteJokeRepo(db),
	}
}

func (j *JokeService) GetJoke(ctx context.Context, id int) (Joke, error) {
	j.log.Info(
		"incoming request",
		slog.String("msg", "joke service called"),
	)

	joke, err := j.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Joke{}, fmt.Errorf("no joke with id %d found: %w", id, ErrNotFound)
		}
		return Joke{}, err
	}

	return fromStorage(joke), nil
}

func (j *JokeService) CreateJoke(ctx context.Context, joke string, nsfw bool) error {
	if joke == "" {
		return errors.New("joke string cannot be empty")
	}

	return j.repo.Create(ctx, sqlite.CreateJokeParams{Joke: joke, Nsfw: nsfw})
}

func (j *JokeService) ListJokes(ctx context.Context) ([]Joke, error) {
	jokes, err := j.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	newJokes := make([]Joke, 0, len(jokes))
	for _, joke := range jokes {
		newJokes = append(newJokes, fromStorage(joke))
	}

	return newJokes, nil
}

// fromStorage adapts a db joke struct to a service layer joke
func fromStorage(j sqlite.Joke) Joke {
	return Joke{
		ID:        j.ID,
		Joke:      j.Joke,
		Nsfw:      j.Nsfw,
		CreatedAt: j.CreatedAt,
	}
}
