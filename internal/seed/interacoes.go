package seed

import (
	"time"

	"github.com/Chicosk1/squadr/internal/domain"
)

var probCurtida = [3]float64{0.15, 0.45, 0.70}

var referencia = time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)

func (g *Gerador) GerarInteracoes(perfis []domain.Perfil) []domain.Interacao {
	porJogo := make(map[domain.Jogo][]int)
	for i, p := range perfis {
		for _, pj := range p.Jogos {
			porJogo[pj.Jogo] = append(porJogo[pj.Jogo], i)
		}
	}

	var interacoes []domain.Interacao
	for i := range perfis {
		eu := perfis[i]
		for _, pj := range eu.Jogos {
			pool := porJogo[pj.Jogo]
			alvo := 20 + g.r.IntN(31)
			tentados := map[int]struct{}{i: {}}
			aceitos := 0

			for tentativa := 0; aceitos < alvo && tentativa < alvo*20; tentativa++ {
				j := pool[g.r.IntN(len(pool))]
				if _, jaTentou := tentados[j]; jaTentou {
					continue
				}
				tentados[j] = struct{}{}

				outro := perfis[j]
				if !PassaFiltros(eu, outro, pj.Jogo) {
					continue
				}
				aceitos++

				grau := Grau(eu, outro, pj.Jogo)
				interacoes = append(interacoes, domain.Interacao{
					De:     eu.ID,
					Para:   outro.ID,
					Jogo:   pj.Jogo,
					Curtiu: g.r.Float64() < probCurtida[grau],
					Em:     referencia.Add(-time.Duration(g.r.IntN(30*24*60)) * time.Minute),
				})
			}
		}
	}
	return interacoes
}
