package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log"
	"math/bits"
	"reflect"
	"slices"
	"time"

	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/Chicosk1/squadr/internal/seed"
)

type chaveInteracao struct {
	de   string
	jogo domain.Jogo
}

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	usuarios := flag.Int("usuarios", 10000, "deve ser igual ao usado no seed")
	semente := flag.Uint64("semente", 42, "deve ser igual à usada no seed")
	amostra := flag.Int("amostra", 200, "quantidade de usuários verificados")
	flag.Parse()
	if *amostra < 1 {
		return fmt.Errorf("-amostra deve ser pelo menos 1")
	}

	ctx := context.Background()
	cfg := config.Carregar()
	cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
	if err != nil {
		return err
	}
	defer cliente.Fechar(ctx)

	g := seed.NovoGerador(*semente)
	perfis := g.GerarUsuarios(*usuarios)
	interacoes := g.GerarInteracoes(perfis)

	if err := conferirTamanho(ctx, cliente, len(perfis)); err != nil {
		return err
	}
	escolhidos := escolherAmostra(perfis, *amostra)
	log.Printf("verificando %d usuários de %d", len(escolhidos), len(perfis))

	if err := verificarPerfis(ctx, cliente, escolhidos); err != nil {
		return err
	}
	if err := verificarCandidatos(ctx, cliente, perfis, interacoes, escolhidos); err != nil {
		return err
	}
	medirLatencia(ctx, cliente, escolhidos)
	return nil
}

func conferirTamanho(ctx context.Context, c *graph.Cliente, esperado int) error {
	res, err := c.Executar(ctx, "MATCH (u:Usuario) RETURN count(u) AS n", nil)
	if err != nil {
		return err
	}
	n, _ := res.Records[0].AsMap()["n"].(int64)
	if int(n) != esperado {
		return fmt.Errorf("o banco tem %d usuários, mas a geração local tem %d: rode o seed com os mesmos -usuarios e -semente", n, esperado)
	}
	return nil
}

func escolherAmostra(perfis []domain.Perfil, n int) []domain.Perfil {
	if n >= len(perfis) {
		return perfis
	}
	passo := len(perfis) / n
	var escolhidos []domain.Perfil
	for i := 0; i < len(perfis) && len(escolhidos) < n; i += passo {
		escolhidos = append(escolhidos, perfis[i])
	}
	return escolhidos
}

// ---------- 1. Ida e volta dos perfis ----------

func verificarPerfis(ctx context.Context, c *graph.Cliente, escolhidos []domain.Perfil) error {
	ids := make([]string, len(escolhidos))
	for i, p := range escolhidos {
		ids[i] = p.ID
	}
	lidos, err := c.BuscarPerfis(ctx, ids)
	if err != nil {
		return err
	}
	lidosPorID := make(map[string]domain.Perfil, len(lidos))
	for _, p := range lidos {
		lidosPorID[p.ID] = p
	}

	divergentes := 0
	for _, gerado := range escolhidos {
		gerado.Jogos = slices.Clone(gerado.Jogos)
		slices.SortFunc(gerado.Jogos, func(a, b domain.PerfilJogo) int { return cmp.Compare(a.Jogo, b.Jogo) })

		lido, existe := lidosPorID[gerado.ID]
		if !existe || !reflect.DeepEqual(gerado, lido) {
			divergentes++
			if divergentes <= 3 {
				log.Printf("DIVERGE %s\n  gerado: %+v\n  lido:   %+v", gerado.ID, gerado, lido)
			}
		}
	}
	log.Printf("perfis: %d verificados (ida e volta pelo grafo), %d divergentes", len(escolhidos), divergentes)
	if divergentes > 0 {
		return fmt.Errorf("%d perfis divergentes", divergentes)
	}
	return nil
}

// ---------- 2 a 4. Candidatos, sinais, ordenação ----------

type verificador struct {
	c         *graph.Cliente
	perfis    []domain.Perfil
	porID     map[string]domain.Perfil
	interagiu map[chaveInteracao]map[string]struct{} // likes e recusas
	curtiu    map[chaveInteracao]map[string]struct{} // só likes
}

func adicionar(m map[chaveInteracao]map[string]struct{}, k chaveInteracao, valor string) {
	if m[k] == nil {
		m[k] = map[string]struct{}{}
	}
	m[k][valor] = struct{}{}
}

func verificarCandidatos(ctx context.Context, c *graph.Cliente, perfis []domain.Perfil,
	interacoes []domain.Interacao, escolhidos []domain.Perfil) error {

	v := &verificador{
		c:         c,
		perfis:    perfis,
		porID:     make(map[string]domain.Perfil, len(perfis)),
		interagiu: map[chaveInteracao]map[string]struct{}{},
		curtiu:    map[chaveInteracao]map[string]struct{}{},
	}
	for _, p := range perfis {
		v.porID[p.ID] = p
	}
	for _, x := range interacoes {
		k := chaveInteracao{x.De, x.Jogo}
		adicionar(v.interagiu, k, x.Para)
		if x.Curtiu {
			adicionar(v.curtiu, k, x.Para)
		}
	}

	pares, divergentes := 0, 0
	for _, eu := range escolhidos {
		for _, pj := range eu.Jogos {
			pares++
			problemas, err := v.conferirPar(ctx, eu, pj.Jogo)
			if err != nil {
				return err
			}
			if len(problemas) > 0 {
				divergentes++
				if divergentes <= 5 {
					log.Printf("DIVERGE %s em %s:", eu.ID, pj.Jogo)
					for _, p := range problemas {
						log.Printf("  - %s", p)
					}
				}
			}
		}
	}
	log.Printf("candidatos: %d pares (usuário, jogo) verificados, %d divergentes", pares, divergentes)
	if divergentes > 0 {
		return fmt.Errorf("%d pares divergentes", divergentes)
	}
	return nil
}

func (v *verificador) conferirPar(ctx context.Context, eu domain.Perfil, jogo domain.Jogo) ([]string, error) {
	chave := chaveInteracao{eu.ID, jogo}

	esperado := map[string]struct{}{}
	for _, outro := range v.perfis {
		if outro.ID == eu.ID {
			continue
		}
		if _, jaViu := v.interagiu[chave][outro.ID]; jaViu {
			continue
		}
		if seed.PassaFiltros(eu, outro, jogo) {
			esperado[outro.ID] = struct{}{}
		}
	}

	op := graph.Opcoes{Limite: 1_000_000, MinHorarios: 0, MaxDiffElo: domain.MaxDiffEloDuo, SinalSocial: true}
	obtidos, err := v.c.BuscarCandidatos(ctx, eu.ID, jogo, op)
	if err != nil {
		return nil, err
	}
	obtido := make(map[string]graph.Candidato, len(obtidos))
	for _, cand := range obtidos {
		obtido[cand.ID] = cand
	}

	var problemas []string
	if f := faltantes(esperado, obtido); len(f) > 0 {
		problemas = append(problemas, fmt.Sprintf("%d candidatos do oráculo não vieram do grafo (ex.: %v)", len(f), primeiros(f)))
	}
	if s := faltantes(obtido, esperado); len(s) > 0 {
		problemas = append(problemas, fmt.Sprintf("%d candidatos do grafo não passam no oráculo (ex.: %v)", len(s), primeiros(s)))
	}

	sinais := v.sinalEsperado(eu.ID, jogo)
	errosHorario, errosSocial := 0, 0
	for _, cand := range obtidos {
		outro, ok := v.porID[cand.ID]
		if !ok {
			continue
		}
		if cand.HorariosComuns != bits.OnesCount32(eu.Disponibilidade&outro.Disponibilidade) {
			errosHorario++
		}
		if cand.SinalSocial != sinais[cand.ID] {
			errosSocial++
		}
	}
	if errosHorario > 0 {
		problemas = append(problemas, fmt.Sprintf("%d candidatos com horários em comum errados", errosHorario))
	}
	if errosSocial > 0 {
		problemas = append(problemas, fmt.Sprintf("%d candidatos com sinal social errado", errosSocial))
	}

	curto, err := v.c.BuscarCandidatos(ctx, eu.ID, jogo,
		graph.Opcoes{Limite: 10, MinHorarios: 0, MaxDiffElo: domain.MaxDiffEloDuo})
	if err != nil {
		return nil, err
	}
	if len(curto) != min(10, len(obtidos)) {
		problemas = append(problemas, fmt.Sprintf("LIMIT 10 devolveu %d candidatos", len(curto)))
	} else {
		for i, cand := range curto {
			if cand.ID != obtidos[i].ID {
				problemas = append(problemas, fmt.Sprintf("posição %d: esperado %s, veio %s", i, obtidos[i].ID, cand.ID))
				break
			}
		}
	}
	return problemas, nil
}

func (v *verificador) sinalEsperado(id string, jogo domain.Jogo) map[string]int {
	meus := v.curtiu[chaveInteracao{id, jogo}]
	sinais := map[string]int{}
	for _, sim := range v.perfis {
		if sim.ID == id {
			continue
		}
		deles := v.curtiu[chaveInteracao{sim.ID, jogo}]
		emComum := 0
		for alvo := range meus {
			if _, ok := deles[alvo]; ok {
				emComum++
			}
		}
		if emComum == 0 {
			continue
		}
		for alvo := range deles {
			sinais[alvo] += emComum
		}
	}
	return sinais
}

func faltantes[A, B any](a map[string]A, b map[string]B) []string {
	var saida []string
	for k := range a {
		if _, ok := b[k]; !ok {
			saida = append(saida, k)
		}
	}
	slices.Sort(saida)
	return saida
}

func primeiros(ids []string) []string { return ids[:min(3, len(ids))] }

// ---------- Latência ----------

func medirLatencia(ctx context.Context, c *graph.Cliente, escolhidos []domain.Perfil) {
	op := graph.OpcoesPadrao()
	var duracoes []time.Duration
	for _, eu := range escolhidos {
		for _, pj := range eu.Jogos {
			inicio := time.Now()
			if _, err := c.BuscarCandidatos(ctx, eu.ID, pj.Jogo, op); err != nil {
				log.Printf("erro medindo %s em %s: %v", eu.ID, pj.Jogo, err)
				continue
			}
			duracoes = append(duracoes, time.Since(inicio))
		}
	}
	if len(duracoes) == 0 {
		return
	}
	slices.Sort(duracoes)
	pct := func(p float64) time.Duration { return duracoes[int(float64(len(duracoes)-1)*p)] }
	log.Printf("latência de BuscarCandidatos com as opções padrão (%d consultas): p50=%s p95=%s máx=%s",
		len(duracoes), pct(0.50), pct(0.95), duracoes[len(duracoes)-1])
}
