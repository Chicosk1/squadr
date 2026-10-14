package domain

type Jogo string

const (
	ARK  Jogo = "ark"
	Rust Jogo = "rust"
	LoL  Jogo = "lol"
)

var Jogos = []Jogo{ARK, Rust, LoL}

var Funcoes = map[Jogo][]string{
	ARK:  {"ark:breeder", "ark:builder", "ark:coletor", "ark:pvp", "ark:explorador"},
	Rust: {"rust:builder", "rust:coletor", "rust:pvp", "rust:raider", "rust:eletricista"},
	LoL:  {"lol:top", "lol:jungle", "lol:mid", "lol:adc", "lol:suporte", "lol:preencher"},
}
