package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/viniciusfal/placar/internal/health"
	"github.com/viniciusfal/placar/internal/middleware"
)

type Dependencies struct {
	Logger  *slog.Logger
	Health  *health.Health
	Modules []Module
}

type Module interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func New(deps *Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.Recovery(deps.Logger),
		middleware.RequestID(),
		middleware.StructuredLogger(deps.Logger),
		middleware.Metrics(),
	)

	r.GET("/health/liveness", deps.Health.Liveness)
	r.GET("/health/readiness", deps.Health.Readiness)

	v1 := r.Group("/v1")
	for _, m := range deps.Modules {
		m.RegisterRoutes(v1)
	}
	return r
}
