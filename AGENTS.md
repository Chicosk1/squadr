# Squadr — instruções para agentes de IA

App de matchmaking para gamers brasileiros. **Monorepo**: mobile em Kotlin
Multiplatform, backend em Go.

## ⚠️ A stack mudou em 29/07/2026

O projeto **era** React Native + Expo com Supabase como backend completo. Não é
mais. Se você encontrar qualquer referência a Expo, React Native, `app.json`,
`eas.json`, Expo Router ou SDK JS do Supabase fora de
`docs/decisions/002-stack-mobile-react-native-expo.md`, é resíduo — corrija.

| Camada | Tecnologia |
|---|---|
| Mobile | Kotlin Multiplatform + Compose Multiplatform (`mobile/`) |
| Backend | Go — `cmd/api` (REST) e `cmd/ws` (WebSocket), em `backend/` |
| Banco | Supabase (Postgres gerenciado + Storage), acessado via `pgx` |
| Auth | Discord OAuth no Supabase Auth; **JWT validado no Go** via JWKS |
| Push | Firebase Cloud Messaging (KMPNotifier no app, Admin SDK no backend) |
| Analytics | PostHog (SDK Kotlin Multiplatform + SDK Go) |
| Deploy | Fly.io (backend) · Codemagic (app) |

## Antes de escrever qualquer código

1. **Leia o contexto:** [`docs/context/arquitetura.md`](docs/context/arquitetura.md)
   define onde cada coisa mora e por quê. [`docs/context/stack.md`](docs/context/stack.md)
   lista o que **ainda não foi decidido** — não invente decisão que está aberta lá.
2. **Consulte a documentação oficial versionada** da tecnologia que for tocar.
   Kotlin Multiplatform, Compose Multiplatform e as bibliotecas do ecossistema
   mudam de API com frequência; não escreva de memória:
   - Kotlin Multiplatform: https://kotlinlang.org/docs/multiplatform.html
   - Compose Multiplatform: https://www.jetbrains.com/help/kotlin-multiplatform-dev/
   - Go: https://go.dev/doc/ · pacotes: https://pkg.go.dev
   - Supabase: https://supabase.com/docs
   - PostHog: https://posthog.com/docs
   - Fly.io: https://fly.io/docs · Codemagic: https://docs.codemagic.io
3. **Confirme a versão instalada** antes de assumir uma (`go version`,
   `java -version`, o Gradle do projeto). Não invente número de versão.

## Regras invioláveis

1. **O app nunca fala direto com o banco.** Toda leitura e escrita passa pelo Go.
   Nenhuma credencial de banco no app.
2. **Contrato primeiro:** altere `contracts/openapi.yaml` → depois o Go → depois o
   Kotlin. Nunca a ordem inversa.
3. **Migration já aplicada nunca é editada.** Mudança de schema é arquivo novo, com
   RLS na mesma migration que cria a tabela. Ver
   [`docs/decisions/001-migrations-e-rls-supabase.md`](docs/decisions/001-migrations-e-rls-supabase.md)
   — esse ADR nasceu de um incidente real neste projeto.
4. **`backend/internal/` é organizado por domínio**, não por camada. Não crie
   `handlers/`, `services/` ou `repositories/` no topo.
5. **`internal/platform/*` não contém regra de negócio.**
6. **Toda tela vive em `mobile/shared/src/commonMain/.../ui`** em Compose
   Multiplatform. `androidApp/` e `iosApp/` são cascas finas.
7. **`expect`/`actual` só sob `platform/`.**
8. **Decisão técnica relevante vira ADR** em `docs/decisions/`, numerado em
   sequência. Decisão antiga é marcada como supersedida, **nunca apagada**.
9. **Autorização é responsabilidade do Go** e precisa de teste explícito por
   endpoint — o RLS não é mais a principal linha de defesa.

## Estado atual do projeto

A estrutura de pastas existe, mas **os projetos ainda não foram gerados**: falta
`go mod init` em `backend/` e o projeto Gradle em `mobile/` (Fase 2 do roadmap).
Não presuma que há código compilável. Roadmap:
[`docs/roadmaps/`](docs/roadmaps/).
