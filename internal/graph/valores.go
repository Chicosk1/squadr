package graph

func texto(m map[string]any, chave string) string {
	v, _ := m[chave].(string)
	return v
}

func inteiro(m map[string]any, chave string) int {
	v, _ := m[chave].(int64)
	return int(v)
}

func decimal(m map[string]any, chave string) float64 {
	v, _ := m[chave].(float64)
	return v
}

func booleano(m map[string]any, chave string) bool {
	v, _ := m[chave].(bool)
	return v
}

func lista(m map[string]any, chave string) []string {
	bruto, ok := m[chave].([]any)
	if !ok {
		return nil
	}
	itens := make([]string, 0, len(bruto))
	for _, x := range bruto {
		if s, ok := x.(string); ok {
			itens = append(itens, s)
		}
	}
	return itens
}
