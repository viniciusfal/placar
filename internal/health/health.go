package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Checker func(ctx context.Context) error

type Health struct {
	checks map[string]Checker
}

func New() *Health {
	return &Health{checks: make(map[string]Checker)}
}

// Register adiciona uma dependência a ser verificada no readiness.
func (h *Health) Register(name string, check Checker) {
	h.checks[name] = check
}

// Liveness só confirma que o processo responde — sem checar dependências.
func (h *Health) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readiness checa todas as dependências registradas.
func (h *Health) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	results := gin.H{}
	healthy := true

	for name, check := range h.checks {
		if err := check(ctx); err != nil {
			results[name] = err.Error()
			healthy = false
			continue
		}
		results[name] = "ok"
	}

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{"checks": results})
}

/* Aplicacao
	// cmd/api/main.go
	h := health.New()

	h.Register("database", func(ctx context.Context) error {
		return db.PingContext(ctx)
	})

	h.Register("redis", func(ctx context.Context) error {
		return redisClient.Ping(ctx).Err()
})

#####################################################

	// internal/router/router.go
	r.GET("/health/liveness", h.Liveness)
	r.GET("/health/readiness", h.Readiness)
*/
