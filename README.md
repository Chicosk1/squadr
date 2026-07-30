# Squadr

> O app que conecta gamers brasileiros para jogar juntos — sem depender de grupo
> de WhatsApp, sem spam de Discord e sem esperar o amigo estar online.

Monorepo do Squadr: app mobile em Kotlin Multiplatform e backend em Go.

## Stack

| Camada | Tecnologia |
|---|---|
| Mobile | Kotlin Multiplatform + Compose Multiplatform (Android e iOS) |
| Backend | Go — API REST (`cmd/api`) + serviço WebSocket (`cmd/ws`) |
| Banco / Storage | Supabase (Postgres gerenciado), acessado via `pgx` |
| Autenticação | Discord OAuth via Supabase Auth; JWT validado no Go (JWKS) |
| Push | Firebase Cloud Messaging (KMPNotifier no app, Admin SDK no backend) |
| Analytics | PostHog (SDK Kotlin Multiplatform + SDK Go, mesmo projeto) |
| Hospedagem | Fly.io (backend) · Codemagic (build e publicação do app) |

Detalhes e justificativas: [`docs/context/stack.md`](docs/context/stack.md).

> ℹ️ O projeto migrou de React Native + Expo em 29/07/2026. Ver
> [ADR-003](docs/decisions/003-migracao-kmp-e-backend-go.md).

## Estrutura

```
squadr/
├── backend/            # serviço Go (API REST + WebSocket), organizado por domínio
├── mobile/             # Kotlin Multiplatform + Compose Multiplatform
├── supabase/           # migrations e config do Supabase CLI
├── contracts/          # especificação OpenAPI da API — fonte única de verdade
├── docs/               # context, decisions, roadmaps
├── .github/workflows/  # CI do backend
└── docker-compose.yml  # sobe api + ws localmente
```

## Documentação

Comece por aqui — a documentação é a fonte de verdade do projeto:

| Onde | O que tem |
|---|---|
| [`docs/context/`](docs/context/) | Produto, stack, arquitetura e banco de dados |
| [`docs/decisions/`](docs/decisions/) | ADRs — toda decisão técnica, incluindo as supersedidas |
| [`docs/roadmaps/`](docs/roadmaps/) | Passo a passo de execução, fase por fase |
| [`contracts/`](contracts/) | Contrato da API REST |

## Rodando localmente

> ⚠️ **Estado atual:** a estrutura de pastas existe, mas os projetos ainda não
> foram gerados — falta o `go mod init` do backend e o projeto Gradle do mobile
> (Fase 2 do roadmap). Os comandos abaixo passam a funcionar depois disso.

### Pré-requisitos

Go, JDK, Android Studio (com o plugin Kotlin Multiplatform), Docker e Supabase
CLI. Lista completa com links e verificação:
[`docs/roadmaps/01-fundacao-ambiente-e-ferramentas.md`](docs/roadmaps/01-fundacao-ambiente-e-ferramentas.md).

### 1. Variáveis de ambiente

```bash
cp .env.example .env
```

Preencha o `.env` com os valores do seu projeto Supabase. Ele **nunca** vai para
o Git.

### 2. Banco de dados

```bash
supabase start
```

```bash
supabase db push
```

⚠️ Antes de criar qualquer migration, leia o
[ADR-001](docs/decisions/001-migrations-e-rls-supabase.md) — migration já
aplicada **nunca** é editada.

### 3. Backend

```bash
cd backend && go run ./cmd/api
```

```bash
cd backend && go run ./cmd/ws
```

Ou, em contêiner (os dois de uma vez):

```bash
docker compose up --build
```

### 4. Mobile

Abra a pasta `mobile/` no Android Studio e rode a configuração `androidApp` no
emulador. Para iOS é necessário macOS com Xcode — sem isso, o build iOS sai pelo
Codemagic (Fase 6).

## Testes

```bash
cd backend && go test ./...
```

```bash
cd mobile && ./gradlew :shared:allTests
```

## Convenções que não se negociam

1. **O app nunca fala direto com o banco.** Toda leitura e escrita passa pelo Go.
2. **Contrato primeiro:** alterar `contracts/openapi.yaml` → depois o Go → depois
   o Kotlin.
3. **Migration aplicada não se edita** — mudança é arquivo novo
   ([ADR-001](docs/decisions/001-migrations-e-rls-supabase.md)).
4. **Decisão técnica relevante vira ADR** em
   [`docs/decisions/`](docs/decisions/), numerado em sequência.
5. **Backend organizado por domínio**, não por camada técnica.

## Licença

Ver [LICENSE](LICENSE).
