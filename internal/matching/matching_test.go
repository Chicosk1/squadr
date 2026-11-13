package matching

import (
	"context"
	"errors"
	"math"
	"math/bits"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/seed"
)

// ---------- Apoio ----------

var rePar = regexp.MustCompile(`\{a:\s*'([^']+)',\s*b:\s*'([^']+)',\s*valor:\s*([0-9.]+)\}`)

// paresDoScript lê a complementaridade direto da fonte única (D-018), sem banco.
func paresDoScript(t *testing.T) map[[2]string]float64 {
	t.Helper()
	texto, err := os.ReadFile("../../scripts/referencia.cypher")
	if err != nil {
		t.Fatal(err)
	}
	pares := map[[2]string]float64{}
	for _, m := range rePar.FindAllStringSubmatch(string(texto), -1) {
		v, err := strconv.ParseFloat(m[3], 64)
		if err != nil {
			t.Fatal(err)
		}
		pares[[2]string{m[1], m[2]}] = v
	}
	if len(pares) != 51 {
		t.Fatalf("esperados 51 pares de COMPLEMENTA no script, lidos %d", len(pares))
	}
	return pares
}

func matrizDoScript(t *testing.T) Matriz {
	t.Helper()
	m, err := NovaMatriz(paresDoScript(t))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// fonteMemoria imita o grafo: candidatos são os que passam no oráculo, sem limite.
type fonteMemoria struct {
	perfis []domain.Perfil
	pares  map[[2]string]float64
	social map[string]int
}

func (f *fonteMemoria) BuscarCandidatos(_ context.Context, id string, jogo domain.Jogo, _ domain.OpcoesCandidatos) ([]domain.Candidato, error) {
	var eu domain.Perfil
	for _, p := range f.perfis {
		if p.ID == id {
			eu = p
		}
	}
	var saida []domain.Candidato
	for _, p := range f.perfis {
		if p.ID != id && domain.PassaFiltros(eu, p, jogo) {
			saida = append(saida, domain.Candidato{ID: p.ID, SinalSocial: f.social[p.ID]})
		}
	}
	return saida, nil
}

func (f *fonteMemoria) BuscarPerfis(_ context.Context, ids []string) ([]domain.Perfil, error) {
	quer := map[string]bool{}
	for _, id := range ids {
		quer[id] = true
	}
	var saida []domain.Perfil
	for _, p := range f.perfis {
		if quer[p.ID] {
			saida = append(saida, p)
		}
	}
	return saida, nil
}

func (f *fonteMemoria) BuscarComplementaridade(context.Context) (map[[2]string]float64, error) {
	return f.pares, nil
}

func quase(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// ---------- Componentes ----------

func TestReciproco(t *testing.T) {
	casos := []struct{ ab, ba, esperado float64 }{
		{0, 0, 0},
		{1, 0, 0},
		{0.5, 0.5, 0.5},
		{1, 1, 1},
		{0.9, 0.1, 0.18},
	}
	for _, c := range casos {
		if got := Reciproco(c.ab, c.ba); !quase(got, c.esperado) {
			t.Errorf("Reciproco(%v, %v) = %v, esperado %v", c.ab, c.ba, got, c.esperado)
		}
		if Reciproco(c.ab, c.ba) != Reciproco(c.ba, c.ab) {
			t.Errorf("Reciproco não é simétrico para (%v, %v)", c.ab, c.ba)
		}
	}
}

func TestComponentes(t *testing.T) {
	casos := []struct {
		nome     string
		got, esp float64
	}{
		{"jaccard vazio", jaccardBits(0, 0), 0},
		{"jaccard igual", jaccardBits(0b1011, 0b1011), 1},
		{"jaccard metade", jaccardBits(0b0011, 0b0110), 1.0 / 3},
		{"compromisso fixo×eventual", compromisso("fixo", "eventual"), 0},
		{"compromisso tanto_faz", compromisso("fixo", "tanto_faz"), 1},
		{"escala vizinhos", escala("trio", "quad", []string{"solo_duo", "trio", "quad", "sem_limite"}), 2.0 / 3},
		{"escala tanto_faz", escala("2-3", "tanto_faz", []string{"2-3", "4-6", "7+"}), 1},
		{"escala desconhecido", escala("x", "2-3", []string{"2-3", "4-6", "7+"}), 0.5},
		{"servidor ARK rates", servidorARK("rates_altas", "rates_baixas"), 0.5},
		{"servidor ARK oficial×rates", servidorARK("oficial", "rates_baixas"), 0},
		{"servidor Rust vanilla", servidorRust("oficial", "comunitario"), 0.5},
		{"servidor Rust modificado", servidorRust("oficial", "modificado"), 0},
		{"voz não exige", fatorVoz(false, "nao"), 1},
		{"voz exige e às vezes", fatorVoz(true, "as_vezes"), 0.5},
		{"voz exige e não usa", fatorVoz(true, "nao"), 0},
		{"tolerância B mais tolerante", tolerancia(1, 5), 0},
		{"tolerância B menos tolerante", tolerancia(5, 1), 1},
		{"nível parecido igual", nivel("parecido", 0), 1},
		{"nível parecido 1 degrau", nivel("parecido", 0.25), 0.5},
		{"nível parecido 3 degraus", nivel("parecido", -0.75), 0},
		{"nível mais experiente igual", nivel("mais_experiente", 0), 0.5},
		{"nível mais experiente acima", nivel("mais_experiente", 0.75), 1},
		{"nível mais experiente abaixo", nivel("mais_experiente", -0.25), 0},
		{"nível tanto faz", nivel("tanto_faz", -0.75), 1},
	}
	for _, c := range casos {
		if !quase(c.got, c.esp) {
			t.Errorf("%s: %v, esperado %v", c.nome, c.got, c.esp)
		}
	}
}

func TestFuncaoArranjo(t *testing.T) {
	m := matrizDoScript(t)
	pj := func(p, s string) domain.PerfilJogo {
		return domain.PerfilJogo{Jogo: domain.LoL, FuncaoPrincipal: p, FuncaoSecundaria: s}
	}
	casos := []struct {
		nome   string
		a, b   domain.PerfilJogo
		espera float64
	}{
		{"principais complementares", pj("lol:adc", "lol:mid"), pj("lol:support", "lol:top"), 1.0},
		{"só pela secundária de B", pj("lol:adc", "lol:mid"), pj("lol:top", "lol:support"), 0.75},
		{"só pelas duas secundárias", pj("lol:mid", "lol:adc"), pj("lol:top", "lol:support"), 0.5625},
		{"sem secundária", pj("lol:adc", ""), pj("lol:adc", ""), 0.1},
	}
	for _, c := range casos {
		if got := funcao(c.a, c.b, m); !quase(got, c.espera) {
			t.Errorf("%s: %v, esperado %v", c.nome, got, c.espera)
		}
		if funcao(c.a, c.b, m) != funcao(c.b, c.a, m) {
			t.Errorf("%s: função não é simétrica", c.nome)
		}
	}
}

func TestFilasDuo(t *testing.T) {
	a := domain.PerfilJogo{Jogo: domain.LoL, Elo: 3, Filas: []string{"solo_duo", "flex"}}
	b := domain.PerfilJogo{Jogo: domain.LoL, Elo: 8, Filas: []string{"solo_duo", "flex"}}
	if got := filas(a, b); !quase(got, 0.5) {
		t.Errorf("elo distante: solo_duo não deveria contar, veio %v", got)
	}
	b.Elo = 4
	if got := filas(a, b); !quase(got, 1) {
		t.Errorf("elo próximo: esperado 1, veio %v", got)
	}
}

// ---------- Propriedades sobre a base gerada ----------

func gerar(t *testing.T, n int) []domain.Perfil {
	t.Helper()
	return seed.NovoGerador(42).GerarUsuarios(n)
}

func TestPropriedades(t *testing.T) {
	m := matrizDoScript(t)
	perfis := gerar(t, 600)
	pares := 0
	for i := 0; i < len(perfis); i += 3 {
		a := perfis[i]
		for j := i + 1; j < len(perfis); j += 7 {
			b := perfis[j]
			for _, pa := range a.Jogos {
				pb, ok := b.PerfilDoJogo(pa.Jogo)
				if !ok {
					continue
				}
				pares++
				ab, ca := Unilateral(a, b, pa, pb, m)
				ba, cb := Unilateral(b, a, pb, pa, m)
				for k := range NumComponentes {
					if ca[k] < 0 || ca[k] > 1 || cb[k] < 0 || cb[k] > 1 {
						t.Fatalf("%s×%s: componente %s fora de [0, 1]: %v / %v", a.ID, b.ID, k, ca[k], cb[k])
					}
				}
				if ab < 0 || ab > 1 || ba < 0 || ba > 1 {
					t.Fatalf("%s×%s: score fora de [0, 1]: %v / %v", a.ID, b.ID, ab, ba)
				}
				r := Reciproco(ab, ba)
				if r != Reciproco(ba, ab) || r > max(ab, ba)+1e-12 || r < min(ab, ba)-1e-12 {
					t.Fatalf("%s×%s: recíproco %v inconsistente com %v e %v", a.ID, b.ID, r, ab, ba)
				}
				// Componentes simétricos por construção.
				for _, k := range []Componente{Horario, Funcao} {
					if !quase(ca[k], cb[k]) {
						t.Fatalf("%s×%s: %s deveria ser simétrico: %v × %v", a.ID, b.ID, k, ca[k], cb[k])
					}
				}
			}
		}
	}
	if pares < 1000 {
		t.Fatalf("poucos pares testados: %d", pares)
	}
}

// TestAprendizGrinder é o caso de controle da avaliação (personas.md, grau 0): o aprendiz
// gosta do grinder (quer alguém mais experiente), mas o grinder não gosta do aprendiz.
// A reciprocidade precisa derrubar o par.
func TestAprendizGrinder(t *testing.T) {
	m := matrizDoScript(t)
	perfis := gerar(t, 2000)
	var aprendizes, grinders []domain.Perfil
	for _, p := range perfis {
		switch p.Persona { // o teste lê a persona; o motor não
		case "RUST-APRENDIZ":
			aprendizes = append(aprendizes, p)
		case "RUST-GRINDER":
			grinders = append(grinders, p)
		}
	}

	var somaAG, somaGA, somaRecAG, somaRecAA float64
	var nAG, nAA int
	for _, a := range aprendizes {
		pa, _ := a.PerfilDoJogo(domain.Rust)
		for _, g := range grinders {
			if !domain.PassaFiltros(a, g, domain.Rust) {
				continue
			}
			pg, _ := g.PerfilDoJogo(domain.Rust)
			ag, _ := Unilateral(a, g, pa, pg, m)
			ga, _ := Unilateral(g, a, pg, pa, m)
			somaAG += ag
			somaGA += ga
			somaRecAG += Reciproco(ag, ga)
			nAG++
		}
		for _, b := range aprendizes {
			if a.ID == b.ID || !domain.PassaFiltros(a, b, domain.Rust) {
				continue
			}
			pb, _ := b.PerfilDoJogo(domain.Rust)
			ab, _ := Unilateral(a, b, pa, pb, m)
			ba, _ := Unilateral(b, a, pb, pa, m)
			somaRecAA += Reciproco(ab, ba)
			nAA++
		}
	}
	if nAG == 0 || nAA == 0 {
		t.Fatalf("sem pares para comparar: %d aprendiz×grinder, %d aprendiz×aprendiz", nAG, nAA)
	}
	mAG, mGA := somaAG/float64(nAG), somaGA/float64(nAG)
	mRecAG, mRecAA := somaRecAG/float64(nAG), somaRecAA/float64(nAA)
	t.Logf("s(aprendiz→grinder)=%.3f s(grinder→aprendiz)=%.3f recíproco=%.3f | aprendiz×aprendiz recíproco=%.3f (%d e %d pares)",
		mAG, mGA, mRecAG, mRecAA, nAG, nAA)

	if mGA >= mAG {
		t.Errorf("o grinder deveria gostar menos do aprendiz do que o contrário: %.3f ≥ %.3f", mGA, mAG)
	}
	if mRecAG >= mAG {
		t.Errorf("a reciprocidade deveria derrubar o par: recíproco %.3f ≥ unilateral %.3f", mRecAG, mAG)
	}
	if mRecAG >= mRecAA {
		t.Errorf("aprendiz×grinder (grau 0) deveria ficar abaixo de aprendiz×aprendiz (grau 2): %.3f ≥ %.3f", mRecAG, mRecAA)
	}
}

// ---------- Recomendador ----------

func novoRecomendador(t *testing.T, f *fonteMemoria) *Recomendador {
	t.Helper()
	r, err := Novo(context.Background(), f, domain.OpcoesPadrao())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRecomendar(t *testing.T) {
	perfis := gerar(t, 400)
	f := &fonteMemoria{perfis: perfis, pares: paresDoScript(t)}
	r := novoRecomendador(t, f)
	ctx := context.Background()

	eu := perfis[0]
	jogo := eu.Jogos[0].Jogo
	recs, err := r.Recomendar(ctx, eu.ID, jogo, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) == 0 || len(recs) > 10 {
		t.Fatalf("esperado de 1 a 10 recomendações, vieram %d", len(recs))
	}
	for i, rec := range recs {
		if rec.ID == eu.ID {
			t.Errorf("recomendou o próprio usuário")
		}
		if !quase(rec.Score, PesoReciproco*rec.Reciproco+PesoSocial*rec.Social) {
			t.Errorf("%s: score não bate com a combinação", rec.ID)
		}
		if i > 0 && rec.Score > recs[i-1].Score {
			t.Errorf("ranking fora de ordem na posição %d", i)
		}
		if len(rec.Motivos) == 0 {
			t.Errorf("%s: sem motivo do match", rec.ID)
		}
	}

	if _, err := r.Recomendar(ctx, "nao-existe", jogo, 10); err == nil {
		t.Error("usuário inexistente deveria dar erro")
	}
}

func TestSocialNormalizado(t *testing.T) {
	if normalizarSocial(0, 0) != 0 || normalizarSocial(5, 0) != 0 {
		t.Error("sem sinal na consulta, o social deve ser 0")
	}
	if !quase(normalizarSocial(10, 10), 1) {
		t.Error("o maior sinal deve valer 1")
	}
	if v := normalizarSocial(1, 100); v <= 0.01 || v >= 0.5 {
		t.Errorf("log1p deveria amortecer: 1 de 100 vale %v", v)
	}
}

// TestIgnoraPersona garante a regra inviolável: o motor não lê Persona nem PersonaJogo.
func TestIgnoraPersona(t *testing.T) {
	perfis := gerar(t, 300)
	pares := paresDoScript(t)
	ctx := context.Background()

	trocados := make([]domain.Perfil, len(perfis))
	for i, p := range perfis {
		p.Persona = "X"
		p.Jogos = append([]domain.PerfilJogo(nil), p.Jogos...)
		for j := range p.Jogos {
			p.Jogos[j].PersonaJogo = "X"
		}
		trocados[i] = p
	}
	original := novoRecomendador(t, &fonteMemoria{perfis: perfis, pares: pares})
	alterado := novoRecomendador(t, &fonteMemoria{perfis: trocados, pares: pares})

	for _, eu := range perfis[:30] {
		for _, pj := range eu.Jogos {
			a, err := original.Pontuar(ctx, eu.ID, pj.Jogo)
			if err != nil {
				t.Fatal(err)
			}
			b, err := alterado.Pontuar(ctx, eu.ID, pj.Jogo)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(a, b) {
				t.Fatalf("%s em %s: o ranking mudou ao alterar a persona", eu.ID, pj.Jogo)
			}
		}
	}
}

// ---------- Explicação de um par ----------

func TestExplicar(t *testing.T) {
	perfis := gerar(t, 400)
	f := &fonteMemoria{perfis: perfis, pares: paresDoScript(t)}
	r := novoRecomendador(t, f)
	ctx := context.Background()

	// Um usuário com pelo menos 3 candidatos no primeiro jogo dele.
	var eu domain.Perfil
	var jogo domain.Jogo
	var ranking []Recomendacao
	for _, p := range perfis {
		recs, err := r.Pontuar(ctx, p.ID, p.Jogos[0].Jogo)
		if err != nil {
			t.Fatal(err)
		}
		if len(recs) >= 3 {
			eu, jogo, ranking = p, p.Jogos[0].Jogo, recs
			break
		}
	}
	if len(ranking) < 3 {
		t.Fatal("nenhum usuário com 3 candidatos")
	}

	// Um candidato do ranking: a explicação repete os números do ranking.
	for pos, rec := range ranking[:3] {
		e, err := r.Explicar(ctx, eu.ID, rec.ID, jogo)
		if err != nil {
			t.Fatal(err)
		}
		if e.Posicao != pos+1 || e.Candidatos != len(ranking) || !e.Calculado {
			t.Errorf("%s: posição %d de %d, esperado %d de %d", rec.ID, e.Posicao, e.Candidatos, pos+1, len(ranking))
		}
		if !quase(e.ScoreAB, rec.ScoreAB) || !quase(e.ScoreBA, rec.ScoreBA) ||
			!quase(e.Reciproco, rec.Reciproco) || !quase(e.Score, rec.Score) || e.DeA != rec.DeA {
			t.Errorf("%s: a explicação não bate com o ranking", rec.ID)
		}
		comuns := bits.OnesCount32(e.HorariosComuns)
		if e.PassaFiltros != (comuns >= 1) {
			t.Errorf("%s: passa=%v com %d horários em comum", rec.ID, e.PassaFiltros, comuns)
		}
	}

	// Alguém barrado por um filtro: fora do ranking, com o motivo, mas com score calculado se joga.
	for _, outro := range perfis[1:] {
		if domain.PassaFiltros(eu, outro, jogo) {
			continue
		}
		e, err := r.Explicar(ctx, eu.ID, outro.ID, jogo)
		if err != nil {
			t.Fatal(err)
		}
		if e.PassaFiltros || e.Posicao != 0 {
			t.Fatalf("%s: barrado, mas passa=%v posição=%d", outro.ID, e.PassaFiltros, e.Posicao)
		}
		barrou := false
		for _, fl := range e.Filtros {
			if !fl.Passou {
				barrou = true
				if fl.Motivo == "" {
					t.Errorf("%s: filtro %q barrou sem motivo", outro.ID, fl.Nome)
				}
			}
		}
		_, joga := outro.PerfilDoJogo(jogo)
		if !barrou || e.Calculado != joga {
			t.Errorf("%s: barrou=%v calculado=%v joga=%v", outro.ID, barrou, e.Calculado, joga)
		}
		break
	}

	if _, err := r.Explicar(ctx, eu.ID, "nao-existe", jogo); !errors.Is(err, ErrUsuarioNaoEncontrado) {
		t.Errorf("B inexistente: erro %v", err)
	}
	for _, j := range domain.Jogos {
		if _, joga := eu.PerfilDoJogo(j); !joga {
			if _, err := r.Explicar(ctx, eu.ID, ranking[0].ID, j); !errors.Is(err, ErrNaoJoga) {
				t.Errorf("A não joga %s: erro %v", j, err)
			}
			break
		}
	}
}

// ---------- Fator de função (D-045) ----------

// TestFatorFuncao usa dois jogadores iguais em tudo, menos nas funções: o fator precisa
// separar os encaixes que a média ponderada sozinha deixava em 0,85 a 1,00.
func TestFatorFuncao(t *testing.T) {
	m, err := NovaMatriz(paresDoScript(t))
	if err != nil {
		t.Fatal(err)
	}
	base := domain.Perfil{Disponibilidade: 0xF0F0, Compromisso: "fixo", UsaVoz: "sim", ToleranciaToxicidade: 3,
		Pesos: domain.Pesos{Horario: 3, Intensidade: 3, Modo: 3, Comunicacao: 3, Funcao: 3, Nivel: 3}}
	pj := func(jogo domain.Jogo, principal, secundaria string) domain.PerfilJogo {
		return domain.PerfilJogo{Jogo: jogo, FuncaoPrincipal: principal, FuncaoSecundaria: secundaria,
			Intensidade: 4, Elo: 4, Nivel: 0.4, NivelDesejado: "parecido", Filas: []string{"solo_duo"},
			Servidor: "oficial", TamanhoTribo: "2-3", LimiteGrupo: "trio", Wipe: "sempre"}
	}
	score := func(a, b domain.PerfilJogo) float64 {
		s, _ := Unilateral(base, base, a, b, m)
		return s
	}
	casos := []struct {
		nome        string
		a, b        domain.PerfilJogo
		funcao, esp float64
	}{
		// Um vai de jungle e o outro de top, a secundária dele: 0,8 × 0,75 = 0,6.
		{"LoL Jungle/Top × Jungle/Top", pj(domain.LoL, "lol:jungle", "lol:top"), pj(domain.LoL, "lol:jungle", "lol:top"), 0.6, (5 + 0.6) / 6 * 0.8},
		{"LoL Jungle × Mid", pj(domain.LoL, "lol:jungle", "lol:top"), pj(domain.LoL, "lol:mid", "lol:top"), 0.9, (5 + 0.9) / 6 * 0.95},
		{"LoL ADC × Suporte", pj(domain.LoL, "lol:adc", "lol:mid"), pj(domain.LoL, "lol:support", "lol:mid"), 1, 1},
		{"ARK Builder × Builder", pj(domain.ARK, "ark:builder", ""), pj(domain.ARK, "ark:builder", ""), 0.4, (5 + 0.4) / 6 * 0.7},
		{"ARK Builder × Breeder", pj(domain.ARK, "ark:builder", ""), pj(domain.ARK, "ark:breeder", ""), 0.9, (5 + 0.9) / 6 * 0.95},
		{"Rust Eletricista × Eletricista", pj(domain.Rust, "rust:electrician", ""), pj(domain.Rust, "rust:electrician", ""), 0.3, (5 + 0.3) / 6 * 0.65},
	}
	for _, c := range casos {
		s, comp := Unilateral(base, base, c.a, c.b, m)
		if !quase(comp[Funcao], c.funcao) || !quase(s, c.esp) {
			t.Errorf("%s: função %.3f e score %.3f, esperado %.3f e %.3f", c.nome, comp[Funcao], s, c.funcao, c.esp)
		}
	}
	// O conflito de rota fica bem abaixo do encaixe ideal, e não 7 pontos abaixo, como antes.
	if d := score(casos[2].a, casos[2].b) - score(casos[0].a, casos[0].b); d < 0.2 {
		t.Errorf("ADC × Suporte só %.3f acima de Jungle/Top × Jungle/Top", d)
	}
	if FatorFuncao(0) != PisoFatorFuncao || FatorFuncao(1) != 1 {
		t.Error("o fator deve ir de PisoFatorFuncao (função 0) a 1 (função 1)")
	}
}
