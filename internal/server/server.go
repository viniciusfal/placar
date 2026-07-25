package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, r *gin.Engine, port string, log *slog.Logger) error {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Info("servidor iniciado", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("falha ao iniciar servidor", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.Info("encerrando servidor, aguardando requisições em andamento...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown forçado: %w", err)
	}

	log.Info("servidor encerrado com sucesso")
	return nil
}
