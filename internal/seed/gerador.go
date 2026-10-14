package seed

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/Chicosk1/squadr/internal/domain"
)

const (
	ruidoHorario     = 0.10
	ruidoIntensidade = 0.20
	ruidoFuncao      = 0.15
	ruidoPeso        = 0.30
	probOutlier      = 0.05
)

var (
	regioes       = Opcoes{{"AS", 80}, {"AN", 12}, {"EU", 8}}
	faixasEtarias = Opcoes{{"18-24", 50}, {"25-34", 35}, {"35+", 15}}

	vizinhas = map[string][]string{
		"18-24": {"18-24", "25-34"},
		"25-34": {"18-24", "25-34", "35+"},
		"35+":   {"25-34", "35+"},
	}

	niveisDesejados = igual("parecido", "mais_experiente", "tanto_faz")
	compromissos    = igual("fixo", "eventual", "tanto_faz")

	nivelPorHoras = map[string]float64{"<100": 0.125, "100-500": 0.375, "500-2000": 0.625, "2000+": 0.875}

	adjetivos = []string{"Sombrio", "Veloz", "Calmo", "Feroz", "Astuto", "Noturno", "Mistico", "Bravo",
		"Gelado", "Ardente", "Furtivo", "Lendario", "Tranquilo", "Selvagem", "Letal", "Sereno"}
	substantivos = []string{"Lobo", "Dragao", "Falcao", "Urso", "Raptor", "Corvo", "Tigre", "Golem",
		"Fenix", "Viking", "Ninja", "Mago", "Cacador", "Guerreiro", "Arqueiro", "Piloto"}
)

type Gerador struct {
	r *rand.Rand
}

func NovoGerador(semente uint64) *Gerador {
	return &Gerador{r: rand.New(rand.NewPCG(semente, semente^0x9E3779B97F4A7C15))}
}

func (g *Gerador) GerarUsuarios(total int) []domain.Perfil {
	perfis := make([]domain.Perfil, 0, total)
	for i, n := range distribuir(total) {
		for k := 0; k < n; k++ {
			perfis = append(perfis, g.gerarPerfil(Personas[i]))
		}
	}

	g.r.Shuffle(len(perfis), func(i, j int) { perfis[i], perfis[j] = perfis[j], perfis[i] })
	for i := range perfis {
		perfis[i].ID = fmt.Sprintf("u%06d", i+1)
	}
	return perfis
}

func distribuir(total int) []int {
	contagens := make([]int, len(Personas))
	usado := 0
	for i, p := range Personas {
		if i == len(Personas)-1 {
			contagens[i] = max(0, total-usado)
			break
		}
		contagens[i] = int(math.Round(float64(total) * p.Proporcao / 100))
		usado += contagens[i]
	}
	return contagens
}

func (g *Gerador) gerarPerfil(p Persona) domain.Perfil {
	regiao := regioes.Sortear(g.r)
	faixa := faixasEtarias.Sortear(g.r)

	perfil := domain.Perfil{
		Nome:                 g.apelido(),
		FaixaEtaria:          faixa,
		FaixasAceitas:        g.faixasAceitas(faixa),
		Regiao:               regiao,
		Idiomas:              g.idiomas(regiao),
		UsaVoz:               p.UsaVoz.Sortear(g.r),
		ToleranciaToxicidade: p.Toxicidade.Sortear(g.r),
		Compromisso:          p.Compromisso.Sortear(g.r),
		Pesos:                g.ruidoPesos(p.Pesos),
		Disponibilidade:      g.horarios(p),
		Persona:              p.Nome,
	}
	perfil.ExigeVoz = perfil.UsaVoz != "nao" && g.r.Float64() < p.ExigeVoz

	modelos := p.Jogos
	if p.SortearJogos > 0 {
		indices := g.r.Perm(len(p.Jogos))[:p.SortearJogos]
		slices.Sort(indices)
		modelos = nil
		for _, i := range indices {
			modelos = append(modelos, p.Jogos[i])
		}
	}
	for _, m := range modelos {
		perfil.Jogos = append(perfil.Jogos, g.gerarPerfilJogo(m))
	}

	if g.r.Float64() < probOutlier {
		g.aplicarOutlier(&perfil)
	}
	return perfil
}

func (g *Gerador) gerarPerfilJogo(m ModeloJogo) domain.PerfilJogo {
	pj := domain.PerfilJogo{
		Jogo:            m.Jogo,
		Intensidade:     g.ruidoIntensidade(m.Intensidade.Sortear(g.r)),
		NivelDesejado:   m.NivelDesejado.Sortear(g.r),
		FuncaoPrincipal: m.FuncaoPrincipal,
		Elo:             domain.SemRanked,
		PersonaJogo:     m.Persona,
	}
	pj.FuncaoSecundaria = g.funcaoSecundaria(m)

	switch m.Jogo {
	case domain.ARK:
		pj.Horas = m.Horas.Sortear(g.r)
		pj.Versao = m.Versao.Sortear(g.r)
		pj.Plataforma = m.Plataforma.Sortear(g.r)
		pj.Modo = m.Modo.Sortear(g.r)
		pj.Servidor = m.Servidor.Sortear(g.r)
		pj.TamanhoTribo = m.TamanhoTribo.Sortear(g.r)
	case domain.Rust:
		pj.Horas = m.Horas.Sortear(g.r)
		pj.Versao = m.Versao.Sortear(g.r)
		pj.Servidor = m.Servidor.Sortear(g.r)
		pj.LimiteGrupo = m.LimiteGrupo.Sortear(g.r)
		pj.Wipe = m.Wipe.Sortear(g.r)
	case domain.LoL:
		pj.Elo = m.Elo.Sortear(g.r) // -1 = sem ranked
		pj.Filas = slices.Clone(m.Filas)
	}
	pj.Nivel = calcularNivel(pj)
	return pj
}

func calcularNivel(pj domain.PerfilJogo) float64 {
	if pj.Jogo == domain.LoL {
		if pj.Elo == domain.SemRanked {
			return 0
		}
		return float64(pj.Elo+1) / 10
	}
	return nivelPorHoras[pj.Horas]
}

// ---------- Atributos gerais ----------

func (g *Gerador) idiomas(regiao string) []string {
	var l []string
	switch regiao {
	case "AS":
		if g.r.Float64() < 0.95 {
			l = append(l, "pt")
		}
		if g.r.Float64() < 0.30 {
			l = append(l, "en")
		}
		if g.r.Float64() < 0.15 {
			l = append(l, "es")
		}
		if len(l) == 0 {
			l = []string{"es"}
		}
	case "AN":
		l = []string{"en"}
		if g.r.Float64() < 0.20 {
			l = append(l, "es")
		}
	default:
		l = []string{"en"}
	}
	return l
}

func (g *Gerador) faixasAceitas(propria string) []string {
	n := g.r.IntN(100)
	switch {
	case n < 70:
		return slices.Clone(vizinhas[propria])
	case n < 90:
		return []string{"18-24", "25-34", "35+"}
	default:
		return []string{propria}
	}
}

func (g *Gerador) apelido() string {
	return fmt.Sprintf("%s%s%02d",
		adjetivos[g.r.IntN(len(adjetivos))],
		substantivos[g.r.IntN(len(substantivos))],
		g.r.IntN(100))
}

// ---------- Horários ----------

func (g *Gerador) horarios(p Persona) uint32 {
	if p.ProbHorario > 0 {
		return g.horariosAleatorios(p.ProbHorario)
	}
	var m uint32
	for _, bit := range p.Horarios {
		m |= 1 << bit
	}
	for bit := 0; bit < domain.NumHorarios; bit++ {
		if g.r.Float64() < ruidoHorario {
			m ^= 1 << bit
		}
	}
	if m == 0 {
		m = 1 << g.r.IntN(domain.NumHorarios)
	}
	return m
}

func (g *Gerador) horariosAleatorios(prob float64) uint32 {
	var m uint32
	for bit := 0; bit < domain.NumHorarios; bit++ {
		if g.r.Float64() < prob {
			m |= 1 << bit
		}
	}
	if m == 0 {
		m = 1 << g.r.IntN(domain.NumHorarios)
	}
	return m
}

// ---------- Ruído ----------

func (g *Gerador) sinal() int {
	if g.r.IntN(2) == 0 {
		return -1
	}
	return 1
}

func limitar(v, lo, hi int) int { return min(max(v, lo), hi) }

func (g *Gerador) ruidoIntensidade(v int) int {
	if g.r.Float64() < ruidoIntensidade {
		v += g.sinal()
	}
	return limitar(v, 1, 5)
}

func (g *Gerador) ruidoPesos(p domain.Pesos) domain.Pesos {
	for _, campo := range []*int{&p.Horario, &p.Intensidade, &p.Modo, &p.Comunicacao, &p.Funcao, &p.Nivel} {
		if g.r.Float64() < ruidoPeso {
			*campo = limitar(*campo+g.sinal(), 1, 5)
		}
	}
	return p
}

func (g *Gerador) funcaoSecundaria(m ModeloJogo) string {
	if m.FuncaoSecundaria != "" && g.r.Float64() >= ruidoFuncao {
		return m.FuncaoSecundaria
	}
	return g.funcaoAleatoria(m.Jogo, m.FuncaoPrincipal)
}

func (g *Gerador) funcaoAleatoria(jogo domain.Jogo, exceto string) string {
	var opcoes []string
	for _, f := range domain.Funcoes[jogo] {
		if f != exceto {
			opcoes = append(opcoes, f)
		}
	}
	return opcoes[g.r.IntN(len(opcoes))]
}

func (g *Gerador) aplicarOutlier(p *domain.Perfil) {
	switch g.r.IntN(6) {
	case 0:
		p.Jogos[g.r.IntN(len(p.Jogos))].Intensidade = 1 + g.r.IntN(5)
	case 1:
		p.Disponibilidade = g.horariosAleatorios(0.30)
	case 2:
		p.ToleranciaToxicidade = 1 + g.r.IntN(5)
	case 3:
		p.Jogos[g.r.IntN(len(p.Jogos))].NivelDesejado = niveisDesejados.Sortear(g.r)
	case 4:
		p.Compromisso = compromissos.Sortear(g.r)
	case 5:
		p.Pesos = domain.Pesos{
			Horario: 1 + g.r.IntN(5), Intensidade: 1 + g.r.IntN(5), Modo: 1 + g.r.IntN(5),
			Comunicacao: 1 + g.r.IntN(5), Funcao: 1 + g.r.IntN(5), Nivel: 1 + g.r.IntN(5),
		}
	}
}
