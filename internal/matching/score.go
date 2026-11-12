package matching

import (
	"math/bits"
	"slices"

	"github.com/Chicosk1/squadr/internal/domain"
)

// Componente indexa os seis componentes do score unilateral (questionario.md, P01 a P06).
type Componente int

const (
	Horario Componente = iota
	Intensidade
	Modo
	Comunicacao
	Funcao
	Nivel
	NumComponentes
)

var nomesComponentes = [NumComponentes]string{
	Horario:     "horários em comum",
	Intensidade: "intensidade e compromisso",
	Modo:        "modo e servidor",
	Comunicacao: "comunicação",
	Funcao:      "funções complementares",
	Nivel:       "nível",
}

func (c Componente) String() string { return nomesComponentes[c] }

// Componentes guarda o valor de cada componente, todos em [0, 1].
type Componentes [NumComponentes]float64

const (
	// PesoPadrao é o peso de um critério que o usuário não respondeu (D-008).
	PesoPadrao = 3
	// DescontoSecundaria multiplica o valor do arranjo para cada jogador fora da função principal (D-035).
	DescontoSecundaria = 0.75
)

// Unilateral calcula s(A→B): quanto B serve para A, pelos critérios e pesos de A (D-035).
// a e b precisam jogar o jogo. O motor nunca lê Persona nem PersonaJogo.
func Unilateral(a, b domain.Perfil, pa, pb domain.PerfilJogo, m Matriz) (float64, Componentes) {
	var c Componentes
	c[Horario] = jaccardBits(a.Disponibilidade, b.Disponibilidade)
	c[Intensidade] = (similaridadeOrdinal(pa.Intensidade, pb.Intensidade, 1, 5) + compromisso(a.Compromisso, b.Compromisso)) / 2
	c[Modo] = modo(pa, pb)
	c[Comunicacao] = fatorVoz(a.ExigeVoz, b.UsaVoz) * tolerancia(a.ToleranciaToxicidade, b.ToleranciaToxicidade)
	c[Funcao] = funcao(pa, pb, m)
	c[Nivel] = nivel(pa.NivelDesejado, pb.Nivel-pa.Nivel)

	pesos := [NumComponentes]int{a.Pesos.Horario, a.Pesos.Intensidade, a.Pesos.Modo,
		a.Pesos.Comunicacao, a.Pesos.Funcao, a.Pesos.Nivel}
	soma, total := 0.0, 0.0
	for i, p := range pesos {
		if p <= 0 {
			p = PesoPadrao
		}
		soma += float64(p) * c[i]
		total += float64(p)
	}
	return soma / total, c
}

// Reciproco é a média harmônica de s(A→B) e s(B→A) (D-015). Se um dos lados é zero, o par é zero.
func Reciproco(ab, ba float64) float64 {
	if ab+ba == 0 {
		return 0
	}
	return 2 * ab * ba / (ab + ba)
}

// ---------- Horário ----------

// jaccardBits é |A∩B| / |A∪B| sobre a grade de 28 bits. Duas grades vazias valem 0.
func jaccardBits(a, b uint32) float64 {
	uniao := bits.OnesCount32(a | b)
	if uniao == 0 {
		return 0
	}
	return float64(bits.OnesCount32(a&b)) / float64(uniao)
}

// ---------- Intensidade e compromisso ----------

// similaridadeOrdinal é 1 − |a − b| / (max − min).
func similaridadeOrdinal(a, b, minimo, maximo int) float64 {
	d := a - b
	if d < 0 {
		d = -d
	}
	return max(0, 1-float64(d)/float64(maximo-minimo))
}

// compromisso (G10): iguais ou qualquer lado "tanto_faz" vale 1; fixo × eventual vale 0.
func compromisso(a, b string) float64 {
	if a == b || a == "tanto_faz" || b == "tanto_faz" || a == "" || b == "" {
		return 1
	}
	return 0
}

// ---------- Modo e servidor ----------

// As escalas ordinais seguem a ordem do vocabulário do formulário (D-038).
var (
	ordemTribo  = domain.Valores(domain.TamanhosTribo)
	ordemLimite = domain.Valores(domain.LimitesGrupo)
	ordemWipe   = domain.Valores(domain.Wipes)
)

func modo(pa, pb domain.PerfilJogo) float64 {
	switch pa.Jogo {
	case domain.ARK:
		return (servidorARK(pa.Servidor, pb.Servidor) +
			escala(pa.TamanhoTribo, pb.TamanhoTribo, ordemTribo)) / 2
	case domain.Rust:
		return (servidorRust(pa.Servidor, pb.Servidor) +
			escala(pa.LimiteGrupo, pb.LimiteGrupo, ordemLimite) +
			escala(pa.Wipe, pb.Wipe, ordemWipe)) / 3
	case domain.LoL:
		return filas(pa, pb)
	}
	return 0
}

// escala compara respostas ordinais: 1 − distância / (n − 1). "tanto_faz" vale 1;
// um valor fora da escala vale 0,5 (neutro), para não premiar nem punir dado inesperado.
func escala(a, b string, ordem []string) float64 {
	if a == "tanto_faz" || b == "tanto_faz" {
		return 1
	}
	i, j := slices.Index(ordem, a), slices.Index(ordem, b)
	if i < 0 || j < 0 {
		return 0.5
	}
	return similaridadeOrdinal(i, j, 0, len(ordem)-1)
}

// servidorARK (A04): iguais 1; dois não oficiais com rates diferentes 0,5; oficial × não oficial 0.
func servidorARK(a, b string) float64 {
	switch {
	case a == b || a == "tanto_faz" || b == "tanto_faz":
		return 1
	case a != "oficial" && b != "oficial":
		return 0.5
	}
	return 0
}

// servidorRust (R02): iguais 1; oficial × comunitário (os dois vanilla) 0,5; o resto 0.
func servidorRust(a, b string) float64 {
	switch {
	case a == b:
		return 1
	case (a == "oficial" && b == "comunitario") || (a == "comunitario" && b == "oficial"):
		return 0.5
	}
	return 0
}

// filas (L04) é o Jaccard das filas, contando solo_duo em comum só se os dois podem jogá-la juntos.
func filas(pa, pb domain.PerfilJogo) float64 {
	uniao := map[string]bool{}
	comum := 0
	for _, f := range pa.Filas {
		uniao[f] = true
		if slices.Contains(pb.Filas, f) && (f != "solo_duo" || duoPermitido(pa.Elo, pb.Elo)) {
			comum++
		}
	}
	for _, f := range pb.Filas {
		uniao[f] = true
	}
	if len(uniao) == 0 {
		return 0
	}
	return float64(comum) / float64(len(uniao))
}

func duoPermitido(a, b int) bool {
	if a == domain.SemRanked || b == domain.SemRanked {
		return false
	}
	return max(a-b, b-a) <= domain.MaxDiffEloDuo
}

// ---------- Comunicação ----------

// fatorVoz (G07, G08) é assimétrico: só importa se A exige voz.
func fatorVoz(aExige bool, bUsa string) float64 {
	if !aExige {
		return 1
	}
	switch bUsa {
	case "sim":
		return 1
	case "as_vezes":
		return 0.5
	}
	return 0
}

// tolerancia (G09) é assimétrica (D-035): A só se incomoda se B tolera mais trash talk que A.
// O desconforto de B com A é medido em s(B→A), e a reciprocidade junta os dois lados.
func tolerancia(a, b int) float64 {
	return 1 - float64(max(0, b-a))/4
}

// ---------- Função ----------

// funcao é o melhor arranjo entre as funções dos dois (D-035). Na partida cada um assume
// uma função, então vale o melhor par possível, e não a média. Cada jogador que precisa
// sair da principal multiplica o valor por DescontoSecundaria.
func funcao(pa, pb domain.PerfilJogo, m Matriz) float64 {
	melhor := 0.0
	for i, fa := range [2]string{pa.FuncaoPrincipal, pa.FuncaoSecundaria} {
		for j, fb := range [2]string{pb.FuncaoPrincipal, pb.FuncaoSecundaria} {
			if fa == "" || fb == "" {
				continue
			}
			v := m.Valor(fa, fb)
			if i == 1 {
				v *= DescontoSecundaria
			}
			if j == 1 {
				v *= DescontoSecundaria
			}
			melhor = max(melhor, v)
		}
	}
	return melhor
}

// ---------- Nível ----------

// nivel é assimétrico (A09, R08, L06): depende do que A quer. d = nível de B − nível de A,
// os dois em [0, 1]. Um degrau de horas vale 0,25; um de elo, 0,1.
func nivel(desejado string, d float64) float64 {
	switch desejado {
	case "parecido":
		if d < 0 {
			d = -d
		}
		return max(0, 1-2*d) // zera com dois degraus de horas
	case "mais_experiente":
		return min(1, max(0, 0.5+2*d)) // igual vale 0,5; um degrau acima já vale 1
	}
	return 1 // tanto_faz
}
