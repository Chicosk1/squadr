package api

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

// Rotas da tela de exploração (D-043): a lista de usuários e a comparação de um par.

const (
	LimitePadraoLista = 50
	LimiteMaximoLista = 100
	TamanhoMaxBusca   = 40
)

// prefixosOrigem traduz a origem pedida no prefixo do id (D-039): o formulário gera ids "w...",
// o seed gera "u...".
var prefixosOrigem = map[string]string{"": "", "formulario": "w", "seed": "u"}

// ---------- GET /usuarios?jogo=&regiao=&q=&origem=&depois=&limite= ----------

type respostaUsuarios struct {
	Usuarios []perfilPublico `json:"usuarios"`
	Total    int             `json:"total"`
	Proximo  string          `json:"proximo,omitempty"`
}

func (s *Servidor) listarUsuarios(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := domain.FiltroUsuarios{Limite: LimitePadraoLista, Depois: q.Get("depois")}

	if v := q.Get("jogo"); v != "" {
		jogo, ok := jogoValido(v)
		if !ok {
			responderErro(w, http.StatusBadRequest, "parâmetro jogo deve ser ark, rust ou lol")
			return
		}
		f.Jogo = jogo
	}
	if v := q.Get("regiao"); v != "" {
		if !slices.ContainsFunc(domain.Regioes, func(o domain.Opcao) bool { return o.Valor == v }) {
			responderErro(w, http.StatusBadRequest, "parâmetro regiao deve ser AS, AN ou EU")
			return
		}
		f.Regiao = v
	}
	prefixo, ok := prefixosOrigem[q.Get("origem")]
	if !ok {
		responderErro(w, http.StatusBadRequest, "parâmetro origem deve ser formulario ou seed")
		return
	}
	f.PrefixoID = prefixo
	f.Busca = strings.ToLower(strings.TrimSpace(q.Get("q")))
	if utf8.RuneCountInString(f.Busca) > TamanhoMaxBusca || len(f.Depois) > 64 {
		responderErro(w, http.StatusBadRequest, "busca ou cursor longo demais")
		return
	}
	if v := q.Get("limite"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > LimiteMaximoLista {
			responderErro(w, http.StatusBadRequest, "parâmetro limite deve ser um inteiro de 1 a "+strconv.Itoa(LimiteMaximoLista))
			return
		}
		f.Limite = n
	}

	pagina, err := s.repo.ListarUsuarios(r.Context(), f)
	if err != nil {
		s.erroInterno(w, r, err)
		return
	}
	saida := respostaUsuarios{Usuarios: make([]perfilPublico, len(pagina.Perfis)), Total: pagina.Total, Proximo: pagina.Proximo}
	for i, p := range pagina.Perfis {
		saida.Usuarios[i] = s.paraPublico(p)
	}
	responderJSON(w, http.StatusOK, saida)
}

// ---------- GET /usuarios/{id}/par/{outro}?jogo= ----------

type filtroJSON struct {
	Nome   string `json:"nome"`
	Passou bool   `json:"passou"`
	Motivo string `json:"motivo,omitempty"`
}

type componentePar struct {
	Nome  string  `json:"nome"`
	ParaA float64 `json:"paraA"` // c(A→B): quanto B atende A
	ParaB float64 `json:"paraB"` // c(B→A): quanto A atende B
}

type respostaPar struct {
	Jogo         domain.Jogo     `json:"jogo"`
	A            perfilPublico   `json:"a"`
	B            perfilPublico   `json:"b"`
	Filtros      []filtroJSON    `json:"filtros"`
	PassaFiltros bool            `json:"passaFiltros"`
	Comuns       []string        `json:"horariosComuns"` // ids, como "sex-noite"
	Calculado    bool            `json:"calculado"`      // falso quando B não joga o jogo
	ScoreParaA   float64         `json:"scoreParaA"`     // s(A→B)
	ScoreParaB   float64         `json:"scoreParaB"`     // s(B→A)
	Reciproco    float64         `json:"reciproco"`
	Componentes  []componentePar `json:"componentes"`
	Motivos      []string        `json:"motivos"`
	Posicao      int             `json:"posicao"` // de B no ranking de A; 0 se B não está nele
	Candidatos   int             `json:"candidatos"`
	Score        float64         `json:"score"`
	Social       float64         `json:"social"`
}

func (s *Servidor) par(w http.ResponseWriter, r *http.Request) {
	idA, idB := r.PathValue("id"), r.PathValue("outro")
	jogo, ok := jogoValido(r.URL.Query().Get("jogo"))
	if !ok {
		responderErro(w, http.StatusBadRequest, "parâmetro jogo deve ser ark, rust ou lol")
		return
	}
	if idA == idB {
		responderErro(w, http.StatusBadRequest, "escolha dois usuários diferentes")
		return
	}

	e, err := s.rec.Explicar(r.Context(), idA, idB, jogo)
	switch {
	case errors.Is(err, domain.ErrNaoEncontrado):
		responderErro(w, http.StatusNotFound, "usuário não encontrado")
		return
	case errors.Is(err, matching.ErrNaoJoga):
		responderErro(w, http.StatusUnprocessableEntity, "o usuário "+idA+" não joga "+string(jogo))
		return
	case err != nil:
		s.erroInterno(w, r, err)
		return
	}

	saida := respostaPar{
		Jogo: jogo, A: s.paraPublico(e.A), B: s.paraPublico(e.B),
		Filtros:      make([]filtroJSON, len(e.Filtros)),
		PassaFiltros: e.PassaFiltros,
		Comuns:       domain.Perfil{Disponibilidade: e.HorariosComuns}.HorarioIDs(),
		Calculado:    e.Calculado,
		ScoreParaA:   e.ScoreAB, ScoreParaB: e.ScoreBA, Reciproco: e.Reciproco,
		Componentes: []componentePar{},
		Motivos:     e.Motivos,
		Posicao:     e.Posicao, Candidatos: e.Candidatos, Score: e.Score, Social: e.Social,
	}
	for i, f := range e.Filtros {
		saida.Filtros[i] = filtroJSON{Nome: f.Nome, Passou: f.Passou, Motivo: f.Motivo}
	}
	if saida.Comuns == nil {
		saida.Comuns = []string{}
	}
	if saida.Motivos == nil {
		saida.Motivos = []string{}
	}
	if e.Calculado {
		for c := range matching.NumComponentes {
			saida.Componentes = append(saida.Componentes, componentePar{Nome: c.String(), ParaA: e.DeA[c], ParaB: e.DeB[c]})
		}
	}
	responderJSON(w, http.StatusOK, saida)
}
