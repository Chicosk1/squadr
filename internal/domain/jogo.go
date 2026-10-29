package domain

type Jogo string

const (
	ARK  Jogo = "ark"
	Rust Jogo = "rust"
	LoL  Jogo = "lol"
)

var Jogos = []Jogo{ARK, Rust, LoL}

var Funcoes = map[Jogo][]string{
	ARK:  {"ark:breeder", "ark:builder", "ark:gatherer", "ark:pvp", "ark:explorer"},
	Rust: {"rust:builder", "rust:gatherer", "rust:pvp", "rust:raider", "rust:electrician"},
	LoL:  {"lol:top", "lol:jungle", "lol:mid", "lol:adc", "lol:support", "lol:fill"},
}

const MaxDiffEloDuo = 2
