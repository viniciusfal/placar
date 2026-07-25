package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/viniciusfal/placar/internal/health"
	"github.com/viniciusfal/placar/internal/middleware"
)

type Dependencies struct {
	Logger *slog.Logger
	Health *health.Health
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

	api := r.Group("/api/v1")
	api.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
