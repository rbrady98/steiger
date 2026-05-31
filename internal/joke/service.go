package joke

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type Service struct {
	log  *slog.Logger
	repo Repository
}

func NewService(log *slog.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) GetJoke(ctx context.Context, id int) (Joke, error) {
	s.log.Info(
		"incoming request",
		slog.String("msg", "joke service called"),
	)

	jk, err := s.repo.Get(ctx, id)
	if err != nil {
		return Joke{}, fmt.Errorf("error getting joke with id %d: %w", id, err)
	}

	return jk, nil
}

func (s *Service) CreateJoke(ctx context.Context, joke string, nsfw bool) error {
	if joke == "" {
		return errors.New("joke string cannot be empty")
	}

	return s.repo.Create(ctx, joke, nsfw)
}

func (s *Service) ListJokes(ctx context.Context) ([]Joke, error) {
	jokes, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return jokes, nil
}
