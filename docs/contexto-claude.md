# Contexto para continuar o squadr (TCC)

> **Gerado em:** 09/11/2026, a partir do commit `1dacbab` (branch `main`, árvore limpa).
> **`go build ./...`:** ok (Go 1.26.5 windows/amd64). **`go vet ./...`:** ok.
> **`go test ./...`:** ok. `internal/domain` (2 testes) e `internal/seed` (9 testes) passam; `cmd/*`, `internal/config` e `internal/graph` não têm testes.
> **Não executados:** `cmd/seed`, `cmd/migrar`, `cmd/verificar`, `cmd/api` e Docker (exigem banco; o seed apaga dados). A afirmação "o `cmd/verificar` termina com 0 divergências" não foi reconferida nesta geração.

---

## 1. Como usar este documento

Você é a nova sessão do Claude que vai continuar este TCC. Você não tem acesso à conversa em que o projeto foi planejado; este arquivo resume o que só existia nela.

1. Leia este arquivo inteiro e `docs/decisoes.md` antes de propor qualquer coisa.
2. O repositório é a fonte da verdade. Se algo aqui contradisser o código, vale o código; aponte a divergência ao aluno.
3. Respeite as convenções da seção 5. Elas não são sugestões.
4. Continue pela **Fase 5** (seção 9). Antes de começar, confirme com o aluno as inconsistências da seção 12, principalmente a ausência de `Cliente.ValidarFuncoes`.
5. Entregue no formato da seção 11. O aluno faz os commits.
6. Se você só tiver os arquivos anexados ao Projeto, peça ao aluno o conteúdo atual de qualquer arquivo antes de reescrevê-lo. Nunca reconstrua um arquivo de memória.

---

## 2. O projeto em uma página

- **O quê:** TCC de um motor de recomendação no estilo "match" para encontrar parceiros em jogos online que não têm bom sistema de formação de grupos.
- **Escopo:** só o motor, com dados mockados e um formulário representativo (D-011). Chat, autenticação, interface completa e integração com Steam e Riot são trabalho futuro.
- **Jogos:** ARK (o caso que motiva o trabalho), Rust e League of Legends (D-001).
- **Arquitetura híbrida (D-013):** o orientador recomendou banco de grafos. O Neo4j gera os candidatos (filtros rígidos, horários em comum, complementaridade de funções, sinal social); o Go calcula o score, a reciprocidade e o ranking.
- **Recomendação sempre por jogo (D-002):** o motor responde "quem é um bom parceiro para *este* jogo".
- **Avaliação:** personas com ruído geram os dados e um gabarito de relevância (graus 0, 1 e 2) usado para Precision@K e NDCG@K (`docs/personas.md`).
- **Stack:** Go (`go.mod`: `go 1.26.5`), Neo4j 5.26 Community em Docker Compose, `godotenv` v1.5.1 e `neo4j-go-driver/v5` v5.28.5. Nenhuma outra dependência.
- **Módulo:** `github.com/Chicosk1/squadr`.
- **Estado:** segundo o aluno, Fases 0 a 4 concluídas; a próxima é a 5. O repositório tem a Fase 4 no código, mas com pendências (seções 8 e 12).
- **Postura esperada:** respostas críticas, que apontem quando algo não combina, e focadas em eficiência.

---

## 3. Arquitetura e pipeline

```
Usuário A + jogo
  │
  ▼ Consulta 1 — Cypher (internal/graph/candidatos.go)           [existe]
  │  filtros rígidos gerais: mesma região, faixa etária aceita pelos dois,
  │  pelo menos 1 idioma em comum, sem CURTIU/RECUSOU de A para B neste jogo
  │  + filtro do jogo (fragmento Cypher fixo por jogo, D-030)
  │  → horários em comum = count(Horario compartilhado) ≥ MinHorarios (D-028)
  │  → complemento = max(COMPLEMENTA.valor) entre funções de A e de B
  │  → ORDER BY horarios DESC, complemento DESC, id  LIMIT 300 (D-029)
  ▼ Consulta 2 — Cypher, só para os 300                             [existe]
  │  sinal social (filtragem colaborativa por travessia)
  ▼ Scoring em Go puro (internal/matching)                          [Fase 5]
  │  s(A→B) e s(B→A) → média harmônica → combina com o sinal social
  ▼ Ranking top K com o motivo do match                             [Fase 5]
  ▼ API net/http (internal/api) + formulário HTML                   [Fase 6]
```

**Filtros por jogo** (iguais no Cypher e no oráculo `seed.PassaFiltros`):

| Jogo | Filtro rígido |
|---|---|
| ARK | Mesma versão (ASE ≠ ASA); mesmo grupo de plataforma (PC Steam e Epic juntos, cada console isolado); modo igual ou um dos dois em `ambos` |
| Rust | Mesma versão (`pc` ≠ `console`) |
| LoL | Pelo menos uma fila em comum; `solo_duo` só conta se os dois são ranqueados e a diferença absoluta de elo é no máximo `MaxDiffEloDuo` (2) |

**Sinais devolvidos em `graph.Candidato`:** `HorariosComuns` (int), `Complemento` (float, 0 se nulo), `SinalSocial` (int: para cada usuário que curtiu alguém que A também curtiu no jogo, soma o número de curtidas em comum em cada candidato que ele curtiu).

**Opções (`graph.OpcoesPadrao`)**: `Limite` 300, `MinHorarios` 1, `MaxDiffElo` = `domain.MaxDiffEloDuo`, `SinalSocial` true. O `cmd/verificar` usa `MinHorarios` 0, porque o oráculo não tem filtro de horário.

**Modelo do grafo** (detalhes em `docs/modelo-grafo.md`):

```
(:Usuario)-[:JOGA {intensidade, nivel, nivelDesejado, personaJogo, + campos do jogo}]->(:Jogo)
(:Usuario)-[:EXERCE {prioridade: 1|2}]->(:Funcao {id, jogo, nome})
(:Usuario)-[:DISPONIVEL]->(:Horario {id, dia, periodo, bit 0..27})
(:Funcao)-[:COMPLEMENTA {valor 0..1}]-(:Funcao)      // 51 pares, gravados uma vez, lidos sem direção
(:Usuario)-[:CURTIU|RECUSOU {jogo, em}]->(:Usuario)
```

Dados de referência: 3 jogos, 28 horários, 16 funções, 51 relações `COMPLEMENTA`. Região, idiomas, faixa etária e pesos são propriedades do `Usuario` (pesos achatados em `pesoHorario`, `pesoIntensidade` etc.).

---

## 4. Estrutura do repositório

| Caminho | Conteúdo |
|---|---|
| `cmd/api/main.go` | Servidor `net/http`; hoje só `GET /saude` (ping no Neo4j, responde `{"banco":"ok","status":"ok"}`), com encerramento gracioso |
| `cmd/migrar/main.go` | Executa `scripts/schema.cypher` e `scripts/referencia.cypher` e confere as contagens (3, 28, 16, 51). Idempotente |
| `cmd/seed/main.go` | Gera usuários e interações pela semente, **apaga todos os `Usuario`** e carrega em lote. Flags `-usuarios` (10000, mín. 100) e `-semente` (42) |
| `cmd/verificar/main.go` | Verificação diferencial: regenera os dados pela semente, confere a contagem de usuários no banco, a ida e volta de perfis, os candidatos contra `PassaFiltros`, horários em comum, sinal social e a ordem do `LIMIT`; mede a latência (p50, p95, máx.). Flags `-usuarios`, `-semente` (iguais às do seed) e `-amostra` (200) |
| `internal/config/config.go` | Lê `.env` (godotenv) e `NEO4J_URI`, `NEO4J_USER`, `NEO4J_PASSWORD`, `NEO4J_DATABASE`, `API_PORT` |
| `internal/domain/jogo.go` | `Jogo` (`ark`, `rust`, `lol`), `Funcoes` por jogo, `MaxDiffEloDuo = 2` |
| `internal/domain/horario.go` | 7 dias × 4 períodos, `IDHorario(bit)`, `BitDoHorario(id)` |
| `internal/domain/perfil.go` | `Perfil`, `PerfilJogo`, `Pesos`, `SemRanked = -1`, `PerfilDoJogo`, `HorarioIDs` e conversão para propriedades do grafo |
| `internal/domain/interacao.go` | `Interacao` (de, para, jogo, curtiu, em) |
| `internal/graph/conexao.go` | `Cliente` (driver único, D-019), `Executar`, `ExecutarAutoCommit`, `Ping`, `Fechar` |
| `internal/graph/migrar.go` | `ExecutarScript` (divide por `;`, ignora linhas `//`) |
| `internal/graph/carga.go` | `LimparUsuarios` (`IN TRANSACTIONS OF 5000 ROWS`), `CarregarUsuarios` (lotes de 2.000), `CarregarInteracoes` (lotes de 5.000) |
| `internal/graph/candidatos.go` | `Opcoes`, `OpcoesPadrao`, `Candidato`, `BuscarCandidatos` (as duas consultas) |
| `internal/graph/perfis.go` | `BuscarPerfis` (lote via `UNWIND`), `BuscarPerfil`, `ErrNaoEncontrado` |
| `internal/graph/valores.go` | Conversores de tipos do driver (`int64`, `[]any` etc.) |
| `internal/seed/personas.go` | As 10 personas como dados (`var Personas []Persona`) e os modelos por jogo (`ModeloJogo`) |
| `internal/seed/gerador.go` | `Gerador` (PCG + semente), ruído, outliers e normalização do nível |
| `internal/seed/interacoes.go` | Curtidas e recusas por amostragem por rejeição (D-026) |
| `internal/seed/gabarito.go` | Oráculo `PassaFiltros` e gabarito `Grau` (0, 1, 2) |
| `internal/matching/`, `internal/api/` | Só `.gitkeep`: Fases 5 e 6 |
| `scripts/schema.cypher` | Constraints de unicidade (`Usuario`, `Jogo`, `Horario`, `Funcao`) e índice `usuario_regiao` |
| `scripts/referencia.cypher` | Jogos, horários, funções e `COMPLEMENTA` com `MERGE` |
| `docs/questionario.md` | Perguntas do formulário (G, A, R, L, P) e o uso de cada uma no motor |
| `docs/modelo-grafo.md` | Nós, relações, matrizes de complementaridade, índices e filtros da Fase 4 |
| `docs/personas.md` | Personas, regras de ruído, gabarito e likes mock |
| `docs/decisoes.md` | Registro D-001 a D-033 |
| `docker-compose.yml` | `neo4j:5.26-community`, portas 7474 e 7687, heap 512M a 1G, pagecache 512M, volume `neo4j_data`, healthcheck |
| `.env.example` | Modelo do `.env` (o `.env` real não é versionado e não deve ser lido) |

**Testes existentes:** `TestIDeBitDoHorario`, `TestHorarioIDs` (domain); `TestGrade`, `TestGeracaoDeterministica`, `TestProporcoes`, `TestInvariantes`, `TestPassaFiltros` (12 casos em tabela), `TestPassaFiltrosGerais`, `TestGrau`, `TestGrauMulti`, `TestInteracoes` (seed). Não há benchmarks.

**Ordem de execução** (Git Bash, na raiz): `docker compose up -d` → `go run ./cmd/migrar` → `go run ./cmd/seed` → `go run ./cmd/verificar` (mesmos `-usuarios` e `-semente`) → `go run ./cmd/api`.

---

## 5. Convenções obrigatórias

| Convenção | Detalhe |
|---|---|
| Módulo | `github.com/Chicosk1/squadr` |
| Recomendação por jogo | Toda consulta, curtida e recusa carrega o jogo (D-002) |
| Ids de função em inglês (D-033) | `ark:breeder`, `ark:builder`, `ark:gatherer`, `ark:pvp`, `ark:explorer`; `rust:builder`, `rust:gatherer`, `rust:pvp`, `rust:raider`, `rust:electrician`; `lol:top`, `lol:jungle`, `lol:mid`, `lol:adc`, `lol:support`, `lol:fill`. Nomes de persona e rótulos do formulário ficam em português |
| API só com `net/http` | Roteamento do Go 1.22+ (`"GET /rota/{id}"`). Sem chi, gin, echo ou fiber (D-014): API pequena, gargalo no banco, menos dependências |
| `persona` e `personaJogo` | Existem só para o gabarito da avaliação. **O motor nunca pode lê-los.** Atenção: `graph.BuscarPerfis` preenche esses campos em `domain.Perfil` (para a ida e volta do verificador), então eles estarão na struct que o matching recebe |
| Gerador reproduzível | Nunca iterar sobre `map` em código que consome números aleatórios. Personas são dados (um slice), não código (D-021, D-022) |
| `MaxDiffEloDuo` | Vive em `internal/domain` e é a mesma no oráculo e na consulta Cypher. Passe sempre `domain.MaxDiffEloDuo` em `Opcoes.MaxDiffElo` |
| Scripts Cypher | Versionados, idempotentes (`MERGE`), são a fonte única do grafo de referência (D-018) |
| Consultas Cypher | Nada vindo do usuário entra no texto da query; só parâmetros (D-030) |
| Driver | Um único `graph.Cliente` compartilhado, nunca um por requisição (D-019) |
| Decisões | Toda decisão nova vai para `docs/decisoes.md` como D-XXX, com contexto, opções consideradas, decisão, motivo e consequências (o modelo do arquivo também tem Data e Status). A próxima é a **D-034** |

---

## 6. O que já foi feito

O repositório tem uma história anterior (jul. a set./2026: app mobile e uma "platform" em Go) que foi removida no commit `a8ea4b9` (14/10/2026), "conforme nova proposta em conversa com professores". Ignore essa história.

| Fase | Commits | Entregas que existem no repositório |
|---|---|---|
| 0 — Planejamento | — | Não há registro explícito de uma Fase 0. As decisões D-011 a D-015 ("definidas no planejamento") são o que mais se aproxima |
| 1 — Modelagem | `be430d9` | `docs/questionario.md`, `docs/modelo-grafo.md`, `docs/personas.md`, `docs/decisoes.md` (D-001 a D-010, todas "Proposta") e o esqueleto de pastas |
| 2 — Infraestrutura | `5dfc356`, `cb35183` | `docker-compose.yml`, `.env.example`, `go.mod`, `internal/config`, `graph/conexao.go`, `graph/migrar.go`, `scripts/*.cypher`, `cmd/migrar`, `cmd/api` com `/saude`, README; D-016 a D-020 |
| 3 — Domínio e mock | `6c175d3` | `internal/domain`, `internal/seed` (personas, gerador, interações, oráculo e gabarito, testes), `graph/carga.go`, `cmd/seed`; D-021 a D-027 |
| 4 — Candidatos | `1dacbab` (29/10/2026) | `graph/candidatos.go` (duas consultas, filtros por jogo), `graph/perfis.go`, `graph/valores.go`, `cmd/verificar`; `MaxDiffEloDuo` movida de `seed` para `domain`; ids de função traduzidos para inglês no Go e no script; D-028 a D-033. **Faltam:** D-032 (PROFILE) e `Cliente.ValidarFuncoes` (seção 12) |

---

## 7. Índice das decisões

Lido de `docs/decisoes.md`. **Proposta** = aguardando validação; as demais estão como Aceita.

| ID | Decisão | Status |
|---|---|---|
| D-001 | Escopo de jogos: ARK, Rust e League of Legends | **Proposta** |
| D-002 | Recomendação contextualizada por jogo | **Proposta** |
| D-003 | Atributos dependentes do jogo na relação `JOGA` | **Proposta** |
| D-004 | Público 18+ e idade coletada por faixa | **Proposta** |
| D-005 | Disponibilidade em 7 dias × 4 períodos, horário de Brasília | **Proposta** |
| D-006 | Funções como nós e complementaridade armazenada no grafo | **Proposta** |
| D-007 | Região, idioma e faixa etária como propriedades | **Proposta** |
| D-008 | Pesos opcionais (1 a 5) com padrão 3 | **Proposta** |
| D-009 | ASE e ASA como versões incompatíveis (filtro rígido) | **Proposta** |
| D-010 | Valores de complementaridade como premissas heurísticas | **Proposta** |
| D-011 | Foco exclusivo no motor de recomendação | Aceita |
| D-012 | Go como linguagem do backend | Aceita |
| D-013 | Arquitetura híbrida: Neo4j para candidatos, Go para scoring | Aceita |
| D-014 | API REST com `net/http` em vez de GraphQL | Aceita |
| D-015 | Reciprocidade por média harmônica | Aceita |
| D-016 | Neo4j Community em Docker Compose | Aceita |
| D-017 | Configuração por variáveis de ambiente com `.env` | Aceita |
| D-018 | Scripts Cypher versionados e idempotentes | Aceita |
| D-019 | Driver único compartilhado | Aceita |
| D-020 | Graph Data Science adiado (Fase 8) | Aceita |
| D-021 | Dados mock em Go, reproduzíveis por semente (PCG, padrão 42) | Aceita |
| D-022 | Personas como dados, não código | Aceita |
| D-023 | Ruído e outliers apenas em atributos flexíveis | Aceita |
| D-024 | Filtros rígidos reimplementados em Go como oráculo | Aceita |
| D-025 | `personaJogo` na relação `JOGA` (somente mock) | Aceita |
| D-026 | Likes por amostragem por rejeição, probabilidade por grau (ameaça à validade) | Aceita |
| D-027 | Carga em lote com `UNWIND` | Aceita |
| D-028 | Mínimo de horários em comum como filtro configurável | Aceita |
| D-029 | Candidatos em duas consultas (corte em 300) | Aceita |
| D-030 | Filtros por jogo como fragmentos Cypher fixos | Aceita |
| D-031 | Verificação diferencial contra o oráculo | Aceita |
| D-032 | Resultados do PROFILE (**não preenchida**) | **Proposta** |
| D-033 | Ids de funções no vocabulário da comunidade (inglês) | Aceita |

---

## 8. Pendências

Pendências trazidas da conversa de planejamento, conferidas no repositório:

| # | Pendência | Situação no repositório |
|---|---|---|
| 1 | Validação do orientador para D-001 a D-010 e para D-002, D-028 e D-033 | **Aberta.** D-001 a D-010 seguem "Proposta"; `questionario.md`, `modelo-grafo.md` e `personas.md` têm cabeçalho "Proposta". D-028 e D-033 estão como "Aceita" (seção 12) |
| 2 | Preencher a D-032 com o PROFILE (tabela de db hits, experimentos A, B e C) e decidir sobre o índice `usuario_regiao` | **Aberta.** A D-032 só tem os textos-modelo. Os experimentos e o "Passo 7.3" citado nela não estão descritos no repositório: peça ao aluno |
| 3a | Atualizar `docs/modelo-grafo.md` | **Quase feita.** Já tem `personaJogo` em `JOGA`, normalização do nível, filtros e sinais da Fase 4 e ids em inglês. Falta: cabeçalho ainda diz "Fase 1 / Proposta" e a cópia do script de funções diverge de `scripts/referencia.cypher` nos nomes (seção 12) |
| 3b | README com a tabela de status em dia | **Aberta.** Fase 4 aparece "Em andamento"; a tabela não tem as Fases 0, 8 e 9; a estrutura omite `cmd/verificar`; pede Go 1.22 enquanto o `go.mod` exige 1.26.5 |
| 4 | Crossplay do ARK e diferença de elo no duo do LoL: confirmar as regras reais ou declará-las como premissas na monografia | **Aberta.** Já estão declaradas como premissas em `personas.md`, `questionario.md` e na D-024 |
| 5 | Implementar `Cliente.ValidarFuncoes` e chamá-la no `cmd/seed` antes de `LimparUsuarios` | **Aberta (nova).** A D-033 descreve a função, mas ela não existe no código (seção 12) |
| 6 | Reexecutar `cmd/verificar` e confirmar 0 divergências | Não reconferido nesta geração |

---

## 9. Próximas fases

### Fase 5 — Scoring e reciprocidade (semanas 11 e 12)

- **Pacote `internal/matching` em Go puro**, sem depender do Neo4j. O recomendador depende de uma **interface de fonte de candidatos** (D-013).
- **Matriz de complementaridade** carregada do grafo para a memória na inicialização (D-006).
- **Score unilateral s(A→B)**, ponderado pelos pesos de A, com seis componentes:

| Componente | Peso | Regra planejada |
|---|---|---|
| Horário | P01 `pesoHorario` | Jaccard sobre o bitmask de 28 bits |
| Intensidade e compromisso | P02 `pesoIntensidade` | Intensidade do jogo + G10 |
| Modo e servidor | P03 `pesoModo` | Por jogo: ARK servidor e tamanho de tribo; Rust servidor, limite de grupo e wipe; LoL filas (mapeamento do `questionario.md`) |
| Comunicação | P04 `pesoComunicacao` | Voz assimétrica (se A exige voz e B não usa, zero) + similaridade da tolerância a toxicidade |
| Função | P05 `pesoFuncao` | Complementaridade, com a função principal pesando mais que a secundária |
| Nível | P06 `pesoNivel` | Assimétrico, conforme o `nivelDesejado` de A |

- **Reciprocidade (D-015):** `S(A,B) = 2·s(A→B)·s(B→A) / (s(A→B) + s(B→A))`.
- **Combinação com o sinal social normalizado:** ponto de partida `0,8 · recíproco + 0,2 · social`.
- **`Recomendar(ctx, id, jogo, k)`:** devolve os K melhores com o motivo do match.
- **Testes:** de propriedade (o score recíproco é simétrico e fica entre 0 e 1), em tabela, e o caso de controle **RUST-APRENDIZ × RUST-GRINDER**, que deve cair no ranking com reciprocidade.
- **`cmd/verificar`** ganha a checagem da complementaridade calculada em Go contra a do grafo (prevista na D-031).
- **Ao fim da fase:** uma versão simples de Precision@10.

**O que já existe para apoiar a Fase 5:** `graph.BuscarCandidatos` (candidatos e sinais), `graph.BuscarPerfis` (perfis completos em uma consulta), `graph.BuscarPerfil`, os tipos de `domain`, e `seed.PassaFiltros` e `seed.Grau` para a avaliação.

**Pontos que o plano não especifica** (observações desta geração; decidir com o aluno e registrar como D-034 em diante):

1. `graph.Candidato` e `graph.Opcoes` estão no pacote `graph`, que importa o driver. Se a interface de `matching` usar esses tipos, `matching` passa a depender do Neo4j. Será preciso mover os tipos para `domain` ou definir tipos próprios em `matching` com um adaptador.
2. O `Complemento` do Cypher é o máximo entre todos os pares de funções, sem distinguir principal e secundária. O componente de função da Fase 5 pondera as duas, então não é o mesmo número. A checagem no verificador deve comparar o máximo calculado em Go com o do grafo, e não o componente do score.
3. A normalização do sinal social (hoje um inteiro sem limite) não está definida.
4. Faltam valores concretos para: voz quando B usa "às vezes"; a fórmula da similaridade de tolerância; cada combinação de `nivelDesejado` × diferença de nível; a soma dos pesos (dividir pela soma dos pesos de A mantém s em [0, 1]).
5. Casos de borda: Jaccard com bitmask vazio; média harmônica com os dois scores iguais a zero (divisão por zero).
6. Um teste que altera `Persona`/`PersonaJogo` e confirma que o ranking não muda garante a regra de que o motor não lê esses campos.
7. Precision@10: `Grau` não aplica os filtros. Relevante = `PassaFiltros` e grau ≥ 1? Definir.

### Fase 6 — API e formulário (semana 13)

`GET /usuarios/{id}/recomendacoes?jogo=...`, `POST /usuarios` e `POST /usuarios/{id}/curtidas`; middlewares de log (que também mede a latência) e de recuperação de panic; página HTML simples servida pelo próprio Go como formulário; resposta com o motivo do match; handlers em `internal/api`.

### Fase 7 — Avaliação (semanas 14 a 16)

- Precision@10 e NDCG@10 com o gabarito das personas.
- Ablação com cinco versões: aleatório após os filtros; só o ranking do grafo; grafo + score unilateral; motor completo; motor completo sem sinal social (por causa da ameaça à validade da D-026).
- Medir quanto o corte de 300 candidatos perde de relevantes (D-029).
- Benchmarks do scoring com `go test -bench`.
- Latência com 1 mil, 10 mil, 50 mil e 100 mil usuários, separando o tempo de banco do tempo de scoring.
- Teste de carga no endpoint.
- Opcional: validação com 20 a 40 pessoas reais.

### Fase 8 — Extras opcionais (semana 17)

Graph Data Science como versão extra na ablação (D-020), ajuste automático de pesos a partir das curtidas, interface de swipe.

### Fase 9 — Fechamento (semana 18)

Revisão da monografia, README final, ensaio da apresentação e vídeo da demonstração como plano B.

### Riscos e recomendações

Escrever a monografia em paralelo (os `docs/` já têm material para a fundamentação e o desenvolvimento), avaliar cedo e não investir em extras antes de a Fase 7 estar completa.

---

## 10. Lições aprendidas

- **MATCH silencioso na carga Cypher.** Um `MATCH` que não encontra o nó faz a relação não ser criada, sem nenhum erro. Aconteceu duas vezes: relações faltando na Fase 3 e ids de função traduzidos na Fase 4 (o Go gerava `ark:gatherer`, o grafo tinha `ark:coletor`, e `EXERCE` não era criada). As defesas previstas são:
  - `Cliente.ValidarFuncoes`, chamada no seed antes de apagar a base. **Ainda não existe no código** (seção 12).
  - `cmd/verificar`, que faz verificação diferencial: ida e volta de perfis pelo grafo e comparação dos candidatos do Cypher com o oráculo `PassaFiltros` em Go. Foi ele que detectou o erro da Fase 4.
- **Duas implementações independentes pegam bugs.** O oráculo em Go não depende da query; qualquer divergência aponta um bug em uma das duas (D-024, D-031). Mantenha essa independência ao adicionar regras.
- **Reprodutibilidade é frágil.** Iterar sobre `map` muda a ordem de consumo dos números aleatórios e quebra a semente. `TestGeracaoDeterministica` protege isso.
- **Apagar tudo em uma transação estoura a memória.** Por isso a limpeza usa `CALL ... IN TRANSACTIONS` em sessão auto-commit (D-027).
- **O verificador exige os mesmos `-usuarios` e `-semente` do seed.** Ele confere a contagem de usuários e falha se não bater.

---

## 11. Como trabalhar com o aluno

- Respostas críticas: aponte quando algo não combina, inclusive no plano. Priorize eficiência.
- Cada fase é entregue com:
  1. Objetivo.
  2. Mudanças em relação ao plano.
  3. Passos com **arquivos completos prontos para colar** (caminho do arquivo e conteúdo inteiro, não trechos).
  4. Comandos para rodar.
  5. O que conferir.
  6. Seção **"O que você aprende"**.
  7. **"Critério de pronto"**.
  8. **"Armadilhas comuns"**.
- Documentação vem em markdown pronta para colar, com o nome do arquivo. Decisões novas no formato D-XXX (seção 5).
- **O aluno faz os commits.** Não proponha nem execute operações de Git.
- Ambiente: Windows com Git Bash (MINGW64). Comandos devem funcionar no Git Bash; no PowerShell, `curl` é `curl.exe`.
- Nunca peça nem copie o conteúdo do `.env` (tem a senha do banco).
- Lembre que `go run ./cmd/seed` apaga os usuários do banco.

---

## 12. Inconsistências encontradas

Divergências entre o código, os documentos e o que veio da conversa de planejamento. Nenhuma foi corrigida; todas aguardam decisão do aluno.

| # | Onde | O que se afirma | O que o repositório mostra |
|---|---|---|---|
| 1 | D-033 e conversa de planejamento | `Cliente.ValidarFuncoes` existe e é chamada no seed antes de apagar a base | **Não existe no código** nem em nenhum commit (`git log -S` só encontra a menção na D-033). `cmd/seed/main.go` chama `LimparUsuarios` sem nenhuma validação. A defesa contra o MATCH silencioso hoje é só o `cmd/verificar`, depois da carga |
| 2 | Conversa × README × D-032 | Fases 0 a 4 concluídas | O README marca a Fase 4 como "Em andamento", e a D-032 (resultado da Fase 4) está vazia e como "Proposta" |
| 3 | Conversa × README | Existem as Fases 0, 8 e 9 | A tabela do README vai da Fase 1 à 7. Nenhum documento define uma Fase 0 |
| 4 | Conversa × D-032 | O PROFILE tem os experimentos A, B **e C** | A D-032 cita só os experimentos A e B, e um "Passo 7.3" que só existe na conversa. Nenhum dos três experimentos está descrito no repositório |
| 5 | `scripts/referencia.cypher` × `docs/modelo-grafo.md` × D-033 | D-033: rótulos ficam em português; D-018: o script é a fonte única | O script grava `nome` em inglês em 6 funções (`ark:gatherer` e `rust:gatherer` = Gatherer, `ark:explorer` = Explorer, `rust:electrician` = Electrician, `lol:support` = Support, `lol:fill` = Fill), e a cópia do script em `modelo-grafo.md` usa Coletor, Explorador, Eletricista, Suporte e Preencher. Se a Fase 6 ler os rótulos do grafo, o formulário sai em inglês |
| 6 | README × `go.mod` | README: "Go 1.22 ou superior" | `go.mod` declara `go 1.26.5`; versões anteriores precisam baixar o toolchain |
| 7 | README | A estrutura do projeto e os "Comandos úteis" | A estrutura omite `cmd/verificar` (citado só nos comandos). `go test -bench=. ./...` aparece como "roda os benchmarks", mas não há benchmarks |
| 8 | `docs/modelo-grafo.md` | Cabeçalho "Fase 1 — Modelagem / Status: Proposta" | O conteúdo já inclui material da Fase 4 (filtros, sinais, `personaJogo`) |
| 9 | `docs/personas.md` (gabarito) | Os filtros rígidos são "região, idioma, faixa etária, versão, modo, filas" | A lista omite a **plataforma** (grupo de crossplay do ARK), que `PassaFiltros` e o Cypher aplicam |
| 10 | Conversa × `decisoes.md` | D-028 e D-033 aguardam validação do orientador; D-002 aparece como "mudança posterior" | D-028 e D-033 estão como "Aceita". A D-002 está no bloco da Fase 1 (D-001 a D-010), não entre as mudanças posteriores |
| 11 | Conversa × modelo de `decisoes.md` | Formato: contexto, opções, decisão, motivo, consequências | O modelo do arquivo também tem Data e Status. Todas as entradas preenchem o Status e nenhuma preenche a Data |

**Afirmações não verificadas nesta geração** (exigem o banco): `cmd/verificar` termina com 0 divergências; os resultados de latência e de PROFILE da Fase 4.
