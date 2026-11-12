// Package api expõe o motor por HTTP usando só net/http e o roteamento do Go 1.22+ (D-014).
// Os handlers dependem de interfaces: o cmd/api liga o *graph.Cliente e o *matching.Recomendador,
// e os testes usam implementações em memória.
package api

import (
	"context"
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

type Recomendador interface {
	Recomendar(ctx context.Context, id string, jogo domain.Jogo, k int) ([]matching.Recomendacao, error)
}

type Repositorio interface {
	Ping(ctx context.Context) error
	CriarUsuario(ctx context.Context, p domain.Perfil) error
	RegistrarInteracao(ctx context.Context, x domain.Interacao) (bool, error)
	// BuscarPerfil devolve domain.ErrNaoEncontrado (embrulhado) se o usuário não existe.
	BuscarPerfil(ctx context.Context, id string) (domain.Perfil, error)
}

const (
	KPadrao       = 10
	KMaximo       = 50
	TamanhoMaxReq = 64 << 10 // 64 KiB por corpo de requisição
)

type Servidor struct {
	rec     Recomendador
	repo    Repositorio
	pagina  *template.Template
	dados   dadosPagina
	rotulos map[string]string // id da função → nome, lido do grafo
	log     *log.Logger
	novoID  func() (string, error)
}

//go:embed pagina.html
var paginaHTML string

// Novo monta o servidor. rotulos é id → nome de cada função, lido do grafo; falha se faltar
// o rótulo de alguma função do domínio, para o formulário não sair com opções em branco.
func Novo(rec Recomendador, repo Repositorio, rotulos map[string]string, logger *log.Logger) (*Servidor, error) {
	dados, err := montarDadosPagina(rotulos)
	if err != nil {
		return nil, err
	}
	pagina, err := template.New("pagina").Parse(paginaHTML)
	if err != nil {
		return nil, err
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Servidor{rec: rec, repo: repo, pagina: pagina, dados: dados, rotulos: rotulos,
		log: logger, novoID: novoID}, nil
}

// Handler devolve as rotas com os middlewares. O de log fica por fora do de recuperação,
// para registrar também o 500 que a recuperação escreve depois de um panic.
func (s *Servidor) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.paginaInicial)
	mux.HandleFunc("GET /saude", s.saude)
	mux.HandleFunc("GET /usuarios/{id}", s.perfil)
	mux.HandleFunc("GET /usuarios/{id}/recomendacoes", s.recomendacoes)
	mux.HandleFunc("POST /usuarios", s.criarUsuario)
	mux.HandleFunc("POST /usuarios/{id}/curtidas", s.interagir(true))
	mux.HandleFunc("POST /usuarios/{id}/recusas", s.interagir(false))
	return s.registrarLog(s.recuperar(mux))
}
