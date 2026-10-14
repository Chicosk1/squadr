package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/graph"
)

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	ctx, parar := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer parar()

	cfg := config.Carregar()
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /saude", func(w http.ResponseWriter, r *http.Request) {
		ctxPing, cancelar := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancelar()

		w.Header().Set("Content-Type", "application/json")
		if err := cliente.Ping(ctxPing); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "indisponivel", "banco": "erro"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "banco": "ok"})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.PortaAPI,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	erroServidor := make(chan error, 1)
	go func() { erroServidor <- srv.ListenAndServe() }()
	log.Printf("API ouvindo em http://localhost:%s", cfg.PortaAPI)

	select {
	case err := <-erroServidor:
		return err
	case <-ctx.Done():
		log.Println("encerrando...")
		ctxFim, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelar()
		return srv.Shutdown(ctxFim)
	}
}
