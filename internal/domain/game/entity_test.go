package game

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGameEntity(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateGameInput
		wantErr error
	}{
		{
			name:  "criar partida valida",
			input: CreateGameInput{Jogador1: "vinicius", Jogador2: "jose", PontuacaoNecessaria: 10},
		},
		{
			name:    "sem jogador 1",
			input:   CreateGameInput{Jogador2: "jose", PontuacaoNecessaria: 10},
			wantErr: ErrJogador1Obrigatorio,
		},
		{
			name:    "sem jogador 2",
			input:   CreateGameInput{Jogador1: "vinicius", PontuacaoNecessaria: 10},
			wantErr: ErrJogador2Obrigatorio,
		},
		{
			name:    "pontuacao necessaria invalida",
			input:   CreateGameInput{Jogador1: "vinicius", Jogador2: "jose", PontuacaoNecessaria: -1},
			wantErr: ErrPontuacaoNecessariaInvalida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			game, err := CreateGame(tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(err, tt.wantErr)
				return
			}

			assert.NoError(err)
			assert.NotEqual(uuid.Nil, game.ID)
			assert.Equal(tt.input.Jogador1, game.Jogador1)
			assert.Equal(tt.input.Jogador2, game.Jogador2)
			assert.Equal(tt.input.PontuacaoNecessaria, game.PontuacaoNecessaria)
			assert.Equal("em_andamento", game.StatusPartida)
			assert.Equal(0, game.PontuacaoJogador1)
			assert.Equal(0, game.PontuacaoJogador2)
			assert.WithinDuration(time.Now(), game.InicioPartida, time.Second)
			assert.Nil(game.Vencedor)
			assert.Nil(game.AtualizadoEm)
		})
	}
}
