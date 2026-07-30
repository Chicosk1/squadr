# backend — serviço Go do Squadr

Dois binários, um conjunto de pacotes de domínio.

| Entrypoint | O que é | Porta padrão |
|---|---|---|
| `cmd/api` | API REST: regras de negócio (perfil, matching, bloqueio, avaliação, grupos) | 8080 |
| `cmd/ws` | Serviço WebSocket: chat em tempo real, hub de conexões | 8081 |

Os dois importam os **mesmos** pacotes de `internal/`. Nenhuma regra de negócio é
duplicada. Serviços separados para que um deploy da API não derrube as conversas
abertas, e um pico de conexões no chat não degrade o CRUD.

## Estrutura

```
backend/
├── cmd/
│   ├── api/          # entrypoint da API REST
│   └── ws/           # entrypoint do WebSocket
└── internal/         # organizado por DOMÍNIO, não por camada técnica
    ├── auth/         # JWT do Supabase via JWKS, middleware, auth do handshake
    ├── user/         # perfil, bloqueio, avaliação, sinais de sessão
    ├── matching/     # feed, compatibilidade, swipe, criação de match
    ├── chat/         # mensagens + hub de conexões
    ├── group/        # squads: vagas, pedidos, aprovação
    ├── notification/ # push via Firebase Admin SDK
    ├── analytics/    # eventos de servidor no PostHog
    └── platform/     # infraestrutura, sem regra de negócio
        ├── config/     # leitura e validação de variáveis de ambiente
        ├── database/   # pool pgx, transações
        ├── logger/     # log estruturado
        └── httpserver/ # servidor HTTP, graceful shutdown, middlewares
```

### Por que domínio e não camada

Pastas `handlers/`, `services/` e `repositories/` no topo obrigam a abrir três
diretórios para entender uma funcionalidade, e nada impede que qualquer coisa
importe qualquer coisa. Por domínio, **uma pasta = uma capacidade do produto**.
Dentro de cada pacote a divisão por camada volta a existir, mas em arquivos.

Tabela completa de "o que entra e o que não entra em cada pacote":
[`docs/context/arquitetura.md`](../docs/context/arquitetura.md), seção 3.

## Rodando

> ⚠️ O módulo Go ainda não foi inicializado (`go mod init` — Fase 2 do roadmap).
> Até isso acontecer, os comandos abaixo falham.

```bash
go run ./cmd/api
```

```bash
go run ./cmd/ws
```

```bash
go test ./...
```

Variáveis de ambiente: ver [`.env.example`](../.env.example) na raiz.

## Regras do backend

1. **`internal/platform/*` não contém regra de negócio.** Se apareceu um `if` de
   produto ali, ele está no lugar errado.
2. **Quem decide *quando* notificar é o domínio do evento**; `notification` só
   sabe *como* enviar.
3. **Autorização é responsabilidade daqui.** O RLS do Postgres não é mais a
   principal linha de defesa — ver
   [`docs/context/banco-de-dados.md`](../docs/context/banco-de-dados.md), seção 3.
4. **Operação atômica usa transação explícita** (like mútuo criando match,
   aprovar pedido de squad incrementando `slots_filled`).
5. **O contrato vem antes do código:** altere
   [`contracts/openapi.yaml`](../contracts/openapi.yaml) primeiro.

## Deploy

Fly.io, dois apps (`squadr-api` e `squadr-ws`), a partir deste mesmo `Dockerfile`
com `target` diferente. ⚠️ O app do WebSocket fica fixo em **uma instância** — o
hub é em memória. Ver
[`docs/roadmaps/06-deploy-e-lancamento.md`](../docs/roadmaps/06-deploy-e-lancamento.md),
seção 2.
