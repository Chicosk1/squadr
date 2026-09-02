# 009 — Levar variáveis do `.env` para o código Kotlin com BuildKonfig

**Status:** Aceito
**Data:** 05/08/2026

## Contexto

Primeiro pré-requisito da seção 7.3 do roadmap: `SUPABASE_URL` e
`SUPABASE_PUBLISHABLE_KEY` (mais `API_BASE_URL`/`WS_BASE_URL`, que já existiam
no `.env` desde a Fase 2) precisam chegar ao código Kotlin em `commonMain`
para montar o `SupabaseClient` do [ADR-008](./008-biblioteca-supabase-auth-kmp.md).

O backend Go lê essas variáveis em tempo de execução (`os.Getenv`), mas o
mobile é compilado — os valores precisam ser embutidos no binário durante o
build. O mecanismo nativo do Android (`buildConfigField` do plugin
`com.android.kotlin.multiplatform.library`) só gera a classe `BuildConfig`
para o alvo Android — não alcança `commonMain`, então não alcança iOS.

## Decisão

**`com.codingfeline.buildkonfig` v0.22.0** — plugin Gradle dedicado a
exatamente esse problema em Kotlin Multiplatform, gerando um objeto
`internal object BuildKonfig` em `commonMain`, acessível dos dois alvos sem
`expect`/`actual` manual. Compilado com Kotlin 2.4.10 — a mesma versão do
projeto, sem risco de desalinhamento.

A alternativa seria uma `Task` Gradle escrita à mão, gerando um `.kt` num
diretório de fontes gerado e registrando esse diretório em `commonMain`. Não
escolhi esse caminho porque, no momento desta decisão, o `./gradlew` por
linha de comando está bloqueado neste ambiente (erro de socket loopback,
investigado e não resolvido — não é bug do projeto). Sem conseguir compilar
para validar, o risco de acertar a sintaxe de wiring de source-set na mão era
alto demais; uma biblioteca já testada pela comunidade reduz esse risco.

`shared/build.gradle.kts` lê o `.env` da raiz do repositório (não o
`.env.example`) na fase de configuração do Gradle, valida que todas as
variáveis necessárias estão presentes — falhando com mensagem clara se
faltar alguma, mesmo princípio de fail-fast do `internal/platform/config` no
backend — e alimenta o `buildkonfig { defaultConfigs { ... } }` com os
valores lidos.

`SUPABASE_URL` **não foi duplicado** no `.env`: o mobile reaproveita a mesma
variável que o backend já usa (não é segredo, é a URL pública do projeto).

## Consequências

- `mobile/gradle/libs.versions.toml`: nova entrada de versão e plugin
  (`buildkonfig`).
- `mobile/build.gradle.kts`: plugin registrado com `apply false` (mesmo
  padrão dos outros).
- `mobile/shared/build.gradle.kts`: lê o `.env`, valida presença de
  `API_BASE_URL`, `WS_BASE_URL`, `SUPABASE_URL`, `SUPABASE_PUBLISHABLE_KEY`, e
  configura o `buildkonfig`.
- Qualquer `./gradlew` (mesmo `tasks` ou `help`) falha se `.env` não existir
  ou faltar uma dessas variáveis — intencional, mas quem for rodar CI de
  mobile no futuro (ainda não existe — hoje só o backend tem CI, ver
  [`02-estrutura-do-projeto.md`](../roadmaps/02-estrutura-do-projeto.md)) vai
  precisar prover um `.env` no pipeline.
- Os valores gerados (`BuildKonfig.API_BASE_URL` etc.) hoje são sempre os de
  desenvolvimento local (`.env` tem só um valor por variável, ex:
  `http://10.0.2.2:8080`). Diferenciar debug/release (apontar pra
  `https://squadr-api.fly.dev` em produção) é decisão de build variant ainda
  não tomada — fica para a Fase 6 (deploy), não bloqueia o login agora.
- Não validado com `./gradlew :shared:build` (bloqueio de ambiente já
  registrado) — validado com sucesso via sincronização no Android Studio
  (05/08/2026), que usa um JDK diferente do `gradlew` de linha de comando.
