package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/viniciusfal/placar/internal/config"
	"github.com/viniciusfal/placar/internal/logger"
	"github.com/viniciusfal/placar/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao Carregar.env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	l := logger.New(cfg.Mode)

	deps := &router.Dependencies{
		Logger: l,
	}

	r := router.New(deps)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
