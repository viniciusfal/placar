package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/viniciusfal/placar/internal/apierror"
	"github.com/viniciusfal/placar/internal/metrics"
)

const (
	HeaderRequestID = "X-Request-ID"
	traceIDKey      = "trace_id"
)

type ctxKey string

const loggerCtxKey ctxKey = "logger"

// LoggerFromContext nunca retorna nil — fallback seguro se ninguém injetou logger.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerCtxKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// RequestID gera (ou reaproveita) um trace_id e guarda no gin.Context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(HeaderRequestID)
		if traceID == "" {
			traceID = uuid.NewString()
		}
		c.Set(traceIDKey, traceID)
		c.Header(HeaderRequestID, traceID)
		c.Next()
	}
}

// StructuredLogger cria um logger com trace_id embutido e injeta no
// context.Context padrão (não no gin.Context), pra handlers e services
// usarem sem depender de Gin.
func StructuredLogger(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID, _ := c.Get(traceIDKey)

		reqLogger := base.With(
			slog.String("trace_id", traceID.(string)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		ctx := context.WithValue(c.Request.Context(), loggerCtxKey, reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		reqLogger.Info(
			"request finalizada",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unmatched" // rota não encontrada (404), evita cardinalidade infinita
		}

		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recuperado", slog.Any("panic", r))
				p := apierror.New(http.StatusInternalServerError, "Erro interno do servidor", "")
				apierror.Respond(c, log, p)
			}
		}()
		c.Next()
	}
}
