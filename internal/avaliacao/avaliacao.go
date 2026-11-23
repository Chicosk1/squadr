package avaliacao

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
	"github.com/Chicosk1/squadr/internal/seed"
)

// Variantes da ablação (D-046), da mais simples para o motor completo. Todas recomendam
// entre os mesmos candidatos compatíveis; só muda o critério de ordenação.
const (
	Aleatorio  = "aleatório entre compatíveis"
	OrdemGrafo = "ordem do grafo"
	Unilateral = "score unilateral s(A→B)"
	Reciproco  = "recíproco, sem sinal social"
	Motor      = "motor completo"
)

var Variantes = []string{Aleatorio, OrdemGrafo, Unilateral, Reciproco, Motor}

// semLimite pede à fonte todos os candidatos que passam nos filtros.
const semLimite = 1_000_000

type Config struct {
	K       int
	Amostra int // quantidade de usuários avaliados, espaçados na base
	Opcoes  domain.OpcoesCandidatos
}

// Media acumula uma métrica no geral, por jogo e por persona.
type Media struct {
	soma    float64
	n       int
	porJogo map[domain.Jogo][2]float64 // soma, contagem
	porPers map[string][2]float64
}

func novaMedia() *Media {
	return &Media{porJogo: map[domain.Jogo][2]float64{}, porPers: map[string][2]float64{}}
}

func (m *Media) somar(jogo domain.Jogo, persona string, v float64) {
	m.soma += v
	m.n++
	j := m.porJogo[jogo]
	m.porJogo[jogo] = [2]float64{j[0] + v, j[1] + 1}
	p := m.porPers[persona]
	m.porPers[persona] = [2]float64{p[0] + v, p[1] + 1}
}

func (m *Media) Geral() float64 { return m.soma / float64(max(m.n, 1)) }

func (m *Media) Jogo(j domain.Jogo) float64 {
	v := m.porJogo[j]
	return v[0] / max(v[1], 1)
}

func (m *Media) Persona(p string) float64 {
	v := m.porPers[p]
	return v[0] / max(v[1], 1)
}

type Relatorio struct {
	K            int
	Consultas    int
	SemRelevante int // consultas sem nenhum relevante possível: ficam fora da média do NDCG
	Precisao     map[string]*Media
	NDCG         map[string]*Media
	Personas     []string
	ConsultasPer map[string]int
	// Cobertura do corte (D-029): relevantes compatíveis que entram nos candidatos pontuados.
	Corte                                    int
	RelevantesCompativeis, RelevantesNoCorte int
	Latencia                                 time.Duration // média de Pontuar (candidatos, perfis e scoring)
}

// Executar avalia o motor contra o gabarito. Os perfis são os regenerados pela semente do
// seed; o motor só lê a fonte. O gabarito (seed.Grau) usa a persona, que o motor nunca vê.
func Executar(ctx context.Context, fonte matching.Fonte, perfis []domain.Perfil, cfg Config) (*Relatorio, error) {
	if cfg.K < 1 || cfg.Amostra < 1 {
		return nil, fmt.Errorf("k e amostra devem ser pelo menos 1")
	}
	rec, err := matching.Novo(ctx, fonte, cfg.Opcoes)
	if err != nil {
		return nil, err
	}
	porID := make(map[string]domain.Perfil, len(perfis))
	for _, p := range perfis {
		porID[p.ID] = p
	}
	grau := func(eu domain.Perfil, id string, jogo domain.Jogo) int {
		outro, ok := porID[id]
		if !ok || !domain.PassaFiltros(eu, outro, jogo) {
			return 0
		}
		return seed.Grau(eu, outro, jogo)
	}
	graus := func(eu domain.Perfil, ids []string, jogo domain.Jogo) []int {
		g := make([]int, len(ids))
		for i, id := range ids {
			g[i] = grau(eu, id, jogo)
		}
		return g
	}

	r := &Relatorio{K: cfg.K, Corte: cfg.Opcoes.Limite, Precisao: map[string]*Media{}, NDCG: map[string]*Media{}, ConsultasPer: map[string]int{}}
	for _, v := range Variantes {
		r.Precisao[v], r.NDCG[v] = novaMedia(), novaMedia()
	}
	var tempo time.Duration

	passo := max(1, len(perfis)/cfg.Amostra)
	for i, n := 0, 0; i < len(perfis) && n < cfg.Amostra; i, n = i+passo, n+1 {
		eu := perfis[i]
		for _, pj := range eu.Jogos {
			jogo := pj.Jogo
			persona := pj.PersonaJogo
			if eu.Persona == seed.PersonaMulti {
				persona = seed.PersonaMulti
			}

			// Todos os candidatos que o sistema poderia recomendar: os filtros rígidos e o
			// mínimo de horários, sem corte. Servem ao ideal do NDCG e à linha de base aleatória.
			todos, err := fonte.BuscarCandidatos(ctx, eu.ID, jogo, domain.OpcoesCandidatos{
				Limite: semLimite, MinHorarios: cfg.Opcoes.MinHorarios, MaxDiffElo: cfg.Opcoes.MaxDiffElo})
			if err != nil {
				return nil, err
			}
			possiveis := graus(eu, ids(todos), jogo)

			inicio := time.Now()
			pontuados, err := rec.Pontuar(ctx, eu.ID, jogo)
			if err != nil {
				return nil, err
			}
			tempo += time.Since(inicio)

			grafo, err := fonte.BuscarCandidatos(ctx, eu.ID, jogo, domain.OpcoesCandidatos{
				Limite: cfg.K, MinHorarios: cfg.Opcoes.MinHorarios, MaxDiffElo: cfg.Opcoes.MaxDiffElo})
			if err != nil {
				return nil, err
			}

			listas := map[string][]int{
				OrdemGrafo: graus(eu, ids(grafo), jogo),
				Unilateral: graus(eu, ordenar(pontuados, func(x matching.Recomendacao) float64 { return x.ScoreAB }), jogo),
				Reciproco:  graus(eu, ordenar(pontuados, func(x matching.Recomendacao) float64 { return x.Reciproco }), jogo),
				Motor:      graus(eu, ordenar(pontuados, func(x matching.Recomendacao) float64 { return x.Score }), jogo),
			}

			r.Consultas++
			if r.ConsultasPer[persona] == 0 {
				r.Personas = append(r.Personas, persona)
			}
			r.ConsultasPer[persona]++

			pAle, nAle, temRelevante := AleatorioEsperado(possiveis, cfg.K)
			r.Precisao[Aleatorio].somar(jogo, persona, pAle)
			for _, v := range Variantes[1:] {
				r.Precisao[v].somar(jogo, persona, PrecisaoK(listas[v], cfg.K))
			}
			if !temRelevante {
				r.SemRelevante++
			} else {
				r.NDCG[Aleatorio].somar(jogo, persona, nAle)
				for _, v := range Variantes[1:] {
					ndcg, _ := NDCGK(listas[v], possiveis, cfg.K)
					r.NDCG[v].somar(jogo, persona, ndcg)
				}
			}

			for _, g := range possiveis {
				if g >= 1 {
					r.RelevantesCompativeis++
				}
			}
			for _, x := range pontuados {
				if grau(eu, x.ID, jogo) >= 1 {
					r.RelevantesNoCorte++
				}
			}
		}
	}
	slices.Sort(r.Personas)
	r.Latencia = tempo / time.Duration(max(r.Consultas, 1))
	return r, nil
}

func ids(cands []domain.Candidato) []string {
	saida := make([]string, len(cands))
	for i, c := range cands {
		saida[i] = c.ID
	}
	return saida
}

// ordenar devolve os ids em ordem decrescente do critério, com empate pelo id (como o motor).
func ordenar(recs []matching.Recomendacao, criterio func(matching.Recomendacao) float64) []string {
	copia := slices.Clone(recs)
	slices.SortFunc(copia, func(x, y matching.Recomendacao) int {
		if c := cmp.Compare(criterio(y), criterio(x)); c != 0 {
			return c
		}
		return cmp.Compare(x.ID, y.ID)
	})
	saida := make([]string, len(copia))
	for i, x := range copia {
		saida[i] = x.ID
	}
	return saida
}

// Imprimir escreve o relatório em tabelas Markdown, prontas para a monografia.
func (r *Relatorio) Imprimir(w io.Writer) {
	fmt.Fprintf(w, "Avaliação contra o gabarito das personas: %d consultas (usuário, jogo), K = %d\n\n", r.Consultas, r.K)
	tabela := func(titulo string, medias map[string]*Media) {
		fmt.Fprintf(w, "### %s\n\n| Variante | Geral |", titulo)
		for _, j := range domain.Jogos {
			fmt.Fprintf(w, " %s |", j)
		}
		fmt.Fprint(w, "\n|---|---|")
		for range domain.Jogos {
			fmt.Fprint(w, "---|")
		}
		fmt.Fprintln(w)
		for _, v := range Variantes {
			fmt.Fprintf(w, "| %s | %.3f |", v, medias[v].Geral())
			for _, j := range domain.Jogos {
				fmt.Fprintf(w, " %.3f |", medias[v].Jogo(j))
			}
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w)
	}
	tabela(fmt.Sprintf("NDCG@%d (graus 0, 1 e 2; %d consultas sem relevante possível ficam fora da média)", r.K, r.SemRelevante), r.NDCG)
	tabela(fmt.Sprintf("Precision@%d (relevante = grau ≥ 1)", r.K), r.Precisao)

	fmt.Fprintf(w, "### NDCG@%d por persona\n\n| Persona | Consultas | Aleatório | Ordem do grafo | Motor completo |\n|---|---|---|---|---|\n", r.K)
	for _, p := range r.Personas {
		fmt.Fprintf(w, "| %s | %d | %.3f | %.3f | %.3f |\n", p, r.ConsultasPer[p],
			r.NDCG[Aleatorio].Persona(p), r.NDCG[OrdemGrafo].Persona(p), r.NDCG[Motor].Persona(p))
	}
	fmt.Fprintln(w)
	if r.RelevantesCompativeis > 0 {
		fmt.Fprintf(w, "Corte de %d candidatos (D-029): %.1f%% dos relevantes compatíveis entram nos candidatos pontuados (%d de %d).\n",
			r.Corte, 100*float64(r.RelevantesNoCorte)/float64(r.RelevantesCompativeis), r.RelevantesNoCorte, r.RelevantesCompativeis)
	}
	fmt.Fprintf(w, "Latência média de Pontuar (candidatos, perfis e scoring): %s\n", r.Latencia.Round(time.Microsecond))
}
