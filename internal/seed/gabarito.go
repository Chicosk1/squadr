package seed

import "github.com/Chicosk1/squadr/internal/domain"

const PersonaMulti = "MULTI-CASUAL"

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
