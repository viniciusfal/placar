package game

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGameEntity(t *testing.T) {
	assert := assert.New(t)

	input := CreateGameInput{
		Jogador1:            "vinicius",
		Jogador2:            "jose",
		PontuacaoNecessaria: 10,
	}

	game, err := CreateGame(input)

	assert.NoError(err)
	assert.NotEmpty(game.ID)
	assert.Equal(input.Jogador1, game.Jogador1)
	assert.Equal(input.Jogador2, game.Jogador2)
	assert.Equal(input.PontuacaoNecessaria, game.PontuacaoNecessaria)
	assert.Equal(0, game.PontuacaoJogador1)
	assert.Equal(0, game.PontuacaoJogador2)
	assert.Equal("em_andamento", game.StatusPartida)
	assert.WithinDuration(time.Now(), game.InicioPartida, time.Second)
	assert.Nil(game.Vencedor)
	assert.Nil(game.AtualizadoEm)
}
