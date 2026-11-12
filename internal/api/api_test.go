package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

// ---------- Implementações em memória ----------

type recFalso struct{ joga map[string][]domain.Jogo }

func (r *recFalso) Recomendar(_ context.Context, id string, jogo domain.Jogo, k int) ([]matching.Recomendacao, error) {
	jogos, ok := r.joga[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domain.ErrNaoEncontrado, id)
	}
	for _, j := range jogos {
		if j == jogo {
			recs := []matching.Recomendacao{
				{ID: "u2", Nome: "Dois", Score: 0.9, Motivos: []string{"funções complementares"}},
				{ID: "u3", Nome: "Três", Score: 0.5, Motivos: []string{"nível"}},
			}
			return recs[:min(k, len(recs))], nil
		}
	}
	return nil, fmt.Errorf("%w: %s", matching.ErrNaoJoga, id)
}

type repoFalso struct {
	criados  []domain.Perfil
	curtidas map[[3]string]bool // de, para, jogo
	joga     map[string]bool
	perfis   map[string]domain.Perfil
}

func (r *repoFalso) Ping(context.Context) error { return nil }

func (r *repoFalso) BuscarPerfil(_ context.Context, id string) (domain.Perfil, error) {
	p, ok := r.perfis[id]
	if !ok {
		return domain.Perfil{}, fmt.Errorf("%w: %s", domain.ErrNaoEncontrado, id)
	}
	return p, nil
}

func (r *repoFalso) CriarUsuario(_ context.Context, p domain.Perfil) error {
	r.criados = append(r.criados, p)
	return nil
}

func (r *repoFalso) RegistrarInteracao(_ context.Context, x domain.Interacao) (bool, error) {
	if !r.joga[x.De] || !r.joga[x.Para] {
		return false, domain.ErrNaoEncontrado
	}
	k := [3]string{x.De, x.Para, string(x.Jogo)}
	if _, ok := r.curtidas[k]; ok {
		return false, domain.ErrJaAvaliado
	}
	r.curtidas[k] = x.Curtiu
	return x.Curtiu && r.curtidas[[3]string{x.Para, x.De, string(x.Jogo)}], nil
}

func rotulosCompletos() map[string]string {
	r := map[string]string{}
	for _, ids := range domain.Funcoes {
		for _, id := range ids {
			r[id] = "Rótulo " + id
		}
	}
	r["ark:gatherer"] = "Coletor"
	return r
}

type ambiente struct {
	h    http.Handler
	repo *repoFalso
	log  *bytes.Buffer
}

func novoAmbiente(t *testing.T) ambiente {
	t.Helper()
	repo := &repoFalso{curtidas: map[[3]string]bool{}, joga: map[string]bool{"u1": true, "u2": true},
		perfis: map[string]domain.Perfil{"u2": perfilU2()}}
	rec := &recFalso{joga: map[string][]domain.Jogo{"u1": {domain.LoL}}}
	buf := &bytes.Buffer{}
	s, err := Novo(rec, repo, rotulosCompletos(), log.New(buf, "", 0))
	if err != nil {
		t.Fatal(err)
	}
	s.novoID = func() (string, error) { return "wteste01", nil }
	return ambiente{h: s.Handler(), repo: repo, log: buf}
}

func (a ambiente) chamar(metodo, url, corpo string) *httptest.ResponseRecorder {
	var leitor io.Reader
	if corpo != "" {
		leitor = strings.NewReader(corpo)
	}
	w := httptest.NewRecorder()
	a.h.ServeHTTP(w, httptest.NewRequest(metodo, url, leitor))
	return w
}

func lerJSON(t *testing.T, w *httptest.ResponseRecorder, destino any) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), destino); err != nil {
		t.Fatalf("resposta não é JSON (%v): %s", err, w.Body.String())
	}
}

// ---------- Testes ----------

func TestPagina(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("GET", "/", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	html := w.Body.String()
	if n := strings.Count(html, `name="horarios"`); n != domain.NumHorarios {
		t.Errorf("a grade tem %d horários, esperado %d", n, domain.NumHorarios)
	}
	if n := strings.Count(html, `data-h="`); n != domain.NumHorarios {
		t.Errorf("a mini-grade do perfil tem %d horários, esperado %d", n, domain.NumHorarios)
	}
	for _, trecho := range []string{"Coletor", `value="ark:gatherer"`, `data-campo="tamanhoTribo"`,
		`data-campo="wipe"`, `data-campo="filas"`, "Desafiante"} {
		if !strings.Contains(html, trecho) {
			t.Errorf("a página não tem %q", trecho)
		}
	}
	if a.chamar("GET", "/nao-existe", "").Code != http.StatusNotFound {
		t.Error("rota desconhecida deveria dar 404")
	}
}

func TestNovoSemRotulo(t *testing.T) {
	r := rotulosCompletos()
	delete(r, "lol:fill")
	if _, err := Novo(&recFalso{}, &repoFalso{}, r, nil); err == nil {
		t.Error("deveria falhar sem o rótulo de lol:fill")
	}
}

func TestRecomendacoes(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("GET", "/usuarios/u1/recomendacoes?jogo=lol&k=1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body)
	}
	var r respostaRecomendacoes
	lerJSON(t, w, &r)
	if r.Usuario != "u1" || r.Jogo != domain.LoL || len(r.Recomendacoes) != 1 || r.Recomendacoes[0].Nome != "Dois" {
		t.Errorf("resposta inesperada: %+v", r)
	}

	casos := []struct {
		url    string
		status int
	}{
		{"/usuarios/u1/recomendacoes", http.StatusBadRequest},
		{"/usuarios/u1/recomendacoes?jogo=dota", http.StatusBadRequest},
		{"/usuarios/u1/recomendacoes?jogo=lol&k=0", http.StatusBadRequest},
		{"/usuarios/u1/recomendacoes?jogo=lol&k=51", http.StatusBadRequest},
		{"/usuarios/u9/recomendacoes?jogo=lol", http.StatusNotFound},
		{"/usuarios/u1/recomendacoes?jogo=ark", http.StatusUnprocessableEntity},
	}
	for _, c := range casos {
		if w := a.chamar("GET", c.url, ""); w.Code != c.status {
			t.Errorf("%s: status %d, esperado %d (%s)", c.url, w.Code, c.status, w.Body)
		}
	}
	if a.chamar("POST", "/usuarios/u1/recomendacoes?jogo=lol", "").Code != http.StatusMethodNotAllowed {
		t.Error("POST em recomendações deveria dar 405")
	}
}

func cadastroValido() map[string]any {
	return map[string]any{
		"maiorDeIdade": true, "nome": "Teste", "faixaEtaria": "18-24",
		"faixasAceitas": []string{"18-24", "25-34"}, "regiao": "AS", "idiomas": []string{"pt"},
		"horarios": []string{"sex-noite", "sab-noite"}, "usaVoz": "sim", "exigeVoz": false,
		"toleranciaToxicidade": 3, "compromisso": "fixo",
		"pesos": map[string]int{"horario": 5},
		"jogos": []map[string]any{
			{"jogo": "rust", "intensidade": 4, "nivelDesejado": "parecido",
				"funcaoPrincipal": "rust:builder", "funcaoSecundaria": "",
				"horas": "500-2000", "versao": "pc", "servidor": "oficial",
				"limiteGrupo": "trio", "wipe": "sempre", "plataforma": "ignorado no rust"},
			{"jogo": "lol", "intensidade": 2, "nivelDesejado": "tanto_faz",
				"funcaoPrincipal": "lol:support", "funcaoSecundaria": "lol:mid", "filas": []string{"aram"}},
		},
	}
}

func corpo(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestCriarUsuario(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("POST", "/usuarios", corpo(t, cadastroValido()))
	if w.Code != http.StatusCreated {
		t.Fatalf("status %d: %s", w.Code, w.Body)
	}
	if w.Header().Get("Location") != "/usuarios/wteste01" {
		t.Errorf("Location = %q", w.Header().Get("Location"))
	}
	if len(a.repo.criados) != 1 {
		t.Fatalf("esperado 1 usuário criado, %d", len(a.repo.criados))
	}
	p := a.repo.criados[0]
	if p.ID != "wteste01" || p.Pesos.Horario != 5 || p.Pesos.Nivel != matching.PesoPadrao {
		t.Errorf("id ou pesos errados: %+v", p)
	}
	if p.Disponibilidade != 1<<19|1<<23 {
		t.Errorf("disponibilidade %b", p.Disponibilidade)
	}
	rust, _ := p.PerfilDoJogo(domain.Rust)
	lol, _ := p.PerfilDoJogo(domain.LoL)
	if rust.Nivel != 0.625 || rust.Plataforma != "" || lol.Elo != domain.SemRanked || lol.Nivel != 0 {
		t.Errorf("conversão errada: rust=%+v lol=%+v", rust, lol)
	}
}

func TestCriarUsuarioInvalido(t *testing.T) {
	a := novoAmbiente(t)
	alterar := func(f func(m map[string]any)) string {
		m := cadastroValido()
		f(m)
		return corpo(t, m)
	}
	casos := map[string]string{
		"menor de idade":   alterar(func(m map[string]any) { m["maiorDeIdade"] = false }),
		"campo persona":    alterar(func(m map[string]any) { m["persona"] = "RUST-GRINDER" }),
		"região inválida":  alterar(func(m map[string]any) { m["regiao"] = "BR" }),
		"horário inválido": alterar(func(m map[string]any) { m["horarios"] = []string{"sex-tarde-noite"} }),
		"sem jogos":        alterar(func(m map[string]any) { m["jogos"] = []any{} }),
		"função de outro jogo": alterar(func(m map[string]any) {
			m["jogos"].([]map[string]any)[0]["funcaoPrincipal"] = "lol:adc"
		}),
		"JSON quebrado": `{"nome": `,
		"dois objetos":  corpo(t, cadastroValido()) + "{}",
	}
	for nome, c := range casos {
		w := a.chamar("POST", "/usuarios", c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, esperado 400", nome, w.Code)
			continue
		}
		var e respostaErro
		lerJSON(t, w, &e)
		if e.Erro == "" {
			t.Errorf("%s: resposta sem mensagem de erro", nome)
		}
	}
	if len(a.repo.criados) != 0 {
		t.Errorf("nenhum cadastro inválido deveria chegar ao banco, chegaram %d", len(a.repo.criados))
	}
}

func TestCurtidasERecusas(t *testing.T) {
	a := novoAmbiente(t)
	match := func(w *httptest.ResponseRecorder) bool {
		var r map[string]bool
		lerJSON(t, w, &r)
		return r["match"]
	}

	w := a.chamar("POST", "/usuarios/u1/curtidas", `{"para":"u2","jogo":"lol"}`)
	if w.Code != http.StatusCreated || match(w) {
		t.Fatalf("primeira curtida: status %d, match %v", w.Code, match(w))
	}
	w = a.chamar("POST", "/usuarios/u2/curtidas", `{"para":"u1","jogo":"lol"}`)
	if w.Code != http.StatusCreated || !match(w) {
		t.Fatalf("curtida recíproca deveria dar match: status %d", w.Code)
	}

	casos := []struct {
		nome, url, corpo string
		status           int
	}{
		{"repetida", "/usuarios/u1/recusas", `{"para":"u2","jogo":"lol"}`, http.StatusConflict},
		{"outro jogo é outra avaliação", "/usuarios/u1/recusas", `{"para":"u2","jogo":"ark"}`, http.StatusCreated},
		{"inexistente", "/usuarios/u1/curtidas", `{"para":"u9","jogo":"lol"}`, http.StatusNotFound},
		{"a si mesmo", "/usuarios/u1/curtidas", `{"para":"u1","jogo":"lol"}`, http.StatusBadRequest},
		{"jogo inválido", "/usuarios/u1/curtidas", `{"para":"u2","jogo":"x"}`, http.StatusBadRequest},
		{"campo extra", "/usuarios/u1/curtidas", `{"para":"u2","jogo":"lol","curtiu":true}`, http.StatusBadRequest},
	}
	for _, c := range casos {
		if w := a.chamar("POST", c.url, c.corpo); w.Code != c.status {
			t.Errorf("%s: status %d, esperado %d (%s)", c.nome, w.Code, c.status, w.Body)
		}
	}
}

type repoComPanic struct{ repoFalso }

func (repoComPanic) Ping(context.Context) error { panic("teste") }

type repoQuebrado struct{ repoFalso }

func (repoQuebrado) CriarUsuario(context.Context, domain.Perfil) error {
	return errors.New("senha=segredo no texto do erro do banco")
}

func TestMiddlewares(t *testing.T) {
	buf := &bytes.Buffer{}
	s, err := Novo(&recFalso{}, &repoComPanic{}, rotulosCompletos(), log.New(buf, "", 0))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/saude", nil))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("panic deveria virar 500, veio %d", w.Code)
	}
	if !strings.Contains(buf.String(), "panic em GET /saude") || !strings.Contains(buf.String(), "GET /saude 500") {
		t.Errorf("o log deveria registrar o panic e o 500:\n%s", buf)
	}

	// Erro interno: a causa vai para o log, nunca para o cliente.
	buf.Reset()
	s, _ = Novo(&recFalso{}, &repoQuebrado{}, rotulosCompletos(), log.New(buf, "", 0))
	w = httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("POST", "/usuarios", strings.NewReader(corpo(t, cadastroValido()))))
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), "segredo") {
		t.Errorf("erro interno vazou ou status errado: %d %s", w.Code, w.Body)
	}
	if !strings.Contains(buf.String(), "segredo") {
		t.Error("a causa do erro deveria estar no log")
	}
}

func perfilU2() domain.Perfil {
	return domain.Perfil{
		ID: "u2", Nome: "Dois", FaixaEtaria: "25-34", FaixasAceitas: []string{"25-34"},
		Regiao: "AS", Idiomas: []string{"pt", "en"}, UsaVoz: "as_vezes", ExigeVoz: true,
		ToleranciaToxicidade: 2, Compromisso: "fixo",
		Pesos:           domain.Pesos{Horario: 5, Intensidade: 1, Modo: 1, Comunicacao: 1, Funcao: 1, Nivel: 1},
		Disponibilidade: 1<<19 | 1<<23,
		Persona:         "PERSONA-SECRETA",
		Jogos: []domain.PerfilJogo{
			{Jogo: domain.ARK, Intensidade: 4, NivelDesejado: "parecido", FuncaoPrincipal: "ark:gatherer",
				Horas: "2000+", Versao: "asa", Plataforma: "pc_epic", Modo: "pvp", Servidor: "oficial",
				TamanhoTribo: "7+", Elo: domain.SemRanked, PersonaJogo: "PERSONA-JOGO-SECRETA"},
			{Jogo: domain.LoL, Intensidade: 2, NivelDesejado: "tanto_faz", FuncaoPrincipal: "lol:adc",
				FuncaoSecundaria: "lol:mid", Elo: 3, Filas: []string{"solo_duo", "aram"}},
		},
	}
}

func TestPerfilPublico(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("GET", "/usuarios/u2", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body)
	}
	// Regra inviolável e D-042: nada de persona, pesos ou faixas aceitas.
	for _, proibido := range []string{"PERSONA", "persona", "pesos", "faixasAceitas", "nivel\""} {
		if strings.Contains(w.Body.String(), proibido) {
			t.Errorf("o perfil público expõe %q: %s", proibido, w.Body)
		}
	}
	var p perfilPublico
	lerJSON(t, w, &p)
	if p.Regiao != "América do Sul" || p.UsaVoz != "Às vezes" || !p.ExigeVoz ||
		strings.Join(p.Idiomas, ",") != "Português,Inglês" || strings.Join(p.Horarios, ",") != "sex-noite,sab-noite" {
		t.Errorf("campos gerais errados: %+v", p)
	}
	if len(p.Jogos) != 2 {
		t.Fatalf("esperados 2 jogos, vieram %d", len(p.Jogos))
	}
	ark, lol := p.Jogos[0], p.Jogos[1]
	if ark.FuncaoPrincipal != "Coletor" || ark.FuncaoSecundaria != "" || ark.Detalhes[1] != (campo{"Plataforma", "PC (Epic)"}) {
		t.Errorf("ARK errado: %+v", ark)
	}
	if lol.Detalhes[0] != (campo{"Elo", "Ouro"}) || lol.Detalhes[1] != (campo{"Filas", "Solo/Duo, ARAM"}) {
		t.Errorf("LoL errado: %+v", lol)
	}
	if a.chamar("GET", "/usuarios/u9", "").Code != http.StatusNotFound {
		t.Error("usuário inexistente deveria dar 404")
	}
}

func TestComponentesNaRecomendacao(t *testing.T) {
	a := novoAmbiente(t)
	var r struct {
		Recomendacoes []struct {
			Componentes []componenteJSON `json:"componentes"`
		} `json:"recomendacoes"`
	}
	lerJSON(t, a.chamar("GET", "/usuarios/u1/recomendacoes?jogo=lol", ""), &r)
	if len(r.Recomendacoes) == 0 || len(r.Recomendacoes[0].Componentes) != int(matching.NumComponentes) {
		t.Fatalf("esperados %d componentes: %+v", matching.NumComponentes, r)
	}
	if r.Recomendacoes[0].Componentes[matching.Funcao].Nome != "funções complementares" {
		t.Errorf("ordem ou nome dos componentes errado: %+v", r.Recomendacoes[0].Componentes)
	}
}
