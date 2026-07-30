# Arquitetura

> **Vigente desde 29/07/2026** — ver [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md).
> Antes disso não havia backend próprio: o app falava direto com o Supabase.

---

## 1. Visão geral

```
        ┌─────────────────────────────────────────┐
        │   mobile/  Kotlin Multiplatform         │
        │   UI única em Compose Multiplatform     │
        │         (androidApp + iosApp)           │
        └───────┬──────────────────────┬──────────┘
                │ HTTPS REST           │ WSS
                │ Bearer JWT           │ Bearer JWT no handshake
                ▼                      ▼
        ┌───────────────┐      ┌───────────────────┐
        │  cmd/api (Go) │      │   cmd/ws (Go)     │
        │  regra de     │      │   chat: hub de    │
        │  negócio      │      │   conexões        │
        └───────┬───────┘      └───────────┬───────┘
                │ ambos importam os mesmos │
                │ pacotes de internal/     │
                └──────────────┬───────────┘
                               │ pgx
                               ▼
              ┌────────────────────────────────┐
              │  Supabase                      │
              │  · Postgres gerenciado         │
              │  · Storage (avatares)          │
              │  · Auth (Discord OAuth, JWKS)  │
              └────────────────────────────────┘

     Saídas do backend:  Firebase Admin SDK → FCM → dispositivos
                         SDK Go do PostHog  → PostHog (mesmo projeto do app)
```

**A regra que organiza tudo:** o app **nunca** fala com o banco. Toda leitura e
escrita passa pelo Go. Isso é a mudança central em relação à arquitetura
anterior, e é o que permite garantir que regra de matching, bloqueio e reputação
não sejam contornáveis pelo cliente.

---

## 2. Monorepo

```
squadr/
├── backend/            # serviço Go (API REST + WebSocket)
├── mobile/             # Kotlin Multiplatform + Compose Multiplatform
├── supabase/           # migrations e config do Supabase CLI
├── contracts/          # especificação OpenAPI (fonte única de verdade)
├── docs/               # context, decisions, roadmaps
├── .github/            # workflows de CI do backend
├── docker-compose.yml  # sobe api + ws localmente
└── README.md
```

Por que monorepo: o contrato da API muda em dois lugares ao mesmo tempo
(handler Go + cliente Kotlin). Num único repositório isso é **um** commit e
**um** PR, com o `contracts/openapi.yaml` no meio como árbitro.

---

## 3. `backend/` — organizado por domínio, não por camada

```
backend/
├── cmd/
│   ├── api/                 # entrypoint da API REST
│   └── ws/                  # entrypoint do serviço WebSocket
├── internal/
│   ├── auth/                # validação de JWT do Supabase via JWKS, middleware
│   ├── user/                # perfil, bloqueio, avaliação/reputação, sinais de sessão
│   ├── matching/            # feed, compatibilidade, swipe, criação de match
│   ├── chat/                # mensagens + hub de conexões WebSocket
│   ├── group/               # squads: criação, vagas, pedidos de entrada
│   ├── notification/        # disparo de push via Firebase Admin SDK
│   ├── analytics/           # eventos de servidor via SDK Go do PostHog
│   └── platform/            # infraestrutura, sem regra de negócio
│       ├── config/          # leitura e validação de variáveis de ambiente
│       ├── database/        # pool de conexões pgx, transações
│       ├── logger/          # log estruturado
│       └── httpserver/      # servidor HTTP, graceful shutdown, middlewares comuns
├── Dockerfile
└── README.md
```

### O critério: por que domínio e não por camada

Pastas `handlers/`, `services/`, `repositories/` no topo obrigam a abrir três
diretórios para entender uma funcionalidade, e nada impede que qualquer coisa
importe qualquer coisa. Por domínio, **uma pasta = uma capacidade do produto**,
com sua própria fronteira. Dentro de cada pacote de domínio a divisão por camada
volta a existir, mas em arquivos, não em pastas espalhadas.

| Pacote | Entra aqui | Não entra aqui |
|---|---|---|
| `internal/auth` | Tudo sobre *quem é o usuário desta requisição*: buscar e cachear JWKS, verificar assinatura/`iss`/`aud`/`exp`, extrair o `sub`, middleware HTTP e handshake do WS | Regra de perfil (isso é `user`) |
| `internal/user` | Perfil, disponibilidade, jogos do jogador, bloqueio, elogios, agregação dos sinais de sessão | Decidir *quem aparece no feed* (isso é `matching`) |
| `internal/matching` | Query do feed, compatibilidade de rank/horário, swipe, detecção de like mútuo, criação do match | Enviar a notificação do match (chama `notification`) |
| `internal/chat` | Persistência de mensagem, histórico, o **hub** de conexões WebSocket, presença | Autenticar a conexão (chama `auth`) |
| `internal/group` | Squads/lobby: vagas, pedidos, aprovação, `slots_filled` | — |
| `internal/notification` | Montar e disparar push via Firebase Admin SDK, guardar device tokens | Decidir *quando* notificar — quem decide é o domínio que gerou o evento |
| `internal/analytics` | Envio de eventos de servidor ao PostHog | Eventos de UI (esses saem do app) |
| `internal/platform/*` | Só infraestrutura: config, pool `pgx`, logger, servidor HTTP | **Qualquer** regra de negócio |

### Os dois entrypoints

`cmd/api` e `cmd/ws` são binários diferentes que importam os **mesmos** pacotes
de `internal/`. Nenhuma regra de negócio é duplicada: se o `ws` precisa saber se
um usuário bloqueou outro antes de entregar a mensagem, ele chama
`internal/user` — a mesma função que a API REST chama.

Consequência prática: um deploy da API não derruba as conversas abertas, e um
pico de conexões no chat não degrada o CRUD.

---

## 4. `mobile/` — Kotlin Multiplatform

```
mobile/
├── androidApp/              # casca fina: Activity, manifest, ícone
├── iosApp/                  # casca fina: projeto Xcode, AppDelegate
└── shared/src/
    ├── commonMain/kotlin/com/squadr/
    │   ├── data/
    │   │   ├── remote/          # cliente HTTP e WebSocket contra o backend Go
    │   │   ├── local/           # cache/preferências
    │   │   └── repository/      # implementações das interfaces de domain
    │   ├── domain/
    │   │   ├── model/           # tipos do domínio, sem dependência de framework
    │   │   ├── repository/      # interfaces
    │   │   └── usecase/         # regra que faz sentido no cliente (validação de form, ordenação local)
    │   ├── ui/                  # telas em Compose Multiplatform
    │   │   ├── login/  matching/  chat/  profile/  groups/  theme/
    │   ├── di/                  # composição de dependências
    │   └── platform/            # declarações expect
    ├── androidMain/kotlin/com/squadr/platform/    # implementações actual (Android)
    ├── iosMain/kotlin/com/squadr/platform/        # implementações actual (iOS)
    └── commonMain/composeResources/               # imagens, fontes, strings
```

**`androidApp/` e `iosApp/` são cascas finas por decisão explícita.** Toda tela
vive em `shared/.../ui` em Compose Multiplatform. Se uma tela começar a ser
escrita dentro de `androidApp/`, é sinal de que algo saiu do trilho — a única
coisa que deve morar nas cascas é o que o sistema operacional exige (manifest,
`Info.plist`, ciclo de vida, permissões).

**`platform/` é onde o `expect`/`actual` mora.** Nada de `expect` espalhado pelo
resto do código. Candidatos naturais: armazenamento seguro do token
(EncryptedSharedPreferences no Android, Keychain no iOS), inicialização do
KMPNotifier, abertura do navegador para o OAuth, `Dispatchers` específicos.

**Desenvolvimento primário no Android Studio.** O Xcode entra só para o que é
inevitável: assinar, rodar em device físico iOS e publicar. Isso está assumido
no roadmap — quem não tem macOS consegue desenvolver tudo, menos gerar o build
iOS (que fica com o Codemagic).

---

## 5. Fluxos

### 5.1 Autenticação

```
1. App: usuário toca "Entrar com Discord"
2. App abre o fluxo OAuth do Supabase Auth no navegador do sistema
   (Custom Tabs no Android, ASWebAuthenticationSession no iOS)
3. Discord autentica → Supabase Auth devolve access_token (JWT) + refresh_token
4. App guarda os tokens no armazenamento seguro da plataforma (platform/, expect/actual)
5. Toda chamada ao backend leva:  Authorization: Bearer <access_token>
6. Go (internal/auth):
   · busca o JWKS do Supabase e mantém em cache (com rotação de chave)
   · confere assinatura, iss, aud e exp
   · extrai o sub → é o ID do usuário para o resto da requisição
7. Token expirado → o app renova direto com o Supabase Auth (refresh_token)
```

O **Go nunca emite token** — só valida. Quem emite é o Supabase Auth. Isso
mantém o login em 2 cliques e evita construir gestão de sessão do zero.

> **Ponto aberto:** como o app fala com o Supabase Auth — via biblioteca de
> comunidade de Supabase para Kotlin Multiplatform (confirmar o suporte a iOS na
> versão corrente) ou chamando os endpoints REST de auth direto com o cliente
> HTTP. Decidir na Fase 3.

### 5.2 Requisição REST comum

```
App → HTTPS → cmd/api
                 │ middleware de internal/auth valida o Bearer
                 │ handler do domínio (ex: matching.Feed)
                 │ pgx → Postgres do Supabase
                 └─ resposta JSON conforme contracts/openapi.yaml
```

### 5.3 Chat em tempo real

```
1. App abre WSS contra cmd/ws, com o JWT no header do handshake
2. internal/auth valida o token antes do upgrade — conexão não autenticada não sobe
3. Conexão entra no hub (internal/chat): user_id → conexões ativas
4. Mensagem recebida:
   · internal/user confirma que não há bloqueio entre as partes
   · internal/chat persiste em messages (pgx)
   · destinatário conectado → entrega pelo hub
   · destinatário offline → internal/notification dispara push via FCM
```

> **Limitação aceita no MVP:** o hub é **em memória**, o que significa **uma
> instância** de `cmd/ws`. Duas instâncias sem fan-out fariam a mensagem chegar
> só a quem estivesse conectado na mesma instância. Antes de subir a segunda
> instância é obrigatório escolher um mecanismo de distribuição (`LISTEN/NOTIFY`
> do Postgres ou Redis pub/sub) e registrar como ADR. Não é dívida escondida —
> é escopo consciente.

### 5.4 Push notification

Quem decide *que* algo merece push é o domínio que gerou o evento
(`matching` criou um match, `chat` recebeu mensagem para alguém offline). Quem
sabe *como* enviar é `internal/notification`, via Firebase Admin SDK. No app, a
recepção e a exibição ficam com o KMPNotifier, com o token do dispositivo
registrado no backend no login.

### 5.5 Analytics

Dois SDKs, **um único projeto PostHog**: o SDK de Kotlin Multiplatform manda
eventos de interface (tela vista, like dado) e o SDK de Go manda eventos de
servidor (match criado, push enviado, erro de regra). Como o `distinct_id` é o
mesmo (o `sub` do JWT), um funil pode cruzar as duas origens.

---

## 6. `contracts/` — o contrato é a fonte de verdade

`contracts/openapi.yaml` descreve a API REST. A ordem de trabalho é sempre:
**alterar o contrato → alterar o Go → alterar o Kotlin**. Nunca o inverso.

Isso existe porque cliente e servidor agora são linguagens diferentes,
compiladas separadamente: sem contrato explícito, a divergência só aparece em
runtime, no dispositivo do usuário.

---

## 7. Topologia de deploy

| O que | Onde | Observação |
|---|---|---|
| `cmd/api` | Fly.io, app próprio | Região preferencial GRU (São Paulo) |
| `cmd/ws` | Fly.io, app próprio, separado | Escala e reinicia independente da API |
| Postgres + Storage | Supabase | Migrations aplicadas via Supabase CLI — ver [ADR-001](../decisions/001-migrations-e-rls-supabase.md) |
| Push | Firebase Cloud Messaging | Credencial de service account só no backend, nunca no app |
| Build e publicação do app | Codemagic | Assina e envia para Google Play e App Store |
| CI do backend | GitHub Actions | `.github/workflows/` — lint, teste e build a cada PR |

Localmente, `docker-compose.yml` sobe `api` e `ws`; o Postgres vem do
`supabase start` (stack local do Supabase CLI). Ver o
[README](../../README.md) da raiz.
