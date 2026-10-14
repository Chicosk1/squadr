package seed

import (
	"slices"

	"github.com/Chicosk1/squadr/internal/domain"
)

const MaxDiffEloDuo = 2

const PersonaMulti = "MULTI-CASUAL"

func PassaFiltros(a, b domain.Perfil, jogo domain.Jogo) bool {
	pa, okA := a.PerfilDoJogo(jogo)
	pb, okB := b.PerfilDoJogo(jogo)
	if !okA || !okB {
		return false
	}
	if a.Regiao != b.Regiao {
		return false
	}
	if !slices.Contains(a.FaixasAceitas, b.FaixaEtaria) || !slices.Contains(b.FaixasAceitas, a.FaixaEtaria) {
		return false
	}
	if !temEmComum(a.Idiomas, b.Idiomas) {
		return false
	}

	switch jogo {
	case domain.ARK:
		return pa.Versao == pb.Versao &&
			grupoPlataforma(pa.Plataforma) == grupoPlataforma(pb.Plataforma) &&
			modoCompativel(pa.Modo, pb.Modo)
	case domain.Rust:
		return pa.Versao == pb.Versao
	case domain.LoL:
		return filasCompativeis(pa, pb)
	}
	return false
}

func temEmComum(a, b []string) bool {
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}

func grupoPlataforma(p string) string {
	if p == "pc_steam" || p == "pc_epic" {
		return "pc"
	}
	return p
}

func modoCompativel(a, b string) bool {
	return a == b || a == "ambos" || b == "ambos"
}

func filasCompativeis(a, b domain.PerfilJogo) bool {
	for _, fila := range a.Filas {
		if !slices.Contains(b.Filas, fila) {
			continue
		}
		if fila == "solo_duo" && !eloCompativel(a.Elo, b.Elo) {
			continue
		}
		return true
	}
	return false
}

func eloCompativel(a, b int) bool {
	if a == domain.SemRanked || b == domain.SemRanked {
		return false
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= MaxDiffEloDuo
}

// ---------- Gabarito de relevância ----------

type par [2]string

func chave(a, b string) par {
	if a > b {
		a, b = b, a
	}
	return par{a, b}
}

var grausPares = map[par]int{
	chave("ARK-PVP-RAIDER", "ARK-PVP-RAIDER"):   2,
	chave("ARK-PVE-BREEDER", "ARK-PVE-BREEDER"): 2,
	chave("ARK-PVE-BUILDER", "ARK-PVE-BUILDER"): 2,
	chave("ARK-PVE-BREEDER", "ARK-PVE-BUILDER"): 2,

	chave("RUST-GRINDER", "RUST-GRINDER"):   2,
	chave("RUST-CASUAL", "RUST-CASUAL"):     2,
	chave("RUST-APRENDIZ", "RUST-APRENDIZ"): 2,
	chave("RUST-APRENDIZ", "RUST-CASUAL"):   1,

	chave("LOL-SUPORTE", "LOL-SUPORTE"): 1,
	chave("LOL-ADC", "LOL-ADC"):         1,
	chave("LOL-ARAM", "LOL-ARAM"):       1,
	chave("LOL-ADC", "LOL-SUPORTE"):     2,
}

func Grau(a, b domain.Perfil, jogo domain.Jogo) int {
	pa, okA := a.PerfilDoJogo(jogo)
	pb, okB := b.PerfilDoJogo(jogo)
	if !okA || !okB {
		return 0
	}

	multiA, multiB := a.Persona == PersonaMulti, b.Persona == PersonaMulti
	switch {
	case multiA && multiB:
		return 1
	case multiA:
		return grauMulti(pa, pb)
	case multiB:
		return grauMulti(pb, pa)
	}
	return grausPares[chave(pa.PersonaJogo, pb.PersonaJogo)]
}

func grauMulti(multi, outro domain.PerfilJogo) int {
	if multi.PersonaJogo == outro.PersonaJogo {
		return 1
	}
	return 0
}
