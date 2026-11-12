# Roadmap

> **Atualizado em:** 10/11/2026, a partir do commit `309c1b9` mais as alterações ainda não commitadas (fechamento da Fase 4).
> **Build, vet e testes:** ok (`internal/domain` e `internal/seed`; os demais pacotes não têm testes).
> Detalhes de cada fase estão em `docs/contexto-claude.md` (seção 9). Se este arquivo divergir do código, vale o código.

## Visão geral

| Fase | Descrição | Semanas | Situação |
|---|---|---|---|
| 0 | Planejamento (escopo e arquitetura) | — | Concluída |
| 1 | Modelagem (questionário, grafo, personas) | — | Concluída, **sem validação do orientador** |
| 2 | Infraestrutura e esqueleto | — | Concluída |
| 3 | Domínio e gerador de dados mock | — | Concluída |
| 4 | Geração de candidatos no grafo | — | Concluída, **falta commitar o fechamento** |
| 5 | Scoring e reciprocidade | 11–12 | Concluída |
| 6 | API e formulário | 13 | Concluída |
| 7 | Avaliação e benchmarks | 14–16 | **Próxima** |
| 8 | Extras opcionais | 17 | Pendente (só se a 7 estiver completa) |
| 9 | Fechamento (monografia e apresentação) | 18 | Pendente |

```
[0]──[1]──[2]──[3]──[4]──[5]──[6]──▶[7]──[8?]──[9]
 ✔    ✔    ✔    ✔    ✔    ✔    ✔    agora           monografia em paralelo desde já
```

---

## O que foi feito

### Fase 0 — Planejamento
- Escopo restrito ao motor (D-011), Go (D-012), arquitetura híbrida Neo4j + Go (D-013), REST com `net/http` (D-014), reciprocidade por média harmônica (D-015).

### Fase 1 — Modelagem (`be430d9`)
- `docs/questionario.md`, `docs/modelo-grafo.md`, `docs/personas.md`.
- D-001 a D-010 (jogos, recomendação por jogo, atributos em `JOGA`, 18+, grade 7×4, funções como nós, propriedades, pesos, ASE ≠ ASA, complementaridade heurística).

### Fase 2 — Infraestrutura (`5dfc356`, `cb35183`)
- Docker Compose com Neo4j 5.26, `.env.example`, `internal/config`.
- `graph.Cliente` único (D-019), scripts Cypher idempotentes (D-018), `cmd/migrar`, `cmd/api` com `GET /saude`.

### Fase 3 — Domínio e mock (`6c175d3`)
- `internal/domain`: jogos, horários (bitmask de 28 bits), perfil, pesos, interações.
- `internal/seed`: 10 personas como dados, gerador reproduzível por semente (PCG), ruído e outliers, curtidas por amostragem por rejeição.
- Oráculo `PassaFiltros` e gabarito `Grau` (0, 1, 2).
- Carga em lote com `UNWIND` (`graph/carga.go`) e `cmd/seed`.
- 9 testes, incluindo `TestGeracaoDeterministica`. D-021 a D-027.

### Fase 4 — Candidatos (`1dacbab` + alterações não commitadas)
- `graph.BuscarCandidatos`: duas consultas (filtros + horários + complemento, corte em 300; depois sinal social). Filtros por jogo como fragmentos fixos.
- `graph.BuscarPerfis` e `BuscarPerfil`.
- `cmd/verificar`: verificação diferencial contra o oráculo e medição de latência.
- Ids de função em inglês (D-033). D-028 a D-033.
- **Fechamento (ainda não commitado):**
  - `Cliente.ValidarFuncoes` chamada no `cmd/seed` antes de apagar a base (defesa contra o MATCH silencioso).
  - D-032 preenchida: índice `usuario_regiao` removido; horários e avaliados coletados uma vez (−64% de db hits, −44% de latência p50).
  - Rótulos das funções no script de volta ao português, iguais à cópia em `modelo-grafo.md`.
  - README com a tabela de status completa e Go 1.26.5; `personas.md` passa a citar a plataforma nos filtros.

---

## O que falta

### Antes da Fase 5 (pendências curtas)

| # | Item | Por quê |
|---|---|---|
| 1 | Commitar o fechamento da Fase 4 e as Fases 5 e 6 | Tudo está fora do histórico |
| 2 | ~~Rodar `go run ./cmd/migrar`~~ | Feito em 10/11/2026 |
| 3 | Atualizar ou arquivar as seções 8 e 12 de `docs/contexto-claude.md` | Quase todas as pendências e inconsistências ali listadas foram resolvidas; o documento agora engana quem o lê |
| 4 | Levar ao orientador D-001 a D-010 e os três documentos da Fase 1 | Continuam "Proposta". É o maior risco do cronograma: uma mudança na modelagem agora repercute em todas as fases seguintes |

### Fase 5 — Scoring e reciprocidade (semanas 11–12)

**Decisões:** tomadas em D-034 a D-037 (10/11/2026). Eram:

1. Onde ficam `Candidato` e `Opcoes`: hoje em `graph`, que importa o driver. `matching` não pode depender disso. Mover para `domain` ou criar tipos próprios em `matching` com adaptador.
2. Normalização do sinal social (hoje inteiro sem limite).
3. Valores concretos do score: voz "às vezes", similaridade de tolerância à toxicidade, tabela `nivelDesejado` × diferença de nível, divisão pela soma dos pesos.
4. Casos de borda: Jaccard com bitmask vazio; média harmônica com os dois scores zero.
5. Definição de "relevante" no Precision@10 (`PassaFiltros` e grau ≥ 1?).

**Entregas:**
- [x] Interface de fonte de candidatos em `internal/matching` (D-013).
- [x] Matriz de complementaridade carregada do grafo para a memória.
- [x] Score unilateral s(A→B) com seis componentes: horário, intensidade, modo/servidor, comunicação, função, nível.
- [x] Reciprocidade `2·s(A→B)·s(B→A) / (s(A→B) + s(B→A))` (D-015).
- [x] Combinação com o sinal social (ponto de partida: `0,8 · recíproco + 0,2 · social`).
- [x] `Recomendar(ctx, id, jogo, k)` com o motivo do match.
- [x] Testes: propriedade (simetria, intervalo [0, 1]), tabela, caso RUST-APRENDIZ × RUST-GRINDER, e um teste que altera `Persona`/`PersonaJogo` e confirma que o ranking não muda.
- [x] `cmd/verificar` confere a complementaridade máxima calculada em Go contra a do grafo (não o componente do score).
- [x] Precision@10 simples (`cmd/avaliar`). **Saturado (D-037)**: rever a métrica na Fase 7.
- [x] `cmd/verificar` com 0 divergências (inclusive complementaridade) e `cmd/avaliar`: motor 0,981 × grafo 0,942; `Recomendar` em 71,7 ms em média (10/11/2026).

### Fase 6 — API e formulário (semana 13)
- [x] `GET /usuarios/{id}/recomendacoes?jogo=...&k=...`
- [x] `POST /usuarios` (transação única com conferência, D-040), `POST /usuarios/{id}/curtidas` e `/recusas` com match mútuo (D-039)
- [x] Middlewares de log (com latência) e de recuperação de panic.
- [x] Formulário HTML servido pelo próprio Go, com os rótulos em português lidos do grafo.
- [x] Resposta com o motivo do match.
- [x] Perfil público do candidato (`GET /usuarios/{id}`) e explicação do score na página (D-042).
- [x] Vocabulário do formulário e validação no `domain` (D-038); testado ponta a ponta contra o banco em 10/11/2026.

### Fase 7 — Avaliação (semanas 14–16)
- [ ] NDCG@10 graduado como métrica principal (o Precision@10 binário satura, D-037), com resultado por persona.
- [ ] Ablação: aleatório pós-filtro; só grafo; grafo + score unilateral; motor completo; motor sem sinal social (ameaça à validade da D-026).
- [ ] Quanto o corte de 300 perde de relevantes (D-029).
- [ ] Benchmarks do scoring (`go test -bench`).
- [ ] Latência com 1 mil, 10 mil, 50 mil e 100 mil usuários, separando banco e scoring. Reavaliar a complementaridade, hoje o trecho mais caro da consulta (D-032).
- [ ] Teste de carga no endpoint.
- [ ] Opcional: validação com 20 a 40 pessoas reais.

### Fase 8 — Extras opcionais (semana 17)
- [ ] Graph Data Science como versão extra na ablação (D-020).
- [ ] Ajuste automático de pesos a partir das curtidas.
- [ ] Interface de swipe.

### Fase 9 — Fechamento (semana 18)
- [ ] Revisão da monografia.
- [ ] README final.
- [ ] Ensaio da apresentação e vídeo da demonstração (plano B).

---

## Riscos

| Risco | Mitigação |
|---|---|
| Orientador pedir mudanças na modelagem (D-001 a D-010) depois da Fase 5 | Validar antes de fechar as fórmulas do score |
| Monografia acumulada para a semana 18 | Escrever em paralelo: fundamentação e desenvolvimento já têm material em `docs/` |
| Likes mock gerados a partir do próprio gabarito inflarem o sinal social (D-026) | Ablação sem sinal social na Fase 7 |
| Extras consumirem o tempo da avaliação | Fase 8 só começa com a Fase 7 completa |
| MATCH silencioso em novas cargas | `ValidarFuncoes` antes de apagar; `cmd/verificar` depois |
