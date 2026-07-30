package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJogador1Obrigatorio         = errors.New("Jogador 1 é obrigatório")
	ErrJogador2Obrigatorio         = errors.New("Jogador 2 é obrigatório")
	ErrPontuacaoNecessariaInvalida = errors.New("Pontuacao necessaria deve ser maior que 0")
	ErrGameNaoEncontrado           = errors.New("Game não encontrado")
)

type Game struct {
	ID                  uuid.UUID `bson:"_id"` // só esse campo tem tag, por NECESSIDADE técnica real (_id é especial no Mongo)
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

func CreateGame(input CreateGameInput) (*Game, error) {
	if input.Jogador1 == "" {
		return nil, ErrJogador1Obrigatorio
	}

	if input.Jogador2 == "" {
		return nil, ErrJogador2Obrigatorio
	}

	if input.PontuacaoNecessaria < 0 {
		return nil, ErrPontuacaoNecessariaInvalida
	}

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

	return &game, nil
}
