package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	Mode        string
	DatabaseURL string
	RedisAddr   string
	CORSOrigins []string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:        getEnv("PORT", "8000"),
		Mode:        getEnv("APP_ENV", "development"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("Configuração inválida: %w", err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL é obrigatório")
	}

	if c.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR é obrigatório")
	}

	return nil
}
