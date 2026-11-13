package domain

import (
	"fmt"
	"slices"
	"strings"
)

// Filtro é o resultado de um filtro rígido para um par em um jogo (D-043). Motivo diz por que
// o filtro barrou o par e fica vazio quando ele passa.
type Filtro struct {
	Nome   string
	Passou bool
	Motivo string
}

// Filtros aplica os filtros rígidos ao par (a, b) no jogo, um por um, com as mesmas regras da
// consulta de candidatos (graph/candidatos.go), e devolve todos, inclusive os que passaram.
// É o oráculo do cmd/verificar (D-030): continua escrito em Go, sem nada em comum com o Cypher.
// O mínimo de horários em comum fica de fora, porque na consulta é um parâmetro (OpcoesCandidatos).
func Filtros(a, b Perfil, jogo Jogo) []Filtro {
	var saida []Filtro
	add := func(nome string, passou bool, motivo string) {
		if passou {
			motivo = ""
		}
		saida = append(saida, Filtro{Nome: nome, Passou: passou, Motivo: motivo})
	}

	pa, okA := a.PerfilDoJogo(jogo)
	pb, okB := b.PerfilDoJogo(jogo)
	var fora []string
	if !okA {
		fora = append(fora, nomeOuID(a))
	}
	if !okB {
		fora = append(fora, nomeOuID(b))
	}
	add("Os dois jogam "+RotulosJogos[jogo], okA && okB, strings.Join(fora, " e ")+" não joga "+RotulosJogos[jogo])

	add("Mesma região", a.Regiao == b.Regiao,
		Rotulo(Regioes, a.Regiao)+" × "+Rotulo(Regioes, b.Regiao))

	var faixas []string
	if !slices.Contains(a.FaixasAceitas, b.FaixaEtaria) {
		faixas = append(faixas, fmt.Sprintf("%s não aceita jogar com %s", nomeOuID(a), Rotulo(FaixasEtarias, b.FaixaEtaria)))
	}
	if !slices.Contains(b.FaixasAceitas, a.FaixaEtaria) {
		faixas = append(faixas, fmt.Sprintf("%s não aceita jogar com %s", nomeOuID(b), Rotulo(FaixasEtarias, a.FaixaEtaria)))
	}
	add("Faixa etária aceita pelos dois", len(faixas) == 0, strings.Join(faixas, "; "))

	add("Idioma em comum", temEmComum(a.Idiomas, b.Idiomas),
		rotulos(Idiomas, a.Idiomas)+" × "+rotulos(Idiomas, b.Idiomas))

	if !okA || !okB {
		return saida // sem o perfil do jogo dos dois, os filtros do jogo não se aplicam
	}
	switch jogo {
	case ARK:
		add("Mesma versão", pa.Versao == pb.Versao,
			Rotulo(VersoesARK, pa.Versao)+" × "+Rotulo(VersoesARK, pb.Versao))
		add("Plataforma compatível", grupoPlataforma(pa.Plataforma) == grupoPlataforma(pb.Plataforma),
			Rotulo(PlataformasARK, pa.Plataforma)+" × "+Rotulo(PlataformasARK, pb.Plataforma))
		add("Modo compatível", modoCompativel(pa.Modo, pb.Modo),
			Rotulo(ModosARK, pa.Modo)+" × "+Rotulo(ModosARK, pb.Modo))
	case Rust:
		add("Mesma versão", pa.Versao == pb.Versao,
			Rotulo(VersoesRust, pa.Versao)+" × "+Rotulo(VersoesRust, pb.Versao))
	case LoL:
		passou, motivo := filasCompativeis(pa, pb)
		add("Fila em comum", passou, motivo)
	}
	return saida
}

// PassaFiltros é verdadeiro quando o par passa em todos os filtros rígidos do jogo.
func PassaFiltros(a, b Perfil, jogo Jogo) bool {
	for _, f := range Filtros(a, b, jogo) {
		if !f.Passou {
			return false
		}
	}
	return true
}

func nomeOuID(p Perfil) string {
	if p.Nome != "" {
		return p.Nome
	}
	return p.ID
}

func rotulos(opcoes []Opcao, valores []string) string {
	saida := make([]string, len(valores))
	for i, v := range valores {
		saida[i] = Rotulo(opcoes, v)
	}
	return strings.Join(saida, ", ")
}

func temEmComum(a, b []string) bool {
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}

// grupoPlataforma junta Steam e Epic: no ARK, as duas jogam nos mesmos servidores de PC.
func grupoPlataforma(p string) string {
	if p == "pc_steam" || p == "pc_epic" {
		return "pc"
	}
	return p
}

func modoCompativel(a, b string) bool {
	return a == b || a == "ambos" || b == "ambos"
}

// filasCompativeis exige uma fila em comum. Solo/Duo só conta se os dois têm elo e a
// diferença é de no máximo MaxDiffEloDuo (a regra da Riot para jogar em dupla).
func filasCompativeis(a, b PerfilJogo) (bool, string) {
	duoBarrado := ""
	for _, fila := range a.Filas {
		if !slices.Contains(b.Filas, fila) {
			continue
		}
		if fila == "solo_duo" && !eloCompativel(a.Elo, b.Elo) {
			duoBarrado = motivoDuo(a.Elo, b.Elo)
			continue
		}
		return true, ""
	}
	if duoBarrado != "" {
		return false, "a única fila em comum é Solo/Duo, e " + duoBarrado
	}
	return false, rotulos(Filas, a.Filas) + " × " + rotulos(Filas, b.Filas)
}

func eloCompativel(a, b int) bool {
	if a == SemRanked || b == SemRanked {
		return false
	}
	return max(a-b, b-a) <= MaxDiffEloDuo
}

func motivoDuo(a, b int) string {
	if a == SemRanked || b == SemRanked {
		return "ela exige elo nos dois"
	}
	return fmt.Sprintf("a distância de elo é %d (%s × %s), acima do máximo de %d",
		max(a-b, b-a), RotuloElo(a), RotuloElo(b), MaxDiffEloDuo)
}
