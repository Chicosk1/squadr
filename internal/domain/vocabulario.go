package domain

// Vocabulário do formulário (docs/questionario.md): o valor gravado no grafo e o rótulo
// mostrado ao usuário. É a fonte única para a validação da API, para as opções do
// formulário e para o teste que confere o gerador (D-038). A ordem importa: nas
// perguntas ordinais, o score usa a posição na lista.

type Opcao struct {
	Valor  string
	Rotulo string
}

const TantoFaz = "tanto_faz"

var (
	// Seção G
	FaixasEtarias   = []Opcao{{"18-24", "18 a 24"}, {"25-34", "25 a 34"}, {"35+", "35 ou mais"}}
	Regioes         = []Opcao{{"AS", "América do Sul"}, {"AN", "América do Norte"}, {"EU", "Europa"}}
	Idiomas         = []Opcao{{"pt", "Português"}, {"en", "Inglês"}, {"es", "Espanhol"}}
	UsosVoz         = []Opcao{{"sim", "Sim"}, {"as_vezes", "Às vezes"}, {"nao", "Não"}}
	Compromissos    = []Opcao{{"fixo", "Parceiro ou grupo fixo"}, {"eventual", "Partidas eventuais"}, {TantoFaz, "Tanto faz"}}
	RotulosDias     = []string{"Seg", "Ter", "Qua", "Qui", "Sex", "Sáb", "Dom"} // mesma ordem de Dias
	RotulosPeriodos = []string{"Madrugada (0h–6h)", "Manhã (6h–12h)", "Tarde (12h–18h)", "Noite (18h–24h)"}

	// Comuns aos jogos
	NiveisDesejados = []Opcao{{"parecido", "Parecido"}, {"mais_experiente", "Mais experiente"}, {TantoFaz, "Tanto faz"}}
	FaixasHoras     = []Opcao{{"<100", "Menos de 100"}, {"100-500", "100 a 500"}, {"500-2000", "500 a 2.000"}, {"2000+", "Mais de 2.000"}}

	// Seção A — ARK
	VersoesARK     = []Opcao{{"ase", "Survival Evolved (ASE)"}, {"asa", "Survival Ascended (ASA)"}}
	PlataformasARK = []Opcao{{"pc_steam", "PC (Steam)"}, {"pc_epic", "PC (Epic)"}, {"xbox", "Xbox"}, {"playstation", "PlayStation"}}
	ModosARK       = []Opcao{{"pve", "PvE"}, {"pvp", "PvP"}, {"ambos", "Ambos"}}
	ServidoresARK  = []Opcao{{"oficial", "Oficial"}, {"rates_baixas", "Não oficial, rates baixas"}, {"rates_altas", "Não oficial, rates altas"}, {TantoFaz, "Tanto faz"}}
	TamanhosTribo  = []Opcao{{"2-3", "2 a 3"}, {"4-6", "4 a 6"}, {"7+", "7 ou mais"}, {TantoFaz, "Tanto faz"}}

	// Seção R — Rust
	VersoesRust    = []Opcao{{"pc", "PC"}, {"console", "Console Edition"}}
	ServidoresRust = []Opcao{{"oficial", "Oficial"}, {"comunitario", "Comunitário vanilla"}, {"modificado", "Modificado"}, {"pve", "PvE"}}
	LimitesGrupo   = []Opcao{{"solo_duo", "Solo/Duo"}, {"trio", "Trio"}, {"quad", "Quad"}, {"sem_limite", "Sem limite"}}
	Wipes          = []Opcao{{"sempre", "Sempre"}, {"as_vezes", "Às vezes"}, {"raramente", "Raramente"}}

	// Seção L — League of Legends
	Filas = []Opcao{{"solo_duo", "Solo/Duo"}, {"flex", "Flex"}, {"normal", "Normal"}, {"aram", "ARAM"}, {"clash", "Clash"}}
	// Elos tem o rótulo de cada valor de 0 (Ferro) a 9 (Desafiante); SemRanked fica fora.
	Elos = []string{"Ferro", "Bronze", "Prata", "Ouro", "Platina", "Esmeralda", "Diamante", "Mestre", "Grão-Mestre", "Desafiante"}

	RotulosJogos = map[Jogo]string{ARK: "ARK", Rust: "Rust", LoL: "League of Legends"}
)

// Valores devolve os valores das opções, na ordem, sem "tanto_faz".
// O score usa essa ordem nas perguntas ordinais.
func Valores(opcoes []Opcao) []string {
	saida := make([]string, 0, len(opcoes))
	for _, o := range opcoes {
		if o.Valor != TantoFaz {
			saida = append(saida, o.Valor)
		}
	}
	return saida
}

func valido(opcoes []Opcao, v string) bool {
	for _, o := range opcoes {
		if o.Valor == v {
			return true
		}
	}
	return false
}

// ---------- Nível ----------

var nivelPorHoras = map[string]float64{"<100": 0.125, "100-500": 0.375, "500-2000": 0.625, "2000+": 0.875}

// CalcularNivel normaliza a experiência para [0, 1]: o centro da faixa de horas no ARK
// e no Rust; (elo + 1) / 10 no LoL, com 0 para quem não joga ranqueada.
func (pj PerfilJogo) CalcularNivel() float64 {
	if pj.Jogo == LoL {
		if pj.Elo == SemRanked {
			return 0
		}
		return float64(pj.Elo+1) / 10
	}
	return nivelPorHoras[pj.Horas]
}

// Rotulo devolve o rótulo de um valor, ou o próprio valor se ele não estiver na lista.
func Rotulo(opcoes []Opcao, valor string) string {
	for _, o := range opcoes {
		if o.Valor == valor {
			return o.Rotulo
		}
	}
	return valor
}

// RotuloElo devolve o nome do elo, ou "Sem ranked".
func RotuloElo(elo int) string {
	if elo < 0 || elo >= len(Elos) {
		return "Sem ranked"
	}
	return Elos[elo]
}
