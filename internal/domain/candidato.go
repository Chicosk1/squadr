package domain

// Candidato é um usuário que passou nos filtros rígidos para outro usuário em um jogo,
// com os sinais calculados pela fonte de candidatos (D-029). Fica no domain para que
// matching dependa só da interface, e não do pacote graph (D-034).
type Candidato struct {
	ID             string
	HorariosComuns int
	Complemento    float64 // máximo de COMPLEMENTA entre as funções dos dois; 0 se não houver
	SinalSocial    int
}

// OpcoesCandidatos controla a busca de candidatos.
type OpcoesCandidatos struct {
	Limite      int
	MinHorarios int
	MaxDiffElo  int
	SinalSocial bool
}

func OpcoesPadrao() OpcoesCandidatos {
	return OpcoesCandidatos{Limite: 300, MinHorarios: 1, MaxDiffElo: MaxDiffEloDuo, SinalSocial: true}
}
