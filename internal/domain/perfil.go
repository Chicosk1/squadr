package domain

const SemRanked = -1

type Pesos struct {
	Horario, Intensidade, Modo, Comunicacao, Funcao, Nivel int
}

type PerfilJogo struct {
	Jogo          Jogo
	Intensidade   int     // 1 a 5
	Nivel         float64 // 0 a 1
	NivelDesejado string  // parecido | mais_experiente | tanto_faz

	FuncaoPrincipal  string
	FuncaoSecundaria string

	// ARK e Rust
	Horas    string
	Versao   string
	Servidor string

	// ARK
	Plataforma   string
	Modo         string
	TamanhoTribo string

	// Rust
	LimiteGrupo string
	Wipe        string

	// LoL
	Elo   int // 0 (Ferro) a 9 (Desafiante), ou SemRanked
	Filas []string

	PersonaJogo string
}

type Perfil struct {
	ID                   string
	Nome                 string
	FaixaEtaria          string
	FaixasAceitas        []string
	Regiao               string
	Idiomas              []string
	UsaVoz               string // sim | as_vezes | nao
	ExigeVoz             bool
	ToleranciaToxicidade int // 1 a 5
	Compromisso          string
	Pesos                Pesos
	Disponibilidade      uint32 // 28 bits, um por faixa de horário
	Jogos                []PerfilJogo

	Persona string
}

func (p Perfil) PerfilDoJogo(j Jogo) (PerfilJogo, bool) {
	for _, pj := range p.Jogos {
		if pj.Jogo == j {
			return pj, true
		}
	}
	return PerfilJogo{}, false
}

func (p Perfil) HorarioIDs() []string {
	var ids []string
	for bit := 0; bit < NumHorarios; bit++ {
		if p.Disponibilidade&(1<<bit) != 0 {
			ids = append(ids, IDHorario(bit))
		}
	}
	return ids
}

func (p Perfil) Propriedades() map[string]any {
	return map[string]any{
		"id":                   p.ID,
		"nome":                 p.Nome,
		"faixaEtaria":          p.FaixaEtaria,
		"faixasAceitas":        p.FaixasAceitas,
		"regiao":               p.Regiao,
		"idiomas":              p.Idiomas,
		"usaVoz":               p.UsaVoz,
		"exigeVoz":             p.ExigeVoz,
		"toleranciaToxicidade": p.ToleranciaToxicidade,
		"compromisso":          p.Compromisso,
		"pesoHorario":          p.Pesos.Horario,
		"pesoIntensidade":      p.Pesos.Intensidade,
		"pesoModo":             p.Pesos.Modo,
		"pesoComunicacao":      p.Pesos.Comunicacao,
		"pesoFuncao":           p.Pesos.Funcao,
		"pesoNivel":            p.Pesos.Nivel,
		"persona":              p.Persona,
	}
}

func (pj PerfilJogo) Propriedades() map[string]any {
	props := map[string]any{
		"intensidade":   pj.Intensidade,
		"nivel":         pj.Nivel,
		"nivelDesejado": pj.NivelDesejado,
		"personaJogo":   pj.PersonaJogo,
	}
	switch pj.Jogo {
	case ARK:
		props["horas"] = pj.Horas
		props["versao"] = pj.Versao
		props["plataforma"] = pj.Plataforma
		props["modo"] = pj.Modo
		props["servidor"] = pj.Servidor
		props["tamanhoTribo"] = pj.TamanhoTribo
	case Rust:
		props["horas"] = pj.Horas
		props["versao"] = pj.Versao
		props["servidor"] = pj.Servidor
		props["limiteGrupo"] = pj.LimiteGrupo
		props["wipe"] = pj.Wipe
	case LoL:
		if pj.Elo != SemRanked {
			props["elo"] = pj.Elo
		}
		props["filas"] = pj.Filas
	}
	return props
}
