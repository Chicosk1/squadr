package graph

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Chicosk1/squadr/internal/domain"
)

const (
	loteUsuarios   = 2000
	loteInteracoes = 5000
)

var queriesUsuarios = []string{
	`UNWIND $lote AS p
	 CREATE (u:Usuario)
	 SET u = p.props`,

	`UNWIND $lote AS p
	 UNWIND p.jogos AS jg
	 MATCH (u:Usuario {id: p.id})
	 MATCH (j:Jogo {id: jg.jogo})
	 CREATE (u)-[r:JOGA]->(j)
	 SET r = jg.props`,

	`UNWIND $lote AS p
	 UNWIND p.funcoes AS fn
	 MATCH (u:Usuario {id: p.id})
	 MATCH (f:Funcao {id: fn.id})
	 CREATE (u)-[:EXERCE {prioridade: fn.prioridade}]->(f)`,

	`UNWIND $lote AS p
	 UNWIND p.horarios AS hid
	 MATCH (u:Usuario {id: p.id})
	 MATCH (h:Horario {id: hid})
	 CREATE (u)-[:DISPONIVEL]->(h)`,
}

const queryCurtidas = `
	UNWIND $lote AS x
	MATCH (a:Usuario {id: x.de})
	MATCH (b:Usuario {id: x.para})
	CREATE (a)-[:CURTIU {jogo: x.jogo, em: datetime(x.em)}]->(b)`

const queryRecusas = `
	UNWIND $lote AS x
	MATCH (a:Usuario {id: x.de})
	MATCH (b:Usuario {id: x.para})
	CREATE (a)-[:RECUSOU {jogo: x.jogo, em: datetime(x.em)}]->(b)`

func (c *Cliente) LimparUsuarios(ctx context.Context) error {
	return c.ExecutarAutoCommit(ctx, `
		MATCH (u:Usuario)
		CALL { WITH u DETACH DELETE u } IN TRANSACTIONS OF 5000 ROWS`, nil)
}

func (c *Cliente) CarregarUsuarios(ctx context.Context, perfis []domain.Perfil) error {
	for inicio := 0; inicio < len(perfis); inicio += loteUsuarios {
		fim := min(inicio+loteUsuarios, len(perfis))
		params := map[string]any{"lote": montarLoteUsuarios(perfis[inicio:fim])}

		for _, query := range queriesUsuarios {
			if _, err := c.Executar(ctx, query, params); err != nil {
				return fmt.Errorf("carregando usuários %d a %d: %w", inicio, fim, err)
			}
		}
		log.Printf("usuários: %d/%d", fim, len(perfis))
	}
	return nil
}

func montarLoteUsuarios(perfis []domain.Perfil) []map[string]any {
	lote := make([]map[string]any, 0, len(perfis))
	for _, p := range perfis {
		var jogos, funcoes []map[string]any
		for _, pj := range p.Jogos {
			jogos = append(jogos, map[string]any{"jogo": string(pj.Jogo), "props": pj.Propriedades()})
			funcoes = append(funcoes,
				map[string]any{"id": pj.FuncaoPrincipal, "prioridade": 1},
				map[string]any{"id": pj.FuncaoSecundaria, "prioridade": 2})
		}
		lote = append(lote, map[string]any{
			"id":       p.ID,
			"props":    p.Propriedades(),
			"jogos":    jogos,
			"funcoes":  funcoes,
			"horarios": p.HorarioIDs(),
		})
	}
	return lote
}

func (c *Cliente) CarregarInteracoes(ctx context.Context, interacoes []domain.Interacao) error {
	var curtidas, recusas []map[string]any
	for _, x := range interacoes {
		item := map[string]any{
			"de": x.De, "para": x.Para, "jogo": string(x.Jogo),
			"em": x.Em.Format(time.RFC3339),
		}
		if x.Curtiu {
			curtidas = append(curtidas, item)
		} else {
			recusas = append(recusas, item)
		}
	}
	if err := c.carregarLotes(ctx, queryCurtidas, curtidas, "curtidas"); err != nil {
		return err
	}
	return c.carregarLotes(ctx, queryRecusas, recusas, "recusas")
}

func (c *Cliente) carregarLotes(ctx context.Context, query string, itens []map[string]any, rotulo string) error {
	for inicio := 0; inicio < len(itens); inicio += loteInteracoes {
		fim := min(inicio+loteInteracoes, len(itens))
		params := map[string]any{"lote": itens[inicio:fim]}
		if _, err := c.Executar(ctx, query, params); err != nil {
			return fmt.Errorf("carregando %s %d a %d: %w", rotulo, inicio, fim, err)
		}
	}
	log.Printf("%s: %d carregadas", rotulo, len(itens))
	return nil
}
