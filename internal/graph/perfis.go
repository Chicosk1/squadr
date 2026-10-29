package graph

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Chicosk1/squadr/internal/domain"
)

var ErrNaoEncontrado = errors.New("usuário não encontrado")

const queryPerfis = `
UNWIND $ids AS uid
MATCH (u:Usuario {id: uid})
RETURN u {.*} AS perfil,
       [(u)-[r:JOGA]->(j:Jogo) | {jogo: j.id, props: properties(r)}] AS jogos,
       [(u)-[r:EXERCE]->(f:Funcao) | {id: f.id, jogo: f.jogo, prioridade: r.prioridade}] AS funcoes,
       [(u)-[:DISPONIVEL]->(h:Horario) | h.bit] AS bits`

func (c *Cliente) BuscarPerfis(ctx context.Context, ids []string) ([]domain.Perfil, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	res, err := c.Executar(ctx, queryPerfis, map[string]any{"ids": ids})
	if err != nil {
		return nil, fmt.Errorf("buscando perfis: %w", err)
	}
	perfis := make([]domain.Perfil, 0, len(res.Records))
	for _, rec := range res.Records {
		p, err := perfilDoRegistro(rec.AsMap())
		if err != nil {
			return nil, err
		}
		perfis = append(perfis, p)
	}
	return perfis, nil
}

func (c *Cliente) BuscarPerfil(ctx context.Context, id string) (domain.Perfil, error) {
	perfis, err := c.BuscarPerfis(ctx, []string{id})
	if err != nil {
		return domain.Perfil{}, err
	}
	if len(perfis) == 0 {
		return domain.Perfil{}, fmt.Errorf("%w: %s", ErrNaoEncontrado, id)
	}
	return perfis[0], nil
}

func perfilDoRegistro(m map[string]any) (domain.Perfil, error) {
	u, ok := m["perfil"].(map[string]any)
	if !ok {
		return domain.Perfil{}, errors.New("registro sem o campo perfil")
	}

	p := domain.Perfil{
		ID:                   texto(u, "id"),
		Nome:                 texto(u, "nome"),
		FaixaEtaria:          texto(u, "faixaEtaria"),
		FaixasAceitas:        lista(u, "faixasAceitas"),
		Regiao:               texto(u, "regiao"),
		Idiomas:              lista(u, "idiomas"),
		UsaVoz:               texto(u, "usaVoz"),
		ExigeVoz:             booleano(u, "exigeVoz"),
		ToleranciaToxicidade: inteiro(u, "toleranciaToxicidade"),
		Compromisso:          texto(u, "compromisso"),
		Persona:              texto(u, "persona"),
		Pesos: domain.Pesos{
			Horario:     inteiro(u, "pesoHorario"),
			Intensidade: inteiro(u, "pesoIntensidade"),
			Modo:        inteiro(u, "pesoModo"),
			Comunicacao: inteiro(u, "pesoComunicacao"),
			Funcao:      inteiro(u, "pesoFuncao"),
			Nivel:       inteiro(u, "pesoNivel"),
		},
	}

	bits, _ := m["bits"].([]any)
	for _, b := range bits {
		if n, ok := b.(int64); ok {
			p.Disponibilidade |= 1 << n
		}
	}

	principal := map[domain.Jogo]string{}
	secundaria := map[domain.Jogo]string{}
	funcoes, _ := m["funcoes"].([]any)
	for _, f := range funcoes {
		fm, _ := f.(map[string]any)
		jogo := domain.Jogo(texto(fm, "jogo"))
		switch inteiro(fm, "prioridade") {
		case 1:
			principal[jogo] = texto(fm, "id")
		case 2:
			secundaria[jogo] = texto(fm, "id")
		}
	}

	jogos, _ := m["jogos"].([]any)
	for _, j := range jogos {
		jm, _ := j.(map[string]any)
		props, _ := jm["props"].(map[string]any)
		jogo := domain.Jogo(texto(jm, "jogo"))

		pj := domain.PerfilJogo{
			Jogo:             jogo,
			Intensidade:      inteiro(props, "intensidade"),
			Nivel:            decimal(props, "nivel"),
			NivelDesejado:    texto(props, "nivelDesejado"),
			FuncaoPrincipal:  principal[jogo],
			FuncaoSecundaria: secundaria[jogo],
			Horas:            texto(props, "horas"),
			Versao:           texto(props, "versao"),
			Servidor:         texto(props, "servidor"),
			Plataforma:       texto(props, "plataforma"),
			Modo:             texto(props, "modo"),
			TamanhoTribo:     texto(props, "tamanhoTribo"),
			LimiteGrupo:      texto(props, "limiteGrupo"),
			Wipe:             texto(props, "wipe"),
			Elo:              domain.SemRanked,
			Filas:            lista(props, "filas"),
			PersonaJogo:      texto(props, "personaJogo"),
		}
		if elo, ok := props["elo"].(int64); ok {
			pj.Elo = int(elo)
		}
		p.Jogos = append(p.Jogos, pj)
	}
	slices.SortFunc(p.Jogos, func(a, b domain.PerfilJogo) int { return cmp.Compare(a.Jogo, b.Jogo) })
	return p, nil
}
