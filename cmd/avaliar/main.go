// Comando avaliar mede a qualidade do ranking contra o gabarito das personas.
// Versão simples da Fase 5: Precision@K do motor completo e da ordem do grafo
// (a ablação completa e o NDCG ficam para a Fase 7).
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/Chicosk1/squadr/internal/matching"
	"github.com/Chicosk1/squadr/internal/seed"
)

// O *graph.Cliente precisa satisfazer a fonte do matching sem adaptador (D-034).
var _ matching.Fonte = (*graph.Cliente)(nil)

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

type acumulado struct {
	soma  float64
	n     int
	porJ  map[domain.Jogo]float64
	nPorJ map[domain.Jogo]int
}

func novoAcumulado() *acumulado {
	return &acumulado{porJ: map[domain.Jogo]float64{}, nPorJ: map[domain.Jogo]int{}}
}

func (a *acumulado) somar(jogo domain.Jogo, v float64) {
	a.soma += v
	a.n++
	a.porJ[jogo] += v
	a.nPorJ[jogo]++
}

func (a *acumulado) linha(nome string) string {
	s := fmt.Sprintf("%-16s geral=%.3f", nome, a.soma/float64(max(a.n, 1)))
	for _, j := range domain.Jogos {
		s += fmt.Sprintf("  %s=%.3f", j, a.porJ[j]/float64(max(a.nPorJ[j], 1)))
	}
	return s
}

func executar() error {
	usuarios := flag.Int("usuarios", 10000, "deve ser igual ao usado no seed")
	semente := flag.Uint64("semente", 42, "deve ser igual à usada no seed")
	amostra := flag.Int("amostra", 200, "quantidade de usuários avaliados")
	k := flag.Int("k", 10, "tamanho do ranking avaliado")
	flag.Parse()
	if *amostra < 1 || *k < 1 {
		return fmt.Errorf("-amostra e -k devem ser pelo menos 1")
	}

	ctx := context.Background()
	cfg := config.Carregar()
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(ctx)

	// O gabarito vem dos perfis regenerados pela semente; o motor lê só o grafo.
	perfis := seed.NovoGerador(*semente).GerarUsuarios(*usuarios)
	porID := make(map[string]domain.Perfil, len(perfis))
	for _, p := range perfis {
		porID[p.ID] = p
	}

	op := domain.OpcoesPadrao()
	rec, err := matching.Novo(ctx, cliente, op)
	if err != nil {
		return err
	}

	relevante := func(eu domain.Perfil, id string, jogo domain.Jogo) bool {
		outro, ok := porID[id]
		return ok && domain.PassaFiltros(eu, outro, jogo) && seed.Grau(eu, outro, jogo) >= 1
	}
	precisao := func(eu domain.Perfil, ids []string, jogo domain.Jogo) float64 {
		acertos := 0
		for _, id := range ids {
			if relevante(eu, id, jogo) {
				acertos++
			}
		}
		return float64(acertos) / float64(*k)
	}

	motor, grafo := novoAcumulado(), novoAcumulado()
	var tempo time.Duration
	consultas := 0
	passo := max(1, len(perfis) / *amostra)
	for i, n := 0, 0; i < len(perfis) && n < *amostra; i, n = i+passo, n+1 {
		eu := perfis[i]
		for _, pj := range eu.Jogos {
			inicio := time.Now()
			recs, err := rec.Recomendar(ctx, eu.ID, pj.Jogo, *k)
			if err != nil {
				return err
			}
			tempo += time.Since(inicio)
			consultas++

			ids := make([]string, len(recs))
			for j, r := range recs {
				ids[j] = r.ID
			}
			motor.somar(pj.Jogo, precisao(eu, ids, pj.Jogo))

			// Linha de base: os K primeiros na ordem do grafo (horários, complemento, id).
			cands, err := cliente.BuscarCandidatos(ctx, eu.ID, pj.Jogo, domain.OpcoesCandidatos{
				Limite: *k, MinHorarios: op.MinHorarios, MaxDiffElo: op.MaxDiffElo})
			if err != nil {
				return err
			}
			idsGrafo := make([]string, len(cands))
			for j, c := range cands {
				idsGrafo[j] = c.ID
			}
			grafo.somar(pj.Jogo, precisao(eu, idsGrafo, pj.Jogo))
		}
	}

	log.Printf("Precision@%d em %d consultas (usuário, jogo):", *k, consultas)
	log.Print(grafo.linha("ordem do grafo"))
	log.Print(motor.linha("motor completo"))
	log.Printf("latência média de Recomendar: %s", (tempo / time.Duration(max(consultas, 1))).Round(time.Microsecond))
	return nil
}
