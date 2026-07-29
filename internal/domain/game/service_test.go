package game_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/viniciusfal/placar/internal/domain/game"
	"github.com/viniciusfal/placar/internal/domain/game/mocks"
)

func TestService_Create(t *testing.T) {
	t.Run("cria e salva o game com sucesso", func(t *testing.T) {
		repoMock := mocks.NewMockRepository(t)
		repoMock.EXPECT().
			Salvar(mock.Anything, mock.AnythingOfType("*game.Game")).
			Return(nil)

		svc := game.NewService(repoMock)

		game, err := svc.Create(context.Background(), game.CreateGameInput{
			Jogador1:            "vinicius",
			Jogador2:            "jose",
			PontuacaoNecessaria: 10,
		})

		assert.NoError(t, err)
		assert.NotNil(t, game)
	})

	t.Run("nao salva se domain rejeitar", func(t *testing.T) {
		repoMock := mocks.NewMockRepository(t)

		svc := game.NewService(repoMock)

		_, err := svc.Create(context.Background(), game.CreateGameInput{
			Jogador2:            "jose",
			PontuacaoNecessaria: 10,
			// jogador1 Omitido de proposito
		})

		assert.ErrorIs(t, err, game.ErrJogador1Obrigatorio)
	})
}
