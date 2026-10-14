package domain

import (
	"slices"
	"strings"
)

const (
	PeriodosPorDia = 4
	NumHorarios    = 7 * PeriodosPorDia
)

var (
	Dias     = []string{"seg", "ter", "qua", "qui", "sex", "sab", "dom"}
	Periodos = []string{"madrugada", "manha", "tarde", "noite"}
)

func IDHorario(bit int) string {
	return Dias[bit/PeriodosPorDia] + "-" + Periodos[bit%PeriodosPorDia]
}

func BitDoHorario(id string) (int, bool) {
	dia, periodo, ok := strings.Cut(id, "-")
	if !ok {
		return 0, false
	}
	d, p := slices.Index(Dias, dia), slices.Index(Periodos, periodo)
	if d < 0 || p < 0 {
		return 0, false
	}
	return d*PeriodosPorDia + p, true
}
