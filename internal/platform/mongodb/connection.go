package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao mongodb: %w", err)
	}

	if err := client.Ping(c, nil); err != nil {
		return nil, fmt.Errorf("mongodb não respondeu ao ping: %w", err)
	}

	return client, nil
}

// Disconnect fecha a conexão graciosamente. Chamado com defer no main.go.
// Loga em vez de propagar erro — falha ao desconectar não deveria
// impedir o resto do shutdown de acontecer.
func Disconnect(ctx context.Context, client *mongo.Client, log *slog.Logger) {
	c, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(c); err != nil {
		log.Error("falha ao desconectar do mongodb", slog.String("error", err.Error()))
	}
}
