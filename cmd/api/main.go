package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/viniciusfal/placar/internal/config"
	"github.com/viniciusfal/placar/internal/health"
	"github.com/viniciusfal/placar/internal/logger"
	"github.com/viniciusfal/placar/internal/platform/mongodb"
	"github.com/viniciusfal/placar/internal/platform/redis"
	"github.com/viniciusfal/placar/internal/router"
	"github.com/viniciusfal/placar/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("aviso: .env não encontrado, usando variáveis de ambiente do sistema")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l := logger.New(cfg.Mode)

	mongoClient, err := mongodb.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer mongodb.Disconnect(context.Background(), mongoClient, l)
	// db := mongoClient.Database("")

	redisClient, err := redis.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Disconnect(redisClient, l)

	h := health.New()
	h.Register("mongodb", func(ctx context.Context) error {
		return mongoClient.Ping(ctx, nil)
	})
	h.Register("redis", func(ctx context.Context) error {
		return redisClient.Ping(ctx).Err()
	})

	deps := &router.Dependencies{
		Logger: l,
		Health: h}
	r := router.New(deps)

	if err := server.Run(ctx, r, cfg.Port, l); err != nil {
		l.Error("erro ao encerrar servidor", "error", err)
	}
}
