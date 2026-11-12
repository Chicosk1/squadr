package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Chicosk1/squadr/internal/domain"
)

// perfilPublico é o que outro usuário pode ver (D-042). É um tipo separado de domain.Perfil
// de propósito: um campo novo no domínio não aparece aqui sem alguém decidir que ele é público.
// Ficam de fora: persona (só existe para o gabarito), pesos e faixas aceitas (preferências
// privadas, que só servem ao motor) e o nível normalizado (aparece como horas ou elo).
type perfilPublico struct {
	ID                   string        `json:"id"`
	Nome                 string        `json:"nome"`
	FaixaEtaria          string        `json:"faixaEtaria"`
	Regiao               string        `json:"regiao"`
	Idiomas              []string      `json:"idiomas"`
	UsaVoz               string        `json:"usaVoz"`
	ExigeVoz             bool          `json:"exigeVoz"`
	ToleranciaToxicidade int           `json:"toleranciaToxicidade"`
	Compromisso          string        `json:"compromisso"`
	Horarios             []string      `json:"horarios"` // ids, como "sex-noite"
	Jogos                []jogoPublico `json:"jogos"`
}

type jogoPublico struct {
	Jogo             domain.Jogo `json:"jogo"`
	Rotulo           string      `json:"rotulo"`
	Intensidade      int         `json:"intensidade"`
	NivelDesejado    string      `json:"nivelDesejado"`
	FuncaoPrincipal  string      `json:"funcaoPrincipal"`
	FuncaoSecundaria string      `json:"funcaoSecundaria,omitempty"`
	Detalhes         []campo     `json:"detalhes"` // as respostas específicas do jogo, na ordem do formulário
}

type campo struct {
	Nome  string `json:"nome"`
	Valor string `json:"valor"`
}

func (s *Servidor) perfil(w http.ResponseWriter, r *http.Request) {
	p, err := s.repo.BuscarPerfil(r.Context(), r.PathValue("id"))
	switch {
	case errors.Is(err, domain.ErrNaoEncontrado):
		responderErro(w, http.StatusNotFound, "usuário não encontrado")
		return
	case err != nil:
		s.erroInterno(w, r, err)
		return
	}
	responderJSON(w, http.StatusOK, s.paraPublico(p))
}

func (s *Servidor) paraPublico(p domain.Perfil) perfilPublico {
	out := perfilPublico{
		ID:                   p.ID,
		Nome:                 p.Nome,
		FaixaEtaria:          domain.Rotulo(domain.FaixasEtarias, p.FaixaEtaria),
		Regiao:               domain.Rotulo(domain.Regioes, p.Regiao),
		UsaVoz:               domain.Rotulo(domain.UsosVoz, p.UsaVoz),
		ExigeVoz:             p.ExigeVoz,
		ToleranciaToxicidade: p.ToleranciaToxicidade,
		Compromisso:          domain.Rotulo(domain.Compromissos, p.Compromisso),
		Horarios:             p.HorarioIDs(),
		Idiomas:              make([]string, len(p.Idiomas)),
		Jogos:                make([]jogoPublico, 0, len(p.Jogos)),
	}
	for i, idioma := range p.Idiomas {
		out.Idiomas[i] = domain.Rotulo(domain.Idiomas, idioma)
	}
	for _, pj := range p.Jogos {
		out.Jogos = append(out.Jogos, s.jogoPublico(pj))
	}
	return out
}

func (s *Servidor) jogoPublico(pj domain.PerfilJogo) jogoPublico {
	j := jogoPublico{
		Jogo:            pj.Jogo,
		Rotulo:          domain.RotulosJogos[pj.Jogo],
		Intensidade:     pj.Intensidade,
		NivelDesejado:   domain.Rotulo(domain.NiveisDesejados, pj.NivelDesejado),
		FuncaoPrincipal: s.rotuloFuncao(pj.FuncaoPrincipal),
	}
	if pj.FuncaoSecundaria != "" {
		j.FuncaoSecundaria = s.rotuloFuncao(pj.FuncaoSecundaria)
	}
	add := func(nome, valor string) { j.Detalhes = append(j.Detalhes, campo{nome, valor}) }
	switch pj.Jogo {
	case domain.ARK:
		add("Versão", domain.Rotulo(domain.VersoesARK, pj.Versao))
		add("Plataforma", domain.Rotulo(domain.PlataformasARK, pj.Plataforma))
		add("Modo", domain.Rotulo(domain.ModosARK, pj.Modo))
		add("Servidor", domain.Rotulo(domain.ServidoresARK, pj.Servidor))
		add("Tamanho de tribo", domain.Rotulo(domain.TamanhosTribo, pj.TamanhoTribo))
		add("Horas de jogo", domain.Rotulo(domain.FaixasHoras, pj.Horas))
	case domain.Rust:
		add("Versão", domain.Rotulo(domain.VersoesRust, pj.Versao))
		add("Servidor", domain.Rotulo(domain.ServidoresRust, pj.Servidor))
		add("Limite de grupo", domain.Rotulo(domain.LimitesGrupo, pj.LimiteGrupo))
		add("Joga no wipe", domain.Rotulo(domain.Wipes, pj.Wipe))
		add("Horas de jogo", domain.Rotulo(domain.FaixasHoras, pj.Horas))
	case domain.LoL:
		add("Elo", domain.RotuloElo(pj.Elo))
		filas := make([]string, len(pj.Filas))
		for i, f := range pj.Filas {
			filas[i] = domain.Rotulo(domain.Filas, f)
		}
		add("Filas", strings.Join(filas, ", "))
	}
	add("Intensidade", strconv.Itoa(pj.Intensidade)+" de 5")
	return j
}

func (s *Servidor) rotuloFuncao(id string) string {
	if nome, ok := s.rotulos[id]; ok {
		return nome
	}
	return id
}
