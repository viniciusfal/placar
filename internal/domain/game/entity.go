package game

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID                  uuid.UUID
	Jogador1            string
	Jogador2            string
	PontuacaoJogador1   int
	PontuacaoJogador2   int
	PontuacaoNecessaria int
	Vencedor            *string
	StatusPartida       string
	InicioPartida       time.Time
	AtualizadoEm        *time.Time
}

type CreateGameInput struct {
	Jogador1            string
	Jogador2            string
	PontuacaoNecessaria int
}

func CreateGame(input CreateGameInput) (Game, error) {
	game := Game{
		ID:                  uuid.New(),
		Jogador1:            input.Jogador1,
		Jogador2:            input.Jogador2,
		PontuacaoNecessaria: input.PontuacaoNecessaria,
		PontuacaoJogador1:   0,
		PontuacaoJogador2:   0,
		StatusPartida:       "em_andamento",
		InicioPartida:       time.Now(),
	}

	return game, nil
}
