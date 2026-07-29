package game

import "context"

type Repository interface {
	Salvar(ctx context.Context, g *Game) error
	BuscarPorID(ctx context.Context, id string) (*Game, error)
}
