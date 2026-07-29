package game

import (
	"context"
	"log/slog"
)

type Service struct {
	repo Repository
	_    *slog.Logger
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateGameInput) (*Game, error) {
	g, err := CreateGame(input)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Salvar(ctx, g); err != nil {
		return nil, err
	}

	return &Game{}, nil
}
