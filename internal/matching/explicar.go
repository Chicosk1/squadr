package matching

import (
	"context"
	"fmt"
	"math/bits"

	"github.com/Chicosk1/squadr/internal/domain"
)

// Explicacao detalha um par para a tela de comparação (D-043): os filtros rígidos com o
// motivo de cada um, o score nos dois sentidos e a posição de B no ranking de A.
type Explicacao struct {
	A, B domain.Perfil
	// Filtros traz os de domain.Filtros e, por último, o mínimo de horários em comum
	// que a consulta de candidatos exige (OpcoesCandidatos.MinHorarios).
	Filtros        []domain.Filtro
	PassaFiltros   bool
	HorariosComuns uint32 // grade de 28 bits com os horários dos dois

	// Calculado é falso quando B não joga o jogo: aí não existe score.
	// Quando B joga mas é barrado por um filtro, o score é calculado mesmo assim,
	// para mostrar o que o motor daria se o par chegasse até ele.
	Calculado bool
	ScoreAB   float64 // s(A→B)
	ScoreBA   float64 // s(B→A)
	Reciproco float64
	DeA, DeB  Componentes
	Motivos   []string

	// Posicao é a posição de B no ranking de A (1 = primeiro), ou 0 quando B não está nele:
	// foi barrado por um filtro, já foi avaliado por A neste jogo ou ficou fora do corte
	// de candidatos (D-029).
	Posicao    int
	Candidatos int     // tamanho do ranking de A
	Score      float64 // score final no ranking, com o sinal social; 0 se B não está nele
	Social     float64
}

// Explicar monta a explicação do par (idA, idB) no jogo. A precisa jogar o jogo
// (ErrNaoJoga); os dois precisam existir (ErrUsuarioNaoEncontrado).
func (r *Recomendador) Explicar(ctx context.Context, idA, idB string, jogo domain.Jogo) (Explicacao, error) {
	lidos, err := r.fonte.BuscarPerfis(ctx, []string{idA, idB})
	if err != nil {
		return Explicacao{}, err
	}
	var e Explicacao
	var achouA, achouB bool
	for _, p := range lidos {
		switch p.ID {
		case idA:
			e.A, achouA = p, true
		case idB:
			e.B, achouB = p, true
		}
	}
	if !achouA || !achouB {
		ausente := idA
		if achouA {
			ausente = idB
		}
		return Explicacao{}, fmt.Errorf("%w: %s", ErrUsuarioNaoEncontrado, ausente)
	}
	pa, ok := e.A.PerfilDoJogo(jogo)
	if !ok {
		return Explicacao{}, fmt.Errorf("%w: %s em %s", ErrNaoJoga, idA, jogo)
	}

	e.Filtros = domain.Filtros(e.A, e.B, jogo)
	e.HorariosComuns = e.A.Disponibilidade & e.B.Disponibilidade
	if minimo := r.opcoes.MinHorarios; minimo > 0 {
		comuns := bits.OnesCount32(e.HorariosComuns)
		motivo := "nenhum horário em comum"
		if comuns > 0 {
			motivo = fmt.Sprintf("%d horários em comum, o mínimo é %d", comuns, minimo)
		}
		f := domain.Filtro{Nome: "Horário em comum", Passou: comuns >= minimo}
		if !f.Passou {
			f.Motivo = motivo
		}
		e.Filtros = append(e.Filtros, f)
	}
	e.PassaFiltros = true
	for _, f := range e.Filtros {
		e.PassaFiltros = e.PassaFiltros && f.Passou
	}

	if pb, ok := e.B.PerfilDoJogo(jogo); ok {
		e.Calculado = true
		e.ScoreAB, e.DeA = Unilateral(e.A, e.B, pa, pb, r.matriz)
		e.ScoreBA, e.DeB = Unilateral(e.B, e.A, pb, pa, r.matriz)
		e.Reciproco = Reciproco(e.ScoreAB, e.ScoreBA)
		e.Motivos = motivos(e.DeA, e.DeB)
	}

	ranking, err := r.Pontuar(ctx, idA, jogo)
	if err != nil {
		return Explicacao{}, err
	}
	e.Candidatos = len(ranking)
	for i, rec := range ranking {
		if rec.ID == idB {
			e.Posicao, e.Score, e.Social = i+1, rec.Score, rec.Social
			break
		}
	}
	return e, nil
}
