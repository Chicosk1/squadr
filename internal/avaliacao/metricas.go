// Package avaliacao mede a qualidade do ranking contra o gabarito das personas (Fase 7, D-046).
//
// As métricas recebem a lista de graus de relevância na ordem em que o sistema recomendou
// (grau 0, 1 ou 2, vindo de seed.Grau). Assim elas não dependem do motor nem do banco.
package avaliacao

import (
	"cmp"
	"math"
	"slices"
)

// PrecisaoK é a fração das k primeiras posições ocupadas por um relevante (grau ≥ 1).
// Divide sempre por k: uma lista com menos de k itens conta as posições vazias como erro.
func PrecisaoK(graus []int, k int) float64 {
	acertos := 0
	for _, g := range graus[:min(k, len(graus))] {
		if g >= 1 {
			acertos++
		}
	}
	return float64(acertos) / float64(k)
}

// desconto é o peso da posição i (contada a partir de 0): 1/log2(i+2).
// A 1ª posição vale 1; a 2ª, 0,63; a 3ª, 0,50; a 10ª, 0,29.
func desconto(i int) float64 { return 1 / math.Log2(float64(i)+2) }

// DCGK é o ganho acumulado com desconto: DCG@k = Σ grau_i / log2(i+1), com i de 1 a k
// (Järvelin; Kekäläinen, 2002). O ganho é o próprio grau, então um par ideal (grau 2)
// vale o dobro de um par aceitável (grau 1), e acertar no topo vale mais que no fim.
func DCGK(graus []int, k int) float64 {
	dcg := 0.0
	for i, g := range graus[:min(k, len(graus))] {
		dcg += float64(g) * desconto(i)
	}
	return dcg
}

// IDCGK é o DCG da melhor lista possível: os graus de todos os candidatos que o sistema
// poderia recomendar, em ordem decrescente.
func IDCGK(possiveis []int, k int) float64 {
	ideal := slices.Clone(possiveis)
	slices.SortFunc(ideal, func(a, b int) int { return cmp.Compare(b, a) })
	return DCGK(ideal, k)
}

// NDCGK normaliza o DCG pelo ideal e fica em [0, 1]. Devolve ok = false quando não há
// nenhum relevante possível (IDCG = 0): a consulta não tem o que acertar e fica fora da média.
func NDCGK(graus, possiveis []int, k int) (ndcg float64, ok bool) {
	idcg := IDCGK(possiveis, k)
	if idcg == 0 {
		return 0, false
	}
	return DCGK(graus, k) / idcg, true
}

// AleatorioEsperado devolve o Precision@k e o NDCG@k esperados de uma lista de k itens
// sorteada sem reposição entre os possíveis, ou seja, a média de todos os sorteios.
// Em um sorteio uniforme cada posição recebe, em média, o grau médio do conjunto;
// por isso a esperança sai exata, sem depender de semente nem de repetições.
func AleatorioEsperado(possiveis []int, k int) (precisao, ndcg float64, ok bool) {
	n := len(possiveis)
	if n == 0 {
		return 0, 0, false
	}
	soma, relevantes := 0, 0
	for _, g := range possiveis {
		soma += g
		if g >= 1 {
			relevantes++
		}
	}
	m := min(k, n)
	precisao = float64(relevantes) / float64(n) * float64(m) / float64(k)

	idcg := IDCGK(possiveis, k)
	if idcg == 0 {
		return precisao, 0, false
	}
	descontos := 0.0
	for i := range m {
		descontos += desconto(i)
	}
	return precisao, float64(soma) / float64(n) * descontos / idcg, true
}
