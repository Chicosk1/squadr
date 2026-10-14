package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/Chicosk1/squadr/internal/seed"
)

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	usuarios := flag.Int("usuarios", 10000, "quantidade de usuários a gerar (mínimo 100)")
	semente := flag.Uint64("semente", 42, "semente do gerador: a mesma semente gera sempre os mesmos dados")
	flag.Parse()

	if *usuarios < 100 {
		return fmt.Errorf("-usuarios deve ser pelo menos 100, recebido %d", *usuarios)
	}

	ctx := context.Background()
	cfg := config.Carregar()
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(ctx)

	inicio := time.Now()
	g := seed.NovoGerador(*semente)
	perfis := g.GerarUsuarios(*usuarios)
	interacoes := g.GerarInteracoes(perfis)
	log.Printf("gerados %d usuários e %d interações em %s (semente %d)",
		len(perfis), len(interacoes), time.Since(inicio).Round(time.Millisecond), *semente)

	contagem := map[string]int{}
	for _, p := range perfis {
		contagem[p.Persona]++
	}
	for _, persona := range seed.Personas {
		log.Printf("  %-20s %d", persona.Nome, contagem[persona.Nome])
	}

	log.Println("apagando usuários anteriores (os dados de referência são mantidos)...")
	if err := cliente.LimparUsuarios(ctx); err != nil {
		return err
	}

	inicioCarga := time.Now()
	if err := cliente.CarregarUsuarios(ctx, perfis); err != nil {
		return err
	}
	if err := cliente.CarregarInteracoes(ctx, interacoes); err != nil {
		return err
	}
	log.Printf("carga concluída em %s", time.Since(inicioCarga).Round(time.Millisecond))
	return nil
}
