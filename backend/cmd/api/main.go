package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chicosk1/squadr/backend/internal/platform/config"
	"github.com/Chicosk1/squadr/backend/internal/platform/database"
	"github.com/Chicosk1/squadr/backend/internal/platform/httpserver"
	"github.com/Chicosk1/squadr/backend/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	appLogger := logger.New(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("erro ao conectar no banco", "erro", err)
		os.Exit(1)
	}
	defer pool.Close()

	server := httpserver.New(fmt.Sprintf(":%d", cfg.APIPort), appLogger)

	server.Router().Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			appLogger.Error("healthz: banco inacessível", "erro", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	if err := server.Run(ctx); err != nil {
		appLogger.Error("erro no servidor http", "erro", err)
		os.Exit(1)
	}

	appLogger.Info("servidor encerrado")
}
