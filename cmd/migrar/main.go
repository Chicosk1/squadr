package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	ctx, cancelar := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelar()

	cfg := config.Carregar()
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(ctx)

	for _, script := range []string{"scripts/schema.cypher", "scripts/referencia.cypher"} {
		log.Printf("executando %s", script)
		if err := cliente.ExecutarScript(ctx, script); err != nil {
			return err
		}
	}
	return verificar(ctx, cliente)
}

func verificar(ctx context.Context, c *graph.Cliente) error {
	const query = `
		CALL { MATCH (j:Jogo) RETURN count(j) AS jogos }
		CALL { MATCH (h:Horario) RETURN count(h) AS horarios }
		CALL { MATCH (f:Funcao) RETURN count(f) AS funcoes }
		CALL { MATCH ()-[c:COMPLEMENTA]->() RETURN count(c) AS complementa }
		RETURN jogos, horarios, funcoes, complementa`

	res, err := c.Executar(ctx, query, nil)
	if err != nil {
		return err
	}
	if len(res.Records) != 1 {
		return fmt.Errorf("verificação devolveu %d registros, esperado 1", len(res.Records))
	}

	esperado := map[string]int64{"jogos": 3, "horarios": 28, "funcoes": 16, "complementa": 51}
	for chave, quero := range esperado {
		obtido, _, err := neo4j.GetRecordValue[int64](res.Records[0], chave)
		if err != nil {
			return err
		}
		if obtido != quero {
			return fmt.Errorf("%s: esperado %d, obtido %d", chave, quero, obtido)
		}
		log.Printf("ok: %s = %d", chave, obtido)
	}
	return nil
}
