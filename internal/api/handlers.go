package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

// ---------- Respostas ----------

type respostaErro struct {
	Erro     string   `json:"erro"`
	Detalhes []string `json:"detalhes,omitempty"`
}

func responderJSON(w http.ResponseWriter, status int, corpo any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(corpo)
}

func responderErro(w http.ResponseWriter, status int, msg string, detalhes ...string) {
	responderJSON(w, status, respostaErro{Erro: msg, Detalhes: detalhes})
}

// erroInterno registra a causa e devolve uma mensagem genérica: detalhes do banco não vão ao cliente.
func (s *Servidor) erroInterno(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Printf("erro em %s %s: %v", r.Method, r.URL.Path, err)
	responderErro(w, http.StatusInternalServerError, "erro interno")
}

// decodificar lê um JSON com limite de tamanho e rejeita campos desconhecidos
// (por exemplo, "persona": o motor nunca recebe esse campo de fora).
func decodificar(w http.ResponseWriter, r *http.Request, destino any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, TamanhoMaxReq))
	dec.DisallowUnknownFields()
	if err := dec.Decode(destino); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("o corpo deve ter um único objeto JSON")
	}
	return nil
}

func jogoValido(v string) (domain.Jogo, bool) {
	j := domain.Jogo(v)
	_, ok := domain.Funcoes[j]
	return j, ok
}

// ---------- GET /saude ----------

func (s *Servidor) saude(w http.ResponseWriter, r *http.Request) {
	ctx, cancelar := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancelar()
	if err := s.repo.Ping(ctx); err != nil {
		responderJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "indisponivel", "banco": "erro"})
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"status": "ok", "banco": "ok"})
}

// ---------- GET /usuarios/{id}/recomendacoes?jogo=...&k=... ----------

type recomendacaoJSON struct {
	ID        string   `json:"id"`
	Nome      string   `json:"nome"`
	Score     float64  `json:"score"`
	Reciproco float64  `json:"reciproco"`
	ScoreAB   float64  `json:"scoreParaVoce"`
	ScoreBA   float64  `json:"scoreParaCandidato"`
	Social    float64  `json:"social"`
	Motivos   []string `json:"motivos"`
	// Componentes explica o score: os seis critérios nos dois sentidos, todos em [0, 1].
	Componentes []componenteJSON `json:"componentes"`
}

type componenteJSON struct {
	Nome          string  `json:"nome"`
	ParaVoce      float64 `json:"paraVoce"`      // c(A→B): quanto o candidato atende você
	ParaCandidato float64 `json:"paraCandidato"` // c(B→A): quanto você atende o candidato
}

func componentes(deA, deB matching.Componentes) []componenteJSON {
	saida := make([]componenteJSON, matching.NumComponentes)
	for c := range matching.NumComponentes {
		saida[c] = componenteJSON{Nome: c.String(), ParaVoce: deA[c], ParaCandidato: deB[c]}
	}
	return saida
}

type respostaRecomendacoes struct {
	Usuario       string             `json:"usuario"`
	Jogo          domain.Jogo        `json:"jogo"`
	Recomendacoes []recomendacaoJSON `json:"recomendacoes"`
}

func (s *Servidor) recomendacoes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	jogo, ok := jogoValido(r.URL.Query().Get("jogo"))
	if !ok {
		responderErro(w, http.StatusBadRequest, "parâmetro jogo deve ser ark, rust ou lol")
		return
	}
	k := KPadrao
	if v := r.URL.Query().Get("k"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > KMaximo {
			responderErro(w, http.StatusBadRequest, "parâmetro k deve ser um inteiro de 1 a "+strconv.Itoa(KMaximo))
			return
		}
		k = n
	}

	recs, err := s.rec.Recomendar(r.Context(), id, jogo, k)
	switch {
	case errors.Is(err, domain.ErrNaoEncontrado):
		responderErro(w, http.StatusNotFound, "usuário não encontrado")
		return
	case errors.Is(err, matching.ErrNaoJoga):
		responderErro(w, http.StatusUnprocessableEntity, "o usuário não joga "+string(jogo))
		return
	case err != nil:
		s.erroInterno(w, r, err)
		return
	}

	saida := respostaRecomendacoes{Usuario: id, Jogo: jogo, Recomendacoes: make([]recomendacaoJSON, len(recs))}
	for i, rec := range recs {
		saida.Recomendacoes[i] = recomendacaoJSON{
			ID: rec.ID, Nome: rec.Nome, Score: rec.Score, Reciproco: rec.Reciproco,
			ScoreAB: rec.ScoreAB, ScoreBA: rec.ScoreBA, Social: rec.Social, Motivos: rec.Motivos,
			Componentes: componentes(rec.DeA, rec.DeB),
		}
	}
	responderJSON(w, http.StatusOK, saida)
}

// ---------- POST /usuarios ----------

func novoID() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// "w" separa os usuários do formulário dos do seed ("u000001"...).
	return "w" + hex.EncodeToString(b), nil
}

func (s *Servidor) criarUsuario(w http.ResponseWriter, r *http.Request) {
	var c cadastro
	if err := decodificar(w, r, &c); err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido", err.Error())
		return
	}
	perfil, erros := c.paraPerfil()
	if len(erros) > 0 {
		responderErro(w, http.StatusBadRequest, "cadastro inválido", erros...)
		return
	}
	id, err := s.novoID()
	if err != nil {
		s.erroInterno(w, r, err)
		return
	}
	perfil.ID = id
	if err := s.repo.CriarUsuario(r.Context(), perfil); err != nil {
		s.erroInterno(w, r, err)
		return
	}
	w.Header().Set("Location", "/usuarios/"+id)
	responderJSON(w, http.StatusCreated, map[string]string{"id": id})
}

// ---------- POST /usuarios/{id}/curtidas e /recusas ----------

type pedidoInteracao struct {
	Para string `json:"para"`
	Jogo string `json:"jogo"`
}

func (s *Servidor) interagir(curtiu bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p pedidoInteracao
		if err := decodificar(w, r, &p); err != nil {
			responderErro(w, http.StatusBadRequest, "JSON inválido", err.Error())
			return
		}
		jogo, ok := jogoValido(p.Jogo)
		if !ok {
			responderErro(w, http.StatusBadRequest, "jogo deve ser ark, rust ou lol")
			return
		}
		de := r.PathValue("id")
		if p.Para == "" || p.Para == de {
			responderErro(w, http.StatusBadRequest, "para deve ser o id de outro usuário")
			return
		}

		match, err := s.repo.RegistrarInteracao(r.Context(), domain.Interacao{
			De: de, Para: p.Para, Jogo: jogo, Curtiu: curtiu, Em: time.Now(),
		})
		switch {
		case errors.Is(err, domain.ErrNaoEncontrado):
			responderErro(w, http.StatusNotFound, "usuário não encontrado ou não joga "+string(jogo))
			return
		case errors.Is(err, domain.ErrJaAvaliado):
			responderErro(w, http.StatusConflict, "este usuário já foi avaliado neste jogo")
			return
		case err != nil:
			s.erroInterno(w, r, err)
			return
		}
		responderJSON(w, http.StatusCreated, map[string]bool{"match": match})
	}
}

// ---------- GET / ----------

func (s *Servidor) paginaInicial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.pagina.Execute(w, s.dados); err != nil {
		s.log.Printf("renderizando a página: %v", err)
	}
}
