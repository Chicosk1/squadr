package graph

import (
	"context"
	"fmt"
)

const queryComplementaridade = `
MATCH (a:Funcao)-[c:COMPLEMENTA]->(b:Funcao)
RETURN a.id AS a, b.id AS b, c.valor AS valor`

// BuscarComplementaridade devolve os pares de COMPLEMENTA na direção em que foram gravados.
// A relação é lida sem direção, então quem consome deve tratar (a, b) e (b, a) como iguais.
func (c *Cliente) BuscarComplementaridade(ctx context.Context) (map[[2]string]float64, error) {
	res, err := c.Executar(ctx, queryComplementaridade, nil)
	if err != nil {
		return nil, fmt.Errorf("buscando complementaridade: %w", err)
	}
	pares := make(map[[2]string]float64, len(res.Records))
	for _, rec := range res.Records {
		m := rec.AsMap()
		pares[[2]string{texto(m, "a"), texto(m, "b")}] = decimal(m, "valor")
	}
	return pares, nil
}
