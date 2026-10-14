package domain

import (
	"slices"
	"testing"
)

func TestIDeBitDoHorario(t *testing.T) {
	for bit := 0; bit < NumHorarios; bit++ {
		volta, ok := BitDoHorario(IDHorario(bit))
		if !ok || volta != bit {
			t.Errorf("bit %d -> %q -> %d (ok=%v)", bit, IDHorario(bit), volta, ok)
		}
	}
	if got := IDHorario(19); got != "sex-noite" {
		t.Errorf("IDHorario(19) = %q, esperado sex-noite", got)
	}
	if _, ok := BitDoHorario("xyz-noite"); ok {
		t.Error("id inválido deveria ser rejeitado")
	}
}

func TestHorarioIDs(t *testing.T) {
	p := Perfil{Disponibilidade: 1<<19 | 1<<0}
	esperado := []string{"seg-madrugada", "sex-noite"}
	if got := p.HorarioIDs(); !slices.Equal(got, esperado) {
		t.Errorf("HorarioIDs() = %v, esperado %v", got, esperado)
	}
}
