package avaliacao

import (
	"context"
	"math"
	"math/bits"
	"testing"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/seed"
)

func quase(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestPrecisaoK(t *testing.T) {
	casos := []struct {
		graus []int
		k     int
		quer  float64
	}{
		{[]int{2, 1, 0, 0}, 4, 0.5},
		{[]int{2, 2}, 4, 0.5}, // posições vazias contam como erro
		{nil, 10, 0},
		{[]int{1, 1, 1, 0, 0}, 3, 1},
	}
	for _, c := range casos {
		if got := PrecisaoK(c.graus, c.k); !quase(got, c.quer) {
			t.Errorf("PrecisaoK(%v, %d) = %v, quer %v", c.graus, c.k, got, c.quer)
		}
	}
}

func TestNDCG(t *testing.T) {
	// DCG de [2, 0, 1]: 2/log2(2) + 0 + 1/log2(4) = 2 + 0,5
	if got := DCGK([]int{2, 0, 1}, 3); !quase(got, 2.5) {
		t.Errorf("DCG = %v, quer 2,5", got)
	}
	possiveis := []int{0, 1, 2, 0, 1}
	// A ordem ideal vale 1; a pior ordem dos mesmos itens vale menos; a lista vazia vale 0.
	if n, ok := NDCGK([]int{2, 1, 1}, possiveis, 3); !ok || !quase(n, 1) {
		t.Errorf("ideal: NDCG = %v (ok=%v), quer 1", n, ok)
	}
	if n, _ := NDCGK([]int{0, 0, 1}, possiveis, 3); n <= 0 || n >= 1 {
		t.Errorf("lista ruim: NDCG = %v, quer entre 0 e 1", n)
	}
	if n, _ := NDCGK(nil, possiveis, 3); n != 0 {
		t.Errorf("lista vazia: NDCG = %v, quer 0", n)
	}
	// O grau 2 no topo vale mais que o grau 1 no topo: a relevância é graduada.
	a, _ := NDCGK([]int{2, 1}, possiveis, 2)
	b, _ := NDCGK([]int{1, 2}, possiveis, 2)
	if a <= b {
		t.Errorf("grau 2 no topo (%v) deveria valer mais que grau 1 no topo (%v)", a, b)
	}
	if _, ok := NDCGK([]int{0, 0}, []int{0, 0, 0}, 2); ok {
		t.Error("sem relevante possível, NDCG deveria ficar fora da média (ok = false)")
	}
}

// A esperança analítica do sorteio tem de bater com a média de todos os sorteios possíveis.
func TestAleatorioEsperado(t *testing.T) {
	possiveis := []int{2, 1, 0, 0, 1}
	k := 3
	var somaP, somaN float64
	total := 0
	var sortear func(lista []int, usados uint)
	sortear = func(lista []int, usados uint) {
		if len(lista) == k {
			somaP += PrecisaoK(lista, k)
			n, _ := NDCGK(lista, possiveis, k)
			somaN += n
			total++
			return
		}
		for i, g := range possiveis {
			if usados&(1<<i) == 0 {
				sortear(append(lista, g), usados|1<<i)
			}
		}
	}
	sortear(nil, 0)
	p, n, ok := AleatorioEsperado(possiveis, k)
	if !ok || !quase(p, somaP/float64(total)) || !quase(n, somaN/float64(total)) {
		t.Errorf("esperado P=%v N=%v; média dos %d sorteios P=%v N=%v", p, n, total, somaP/float64(total), somaN/float64(total))
	}
	// Com menos candidatos que k, as posições vazias contam como erro no Precision.
	if p, _, _ := AleatorioEsperado([]int{1, 1}, 4); !quase(p, 0.5) {
		t.Errorf("dois relevantes para k = 4: P = %v, quer 0,5", p)
	}
}

func novaFonte(t *testing.T, n int) (*FonteMemoria, []domain.Perfil, []domain.Interacao) {
	t.Helper()
	g := seed.NovoGerador(7)
	perfis := g.GerarUsuarios(n)
	interacoes := g.GerarInteracoes(perfis)
	pares, err := LerComplementaridade("../../scripts/referencia.cypher")
	if err != nil {
		t.Fatal(err)
	}
	f, err := NovaFonteMemoria(perfis, interacoes, pares)
	if err != nil {
		t.Fatal(err)
	}
	return f, perfis, interacoes
}

// A fonte em memória segue as regras do grafo: só quem passa no oráculo e não foi avaliado,
// com o mínimo de horários, em ordem, com corte, e o sinal social igual ao do cmd/verificar.
func TestFonteMemoria(t *testing.T) {
	f, perfis, interacoes := novaFonte(t, 600)
	ctx := context.Background()
	avaliados := map[[3]string]bool{}
	curtiu := map[[2]string]map[string]bool{}
	for _, x := range interacoes {
		avaliados[[3]string{x.De, x.Para, string(x.Jogo)}] = true
		if x.Curtiu {
			k := [2]string{x.De, string(x.Jogo)}
			if curtiu[k] == nil {
				curtiu[k] = map[string]bool{}
			}
			curtiu[k][x.Para] = true
		}
	}
	for _, eu := range perfis[:40] {
		for _, pj := range eu.Jogos {
			op := domain.OpcoesPadrao()
			cands, err := f.BuscarCandidatos(ctx, eu.ID, pj.Jogo, op)
			if err != nil {
				t.Fatal(err)
			}
			if len(cands) > op.Limite {
				t.Fatalf("%d candidatos, acima do limite %d", len(cands), op.Limite)
			}
			// sinal social por força bruta, como sinalEsperado do cmd/verificar
			sinais := map[string]int{}
			meus := curtiu[[2]string{eu.ID, string(pj.Jogo)}]
			for _, sim := range perfis {
				if sim.ID == eu.ID {
					continue
				}
				deles := curtiu[[2]string{sim.ID, string(pj.Jogo)}]
				comum := 0
				for alvo := range meus {
					if deles[alvo] {
						comum++
					}
				}
				for alvo := range deles {
					if comum > 0 {
						sinais[alvo] += comum
					}
				}
			}
			porID := map[string]domain.Perfil{}
			for _, p := range perfis {
				porID[p.ID] = p
			}
			for i, c := range cands {
				outro := porID[c.ID]
				if !domain.PassaFiltros(eu, outro, pj.Jogo) || avaliados[[3]string{eu.ID, c.ID, string(pj.Jogo)}] || c.ID == eu.ID {
					t.Fatalf("%s não deveria ser candidato de %s em %s", c.ID, eu.ID, pj.Jogo)
				}
				if c.HorariosComuns != bits.OnesCount32(eu.Disponibilidade&outro.Disponibilidade) || c.HorariosComuns < op.MinHorarios {
					t.Fatalf("horários em comum errados para %s × %s", eu.ID, c.ID)
				}
				if c.SinalSocial != sinais[c.ID] {
					t.Fatalf("sinal social de %s para %s: %d, quer %d", c.ID, eu.ID, c.SinalSocial, sinais[c.ID])
				}
				if i > 0 {
					a := cands[i-1]
					if a.HorariosComuns < c.HorariosComuns ||
						(a.HorariosComuns == c.HorariosComuns && a.Complemento < c.Complemento) ||
						(a.HorariosComuns == c.HorariosComuns && a.Complemento == c.Complemento && a.ID > c.ID) {
						t.Fatalf("ordem errada nas posições %d e %d", i-1, i)
					}
				}
			}
		}
	}
}

func TestExecutar(t *testing.T) {
	f, perfis, _ := novaFonte(t, 800)
	r, err := Executar(context.Background(), f, perfis, Config{K: 10, Amostra: 40, Opcoes: domain.OpcoesPadrao()})
	if err != nil {
		t.Fatal(err)
	}
	if r.Consultas == 0 {
		t.Fatal("nenhuma consulta avaliada")
	}
	for _, v := range Variantes {
		for _, m := range []*Media{r.NDCG[v], r.Precisao[v]} {
			if g := m.Geral(); g < 0 || g > 1 {
				t.Errorf("%s: média %v fora de [0, 1]", v, g)
			}
		}
	}
}
