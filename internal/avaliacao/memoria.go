package avaliacao

import (
	"cmp"
	"context"
	"fmt"
	"math/bits"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

// FonteMemoria reproduz em Go a fonte de candidatos do grafo (D-029), para avaliar sem o Neo4j:
// filtros rígidos pelo oráculo domain.PassaFiltros, mínimo de horários em comum, exclusão de
// quem o usuário já avaliou no jogo, ordem por horários em comum, complemento e id, corte no
// limite e sinal social por travessia das curtidas. São as mesmas regras que o cmd/verificar
// confere contra o Cypher (D-031), com 0 divergências. Satisfaz matching.Fonte.
type FonteMemoria struct {
	perfis  []domain.Perfil
	porID   map[string]int
	pares   map[[2]string]float64
	matriz  matching.Matriz
	avaliou map[chave]map[string]struct{}
	curtiu  map[chave]map[string]struct{}
	curtido map[chave]map[string]struct{} // quem curtiu cada usuário, por jogo
}

type chave struct {
	id   string
	jogo domain.Jogo
}

var _ matching.Fonte = (*FonteMemoria)(nil)

func NovaFonteMemoria(perfis []domain.Perfil, interacoes []domain.Interacao, pares map[[2]string]float64) (*FonteMemoria, error) {
	m, err := matching.NovaMatriz(pares)
	if err != nil {
		return nil, err
	}
	f := &FonteMemoria{
		perfis: perfis, porID: make(map[string]int, len(perfis)), pares: pares, matriz: m,
		avaliou: map[chave]map[string]struct{}{}, curtiu: map[chave]map[string]struct{}{}, curtido: map[chave]map[string]struct{}{},
	}
	for i, p := range perfis {
		f.porID[p.ID] = i
	}
	for _, x := range interacoes {
		incluir(f.avaliou, chave{x.De, x.Jogo}, x.Para)
		if x.Curtiu {
			incluir(f.curtiu, chave{x.De, x.Jogo}, x.Para)
			incluir(f.curtido, chave{x.Para, x.Jogo}, x.De)
		}
	}
	return f, nil
}

func incluir(m map[chave]map[string]struct{}, k chave, v string) {
	if m[k] == nil {
		m[k] = map[string]struct{}{}
	}
	m[k][v] = struct{}{}
}

func (f *FonteMemoria) BuscarCandidatos(_ context.Context, id string, jogo domain.Jogo, op domain.OpcoesCandidatos) ([]domain.Candidato, error) {
	if op.MaxDiffElo != domain.MaxDiffEloDuo {
		return nil, fmt.Errorf("a fonte em memória usa o oráculo, que só conhece MaxDiffElo = %d", domain.MaxDiffEloDuo)
	}
	i, ok := f.porID[id]
	if !ok {
		return nil, nil // como o Cypher: usuário inexistente não tem candidatos
	}
	eu := f.perfis[i]
	meu, ok := eu.PerfilDoJogo(jogo)
	if !ok {
		return nil, nil
	}
	vistos := f.avaliou[chave{id, jogo}]

	var saida []domain.Candidato
	for _, outro := range f.perfis {
		if outro.ID == id {
			continue
		}
		if _, jaViu := vistos[outro.ID]; jaViu {
			continue
		}
		if !domain.PassaFiltros(eu, outro, jogo) {
			continue
		}
		horarios := bits.OnesCount32(eu.Disponibilidade & outro.Disponibilidade)
		if horarios < op.MinHorarios {
			continue
		}
		dele, _ := outro.PerfilDoJogo(jogo)
		saida = append(saida, domain.Candidato{ID: outro.ID, HorariosComuns: horarios, Complemento: f.complemento(meu, dele)})
	}
	slices.SortFunc(saida, func(a, b domain.Candidato) int {
		if c := cmp.Compare(b.HorariosComuns, a.HorariosComuns); c != 0 {
			return c
		}
		if c := cmp.Compare(b.Complemento, a.Complemento); c != 0 {
			return c
		}
		return cmp.Compare(a.ID, b.ID)
	})
	if op.Limite > 0 && len(saida) > op.Limite {
		saida = saida[:op.Limite]
	}
	if op.SinalSocial {
		sinais := f.sinalSocial(id, jogo)
		for k := range saida {
			saida[k].SinalSocial = sinais[saida[k].ID]
		}
	}
	return saida, nil
}

// complemento é o maior valor de COMPLEMENTA entre qualquer função de um e do outro, sem desconto.
func (f *FonteMemoria) complemento(a, b domain.PerfilJogo) float64 {
	melhor := 0.0
	for _, fa := range []string{a.FuncaoPrincipal, a.FuncaoSecundaria} {
		for _, fb := range []string{b.FuncaoPrincipal, b.FuncaoSecundaria} {
			if fa != "" && fb != "" {
				melhor = max(melhor, f.matriz.Valor(fa, fb))
			}
		}
	}
	return melhor
}

// sinalSocial segue a querySocial: para cada usuário parecido (que curtiu alguém que eu curti),
// soma ao alvo o número de curtidas em comum comigo, para cada alvo que ele curtiu.
func (f *FonteMemoria) sinalSocial(id string, jogo domain.Jogo) map[string]int {
	emComum := map[string]int{}
	for alvo := range f.curtiu[chave{id, jogo}] {
		for sim := range f.curtido[chave{alvo, jogo}] {
			if sim != id {
				emComum[sim]++
			}
		}
	}
	sinais := map[string]int{}
	for sim, n := range emComum {
		for alvo := range f.curtiu[chave{sim, jogo}] {
			sinais[alvo] += n
		}
	}
	return sinais
}

func (f *FonteMemoria) BuscarPerfis(_ context.Context, ids []string) ([]domain.Perfil, error) {
	saida := make([]domain.Perfil, 0, len(ids))
	for _, id := range ids {
		if i, ok := f.porID[id]; ok {
			saida = append(saida, f.perfis[i])
		}
	}
	return saida, nil
}

func (f *FonteMemoria) BuscarComplementaridade(context.Context) (map[[2]string]float64, error) {
	return f.pares, nil
}

var rePar = regexp.MustCompile(`\{a:\s*'([^']+)',\s*b:\s*'([^']+)',\s*valor:\s*([0-9.]+)\}`)

// LerComplementaridade lê os pares de COMPLEMENTA direto de scripts/referencia.cypher,
// a fonte única do grafo de referência (D-018), para a avaliação sem banco.
func LerComplementaridade(caminho string) (map[[2]string]float64, error) {
	texto, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	pares := map[[2]string]float64{}
	for _, m := range rePar.FindAllStringSubmatch(string(texto), -1) {
		v, err := strconv.ParseFloat(m[3], 64)
		if err != nil {
			return nil, fmt.Errorf("valor de complementaridade inválido %q: %w", m[3], err)
		}
		pares[[2]string{m[1], m[2]}] = v
	}
	if len(pares) == 0 {
		return nil, fmt.Errorf("nenhum par de COMPLEMENTA em %s", caminho)
	}
	return pares, nil
}
