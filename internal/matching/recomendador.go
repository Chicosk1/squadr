package matching

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/Chicosk1/squadr/internal/domain"
)

// Fonte é tudo que o recomendador precisa do banco (D-013, D-034). O *graph.Cliente
// satisfaz esta interface sem adaptador; os testes usam uma fonte em memória.
type Fonte interface {
	BuscarCandidatos(ctx context.Context, id string, jogo domain.Jogo, op domain.OpcoesCandidatos) ([]domain.Candidato, error)
	BuscarPerfis(ctx context.Context, ids []string) ([]domain.Perfil, error)
	BuscarComplementaridade(ctx context.Context) (map[[2]string]float64, error)
}

var (
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
	ErrNaoJoga              = errors.New("usuário não joga este jogo")
)

const (
	// PesoReciproco e PesoSocial combinam o score recíproco com o sinal social normalizado (D-036).
	PesoReciproco = 0.8
	PesoSocial    = 0.2
	// NumMotivos é quantos componentes aparecem como motivo do match.
	NumMotivos = 2
)

type Recomendador struct {
	fonte  Fonte
	matriz Matriz
	opcoes domain.OpcoesCandidatos
}

// Novo carrega a matriz de complementaridade uma vez. Crie um Recomendador por processo.
func Novo(ctx context.Context, fonte Fonte, op domain.OpcoesCandidatos) (*Recomendador, error) {
	pares, err := fonte.BuscarComplementaridade(ctx)
	if err != nil {
		return nil, err
	}
	if len(pares) == 0 {
		return nil, errors.New("matriz de complementaridade vazia: rode go run ./cmd/migrar")
	}
	m, err := NovaMatriz(pares)
	if err != nil {
		return nil, err
	}
	return &Recomendador{fonte: fonte, matriz: m, opcoes: op}, nil
}

func (r *Recomendador) Matriz() Matriz { return r.matriz }

type Recomendacao struct {
	ID        string
	Score     float64 // PesoReciproco·Reciproco + PesoSocial·Social
	Reciproco float64
	ScoreAB   float64 // s(A→B): quanto B serve para A
	ScoreBA   float64 // s(B→A): quanto A serve para B
	Social    float64 // sinal social normalizado em [0, 1]
	DeA, DeB  Componentes
	Motivos   []string
}

// Recomendar devolve os k melhores parceiros para o usuário id no jogo, em ordem decrescente de score.
func (r *Recomendador) Recomendar(ctx context.Context, id string, jogo domain.Jogo, k int) ([]Recomendacao, error) {
	todas, err := r.Pontuar(ctx, id, jogo)
	if err != nil {
		return nil, err
	}
	return todas[:min(k, len(todas))], nil
}

// Pontuar devolve todos os candidatos pontuados e ordenados. A avaliação usa o ranking inteiro.
func (r *Recomendador) Pontuar(ctx context.Context, id string, jogo domain.Jogo) ([]Recomendacao, error) {
	candidatos, err := r.fonte.BuscarCandidatos(ctx, id, jogo, r.opcoes)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(candidatos)+1)
	ids = append(ids, id)
	for _, c := range candidatos {
		ids = append(ids, c.ID)
	}
	lidos, err := r.fonte.BuscarPerfis(ctx, ids)
	if err != nil {
		return nil, err
	}
	perfis := make(map[string]domain.Perfil, len(lidos))
	for _, p := range lidos {
		perfis[p.ID] = p
	}

	a, ok := perfis[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUsuarioNaoEncontrado, id)
	}
	pa, ok := a.PerfilDoJogo(jogo)
	if !ok {
		return nil, fmt.Errorf("%w: %s em %s", ErrNaoJoga, id, jogo)
	}

	maxSocial := 0
	for _, c := range candidatos {
		maxSocial = max(maxSocial, c.SinalSocial)
	}

	recs := make([]Recomendacao, 0, len(candidatos))
	for _, c := range candidatos {
		b, ok := perfis[c.ID]
		if !ok {
			return nil, fmt.Errorf("perfil do candidato %s não encontrado", c.ID)
		}
		pb, ok := b.PerfilDoJogo(jogo)
		if !ok {
			return nil, fmt.Errorf("candidato %s não joga %s", c.ID, jogo)
		}
		rec := Recomendacao{ID: c.ID, Social: normalizarSocial(c.SinalSocial, maxSocial)}
		rec.ScoreAB, rec.DeA = Unilateral(a, b, pa, pb, r.matriz)
		rec.ScoreBA, rec.DeB = Unilateral(b, a, pb, pa, r.matriz)
		rec.Reciproco = Reciproco(rec.ScoreAB, rec.ScoreBA)
		rec.Score = PesoReciproco*rec.Reciproco + PesoSocial*rec.Social
		rec.Motivos = motivos(rec.DeA, rec.DeB)
		recs = append(recs, rec)
	}

	slices.SortFunc(recs, func(x, y Recomendacao) int {
		if c := cmp.Compare(y.Score, x.Score); c != 0 {
			return c
		}
		return cmp.Compare(x.ID, y.ID)
	})
	return recs, nil
}

// normalizarSocial é ln(1+s) / ln(1+máximo) entre os candidatos da mesma consulta (D-036).
// O log amortece candidatos com sinal muito alto. Sem nenhum sinal, vale 0 para todos.
func normalizarSocial(s, maximo int) float64 {
	if maximo <= 0 || s <= 0 {
		return 0
	}
	return math.Log1p(float64(s)) / math.Log1p(float64(maximo))
}

// motivos escolhe os componentes em que o par vai melhor para os dois lados ao mesmo tempo
// (o menor dos dois valores), para não explicar o match com algo que só agrada a um lado.
func motivos(deA, deB Componentes) []string {
	ordem := make([]Componente, NumComponentes)
	for i := range ordem {
		ordem[i] = Componente(i)
	}
	slices.SortStableFunc(ordem, func(x, y Componente) int {
		return cmp.Compare(min(deA[y], deB[y]), min(deA[x], deB[x]))
	})
	saida := make([]string, 0, NumMotivos)
	for _, c := range ordem[:NumMotivos] {
		if min(deA[c], deB[c]) > 0 {
			saida = append(saida, c.String())
		}
	}
	return saida
}
