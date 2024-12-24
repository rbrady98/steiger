package joke

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rbrady98/steiger/internal/storage"
)

var ErrNotFound = errors.New("not found")

type Joke struct {
	ID        int       `json:"id"`
	Joke      string    `json:"joke"`
	Nsfw      bool      `json:"nsfw"`
	CreatedAt time.Time `json:"createdAt"`
}

type JokeRepo interface {
	Get(ctx context.Context, id int) (Joke, error)
	Create(ctx context.Context, content string, nsfw bool) error
	List(ctx context.Context) ([]Joke, error)
}

type JokeService struct {
	log  *slog.Logger
	repo JokeRepo
}

func NewJokeService(log *slog.Logger, repo JokeRepo) *JokeService {
	return &JokeService{
		log:  log,
		repo: repo,
	}
}

func (j *JokeService) GetJoke(ctx context.Context, id int) (Joke, error) {
	j.log.Info(
		"incoming request",
		slog.String("msg", "joke service called"),
	)

	joke, err := j.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			return Joke{}, fmt.Errorf("no joke with id %d found: %w", id, ErrNotFound)
		}
		return Joke{}, err
	}

	return joke, nil
}

func (j *JokeService) CreateJoke(ctx context.Context, joke string, nsfw bool) error {
	if joke == "" {
		return errors.New("joke string cannot be empty")
	}

	return j.repo.Create(ctx, joke, nsfw)
}

func (j *JokeService) ListJokes(ctx context.Context) ([]Joke, error) {
	jokes, err := j.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return jokes, nil
}
