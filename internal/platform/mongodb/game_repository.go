package mongodb

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/viniciusfal/placar/internal/domain/game"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameRepository struct {
	collection *mongo.Collection
}

func NewGameRepository(db *mongo.Database) *GameRepository {
	return &GameRepository{
		collection: db.Collection("games"),
	}
}

func (r *GameRepository) Salvar(ctx context.Context, g *game.Game) error {
	_, err := r.collection.InsertOne(ctx, g)
	return err
}

func (r *GameRepository) BuscarPorID(ctx context.Context, id string) (*game.Game, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, game.ErrGameNaoEncontrado
	}

	var g game.Game
	err = r.collection.FindOne(ctx, bson.M{"_id": uid}).Decode(&g)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, game.ErrGameNaoEncontrado
	}
	if err != nil {
		return nil, err
	}

	return &g, nil
}
