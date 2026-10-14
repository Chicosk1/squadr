package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Cliente struct {
	driver neo4j.DriverWithContext
	banco  string
}

func Conectar(ctx context.Context, uri, usuario, senha, banco string) (*Cliente, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(usuario, senha, ""))
	if err != nil {
		return nil, fmt.Errorf("criando driver: %w", err)
	}

	ctxVerif, cancelar := context.WithTimeout(ctx, 5*time.Second)
	defer cancelar()
	if err := driver.VerifyConnectivity(ctxVerif); err != nil {
		_ = driver.Close(ctx)
		return nil, fmt.Errorf("conectando ao Neo4j em %s: %w", uri, err)
	}

	return &Cliente{driver: driver, banco: banco}, nil
}

func (c *Cliente) Fechar(ctx context.Context) error {
	return c.driver.Close(ctx)
}

func (c *Cliente) Ping(ctx context.Context) error {
	return c.driver.VerifyConnectivity(ctx)
}

func (c *Cliente) Executar(ctx context.Context, query string, params map[string]any) (*neo4j.EagerResult, error) {
	return neo4j.ExecuteQuery(ctx, c.driver, query, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(c.banco))
}
