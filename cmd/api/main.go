package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chicosk1/squadr/internal/api"
	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/Chicosk1/squadr/internal/matching"
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
	// Um único cliente para todo o processo (D-019).
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(context.Background())

	// A matriz de complementaridade e os rótulos são lidos uma vez, na inicialização.
	rec, err := matching.Novo(ctx, cliente, domain.OpcoesPadrao())
	if err != nil {
		return err
	}
	rotulos, err := cliente.BuscarRotulosFuncoes(ctx)
	if err != nil {
		return err
	}
	servidor, err := api.Novo(rec, cliente, rotulos, log.Default())
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              ":" + cfg.PortaAPI,
		Handler:           servidor.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	erroServidor := make(chan error, 1)
	go func() { erroServidor <- srv.ListenAndServe() }()
	log.Printf("API ouvindo em http://localhost:%s (formulário em /)", cfg.PortaAPI)

	select {
	case err := <-erroServidor:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		log.Println("encerrando...")
		ctxFim, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelar()
		return srv.Shutdown(ctxFim)
	}
}
