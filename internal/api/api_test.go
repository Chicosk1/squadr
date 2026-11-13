package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"slices"
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

func (r *recFalso) Explicar(_ context.Context, idA, idB string, jogo domain.Jogo) (matching.Explicacao, error) {
	jogos, ok := r.joga[idA]
	if !ok || idB == "u9" {
		return matching.Explicacao{}, fmt.Errorf("%w: %s", domain.ErrNaoEncontrado, idA)
	}
	if !slices.Contains(jogos, jogo) {
		return matching.Explicacao{}, fmt.Errorf("%w: %s", matching.ErrNaoJoga, idA)
	}
	a := perfilU2()
	a.ID, a.Nome = idA, "Um"
	return matching.Explicacao{
		A: a, B: perfilU2(),
		Filtros: []domain.Filtro{
			{Nome: "Mesma região", Passou: true},
			{Nome: "Fila em comum", Passou: false, Motivo: "ARAM × Flex"},
		},
		HorariosComuns: 1 << 19,
		Calculado:      true,
		ScoreAB:        0.7, ScoreBA: 0.5, Reciproco: 2 * 0.7 * 0.5 / 1.2,
		DeA:        matching.Componentes{0.5, 1, 0, 1, 0.9, 0.25},
		DeB:        matching.Componentes{0.5, 0.75, 0, 1, 0.9, 1},
		Motivos:    []string{"comunicação"},
		Candidatos: 2,
	}, nil
}

type repoFalso struct {
	criados  []domain.Perfil
	curtidas map[[3]string]bool // de, para, jogo
	joga     map[string]bool
	perfis   map[string]domain.Perfil
	filtro   domain.FiltroUsuarios // o último filtro recebido por ListarUsuarios
}

func (r *repoFalso) ListarUsuarios(_ context.Context, f domain.FiltroUsuarios) (domain.PaginaUsuarios, error) {
	r.filtro = f
	return domain.PaginaUsuarios{Perfis: []domain.Perfil{perfilU2()}, Total: 1}, nil
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
	if n := strings.Count(html, `name="horarios" value="`); n != domain.NumHorarios {
		t.Errorf("a grade do cadastro tem %d horários, esperado %d", n, domain.NumHorarios)
	}
	if n := strings.Count(html, `<td data-h="`); n != domain.NumHorarios {
		t.Errorf("a grade da comparação tem %d horários, esperado %d", n, domain.NumHorarios)
	}
	if !strings.Contains(html, `data-horarios="seg-madrugada seg-manha`) || !strings.Contains(html, `data-limite="300"`) {
		t.Error("a página deveria levar a ordem dos horários e o corte de candidatos para o script")
	}
	for _, trecho := range []string{"Coletor", `value="ark:gatherer"`, `data-campo="tamanhoTribo"`,
		`data-campo="wipe"`, `data-campo="filas"`, "Desafiante", `id="painelPar"`,
		"Madrugada<small>0h–6h</small>", `data-peso="comunicacao"`} {
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

// TestSecundariaPreencher: no LoL, só quem escolhe Preencher pode ficar sem secundária (D-045).
func TestSecundariaPreencher(t *testing.T) {
	a := novoAmbiente(t)
	m := cadastroValido()
	lol := m["jogos"].([]map[string]any)[1]
	lol["funcaoPrincipal"], lol["funcaoSecundaria"] = "lol:fill", ""
	if w := a.chamar("POST", "/usuarios", corpo(t, m)); w.Code != http.StatusCreated {
		t.Fatalf("Preencher sem secundária: status %d: %s", w.Code, w.Body)
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
		"LoL sem secundária": alterar(func(m map[string]any) {
			m["jogos"].([]map[string]any)[1]["funcaoSecundaria"] = ""
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

// ---------- Exploração (D-043) ----------

func TestListarUsuarios(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("GET", "/usuarios", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body)
	}
	for _, proibido := range []string{"PERSONA", "persona", "pesos", "faixasAceitas"} {
		if strings.Contains(w.Body.String(), proibido) {
			t.Errorf("a lista expõe %q: %s", proibido, w.Body)
		}
	}
	var r respostaUsuarios
	lerJSON(t, w, &r)
	if r.Total != 1 || len(r.Usuarios) != 1 || r.Usuarios[0].Nome != "Dois" || r.Proximo != "" {
		t.Errorf("resposta inesperada: %+v", r)
	}
	if a.repo.filtro != (domain.FiltroUsuarios{Limite: LimitePadraoLista}) {
		t.Errorf("sem parâmetros, o filtro deveria ser vazio com o limite padrão: %+v", a.repo.filtro)
	}

	w = a.chamar("GET", "/usuarios?jogo=lol&regiao=EU&q=%20AnA%20&origem=formulario&depois=u000010&limite=20", "")
	esperado := domain.FiltroUsuarios{Jogo: domain.LoL, Regiao: "EU", Busca: "ana", PrefixoID: "w", Depois: "u000010", Limite: 20}
	if w.Code != http.StatusOK || a.repo.filtro != esperado {
		t.Errorf("status %d, filtro %+v, esperado %+v", w.Code, a.repo.filtro, esperado)
	}

	for _, url := range []string{
		"/usuarios?jogo=dota", "/usuarios?regiao=BR", "/usuarios?origem=x",
		"/usuarios?limite=0", "/usuarios?limite=101", "/usuarios?limite=dez",
		"/usuarios?q=" + strings.Repeat("a", TamanhoMaxBusca+1),
	} {
		if w := a.chamar("GET", url, ""); w.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, esperado 400", url, w.Code)
		}
	}
}

func TestPar(t *testing.T) {
	a := novoAmbiente(t)
	w := a.chamar("GET", "/usuarios/u1/par/u2?jogo=lol", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body)
	}
	for _, proibido := range []string{"PERSONA", "persona", "pesos", "faixasAceitas"} {
		if strings.Contains(w.Body.String(), proibido) {
			t.Errorf("a comparação expõe %q: %s", proibido, w.Body)
		}
	}
	var r respostaPar
	lerJSON(t, w, &r)
	if r.A.ID != "u1" || r.B.ID != "u2" || r.PassaFiltros || !r.Calculado || r.Posicao != 0 || r.Candidatos != 2 {
		t.Errorf("cabeçalho errado: %+v", r)
	}
	if len(r.Filtros) != 2 || r.Filtros[0].Motivo != "" || r.Filtros[1] != (filtroJSON{"Fila em comum", false, "ARAM × Flex"}) {
		t.Errorf("filtros errados: %+v", r.Filtros)
	}
	if !quaseIgual(r.FatorFuncao, 0.95) {
		t.Errorf("fator de função %v, esperado 0,95 (função 0,9)", r.FatorFuncao)
	}
	if strings.Join(r.Comuns, ",") != "sex-noite" {
		t.Errorf("horários em comum: %v", r.Comuns)
	}
	if len(r.Componentes) != int(matching.NumComponentes) ||
		r.Componentes[matching.Nivel] != (componentePar{"nível", 0.25, 1}) {
		t.Errorf("componentes errados: %+v", r.Componentes)
	}

	casos := []struct {
		url    string
		status int
	}{
		{"/usuarios/u1/par/u2", http.StatusBadRequest},
		{"/usuarios/u1/par/u1?jogo=lol", http.StatusBadRequest},
		{"/usuarios/u9/par/u2?jogo=lol", http.StatusNotFound},
		{"/usuarios/u1/par/u9?jogo=lol", http.StatusNotFound},
		{"/usuarios/u1/par/u2?jogo=ark", http.StatusUnprocessableEntity},
	}
	for _, c := range casos {
		if w := a.chamar("GET", c.url, ""); w.Code != c.status {
			t.Errorf("%s: status %d, esperado %d (%s)", c.url, w.Code, c.status, w.Body)
		}
	}
}

func quaseIgual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
