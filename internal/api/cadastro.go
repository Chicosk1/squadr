package api

import (
	"fmt"
	"strconv"

	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/matching"
)

// cadastro é o corpo de POST /usuarios: as respostas do formulário (questionario.md).
// Não tem id, nível nem persona: o id é gerado pelo servidor, o nível é calculado
// e a persona só existe nos dados mock.
type cadastro struct {
	MaiorDeIdade         bool           `json:"maiorDeIdade"` // G01
	Nome                 string         `json:"nome"`
	FaixaEtaria          string         `json:"faixaEtaria"`          // G02
	FaixasAceitas        []string       `json:"faixasAceitas"`        // G03
	Regiao               string         `json:"regiao"`               // G04
	Idiomas              []string       `json:"idiomas"`              // G05
	Horarios             []string       `json:"horarios"`             // G06, ids como "sex-noite"
	UsaVoz               string         `json:"usaVoz"`               // G07
	ExigeVoz             bool           `json:"exigeVoz"`             // G08
	ToleranciaToxicidade int            `json:"toleranciaToxicidade"` // G09
	Compromisso          string         `json:"compromisso"`          // G10
	Pesos                *pesosJSON     `json:"pesos"`                // P01 a P06, opcional
	Jogos                []jogoCadastro `json:"jogos"`                // G11 e seções A, R, L
}

// pesosJSON: um peso ausente ou 0 assume o padrão 3 (D-008).
type pesosJSON struct {
	Horario     int `json:"horario"`
	Intensidade int `json:"intensidade"`
	Modo        int `json:"modo"`
	Comunicacao int `json:"comunicacao"`
	Funcao      int `json:"funcao"`
	Nivel       int `json:"nivel"`
}

type jogoCadastro struct {
	Jogo             string `json:"jogo"`
	Intensidade      int    `json:"intensidade"`
	NivelDesejado    string `json:"nivelDesejado"`
	FuncaoPrincipal  string `json:"funcaoPrincipal"`
	FuncaoSecundaria string `json:"funcaoSecundaria"`

	Horas    string `json:"horas"`    // ARK e Rust
	Versao   string `json:"versao"`   // ARK e Rust
	Servidor string `json:"servidor"` // ARK e Rust

	Plataforma   string `json:"plataforma"`   // ARK
	Modo         string `json:"modo"`         // ARK
	TamanhoTribo string `json:"tamanhoTribo"` // ARK

	LimiteGrupo string `json:"limiteGrupo"` // Rust
	Wipe        string `json:"wipe"`        // Rust

	Elo   *int     `json:"elo"`   // LoL: ausente ou -1 = sem ranked
	Filas []string `json:"filas"` // LoL
}

// paraPerfil converte e valida. Os campos que não pertencem ao jogo são ignorados.
func (c cadastro) paraPerfil() (domain.Perfil, []string) {
	var erros []string
	if !c.MaiorDeIdade {
		erros = append(erros, "o squadr é só para maiores de 18 anos")
	}

	p := domain.Perfil{
		Nome:                 c.Nome,
		FaixaEtaria:          c.FaixaEtaria,
		FaixasAceitas:        c.FaixasAceitas,
		Regiao:               c.Regiao,
		Idiomas:              c.Idiomas,
		UsaVoz:               c.UsaVoz,
		ExigeVoz:             c.ExigeVoz,
		ToleranciaToxicidade: c.ToleranciaToxicidade,
		Compromisso:          c.Compromisso,
		Pesos:                pesos(c.Pesos),
	}
	for _, h := range c.Horarios {
		bit, ok := domain.BitDoHorario(h)
		if !ok {
			erros = append(erros, fmt.Sprintf("horário inválido: %q", h))
			continue
		}
		p.Disponibilidade |= 1 << bit
	}

	for _, jc := range c.Jogos {
		pj := domain.PerfilJogo{
			Jogo:             domain.Jogo(jc.Jogo),
			Intensidade:      jc.Intensidade,
			NivelDesejado:    jc.NivelDesejado,
			FuncaoPrincipal:  jc.FuncaoPrincipal,
			FuncaoSecundaria: jc.FuncaoSecundaria,
			Elo:              domain.SemRanked,
		}
		switch pj.Jogo {
		case domain.ARK:
			pj.Horas, pj.Versao, pj.Servidor = jc.Horas, jc.Versao, jc.Servidor
			pj.Plataforma, pj.Modo, pj.TamanhoTribo = jc.Plataforma, jc.Modo, jc.TamanhoTribo
		case domain.Rust:
			pj.Horas, pj.Versao, pj.Servidor = jc.Horas, jc.Versao, jc.Servidor
			pj.LimiteGrupo, pj.Wipe = jc.LimiteGrupo, jc.Wipe
		case domain.LoL:
			if jc.Elo != nil {
				pj.Elo = *jc.Elo
			}
			pj.Filas = jc.Filas
		}
		pj.Nivel = pj.CalcularNivel()
		p.Jogos = append(p.Jogos, pj)
	}

	erros = append(erros, p.Validar()...)
	return p, erros
}

func pesos(pj *pesosJSON) domain.Pesos {
	if pj == nil {
		pj = &pesosJSON{}
	}
	padrao := func(v int) int {
		if v == 0 {
			return matching.PesoPadrao
		}
		return v
	}
	return domain.Pesos{
		Horario: padrao(pj.Horario), Intensidade: padrao(pj.Intensidade), Modo: padrao(pj.Modo),
		Comunicacao: padrao(pj.Comunicacao), Funcao: padrao(pj.Funcao), Nivel: padrao(pj.Nivel),
	}
}

// ---------- Dados do formulário ----------

type jogoPagina struct {
	ID      domain.Jogo
	Rotulo  string
	Funcoes []domain.Opcao
}

type diaPagina struct {
	Rotulo   string
	Horarios []string // ids na ordem dos períodos
}

type dadosPagina struct {
	FaixasEtarias, Regioes, Idiomas, UsosVoz, Compromissos []domain.Opcao
	NiveisDesejados, FaixasHoras                           []domain.Opcao
	VersoesARK, PlataformasARK, ModosARK, ServidoresARK    []domain.Opcao
	TamanhosTribo, VersoesRust, ServidoresRust             []domain.Opcao
	LimitesGrupo, Wipes, Filas, Elos, Pesos                []domain.Opcao
	Escala                                                 []int // 1 a 5: intensidade e tolerância
	Periodos                                               []string
	Dias                                                   []diaPagina
	Jogos                                                  []jogoPagina
}

func montarDadosPagina(rotulos map[string]string) (dadosPagina, error) {
	d := dadosPagina{
		FaixasEtarias: domain.FaixasEtarias, Regioes: domain.Regioes, Idiomas: domain.Idiomas,
		UsosVoz: domain.UsosVoz, Compromissos: domain.Compromissos,
		NiveisDesejados: domain.NiveisDesejados, FaixasHoras: domain.FaixasHoras,
		VersoesARK: domain.VersoesARK, PlataformasARK: domain.PlataformasARK, ModosARK: domain.ModosARK,
		ServidoresARK: domain.ServidoresARK, TamanhosTribo: domain.TamanhosTribo,
		VersoesRust: domain.VersoesRust, ServidoresRust: domain.ServidoresRust,
		LimitesGrupo: domain.LimitesGrupo, Wipes: domain.Wipes, Filas: domain.Filas,
		Periodos: domain.RotulosPeriodos,
		Escala:   []int{1, 2, 3, 4, 5},
		Pesos: []domain.Opcao{
			{Valor: "horario", Rotulo: "Horários em comum"},
			{Valor: "intensidade", Rotulo: "Intensidade e compromisso"},
			{Valor: "modo", Rotulo: "Modo e tipo de servidor"},
			{Valor: "comunicacao", Rotulo: "Comunicação (voz e toxicidade)"},
			{Valor: "funcao", Rotulo: "Função complementar"},
			{Valor: "nivel", Rotulo: "Nível do parceiro"},
		},
	}
	for i, nome := range domain.Elos {
		d.Elos = append(d.Elos, domain.Opcao{Valor: strconv.Itoa(i), Rotulo: nome})
	}
	for i, rotulo := range domain.RotulosDias {
		dia := diaPagina{Rotulo: rotulo}
		for p := range domain.PeriodosPorDia {
			dia.Horarios = append(dia.Horarios, domain.IDHorario(i*domain.PeriodosPorDia+p))
		}
		d.Dias = append(d.Dias, dia)
	}
	for _, j := range domain.Jogos {
		jp := jogoPagina{ID: j, Rotulo: domain.RotulosJogos[j]}
		for _, id := range domain.Funcoes[j] {
			nome, ok := rotulos[id]
			if !ok || nome == "" {
				return dadosPagina{}, fmt.Errorf("função %s sem rótulo no grafo: rode go run ./cmd/migrar", id)
			}
			jp.Funcoes = append(jp.Funcoes, domain.Opcao{Valor: id, Rotulo: nome})
		}
		d.Jogos = append(d.Jogos, jp)
	}
	return d, nil
}
