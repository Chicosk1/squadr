package matching

import "fmt"

// Matriz guarda a complementaridade entre funções em memória (D-006), carregada
// uma vez na inicialização. É simétrica: Valor(a, b) == Valor(b, a).
type Matriz struct {
	valores map[[2]string]float64
}

// NovaMatriz recebe os pares como estão no grafo (uma direção só) e os espelha.
// Falha se o mesmo par aparecer com dois valores diferentes.
func NovaMatriz(pares map[[2]string]float64) (Matriz, error) {
	m := Matriz{valores: make(map[[2]string]float64, 2*len(pares))}
	for par, v := range pares {
		if v < 0 || v > 1 {
			return Matriz{}, fmt.Errorf("complementaridade %v = %v fora de [0, 1]", par, v)
		}
		inverso := [2]string{par[1], par[0]}
		if outro, ok := m.valores[inverso]; ok && outro != v {
			return Matriz{}, fmt.Errorf("par %v com valores diferentes: %v e %v", par, v, outro)
		}
		m.valores[par] = v
		m.valores[inverso] = v
	}
	return m, nil
}

// Valor devolve 0 para pares desconhecidos ou funções vazias.
func (m Matriz) Valor(a, b string) float64 {
	return m.valores[[2]string{a, b}]
}

func (m Matriz) Tamanho() int { return len(m.valores) }
