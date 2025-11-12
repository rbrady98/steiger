package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rbrady98/steiger/internal/domain/joke"
)

type JokeService struct {
	log  *slog.Logger
	repo joke.Repository
}

func NewJokeService(log *slog.Logger, repo joke.Repository) *JokeService {
	return &JokeService{
		log:  log,
		repo: repo,
	}
}

func (s *JokeService) GetJoke(ctx context.Context, id int) (joke.Joke, error) {
	s.log.Info(
		"incoming request",
		slog.String("msg", "joke service called"),
	)

	jk, err := s.repo.Get(ctx, id)
	if err != nil {
		return joke.Joke{}, fmt.Errorf("error getting joke with id %d: %w", id, err)
	}

	return jk, nil
}

func (s *JokeService) CreateJoke(ctx context.Context, joke string, nsfw bool) error {
	if joke == "" {
		return errors.New("joke string cannot be empty")
	}

	return s.repo.Create(ctx, joke, nsfw)
}

func (s *JokeService) ListJokes(ctx context.Context) ([]joke.Joke, error) {
	jokes, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return jokes, nil
}
