//go:build integration

package mongodb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcmongodb "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/viniciusfal/placar/internal/domain/game"
	"github.com/viniciusfal/placar/internal/platform/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupTestDB(t *testing.T) *mongo.Database {
	t.Helper()

	ctx := context.Background()
	container, err := tcmongodb.Run(ctx, "mongo:7")
	require.NoError(t, err) // Se tiver erro o container vai ser parado imediatamente (diferente do assert)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("falha ao encerrar container: %v", err)
		}
	})
	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	cliente, err := mongodb.Connect(ctx, uri)
	require.NoError(t, err)

	return cliente.Database("test")
}

func TestGameRepository_Salvar(t *testing.T) {
	db := setupTestDB(t)
	repo := mongodb.NewGameRepository(db)

	g, err := game.CreateGame(game.CreateGameInput{
		Jogador1:            "vinicius",
		Jogador2:            "jose",
		PontuacaoNecessaria: 10,
	})
	require.NoError(t, err)

	err = repo.Salvar(context.Background(), g)
	assert.NoError(t, err)
}

func TestGameRepository_BuscarID(t *testing.T) {
	db := setupTestDB(t)
	repo := mongodb.NewGameRepository(db)

	g, _ := game.CreateGame(game.CreateGameInput{
		Jogador1:            "vinicius",
		Jogador2:            "jose",
		PontuacaoNecessaria: 10,
	})

	err := repo.Salvar(context.Background(), g)
	require.NoError(t, err)

	gameEncontrato, err := repo.BuscarPorID(context.Background(), g.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, g.Jogador1, gameEncontrato.Jogador1)
	assert.Equal(t, g.Jogador2, gameEncontrato.Jogador2)
}

func TestGameRepository_BuscaPorIdNaoEncontrado(t *testing.T) {
	db := setupTestDB(t)
	repo := mongodb.NewGameRepository(db)

	_, err := repo.BuscarPorID(context.Background(), "id_nao_encontrado")

	assert.ErrorIs(t, err, game.ErrGameNaoEncontrado)
}
