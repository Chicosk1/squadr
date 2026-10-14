package seed

import (
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/Chicosk1/squadr/internal/domain"
)

type Faixa struct{ Min, Max int }

func (f Faixa) Sortear(r *rand.Rand) int { return f.Min + r.IntN(f.Max-f.Min+1) }

func entre(lo, hi int) Faixa { return Faixa{lo, hi} }
func exato(n int) Faixa      { return Faixa{n, n} }

type Opcao struct {
	Valor string
	Peso  int
}

type Opcoes []Opcao

func (o Opcoes) Sortear(r *rand.Rand) string {
	total := 0
	for _, op := range o {
		total += op.Peso
	}
	n := r.IntN(total)
	for _, op := range o {
		if n < op.Peso {
			return op.Valor
		}
		n -= op.Peso
	}
	return o[len(o)-1].Valor
}

func unico(valor string) Opcoes { return Opcoes{{valor, 1}} }

func igual(valores ...string) Opcoes {
	o := make(Opcoes, len(valores))
	for i, v := range valores {
		o[i] = Opcao{v, 1}
	}
	return o
}

func pesos(horario, intensidade, modo, comunicacao, funcao, nivel int) domain.Pesos {
	return domain.Pesos{
		Horario: horario, Intensidade: intensidade, Modo: modo,
		Comunicacao: comunicacao, Funcao: funcao, Nivel: nivel,
	}
}

func grade(partes ...string) []int {
	var bits []int
	for _, parte := range partes {
		dias, periodos, ok := strings.Cut(parte, ":")
		if !ok {
			panic("formato esperado dias:periodos, recebido " + parte)
		}
		for _, d := range expandirDias(dias) {
			for _, p := range strings.Split(periodos, ",") {
				bits = append(bits, d*domain.PeriodosPorDia+indice(domain.Periodos, p))
			}
		}
	}
	return bits
}

func expandirDias(texto string) []int {
	var dias []int
	for _, item := range strings.Split(texto, ",") {
		if ini, fim, ehIntervalo := strings.Cut(item, "-"); ehIntervalo {
			for d := indice(domain.Dias, ini); d <= indice(domain.Dias, fim); d++ {
				dias = append(dias, d)
			}
		} else {
			dias = append(dias, indice(domain.Dias, item))
		}
	}
	return dias
}

func indice(lista []string, item string) int {
	i := slices.Index(lista, item)
	if i < 0 {
		panic("valor inválido na definição de persona: " + item)
	}
	return i
}

type ModeloJogo struct {
	Persona          string
	Jogo             domain.Jogo
	Intensidade      Faixa
	NivelDesejado    Opcoes
	FuncaoPrincipal  string
	FuncaoSecundaria string

	// ARK e Rust
	Horas, Versao, Servidor Opcoes
	// ARK
	Plataforma, Modo, TamanhoTribo Opcoes
	// Rust
	LimiteGrupo, Wipe Opcoes
	// LoL
	Elo   Faixa
	Filas []string
}

type Persona struct {
	Nome        string
	Proporcao   float64
	Horarios    []int
	ProbHorario float64
	UsaVoz      Opcoes
	ExigeVoz    float64
	Toxicidade  Faixa
	Compromisso Opcoes
	Pesos       domain.Pesos

	Jogos        []ModeloJogo
	SortearJogos int
}

// ---------- ARK ----------

var arkRaider = ModeloJogo{
	Persona: "ARK-PVP-RAIDER", Jogo: domain.ARK,
	Intensidade: entre(4, 5), NivelDesejado: unico("parecido"),
	FuncaoPrincipal: "ark:pvp", FuncaoSecundaria: "ark:breeder",
	Horas: unico("2000+"), Versao: unico("asa"),
	Plataforma: Opcoes{{"pc_steam", 80}, {"xbox", 10}, {"playstation", 10}},
	Modo:       unico("pvp"), Servidor: unico("oficial"), TamanhoTribo: unico("7+"),
}

var arkBreeder = ModeloJogo{
	Persona: "ARK-PVE-BREEDER", Jogo: domain.ARK,
	Intensidade: entre(1, 2), NivelDesejado: unico("tanto_faz"),
	FuncaoPrincipal: "ark:breeder", FuncaoSecundaria: "ark:coletor",
	Horas: unico("500-2000"), Versao: unico("asa"), Plataforma: unico("pc_steam"),
	Modo: unico("pve"), Servidor: unico("rates_altas"), TamanhoTribo: unico("2-3"),
}

var arkBuilder = ModeloJogo{
	Persona: "ARK-PVE-BUILDER", Jogo: domain.ARK,
	Intensidade: entre(2, 3), NivelDesejado: unico("tanto_faz"),
	FuncaoPrincipal: "ark:builder", FuncaoSecundaria: "ark:coletor",
	Horas: unico("100-500"), Versao: Opcoes{{"asa", 70}, {"ase", 30}}, Plataforma: unico("pc_steam"),
	Modo: unico("pve"), Servidor: unico("rates_baixas"), TamanhoTribo: unico("4-6"),
}

// ---------- Rust ----------

var rustGrinder = ModeloJogo{
	Persona: "RUST-GRINDER", Jogo: domain.Rust,
	Intensidade: exato(5), NivelDesejado: unico("parecido"),
	FuncaoPrincipal: "rust:pvp", FuncaoSecundaria: "rust:raider",
	Horas: unico("2000+"), Versao: unico("pc"), Servidor: unico("oficial"),
	LimiteGrupo: Opcoes{{"quad", 50}, {"sem_limite", 50}}, Wipe: unico("sempre"),
}

var rustCasual = ModeloJogo{
	Persona: "RUST-CASUAL", Jogo: domain.Rust,
	Intensidade: exato(2), NivelDesejado: unico("tanto_faz"),
	FuncaoPrincipal: "rust:builder", FuncaoSecundaria: "rust:eletricista",
	Horas: unico("500-2000"), Versao: unico("pc"), Servidor: unico("comunitario"),
	LimiteGrupo: unico("solo_duo"), Wipe: unico("as_vezes"),
}

var rustAprendiz = ModeloJogo{
	Persona: "RUST-APRENDIZ", Jogo: domain.Rust,
	Intensidade: exato(3), NivelDesejado: unico("mais_experiente"),
	FuncaoPrincipal: "rust:coletor", FuncaoSecundaria: "rust:builder",
	Horas: unico("<100"), Versao: Opcoes{{"pc", 85}, {"console", 15}},
	Servidor:    igual("comunitario", "modificado"),
	LimiteGrupo: igual("solo_duo", "trio"), Wipe: unico("as_vezes"),
}

// ---------- League of Legends ----------

var lolSuporte = ModeloJogo{
	Persona: "LOL-SUPORTE", Jogo: domain.LoL,
	Intensidade: exato(4), NivelDesejado: unico("parecido"),
	FuncaoPrincipal: "lol:suporte", FuncaoSecundaria: "lol:mid",
	Elo: entre(3, 4), Filas: []string{"solo_duo", "flex"},
}

var lolAdc = ModeloJogo{
	Persona: "LOL-ADC", Jogo: domain.LoL,
	Intensidade: exato(4), NivelDesejado: unico("parecido"),
	FuncaoPrincipal: "lol:adc", FuncaoSecundaria: "lol:mid",
	Elo: entre(3, 4), Filas: []string{"solo_duo", "flex"},
}

var lolAram = ModeloJogo{
	Persona: "LOL-ARAM", Jogo: domain.LoL,
	Intensidade: entre(1, 2), NivelDesejado: unico("tanto_faz"),
	FuncaoPrincipal: "lol:preencher",
	Elo:             entre(-1, 2), Filas: []string{"aram", "normal", "clash"},
}

func casual(m ModeloJogo) ModeloJogo {
	m.Intensidade = entre(1, 2)
	return m
}

// ---------- As personas ----------

var Personas = []Persona{
	{
		Nome: "ARK-PVP-RAIDER", Proporcao: 12,
		Horarios: grade("seg-dom:noite,madrugada"),
		UsaVoz:   unico("sim"), ExigeVoz: 1,
		Toxicidade: entre(4, 5), Compromisso: unico("fixo"),
		Pesos: pesos(5, 3, 5, 3, 4, 3),
		Jogos: []ModeloJogo{arkRaider},
	},
	{
		Nome: "ARK-PVE-BREEDER", Proporcao: 12,
		Horarios: grade("sex:noite", "sab,dom:tarde,noite"),
		UsaVoz:   unico("as_vezes"), ExigeVoz: 0,
		Toxicidade: entre(1, 2), Compromisso: unico("eventual"),
		Pesos: pesos(3, 3, 5, 4, 3, 3),
		Jogos: []ModeloJogo{arkBreeder},
	},
	{
		Nome: "ARK-PVE-BUILDER", Proporcao: 11,
		Horarios: grade("qui,sex:noite", "sab,dom:tarde,noite"),
		UsaVoz:   unico("sim"), ExigeVoz: 0,
		Toxicidade: entre(2, 3), Compromisso: unico("fixo"),
		Pesos: pesos(4, 3, 3, 3, 5, 3),
		Jogos: []ModeloJogo{arkBuilder},
	},
	{
		Nome: "RUST-GRINDER", Proporcao: 10,
		Horarios: grade("qui:tarde,noite,madrugada", "sex-dom:noite,madrugada"),
		UsaVoz:   unico("sim"), ExigeVoz: 1,
		Toxicidade: exato(5), Compromisso: unico("fixo"),
		Pesos: pesos(5, 3, 3, 3, 3, 5),
		Jogos: []ModeloJogo{rustGrinder},
	},
	{
		Nome: "RUST-CASUAL", Proporcao: 10,
		Horarios: grade("seg-sex:noite"),
		UsaVoz:   unico("sim"), ExigeVoz: 0,
		Toxicidade: exato(2), Compromisso: unico("fixo"),
		Pesos: pesos(4, 3, 3, 4, 3, 3),
		Jogos: []ModeloJogo{rustCasual},
	},
	{
		Nome: "RUST-APRENDIZ", Proporcao: 10,
		Horarios: grade("ter-sex:noite", "sab:tarde"),
		UsaVoz:   unico("sim"), ExigeVoz: 0,
		Toxicidade: exato(3), Compromisso: unico("fixo"),
		Pesos: pesos(3, 3, 3, 3, 4, 5),
		Jogos: []ModeloJogo{rustAprendiz},
	},
	{
		Nome: "LOL-SUPORTE", Proporcao: 10,
		Horarios: grade("seg-sex:noite", "dom:noite"),
		UsaVoz:   unico("sim"), ExigeVoz: 0,
		Toxicidade: exato(2), Compromisso: unico("fixo"),
		Pesos: pesos(3, 3, 3, 3, 5, 4),
		Jogos: []ModeloJogo{lolSuporte},
	},
	{
		Nome: "LOL-ADC", Proporcao: 10,
		Horarios: grade("seg-sex:noite", "sab:noite"),
		UsaVoz:   unico("sim"), ExigeVoz: 1,
		Toxicidade: entre(3, 4), Compromisso: unico("fixo"),
		Pesos: pesos(4, 3, 3, 3, 5, 3),
		Jogos: []ModeloJogo{lolAdc},
	},
	{
		Nome: "LOL-ARAM", Proporcao: 10,
		Horarios: grade("sex-dom:tarde,noite"),
		UsaVoz:   unico("as_vezes"), ExigeVoz: 0,
		Toxicidade: exato(2), Compromisso: unico("eventual"),
		Pesos: pesos(4, 3, 4, 3, 3, 3),
		Jogos: []ModeloJogo{lolAram},
	},
	{
		Nome: "MULTI-CASUAL", Proporcao: 5,
		ProbHorario: 0.30,
		UsaVoz:      unico("as_vezes"), ExigeVoz: 0,
		Toxicidade: entre(2, 3), Compromisso: unico("eventual"),
		Pesos:        pesos(3, 3, 3, 3, 3, 3),
		Jogos:        []ModeloJogo{casual(arkBreeder), casual(rustCasual), casual(lolAram)},
		SortearJogos: 2,
	},
}
