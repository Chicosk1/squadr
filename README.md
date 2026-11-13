# squadr

Motor de recomendação para formação de grupos em jogos online, no estilo "match". O projeto foca em jogos que não oferecem um bom sistema de criação de grupos. Desenvolvido como Trabalho de Conclusão de Curso.

O escopo é exclusivamente o **motor de recomendação**: os dados são mockados e o formulário de perfil é representativo. Integrações com Steam e Riot, chat e autenticação ficam como trabalhos futuros.

## Como funciona

1. **Neo4j (grafo)** guarda usuários, jogos, horários, funções e curtidas, e gera os candidatos a partir de filtros rígidos e travessias.
2. **Go** calcula o score de compatibilidade, aplica a reciprocidade (média harmônica entre o score de A para B e de B para A) e ordena o resultado final.

A recomendação é sempre contextualizada por jogo: o motor responde "quem é um bom parceiro para este jogo".

## Status

| Fase | Descrição | Situação |
|---|---|---|
| 0 | Planejamento (escopo e arquitetura, D-011 a D-015) | Concluída |
| 1 | Modelagem (questionário, grafo, personas) | Concluída |
| 2 | Infraestrutura e esqueleto do projeto | Concluída |
| 3 | Domínio e gerador de dados mock | Concluída |
| 4 | Geração de candidatos no grafo | Concluída |
| 5 | Scoring e reciprocidade | Concluída |
| 6 | API e formulário (tela de exploração e cadastro em passos, D-043 e D-044; fator de função, D-045) | Concluída |
| 7 | Avaliação e benchmarks | Pendente |
| 8 | Extras opcionais (Graph Data Science, ajuste de pesos) | Pendente |
| 9 | Fechamento (monografia e apresentação) | Pendente |

## Pré-requisitos

- [Go](https://go.dev/dl/) 1.26.5 ou superior
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (no Windows, com WSL2)
- [Git](https://git-scm.com/)

## Como rodar

Execute todos os comandos a partir da raiz do projeto.

**1. Clonar e configurar o ambiente**

```bash
git clone https://github.com/Chicosk1/squadr.git
cd squadr
cp .env.example .env
```

No Windows (PowerShell), use `Copy-Item .env.example .env`. Depois abra o `.env` e defina `NEO4J_PASSWORD` com pelo menos 8 caracteres.

**2. Subir o banco de dados**

```bash
docker compose up -d
docker compose ps
```

Aguarde o status `healthy` (de 30 a 60 segundos). O Neo4j Browser fica disponível em http://localhost:7474.

**3. Criar o schema e carregar os dados de referência**

```bash
go mod tidy
go run ./cmd/migrar
```

A saída deve mostrar quatro linhas `ok:` (jogos, horários, funções e complementaridade). O comando é idempotente: pode ser executado mais de uma vez.

**4. Iniciar a API**

```bash
go run ./cmd/api
```

**5. Verificar**

```bash
curl http://localhost:8080/saude
```

No PowerShell, use `curl.exe` no lugar de `curl`. A resposta esperada é `{"banco":"ok","status":"ok"}`.

Para ter usuários para recomendar, carregue os dados mock com `go run ./cmd/seed` (**apaga todos os usuários**, inclusive os cadastrados pelo formulário) antes de iniciar a API. A tela de exploração e o cadastro ficam em <http://localhost:8080/>.

## API

Detalhes do contrato na D-039 (`docs/decisoes.md`). Erros vêm como `{"erro": "...", "detalhes": [...]}`.

| Rota | Descrição |
|---|---|
| `GET /` | Tela de exploração (lista de jogadores, ranking e comparação de pares) e cadastro em passos (D-044) |
| `GET /saude` | Estado da API e do banco |
| `GET /usuarios?jogo=lol&regiao=AS&q=lobo&origem=seed&limite=50&depois=u000050` | Lista paginada de perfis públicos, com o total do filtro (D-043) |
| `GET /usuarios/{id}` | Perfil público do usuário (sem persona, pesos nem faixas aceitas, D-042) |
| `GET /usuarios/{id}/recomendacoes?jogo=lol&k=10` | Os `k` melhores parceiros (1 a 50, padrão 10), com score, motivo do match e os seis componentes do score nos dois sentidos |
| `GET /usuarios/{id}/par/{outro}?jogo=lol` | Explica um par: filtros rígidos com o motivo de cada um, scores nos dois sentidos, componentes e posição no ranking (D-043) |
| `POST /usuarios` | Cadastra um perfil (respostas do `docs/questionario.md`) e devolve o id gerado |
| `POST /usuarios/{id}/curtidas` | Curte `{"para": "u000002", "jogo": "lol"}`; devolve `{"match": true}` se a curtida for recíproca |
| `POST /usuarios/{id}/recusas` | Recusa, com o mesmo corpo |

```bash
curl "http://localhost:8080/usuarios/u000001/recomendacoes?jogo=lol&k=3"
curl "http://localhost:8080/usuarios/u000001/par/u000002?jogo=lol"
curl -X POST http://localhost:8080/usuarios/u000001/curtidas -d '{"para":"u000002","jogo":"lol"}'
```

Cadastros e curtidas feitos pela API mudam a base: depois deles, rode `go run ./cmd/seed` antes de `cmd/verificar` ou `cmd/avaliar`.

## Variáveis de ambiente

| Variável | Descrição | Padrão |
|---|---|---|
| `NEO4J_URI` | Endereço do Neo4j (protocolo Bolt) | `bolt://localhost:7687` |
| `NEO4J_USER` | Usuário do Neo4j | `neo4j` |
| `NEO4J_PASSWORD` | Senha do Neo4j (mínimo 8 caracteres) | — |
| `NEO4J_DATABASE` | Nome do banco | `neo4j` |
| `API_PORT` | Porta da API | `8080` |

## Estrutura do projeto

```
squadr/
├── cmd/
│   ├── api/          # servidor HTTP
│   ├── avaliar/      # Precision@K do motor contra o gabarito das personas
│   ├── migrar/       # cria constraints e carrega dados de referência
│   ├── seed/         # gerador de dados mock (apaga e recarrega os usuários)
│   └── verificar/    # compara o grafo com o oráculo em Go
├── internal/
│   ├── config/       # leitura de variáveis de ambiente
│   ├── graph/        # acesso ao Neo4j
│   ├── domain/       # structs do domínio
│   ├── matching/     # scoring e reciprocidade (Go puro, sem banco)
│   ├── seed/         # personas e gerador
│   └── api/          # handlers HTTP
├── scripts/          # scripts Cypher (schema e dados de referência)
├── docs/             # documentação e decisões de projeto
├── docker-compose.yml
└── .env.example
```

## Documentação

Toda a documentação de projeto está em [`docs/`](docs/):

- [`questionario.md`](docs/questionario.md): perguntas do formulário e o uso de cada uma no motor
- [`modelo-grafo.md`](docs/modelo-grafo.md): nós, relações, matrizes de complementaridade e índices
- [`personas.md`](docs/personas.md): perfis usados nos dados mock e gabarito de avaliação
- [`decisoes.md`](docs/decisoes.md): registro das decisões de projeto

## Comandos úteis

```bash
docker compose stop          # para o banco, mantendo os dados
docker compose down -v       # apaga o banco e o volume (recomeça do zero)
go test ./...                # roda os testes
go run ./cmd/verificar       # confere o grafo contra o oráculo (use os mesmos -usuarios e -semente do seed)
go run ./cmd/avaliar         # Precision@10 do motor e da ordem do grafo (mesmos -usuarios e -semente do seed)
```

## Tecnologias

- [Go](https://go.dev/)
- [Neo4j](https://neo4j.com/) 5.26 Community
- [Docker Compose](https://docs.docker.com/compose/)

## Licença

Projeto acadêmico.