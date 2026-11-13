package graph

import (
	"context"
	"fmt"

	"github.com/Chicosk1/squadr/internal/domain"
)

// Os parâmetros vazios desligam o filtro correspondente (um prefixo vazio vale para todos).
const filtroListagem = `
WHERE ($jogo = '' OR EXISTS { (u)-[:JOGA]->(:Jogo {id: $jogo}) })
  AND ($regiao = '' OR u.regiao = $regiao)
  AND u.id STARTS WITH $prefixo
  AND ($busca = '' OR toLower(u.nome) CONTAINS $busca OR toLower(u.id) STARTS WITH $busca)`

// A paginação é por chave (id > depois), que usa o índice da constraint usuario_id
// e não fica mais lenta nas últimas páginas, como ficaria com SKIP.
const queryListar = `
MATCH (u:Usuario)` + filtroListagem + `
  AND u.id > $depois
RETURN u.id AS id
ORDER BY id
LIMIT $limite`

const queryContar = `
MATCH (u:Usuario)` + filtroListagem + `
RETURN count(u) AS total`

// ListarUsuarios devolve uma página de usuários em ordem de id, com o total do filtro (D-043).
func (c *Cliente) ListarUsuarios(ctx context.Context, f domain.FiltroUsuarios) (domain.PaginaUsuarios, error) {
	params := map[string]any{
		"jogo":    string(f.Jogo),
		"regiao":  f.Regiao,
		"prefixo": f.PrefixoID,
		"busca":   f.Busca,
		"depois":  f.Depois,
		"limite":  f.Limite + 1, // um a mais diz se existe próxima página
	}
	res, err := c.Executar(ctx, queryListar, params)
	if err != nil {
		return domain.PaginaUsuarios{}, fmt.Errorf("listando usuários: %w", err)
	}
	ids := make([]string, 0, len(res.Records))
	for _, rec := range res.Records {
		ids = append(ids, texto(rec.AsMap(), "id"))
	}

	var pagina domain.PaginaUsuarios
	if len(ids) > f.Limite {
		ids = ids[:f.Limite]
		pagina.Proximo = ids[len(ids)-1]
	}

	res, err = c.Executar(ctx, queryContar, params)
	if err != nil {
		return domain.PaginaUsuarios{}, fmt.Errorf("contando usuários: %w", err)
	}
	if len(res.Records) == 1 {
		pagina.Total = inteiro(res.Records[0].AsMap(), "total")
	}

	perfis, err := c.BuscarPerfis(ctx, ids)
	if err != nil {
		return domain.PaginaUsuarios{}, err
	}
	// BuscarPerfis não garante a ordem: devolve na ordem dos ids.
	porID := make(map[string]domain.Perfil, len(perfis))
	for _, p := range perfis {
		porID[p.ID] = p
	}
	pagina.Perfis = make([]domain.Perfil, 0, len(ids))
	for _, id := range ids {
		if p, ok := porID[id]; ok {
			pagina.Perfis = append(pagina.Perfis, p)
		}
	}
	return pagina, nil
}
