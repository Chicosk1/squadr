package domain

import (
	"errors"
	"fmt"
	"math/bits"
	"slices"
	"unicode/utf8"
)

var (
	ErrNaoEncontrado = errors.New("não encontrado")
	ErrJaAvaliado    = errors.New("já avaliado neste jogo")
)

const TamanhoMaxNome = 40

// Validar devolve a lista de problemas do perfil; vazia quando ele é válido.
// As regras são as do questionario.md, com os valores de vocabulario.go (D-038).
// Não olha Persona nem PersonaJogo.
func (p Perfil) Validar() []string {
	var erros []string
	add := func(formato string, args ...any) { erros = append(erros, fmt.Sprintf(formato, args...)) }

	if n := utf8.RuneCountInString(p.Nome); n == 0 || n > TamanhoMaxNome {
		add("nome deve ter de 1 a %d caracteres", TamanhoMaxNome)
	}
	if !valido(FaixasEtarias, p.FaixaEtaria) {
		add("faixa etária inválida: %q", p.FaixaEtaria)
	}
	conferirLista(&erros, "faixas aceitas", FaixasEtarias, p.FaixasAceitas)
	if !valido(Regioes, p.Regiao) {
		add("região inválida: %q", p.Regiao)
	}
	conferirLista(&erros, "idiomas", Idiomas, p.Idiomas)
	if !valido(UsosVoz, p.UsaVoz) {
		add("uso de voz inválido: %q", p.UsaVoz)
	}
	if p.ToleranciaToxicidade < 1 || p.ToleranciaToxicidade > 5 {
		add("tolerância a toxicidade deve ser de 1 a 5")
	}
	if !valido(Compromissos, p.Compromisso) {
		add("compromisso inválido: %q", p.Compromisso)
	}
	for _, peso := range []int{p.Pesos.Horario, p.Pesos.Intensidade, p.Pesos.Modo,
		p.Pesos.Comunicacao, p.Pesos.Funcao, p.Pesos.Nivel} {
		if peso < 1 || peso > 5 {
			add("pesos devem ser de 1 a 5")
			break
		}
	}
	if p.Disponibilidade == 0 {
		add("marque pelo menos um horário")
	}
	if bits.Len32(p.Disponibilidade) > NumHorarios {
		add("disponibilidade fora da grade de %d horários", NumHorarios)
	}

	if len(p.Jogos) == 0 {
		add("escolha pelo menos um jogo")
	}
	vistos := map[Jogo]bool{}
	for _, pj := range p.Jogos {
		if vistos[pj.Jogo] {
			add("jogo repetido: %s", pj.Jogo)
			continue
		}
		vistos[pj.Jogo] = true
		erros = append(erros, pj.validar()...)
	}
	return erros
}

func (pj PerfilJogo) validar() []string {
	var erros []string
	add := func(formato string, args ...any) {
		erros = append(erros, fmt.Sprintf("%s: "+formato, append([]any{pj.Jogo}, args...)...))
	}
	opcao := func(nome string, opcoes []Opcao, v string) {
		if !valido(opcoes, v) {
			add("%s inválido: %q", nome, v)
		}
	}

	funcoes, ok := Funcoes[pj.Jogo]
	if !ok {
		return []string{fmt.Sprintf("jogo desconhecido: %q", pj.Jogo)}
	}
	if pj.Intensidade < 1 || pj.Intensidade > 5 {
		add("intensidade deve ser de 1 a 5")
	}
	opcao("nível desejado", NiveisDesejados, pj.NivelDesejado)
	if !slices.Contains(funcoes, pj.FuncaoPrincipal) {
		add("função principal inválida: %q", pj.FuncaoPrincipal)
	}
	if pj.FuncaoSecundaria != "" {
		if !slices.Contains(funcoes, pj.FuncaoSecundaria) {
			add("função secundária inválida: %q", pj.FuncaoSecundaria)
		} else if pj.FuncaoSecundaria == pj.FuncaoPrincipal {
			add("a função secundária deve ser diferente da principal")
		}
	}

	switch pj.Jogo {
	case ARK:
		opcao("horas", FaixasHoras, pj.Horas)
		opcao("versão", VersoesARK, pj.Versao)
		opcao("plataforma", PlataformasARK, pj.Plataforma)
		opcao("modo", ModosARK, pj.Modo)
		opcao("servidor", ServidoresARK, pj.Servidor)
		opcao("tamanho de tribo", TamanhosTribo, pj.TamanhoTribo)
	case Rust:
		opcao("horas", FaixasHoras, pj.Horas)
		opcao("versão", VersoesRust, pj.Versao)
		opcao("servidor", ServidoresRust, pj.Servidor)
		opcao("limite de grupo", LimitesGrupo, pj.LimiteGrupo)
		opcao("wipe", Wipes, pj.Wipe)
	case LoL:
		if pj.Elo != SemRanked && (pj.Elo < 0 || pj.Elo >= len(Elos)) {
			add("elo inválido: %d", pj.Elo)
		}
		if len(pj.Filas) == 0 {
			add("escolha pelo menos uma fila")
		}
		for _, f := range pj.Filas {
			opcao("fila", Filas, f)
		}
	}
	return erros
}

func conferirLista(erros *[]string, nome string, opcoes []Opcao, valores []string) {
	if len(valores) == 0 {
		*erros = append(*erros, "escolha pelo menos uma opção em "+nome)
		return
	}
	for i, v := range valores {
		if !valido(opcoes, v) {
			*erros = append(*erros, fmt.Sprintf("%s: valor inválido %q", nome, v))
		} else if slices.Contains(valores[:i], v) {
			*erros = append(*erros, fmt.Sprintf("%s: valor repetido %q", nome, v))
		}
	}
}
