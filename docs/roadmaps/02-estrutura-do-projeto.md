# Fase 2 — Estrutura do Projeto (monorepo)

**Status: 🟨 Em andamento** — estrutura de pastas criada em 29/07/2026; projetos (Go e Gradle) ainda não gerados

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [01 — Fundação: Ambiente e Ferramentas](./01-fundacao-ambiente-e-ferramentas.md)

## O que essa fase entrega

Uma organização de pastas fixa, que todo o resto do projeto vai seguir. Definir
isso **antes** de escrever telas e endpoints evita o problema mais comum em
projetos que crescem sem planejamento: arquivos jogados em qualquer lugar, sem
critério, até ficar difícil achar qualquer coisa.

> ⚠️ **Esta fase foi refeita em 29/07/2026.** A estrutura anterior era de um app
> Expo (`app/` com Expo Router + `src/`). Agora o repositório é um **monorepo**
> com backend Go e mobile Kotlin Multiplatform. As pastas antigas foram removidas
> — nenhum código foi perdido, porque todos aqueles arquivos estavam vazios. Ver
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md).

---

## 1. Estrutura completa do projeto

```
squadr/
├── backend/                          # SERVIÇO GO (API REST + WebSocket)
│   ├── cmd/
│   │   ├── api/                      # entrypoint da API REST
│   │   └── ws/                       # entrypoint do serviço WebSocket
│   ├── internal/                     # organizado por DOMÍNIO, não por camada
│   │   ├── auth/                     # validação de JWT do Supabase (JWKS), middleware
│   │   ├── user/                     # perfil, bloqueio, avaliação, sinais de sessão
│   │   ├── matching/                 # feed, compatibilidade, swipe, match
│   │   ├── chat/                     # mensagens + hub de WebSocket
│   │   ├── group/                    # squads: vagas, pedidos, aprovação
│   │   ├── notification/             # push via Firebase Admin SDK
│   │   ├── analytics/                # eventos de servidor (PostHog)
│   │   └── platform/                 # infraestrutura, sem regra de negócio
│   │       ├── config/  database/  logger/  httpserver/
│   ├── Dockerfile
│   ├── go.mod                        # ⬜ a gerar
│   └── README.md
│
├── mobile/                           # KOTLIN MULTIPLATFORM + COMPOSE MULTIPLATFORM
│   ├── androidApp/                   # casca fina: Activity, manifest, ícone
│   ├── iosApp/                       # casca fina: projeto Xcode
│   ├── shared/src/
│   │   ├── commonMain/kotlin/com/squadr/
│   │   │   ├── data/                 # remote/ local/ repository/
│   │   │   ├── domain/               # model/ repository/ usecase/
│   │   │   ├── ui/                   # login/ matching/ chat/ profile/ groups/ theme/
│   │   │   ├── di/                   # composição de dependências
│   │   │   └── platform/             # declarações expect
│   │   ├── androidMain/kotlin/com/squadr/platform/   # actual (Android)
│   │   ├── iosMain/kotlin/com/squadr/platform/       # actual (iOS)
│   │   └── commonMain/composeResources/              # imagens, fontes, strings
│   ├── settings.gradle.kts           # ⬜ a gerar
│   └── README.md
│
├── supabase/                         # BANCO DE DADOS
│   ├── migrations/                   # histórico de mudanças no schema
│   └── config.toml
│
├── contracts/                        # CONTRATO DA API (fonte única de verdade)
│   ├── openapi.yaml
│   └── README.md
│
├── docs/                             # DOCUMENTAÇÃO DO PROJETO
│   ├── context/                      # produto, stack, arquitetura, banco
│   ├── decisions/                    # ADRs
│   └── roadmaps/                     # este roadmap e as fases
│
├── .github/workflows/                # CI do backend
├── docker-compose.yml                # sobe api + ws localmente
├── .env.example                      # modelo das variáveis de ambiente
├── .env                              # variáveis reais (NUNCA vai pro Git)
├── .gitignore
├── AGENTS.md / CLAUDE.md             # instruções para agentes de IA
├── LICENSE
└── README.md
```

---

## 2. Por que cada pasta existe (o critério de organização)

### Raiz

| Pasta | Critério para algo entrar aqui |
|---|---|
| `backend/` | Código Go. Toda regra de negócio do produto vive aqui |
| `mobile/` | Código Kotlin. Interface e o que é específico de cliente |
| `supabase/` | Só migrations e config do Supabase CLI. Nada de código de aplicação |
| `contracts/` | A especificação da API. Muda **antes** do Go e do Kotlin, nunca depois |
| `docs/` | Documentação viva: contexto, decisões e roadmap |
| `.github/` | CI do backend (lint, teste, build). O CI do app fica no Codemagic |

### `backend/internal/` — por domínio, não por camada

O critério é: **uma pasta = uma capacidade do produto**. Pastas `handlers/`,
`services/` e `repositories/` no topo obrigariam a abrir três diretórios para
entender uma funcionalidade só.

| Pacote | Entra aqui | Não entra aqui |
|---|---|---|
| `auth/` | Quem é o usuário desta requisição: JWKS, verificação de assinatura/`iss`/`aud`/`exp`, extração do `sub`, middleware HTTP e do handshake WS | Regra de perfil (é `user/`) |
| `user/` | Perfil, disponibilidade, jogos do jogador, bloqueio, elogios, agregação dos sinais de sessão | Quem aparece no feed (é `matching/`) |
| `matching/` | Feed, compatibilidade de rank/horário, swipe, like mútuo, criação do match | Enviar o push do match (chama `notification/`) |
| `chat/` | Persistência de mensagem, histórico, hub de conexões, presença | Autenticar a conexão (chama `auth/`) |
| `group/` | Squads: vagas, pedidos, aprovação, `slots_filled` | — |
| `notification/` | Montar e disparar push (Firebase Admin SDK), guardar device tokens | Decidir *quando* notificar — quem decide é o domínio do evento |
| `analytics/` | Eventos de servidor no PostHog | Eventos de UI (saem do app) |
| `platform/*` | Só infraestrutura: config, pool `pgx`, logger, servidor HTTP | **Qualquer** regra de negócio |

### `mobile/shared/src/commonMain/` — por camada, dentro do módulo compartilhado

| Pasta | Critério para algo entrar aqui |
|---|---|
| `ui/` | Telas e componentes em Compose Multiplatform. Uma subpasta por área do produto |
| `ui/theme/` | Cores, tipografia e formas — o equivalente ao antigo `constants/theme.ts` |
| `domain/model/` | Tipos do domínio, sem dependência de framework |
| `domain/repository/` | **Interfaces** de acesso a dados |
| `domain/usecase/` | Regra que faz sentido no cliente: validação de formulário, ordenação local. Regra de negócio de verdade fica no Go |
| `data/remote/` | Cliente HTTP e WebSocket contra o backend Go. **Único** lugar que conhece URL de API |
| `data/local/` | Cache e preferências |
| `data/repository/` | Implementações das interfaces de `domain/repository/` |
| `di/` | Onde as dependências são montadas |
| `platform/` | Só declarações `expect`. As `actual` ficam em `androidMain/` e `iosMain/` |
| `composeResources/` | Imagens, fontes e strings compartilhadas |

**Duas regras que valem a pena repetir:**

1. **`androidApp/` e `iosApp/` são cascas finas.** Se você está escrevendo uma
   tela dentro de uma delas, pare — a tela deveria estar em `shared/.../ui/`. Nas
   cascas só entra o que o sistema operacional exige (manifest, `Info.plist`,
   ciclo de vida, permissões).
2. **`expect` só em `platform/`.** Espalhar `expect`/`actual` pelo código é o
   caminho mais rápido para um projeto Kotlin Multiplatform virar dois projetos.

---

## 3. Regra de nomenclatura (convenção do projeto)

- [ ] **Go — pacotes:** nome curto, minúsculo, sem underscore (`matching`, não `Matching` nem `matching_service`)
- [ ] **Go — arquivos:** `snake_case.go` (ex: `feed_query.go`, `jwks_cache.go`)
- [ ] **Go — testes:** `nome_do_arquivo_test.go`, ao lado do arquivo testado
- [ ] **Kotlin — arquivos e classes:** `PascalCase.kt` (ex: `PlayerCard.kt`, `MatchingScreen.kt`)
- [ ] **Kotlin — funções `@Composable`:** `PascalCase` (ex: `PlayerCard()`), como manda a convenção do Compose
- [ ] **Kotlin — pacotes:** minúsculo, sob `com.squadr`
- [ ] **Migrations do Supabase:** `NNN_descricao_curta.sql`, número sequencial de 3 dígitos (ex: `011_create_commendations_table.sql`)
- [ ] **ADRs:** `NNN-titulo-curto.md`, mesmo padrão de numeração
- [ ] **Endpoints REST:** plural, kebab-case quando composto (`/players`, `/squad-requests`) — definidos em `contracts/openapi.yaml`

---

## 4. O que já está feito

- [ ] Estrutura de pastas do monorepo criada (`backend/`, `mobile/`, `contracts/`, `.github/`)
- [ ] `docs/` reorganizada em `context/`, `decisions/` e `roadmaps/`
- [ ] Pastas e arquivos exclusivos do Expo removidos (`app/`, `src/`, `app.json`, `eas.json`, `package.json`, `package-lock.json`, `tsconfig.json`, `babel.config.js`)
- [ ] `docker-compose.yml`, `backend/Dockerfile`, `.env.example` e `.gitignore` criados para a stack nova
- [ ] `contracts/openapi.yaml` com o esqueleto da API
- [ ] Workflow de CI do backend em `.github/workflows/`
- [ ] `README.md` da raiz explicando como rodar o projeto

## 5. O que falta fazer nesta fase

### 5.1 Inicializar o módulo Go
- [ ] Com o Go instalado (Fase 1, item 1.4), dentro de `backend/`:
  ```bash
  go mod init github.com/Chicosk1/squadr/backend
  ```
- [ ] Criar um `main.go` mínimo em `cmd/api/` e outro em `cmd/ws/`, só para
      confirmar que a compilação funciona:
  ```bash
  go build ./...
  ```
- [ ] Remover os arquivos `.gitkeep` das pastas de `internal/` conforme cada uma
      receber seu primeiro arquivo `.go` de verdade
- [ ] Preencher `GO_VERSION` no `.env` com a versão instalada (usada pelo
      `docker-compose.yml`)

### 5.2 Gerar o projeto Kotlin Multiplatform
- [ ] Gerar o projeto pelo assistente (Android Studio → New Project → Kotlin
      Multiplatform, ou o web wizard da JetBrains), com:
  - package/namespace: `com.squadr`
  - UI compartilhada com **Compose Multiplatform** (não UI nativa por plataforma)
  - alvos: Android e iOS
- [ ] Gerar em uma pasta temporária e depois **encaixar** o resultado em
      `mobile/`, preservando a estrutura da seção 1 (o assistente cria
      `androidApp`/`iosApp`/`shared` e os arquivos Gradle — é isso que precisa
      vir; renomeie o que ele nomear diferente)
- [ ] Confirmar que compila:
  ```bash
  cd mobile && ./gradlew :shared:build
  ```
- [ ] Rodar o app Android no emulador, mesmo em branco
- [ ] Remover os `.gitkeep` das pastas que receberem arquivos reais

### 5.3 Escolher as bibliotecas (e registrar)
- [ ] Decidir cliente HTTP e injeção de dependência do mobile — candidatos em
      [`stack.md`](../context/stack.md), seção 3
- [ ] Decidir roteador HTTP e biblioteca de WebSocket do Go — idem
- [ ] Registrar cada escolha em `docs/decisions/` (ADR novo), não só no código

---

## Antes de avançar

- [ ] Estrutura de pastas do monorepo criada e commitada no Git
- [ ] Documentos organizados em `docs/context/`, `docs/decisions/` e `docs/roadmaps/`
- [ ] `go build ./...` roda sem erro em `backend/`
- [ ] `./gradlew :shared:build` roda sem erro em `mobile/`
- [ ] App Android abre no emulador (tela em branco é suficiente)
- [ ] Critério de "o que vai em cada pasta" (seção 2) entendido — isso evita
      decisões diferentes por quem trabalhar no projeto depois

➡️ Próxima fase: [`03-backend-banco-e-autenticacao.md`](./03-backend-banco-e-autenticacao.md)
