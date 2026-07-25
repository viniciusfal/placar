package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, url string) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parsear REDIS_URL: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis não respondeu ao ping: %w", err)
	}

	return client, nil
}

func Disconnect(client *redis.Client, log *slog.Logger) {
	if err := client.Close(); err != nil {
		log.Error("falha ao desconectar do redis", slog.String("error", err.Error()))
	}
}
