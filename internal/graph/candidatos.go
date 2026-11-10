package graph

import (
	"context"
	"fmt"

	"github.com/Chicosk1/squadr/internal/domain"
)

var filtrosPorJogo = map[domain.Jogo]string{
	domain.ARK: `
		ro.versao = re.versao
		AND (ro.plataforma = re.plataforma
		     OR (ro.plataforma IN ['pc_steam', 'pc_epic'] AND re.plataforma IN ['pc_steam', 'pc_epic']))
		AND (ro.modo = re.modo OR ro.modo = 'ambos' OR re.modo = 'ambos')`,

	domain.Rust: `ro.versao = re.versao`,

	domain.LoL: `
		any(f IN re.filas WHERE f IN ro.filas
		    AND (f <> 'solo_duo'
		         OR (re.elo IS NOT NULL AND ro.elo IS NOT NULL
		             AND abs(re.elo - ro.elo) <= $maxDiffElo)))`,
}

// Os horários de eu e os usuários que eu já avaliou neste jogo são coletados uma vez,
// antes da expansão pelo jogo. Testá-los candidato a candidato com padrões a partir de eu
// (Expand(Into)) custava 64% mais db hits e quase o dobro da latência (D-032).
const queryCandidatos = `
MATCH (eu:Usuario {id: $id})-[re:JOGA]->(j:Jogo {id: $jogo})
WITH eu, re, j,
     [(eu)-[:DISPONIVEL]->(h:Horario) | h] AS meusHorarios,
     [(eu)-[x:CURTIU|RECUSOU]->(v) WHERE x.jogo = $jogo | v] AS vistos
MATCH (j)<-[ro:JOGA]-(outro:Usuario)
WHERE outro <> eu
  AND outro.regiao = eu.regiao
  AND outro.faixaEtaria IN eu.faixasAceitas
  AND eu.faixaEtaria IN outro.faixasAceitas
  AND any(i IN outro.idiomas WHERE i IN eu.idiomas)
  AND (%s)
  AND NOT outro IN vistos
WITH eu, outro, size([(outro)-[:DISPONIVEL]->(h) WHERE h IN meusHorarios | h]) AS horarios
WHERE horarios >= $minHorarios
OPTIONAL MATCH (eu)-[:EXERCE]->(:Funcao {jogo: $jogo})-[c:COMPLEMENTA]-(:Funcao)<-[:EXERCE]-(outro)
WITH outro, horarios, max(c.valor) AS complemento
RETURN outro.id AS id, horarios, complemento
ORDER BY horarios DESC, coalesce(complemento, 0) DESC, id
LIMIT $limite`

const querySocial = `
MATCH (eu:Usuario {id: $id})-[:CURTIU {jogo: $jogo}]->(:Usuario)<-[:CURTIU {jogo: $jogo}]-(sim:Usuario)
WHERE sim <> eu
WITH sim, count(*) AS emComum
MATCH (sim)-[:CURTIU {jogo: $jogo}]->(alvo:Usuario)
WHERE alvo.id IN $ids
RETURN alvo.id AS id, sum(emComum) AS sinal`

func (c *Cliente) BuscarCandidatos(ctx context.Context, id string, jogo domain.Jogo, op domain.OpcoesCandidatos) ([]domain.Candidato, error) {
	filtro, ok := filtrosPorJogo[jogo]
	if !ok {
		return nil, fmt.Errorf("jogo desconhecido: %q", jogo)
	}

	res, err := c.Executar(ctx, fmt.Sprintf(queryCandidatos, filtro), map[string]any{
		"id":          id,
		"jogo":        string(jogo),
		"minHorarios": op.MinHorarios,
		"limite":      op.Limite,
		"maxDiffElo":  op.MaxDiffElo,
	})
	if err != nil {
		return nil, fmt.Errorf("buscando candidatos de %s em %s: %w", id, jogo, err)
	}

	candidatos := make([]domain.Candidato, 0, len(res.Records))
	for _, rec := range res.Records {
		m := rec.AsMap()
		candidatos = append(candidatos, domain.Candidato{
			ID:             texto(m, "id"),
			HorariosComuns: inteiro(m, "horarios"),
			Complemento:    decimal(m, "complemento"),
		})
	}

	if op.SinalSocial && len(candidatos) > 0 {
		ids := make([]string, len(candidatos))
		for i, cand := range candidatos {
			ids[i] = cand.ID
		}
		sinais, err := c.sinalSocial(ctx, id, jogo, ids)
		if err != nil {
			return nil, err
		}
		for i := range candidatos {
			candidatos[i].SinalSocial = sinais[candidatos[i].ID]
		}
	}
	return candidatos, nil
}

func (c *Cliente) sinalSocial(ctx context.Context, id string, jogo domain.Jogo, ids []string) (map[string]int, error) {
	res, err := c.Executar(ctx, querySocial, map[string]any{"id": id, "jogo": string(jogo), "ids": ids})
	if err != nil {
		return nil, fmt.Errorf("calculando sinal social de %s em %s: %w", id, jogo, err)
	}
	sinais := make(map[string]int, len(res.Records))
	for _, rec := range res.Records {
		m := rec.AsMap()
		sinais[texto(m, "id")] = inteiro(m, "sinal")
	}
	return sinais, nil
}
