# squadr

TCC: motor de recomendação no estilo "match" para encontrar parceiros em jogos online (ARK, Rust e League of Legends). O escopo é só o motor, com dados mock e um formulário representativo (D-011). Chat, autenticação e integração com Steam e Riot são trabalho futuro.

O aluno é o autor do TCC e está aprendendo com o projeto. Responda em português, seja crítico (aponte o que não combina, inclusive no plano) e priorize eficiência.

## Arquitetura

- O Neo4j 5.26 gera os candidatos: filtros rígidos, horários em comum, complementaridade de funções e sinal social (`internal/graph/candidatos.go`: duas consultas, corte em 300).
- O Go puro calcula o score unilateral, a reciprocidade (média harmônica, D-015) e o ranking (`internal/matching`).
- A API REST usa só `net/http` (`internal/api`, `cmd/api`).
- A recomendação é sempre por jogo (D-002): toda consulta, curtida e recusa carrega o jogo.

Detalhes:

- `docs/modelo-grafo.md`: grafo, filtros por jogo e sinais.
- `docs/questionario.md`: perguntas do formulário e o uso de cada uma.
- `docs/personas.md`: personas, ruído e gabarito.
- `docs/decisoes.md`: todas as decisões.
- `docs/contexto-claude.md`: histórico do planejamento e plano das fases seguintes. É um retrato de 09/11/2026; se divergir do código, vale o código.

## Comandos

Ambiente: Windows com Git Bash e Go 1.26.5 (`go.mod`). Rode tudo a partir da raiz.

```bash
docker compose up -d && docker compose ps   # esperar "healthy"
go run ./cmd/migrar                          # schema + dados de referência (idempotente)
go run ./cmd/seed                            # APAGA todos os Usuario e recarrega (-usuarios 10000, -semente 42)
go run ./cmd/verificar                       # compara o grafo com o oráculo (mesmos -usuarios e -semente do seed)
go run ./cmd/api                             # GET /saude
go build ./... && go vet ./... && go test ./...
```

- `cmd/seed` apaga a base. Avise antes de rodar.
- `cmd/verificar` falha se `-usuarios` e `-semente` não forem os mesmos usados no seed.
- Os comandos de `cmd/` exigem o Neo4j de pé.

## Regras invioláveis

- **Não leia nem copie o `.env`** (tem a senha do banco). Use `.env.example` e `internal/config`.
- **Não faça commits nem push.** O aluno faz os commits. Leitura do git (`status`, `log`, `diff`, `stash list`) é permitida.
- Módulo `github.com/Chicosk1/squadr`. Dependências permitidas: `godotenv` e `neo4j-go-driver/v5`. Nada novo sem decisão registrada.
- API só com `net/http` e o roteamento do Go 1.22+ (`"GET /rota/{id}"`). Sem chi, gin, echo ou fiber (D-014).
- `persona` e `personaJogo` existem só para o gabarito da avaliação. **O motor nunca lê esses campos**, mesmo vindo preenchidos em `domain.Perfil` por `graph.BuscarPerfis`.
- `internal/matching` não importa `internal/graph` nem o driver do Neo4j. Ele depende de uma interface de fonte de candidatos (D-013).
- Nas consultas Cypher, nada vindo do usuário entra no texto da query, só parâmetros. Os filtros por jogo são fragmentos fixos (D-030).
- Um único `graph.Cliente` compartilhado, nunca um por requisição (D-019).
- `scripts/*.cypher` são a fonte única do grafo de referência e usam `MERGE` (idempotentes, D-018).
- `MaxDiffEloDuo` vive em `internal/domain` e é a mesma no oráculo e no Cypher: passe sempre `domain.MaxDiffEloDuo` em `Opcoes.MaxDiffElo`.
- O gerador é reproduzível: nunca itere sobre `map` em código que consome números aleatórios. Personas são dados, não código (D-021, D-022). `TestGeracaoDeterministica` precisa continuar passando.
- Ids de função em inglês (`ark:gatherer`, `rust:electrician`, `lol:support`...), conforme a D-033. Rótulos do formulário, nomes de persona, identificadores e comentários do código ficam em português, como no código existente.

## Defesas que existem por causa de bugs reais

- **MATCH silencioso.** Na carga Cypher, um `MATCH` que não encontra o nó faz a relação não ser criada, sem nenhum erro. Já aconteceu duas vezes. Ao mexer em carga ou em ids, o seed deve validar contra o grafo antes de apagar a base, e o `cmd/verificar` deve conferir depois.
- **Verificação diferencial.** O oráculo `seed.PassaFiltros` e o Cypher implementam as mesmas regras de forma independente (D-024, D-031). Ao mudar uma regra, mude as duas implementações e mantenha-as independentes: uma nunca chama a outra.

## Decisões

Toda decisão nova vai para `docs/decisoes.md`, no formato do arquivo: título `D-XXX`, Data, Status, Contexto, Opções consideradas, Decisão, Motivo e Consequências. Use o próximo número livre. Decisões que dependem do orientador ficam como "Proposta".

## Ao terminar uma fase ou tarefa grande

Feche com um resumo curto contendo:

- objetivo;
- mudanças em relação ao plano;
- arquivos alterados;
- comandos para conferir;
- **O que você aprende** (conceitos úteis para a monografia);
- **Critério de pronto**;
- **Armadilhas comuns**.

Depois atualize a tabela de status do README.