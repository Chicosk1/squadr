package seed

import (
	"math"
	"reflect"
	"slices"
	"testing"

	"github.com/Chicosk1/squadr/internal/domain"
)

func TestGrade(t *testing.T) {
	if got := grade("sex:noite"); !slices.Equal(got, []int{19}) {
		t.Errorf("grade(sex:noite) = %v, esperado [19]", got)
	}
	if got := len(grade("seg-sex:noite")); got != 5 {
		t.Errorf("seg-sex:noite deveria ter 5 faixas, tem %d", got)
	}
	if got := len(grade("seg-dom:noite,madrugada")); got != 14 {
		t.Errorf("seg-dom:noite,madrugada deveria ter 14 faixas, tem %d", got)
	}
}

func TestGeracaoDeterministica(t *testing.T) {
	a := NovoGerador(7).GerarUsuarios(500)
	b := NovoGerador(7).GerarUsuarios(500)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("a mesma semente gerou dados diferentes")
	}
	c := NovoGerador(8).GerarUsuarios(500)
	if reflect.DeepEqual(a, c) {
		t.Fatal("sementes diferentes geraram dados iguais")
	}
}

func TestProporcoes(t *testing.T) {
	const total = 2000
	perfis := NovoGerador(2).GerarUsuarios(total)
	if len(perfis) != total {
		t.Fatalf("gerou %d perfis, esperado %d", len(perfis), total)
	}
	contagem := map[string]int{}
	for _, p := range perfis {
		contagem[p.Persona]++
	}
	for _, persona := range Personas {
		esperado := int(math.Round(total * persona.Proporcao / 100))
		if d := contagem[persona.Nome] - esperado; d < -1 || d > 1 {
			t.Errorf("%s: %d usuários, esperado cerca de %d", persona.Nome, contagem[persona.Nome], esperado)
		}
	}
}

func TestInvariantes(t *testing.T) {
	perfis := NovoGerador(3).GerarUsuarios(3000)
	ids := map[string]bool{}

	for _, p := range perfis {
		if ids[p.ID] {
			t.Fatalf("ID duplicado: %s", p.ID)
		}
		ids[p.ID] = true

		if len(p.Idiomas) == 0 {
			t.Errorf("%s: sem idiomas", p.ID)
		}
		if !slices.Contains(p.FaixasAceitas, p.FaixaEtaria) {
			t.Errorf("%s: a própria faixa etária não está entre as aceitas", p.ID)
		}
		if p.Disponibilidade == 0 {
			t.Errorf("%s: sem nenhum horário", p.ID)
		}
		if p.ExigeVoz && p.UsaVoz == "nao" {
			t.Errorf("%s: exige voz mas não usa", p.ID)
		}
		if len(p.Jogos) == 0 {
			t.Fatalf("%s: não joga nada", p.ID)
		}
		if p.Persona == PersonaMulti && len(p.Jogos) != 2 {
			t.Errorf("%s: multi-jogo com %d jogos, esperado 2", p.ID, len(p.Jogos))
		}
		for _, peso := range []int{p.Pesos.Horario, p.Pesos.Intensidade, p.Pesos.Modo,
			p.Pesos.Comunicacao, p.Pesos.Funcao, p.Pesos.Nivel} {
			if peso < 1 || peso > 5 {
				t.Errorf("%s: peso %d fora de 1..5", p.ID, peso)
			}
		}
		for _, pj := range p.Jogos {
			if pj.Intensidade < 1 || pj.Intensidade > 5 {
				t.Errorf("%s: intensidade %d fora de 1..5", p.ID, pj.Intensidade)
			}
			if pj.FuncaoPrincipal == pj.FuncaoSecundaria {
				t.Errorf("%s: funções principal e secundária iguais", p.ID)
			}
			if pj.Nivel < 0 || pj.Nivel > 1 {
				t.Errorf("%s: nível %v fora de 0..1", p.ID, pj.Nivel)
			}
		}
	}
}

// ---------- Filtros ----------

func perfilBase(pj domain.PerfilJogo) domain.Perfil {
	return domain.Perfil{
		FaixaEtaria: "18-24", FaixasAceitas: []string{"18-24"},
		Regiao: "AS", Idiomas: []string{"pt"},
		Jogos: []domain.PerfilJogo{pj},
	}
}

func ark(versao, modo, plataforma string) domain.PerfilJogo {
	return domain.PerfilJogo{Jogo: domain.ARK, Versao: versao, Modo: modo, Plataforma: plataforma}
}

func rust(versao string) domain.PerfilJogo {
	return domain.PerfilJogo{Jogo: domain.Rust, Versao: versao}
}

func lol(elo int, filas ...string) domain.PerfilJogo {
	return domain.PerfilJogo{Jogo: domain.LoL, Elo: elo, Filas: filas}
}

func TestPassaFiltros(t *testing.T) {
	casos := []struct {
		nome     string
		a, b     domain.PerfilJogo
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
		{"LoL: duo com sem ranked", lol(domain.SemRanked, "solo_duo"), lol(3, "solo_duo"), false},
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
	if PassaFiltros(a, b, domain.Rust) {
		t.Error("regiões diferentes não deveriam passar")
	}

	b = perfilBase(rust("pc"))
	b.Idiomas = []string{"es"}
	if PassaFiltros(a, b, domain.Rust) {
		t.Error("sem idioma em comum não deveria passar")
	}

	b = perfilBase(rust("pc"))
	b.FaixaEtaria = "35+"
	if PassaFiltros(a, b, domain.Rust) {
		t.Error("faixa etária não aceita por A não deveria passar")
	}

	if PassaFiltros(a, perfilBase(lol(3, "flex")), domain.Rust) {
		t.Error("quem não joga o jogo não deveria passar")
	}
}

// ---------- Gabarito ----------

func primeiro(t *testing.T, perfis []domain.Perfil, persona string) domain.Perfil {
	t.Helper()
	for _, p := range perfis {
		if p.Persona == persona {
			return p
		}
	}
	t.Fatalf("nenhum perfil da persona %s", persona)
	return domain.Perfil{}
}

func TestGrau(t *testing.T) {
	perfis := NovoGerador(1).GerarUsuarios(1000)
	casos := []struct {
		a, b     string
		jogo     domain.Jogo
		esperado int
	}{
		{"LOL-SUPORTE", "LOL-ADC", domain.LoL, 2},
		{"LOL-ADC", "LOL-ADC", domain.LoL, 1},
		{"LOL-ARAM", "LOL-ARAM", domain.LoL, 1},
		{"ARK-PVE-BREEDER", "ARK-PVE-BUILDER", domain.ARK, 2},
		{"ARK-PVE-BREEDER", "ARK-PVE-BREEDER", domain.ARK, 2},
		{"RUST-APRENDIZ", "RUST-CASUAL", domain.Rust, 1},
		{"RUST-APRENDIZ", "RUST-GRINDER", domain.Rust, 0},
	}
	for _, c := range casos {
		a, b := primeiro(t, perfis, c.a), primeiro(t, perfis, c.b)
		if got := Grau(a, b, c.jogo); got != c.esperado {
			t.Errorf("Grau(%s, %s) = %d, esperado %d", c.a, c.b, got, c.esperado)
		}
		if got := Grau(b, a, c.jogo); got != c.esperado {
			t.Errorf("Grau(%s, %s) = %d, esperado %d (simetria)", c.b, c.a, got, c.esperado)
		}
	}
}

func TestGrauMulti(t *testing.T) {
	perfis := NovoGerador(1).GerarUsuarios(1000)
	var multi domain.Perfil
	for _, p := range perfis {
		if _, joga := p.PerfilDoJogo(domain.ARK); p.Persona == PersonaMulti && joga {
			multi = p
			break
		}
	}
	if multi.ID == "" {
		t.Fatal("nenhum multi-jogo joga ARK")
	}
	if got := Grau(multi, primeiro(t, perfis, "ARK-PVE-BREEDER"), domain.ARK); got != 1 {
		t.Errorf("multi com breeder no ARK = %d, esperado 1", got)
	}
	if got := Grau(multi, primeiro(t, perfis, "ARK-PVE-BUILDER"), domain.ARK); got != 0 {
		t.Errorf("multi com builder no ARK = %d, esperado 0", got)
	}
}

// ---------- Interações ----------

func TestInteracoes(t *testing.T) {
	g := NovoGerador(5)
	perfis := g.GerarUsuarios(2000)
	interacoes := g.GerarInteracoes(perfis)
	if len(interacoes) == 0 {
		t.Fatal("nenhuma interação gerada")
	}

	porID := make(map[string]domain.Perfil, len(perfis))
	for _, p := range perfis {
		porID[p.ID] = p
	}

	vistos := map[[3]string]bool{}
	var curtidas, total [3]int
	for _, x := range interacoes {
		if x.De == x.Para {
			t.Fatalf("interação de %s com ele mesmo", x.De)
		}
		k := [3]string{x.De, x.Para, string(x.Jogo)}
		if vistos[k] {
			t.Fatalf("interação duplicada: %v", k)
		}
		vistos[k] = true

		a, b := porID[x.De], porID[x.Para]
		if !PassaFiltros(a, b, x.Jogo) {
			t.Fatalf("interação entre %s e %s não passa nos filtros", x.De, x.Para)
		}
		grau := Grau(a, b, x.Jogo)
		total[grau]++
		if x.Curtiu {
			curtidas[grau]++
		}
	}

	var taxa [3]float64
	for grau := 0; grau < 3; grau++ {
		if total[grau] == 0 {
			t.Fatalf("nenhuma interação de grau %d", grau)
		}
		taxa[grau] = float64(curtidas[grau]) / float64(total[grau])
	}
	if !(taxa[2] > taxa[1] && taxa[1] > taxa[0]) {
		t.Errorf("taxas de like deveriam crescer com o grau, obtido %v", taxa)
	}
}
