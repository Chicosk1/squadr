package domain

import (
	"strings"
	"testing"
)

func perfilBase(pj PerfilJogo) Perfil {
	return Perfil{
		Nome: "Base", FaixaEtaria: "18-24", FaixasAceitas: []string{"18-24"},
		Regiao: "AS", Idiomas: []string{"pt"},
		Jogos: []PerfilJogo{pj},
	}
}

func ark(versao, modo, plataforma string) PerfilJogo {
	return PerfilJogo{Jogo: ARK, Versao: versao, Modo: modo, Plataforma: plataforma}
}

func rust(versao string) PerfilJogo {
	return PerfilJogo{Jogo: Rust, Versao: versao}
}

func lol(elo int, filas ...string) PerfilJogo {
	return PerfilJogo{Jogo: LoL, Elo: elo, Filas: filas}
}

func TestPassaFiltros(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     PerfilJogo
		esperado bool
	}{
		{"ARK: ASE com ASA", ark("ase", "pve", "pc_steam"), ark("asa", "pve", "pc_steam"), false},
		{"ARK: PvE com PvP", ark("asa", "pve", "pc_steam"), ark("asa", "pvp", "pc_steam"), false},
		{"ARK: PvE com Ambos", ark("asa", "pve", "pc_steam"), ark("asa", "ambos", "pc_steam"), true},
		{"ARK: Steam com Epic", ark("asa", "pve", "pc_steam"), ark("asa", "pve", "pc_epic"), true},
		{"ARK: PC com Xbox", ark("asa", "pve", "pc_steam"), ark("asa", "pve", "xbox"), false},
		{"Rust: PC com console", rust("pc"), rust("console"), false},
		{"Rust: PC com PC", rust("pc"), rust("pc"), true},
		{"LoL: sem fila em comum", lol(3, "solo_duo"), lol(3, "aram"), false},
		{"LoL: duo com diferença de elo 3", lol(1, "solo_duo"), lol(4, "solo_duo"), false},
		{"LoL: duo com diferença de elo 2", lol(2, "solo_duo"), lol(4, "solo_duo"), true},
		{"LoL: duo com sem ranked", lol(SemRanked, "solo_duo"), lol(3, "solo_duo"), false},
		{"LoL: flex ignora diferença de elo", lol(0, "flex"), lol(9, "flex"), true},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got := PassaFiltros(perfilBase(c.a), perfilBase(c.b), c.a.Jogo)
			if got != c.esperado {
				t.Errorf("PassaFiltros = %v, esperado %v", got, c.esperado)
			}
		})
	}
}

func TestPassaFiltrosGerais(t *testing.T) {
	a := perfilBase(rust("pc"))

	b := perfilBase(rust("pc"))
	b.Regiao = "EU"
	if PassaFiltros(a, b, Rust) {
		t.Error("regiões diferentes não deveriam passar")
	}

	b = perfilBase(rust("pc"))
	b.Idiomas = []string{"es"}
	if PassaFiltros(a, b, Rust) {
		t.Error("sem idioma em comum não deveria passar")
	}

	b = perfilBase(rust("pc"))
	b.FaixaEtaria = "35+"
	if PassaFiltros(a, b, Rust) {
		t.Error("faixa etária não aceita por A não deveria passar")
	}

	if PassaFiltros(a, perfilBase(lol(3, "flex")), Rust) {
		t.Error("quem não joga o jogo não deveria passar")
	}
}

// TestFiltrosMotivos confere a explicação que a tela de comparação mostra: todos os filtros
// aparecem, cada um com o motivo quando barra o par.
func TestFiltrosMotivos(t *testing.T) {
	a := perfilBase(ark("asa", "pve", "pc_steam"))
	a.Nome = "Ana"
	b := perfilBase(ark("asa", "pvp", "xbox"))
	b.Nome, b.FaixaEtaria, b.Idiomas = "Bia", "35+", []string{"es"}

	filtros := Filtros(a, b, ARK)
	nomes := make([]string, len(filtros))
	barrados := map[string]string{}
	for i, f := range filtros {
		nomes[i] = f.Nome
		if f.Passou && f.Motivo != "" {
			t.Errorf("%s passou, mas tem motivo %q", f.Nome, f.Motivo)
		}
		if !f.Passou {
			barrados[f.Nome] = f.Motivo
		}
	}
	esperados := "Os dois jogam ARK,Mesma região,Faixa etária aceita pelos dois,Idioma em comum," +
		"Mesma versão,Plataforma compatível,Modo compatível"
	if strings.Join(nomes, ",") != esperados {
		t.Errorf("filtros na ordem errada: %v", nomes)
	}
	casos := map[string]string{
		"Faixa etária aceita pelos dois": "Ana não aceita jogar com 35 ou mais",
		"Idioma em comum":                "Português × Espanhol",
		"Plataforma compatível":          "PC (Steam) × Xbox",
		"Modo compatível":                "PvE × PvP",
	}
	if len(barrados) != len(casos) {
		t.Errorf("barrados %v, esperados %v", barrados, casos)
	}
	for nome, motivo := range casos {
		if barrados[nome] != motivo {
			t.Errorf("%s: motivo %q, esperado %q", nome, barrados[nome], motivo)
		}
	}
}

func TestFiltrosQuemNaoJoga(t *testing.T) {
	a := perfilBase(rust("pc"))
	b := perfilBase(lol(3, "flex"))
	b.Nome = "Bia"
	filtros := Filtros(a, b, Rust)
	if filtros[0].Passou || filtros[0].Motivo != "Bia não joga Rust" {
		t.Errorf("primeiro filtro deveria barrar quem não joga: %+v", filtros[0])
	}
	for _, f := range filtros {
		if f.Nome == "Mesma versão" {
			t.Error("os filtros do jogo não se aplicam a quem não joga")
		}
	}
}

func TestFiltrosMotivoDuo(t *testing.T) {
	casos := []struct {
		a, b   PerfilJogo
		trecho string
	}{
		{lol(1, "solo_duo", "aram"), lol(5, "solo_duo", "flex"), "a distância de elo é 4 (Bronze × Esmeralda)"},
		{lol(SemRanked, "solo_duo"), lol(3, "solo_duo"), "exige elo nos dois"},
		{lol(3, "aram"), lol(3, "flex"), "ARAM × Flex"},
	}
	for _, c := range casos {
		f := Filtros(perfilBase(c.a), perfilBase(c.b), LoL)
		ultimo := f[len(f)-1]
		if ultimo.Nome != "Fila em comum" || ultimo.Passou || !strings.Contains(ultimo.Motivo, c.trecho) {
			t.Errorf("motivo %q não contém %q", ultimo.Motivo, c.trecho)
		}
	}
}
