package graph

import (
	"context"
	"fmt"
	"math/bits"
	"time"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ---------- Criação de um usuário (formulário) ----------

const queryConferirUsuario = `
MATCH (u:Usuario {id: $id})
RETURN COUNT { (u)-[:JOGA]->(:Jogo) } AS jogos,
       COUNT { (u)-[:EXERCE]->(:Funcao) } AS funcoes,
       COUNT { (u)-[:DISPONIVEL]->(:Horario) } AS horarios`

// CriarUsuario grava um perfil em uma única transação (D-040). As mesmas consultas da carga
// em lote rodam com um lote de um usuário. Antes do commit, a transação confere quantas
// relações foram criadas: um MATCH que não acha o nó não dá erro, só deixa de criar a relação.
// Se faltar alguma, nada é gravado.
func (c *Cliente) CriarUsuario(ctx context.Context, p domain.Perfil) error {
	params := map[string]any{"lote": montarLoteUsuarios([]domain.Perfil{p})}
	esperado := map[string]int{"jogos": len(p.Jogos), "horarios": bits.OnesCount32(p.Disponibilidade)}
	for _, pj := range p.Jogos {
		esperado["funcoes"]++
		if pj.FuncaoSecundaria != "" {
			esperado["funcoes"]++
		}
	}

	sessao := c.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: c.banco, AccessMode: neo4j.AccessModeWrite})
	defer sessao.Close(ctx)

	_, err := sessao.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		for _, query := range queriesUsuarios {
			res, err := tx.Run(ctx, query, params)
			if err != nil {
				return nil, err
			}
			if _, err := res.Consume(ctx); err != nil {
				return nil, err
			}
		}
		res, err := tx.Run(ctx, queryConferirUsuario, map[string]any{"id": p.ID})
		if err != nil {
			return nil, err
		}
		rec, err := res.Single(ctx)
		if err != nil {
			return nil, err
		}
		m := rec.AsMap()
		for _, chave := range []string{"jogos", "funcoes", "horarios"} {
			if inteiro(m, chave) != esperado[chave] {
				return nil, fmt.Errorf("%s: %d relações criadas, esperadas %d (rode go run ./cmd/migrar)",
					chave, inteiro(m, chave), esperado[chave])
			}
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("criando usuário %s: %w", p.ID, err)
	}
	return nil
}

// ---------- Curtidas e recusas ----------

// Os dois usuários precisam jogar o jogo. vistos conta o que A já fez com B neste jogo;
// só cria a relação se for zero. O tipo da relação não pode ser parâmetro em Cypher,
// então há uma consulta fixa para cada tipo (como os filtros da D-030).
const queryInteracao = `
MATCH (a:Usuario {id: $de})-[:JOGA]->(:Jogo {id: $jogo})<-[:JOGA]-(b:Usuario {id: $para})
WHERE a <> b
WITH a, b, COUNT { (a)-[x:CURTIU|RECUSOU]->(b) WHERE x.jogo = $jogo } AS vistos
FOREACH (_ IN CASE WHEN vistos = 0 THEN [1] ELSE [] END |
  CREATE (a)-[:%s {jogo: $jogo, em: datetime($em)}]->(b))
RETURN vistos,
       EXISTS { (b)-[x:CURTIU]->(a) WHERE x.jogo = $jogo } AS reciproca`

var queriesInteracao = map[bool]string{
	true:  fmt.Sprintf(queryInteracao, "CURTIU"),
	false: fmt.Sprintf(queryInteracao, "RECUSOU"),
}

// RegistrarInteracao grava uma curtida ou uma recusa. Devolve match = true quando é uma
// curtida e B já tinha curtido A neste jogo. Erros: domain.ErrNaoEncontrado se algum dos
// dois não existe ou não joga o jogo (ou se de == para); domain.ErrJaAvaliado se A já curtiu
// ou recusou B neste jogo.
func (c *Cliente) RegistrarInteracao(ctx context.Context, x domain.Interacao) (bool, error) {
	if x.Em.IsZero() {
		x.Em = time.Now()
	}
	res, err := c.Executar(ctx, queriesInteracao[x.Curtiu], map[string]any{
		"de": x.De, "para": x.Para, "jogo": string(x.Jogo), "em": x.Em.Format(time.RFC3339),
	})
	if err != nil {
		return false, fmt.Errorf("registrando interação de %s para %s: %w", x.De, x.Para, err)
	}
	if len(res.Records) == 0 {
		return false, fmt.Errorf("%w: %s ou %s em %s", domain.ErrNaoEncontrado, x.De, x.Para, x.Jogo)
	}
	m := res.Records[0].AsMap()
	if inteiro(m, "vistos") > 0 {
		return false, fmt.Errorf("%w: %s → %s", domain.ErrJaAvaliado, x.De, x.Para)
	}
	return x.Curtiu && booleano(m, "reciproca"), nil
}

// ---------- Rótulos das funções ----------

// BuscarRotulosFuncoes devolve id → nome de cada função, para o formulário (os rótulos
// em português vivem no grafo, D-018 e D-033).
func (c *Cliente) BuscarRotulosFuncoes(ctx context.Context) (map[string]string, error) {
	res, err := c.Executar(ctx, `MATCH (f:Funcao) RETURN f.id AS id, f.nome AS nome`, nil)
	if err != nil {
		return nil, fmt.Errorf("buscando funções: %w", err)
	}
	rotulos := make(map[string]string, len(res.Records))
	for _, rec := range res.Records {
		m := rec.AsMap()
		rotulos[texto(m, "id")] = texto(m, "nome")
	}
	return rotulos, nil
}
