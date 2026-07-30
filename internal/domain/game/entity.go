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
	ID                  uuid.UUID  `bson:"_id"`
	Jogador1            string     `bson:"jogador1"`
	Jogador2            string     `bson:"jogador2"`
	PontuacaoJogador1   int        `bson:"pontuacao_jogador1"`
	PontuacaoJogador2   int        `bson:"pontuacao_jogador2"`
	StatusPartida       string     `bson:"status_partida"`
	Vencedor            *string    `bson:"vencedor"`
	PontuacaoNecessaria int        `bson:"pontuacao_necessaria"`
	InicioPartida       time.Time  `bson:"inicio_partida"`
	AtualizadoEm        *time.Time `bson:"atualizado_em"`
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
