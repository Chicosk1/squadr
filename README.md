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
| 1 | Modelagem (questionário, grafo, personas) | Concluída |
| 2 | Infraestrutura e esqueleto do projeto | Em andamento |
| 3 | Domínio e gerador de dados mock | Pendente |
| 4 | Geração de candidatos no grafo | Pendente |
| 5 | Scoring e reciprocidade | Pendente |
| 6 | API e formulário | Pendente |
| 7 | Avaliação e benchmarks | Pendente |

## Pré-requisitos

- [Go](https://go.dev/dl/) 1.22 ou superior
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
│   ├── migrar/       # cria constraints e carrega dados de referência
│   └── seed/         # gerador de dados mock
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
go test -bench=. ./...       # roda os benchmarks
```

## Tecnologias

- [Go](https://go.dev/)
- [Neo4j](https://neo4j.com/) 5.26 Community
- [Docker Compose](https://docs.docker.com/compose/)

## Licença

Projeto acadêmico.