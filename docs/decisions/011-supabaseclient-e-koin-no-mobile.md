# 011 — Primeira instalação do Koin: `SupabaseClient` como singleton

**Status:** Aceito
**Data:** 05/08/2026

## Contexto

O [ADR-004](./004-bibliotecas-do-mobile-e-backend.md) escolheu Koin para
injeção de dependência no mobile, mas — como o Ktor no mesmo ADR — isso era
só decisão no papel: nenhum código do projeto usava Koin até agora. O
[ADR-008](./008-biblioteca-supabase-auth-kmp.md) e o
[ADR-010](./010-esquema-deep-link-oauth.md) já definiram a biblioteca de auth
e o esquema de deep link; faltava montar o `SupabaseClient` de verdade e
decidir onde ele vive no grafo de dependências do app.

Antes de escrever qualquer import, confirmei direto no código-fonte do
`supabase-kt` (não presumido, para não escrever import errado num ambiente
onde não consigo compilar para checar — ver bloqueio de `./gradlew`
registrado no ADR-009):

- `createSupabaseClient` e a interface `SupabaseClient` vivem no pacote
  `io.github.jan.supabase`.
- `Auth` (o plugin) vive em `io.github.jan.supabase.auth`, com `scheme` e
  `host` como propriedades de `AuthConfig`.
- `handleDeeplinks(intent: Intent, ...)` é uma **extensão Android-only** de
  `SupabaseClient`, definida só no `androidMain` do módulo `auth-kt`.
- A inicialização do `Context` do Android é automática, via
  `androidx.startup` (o `auth-kt` já registra um `Initializer` próprio no seu
  `AndroidManifest.xml`, mesclado no do app) — não precisa de `Application`
  customizada nem de passar `Context` manualmente.

## Decisão

Um módulo Koin (`appModule`, em `shared/src/commonMain/.../di/AppModule.kt`)
expõe **um único singleton**: o `SupabaseClient`, já com o plugin `Auth`
instalado (`scheme`/`host` do ADR-010). `initKoin()` (mesmo arquivo, `Koin.kt`)
chama `startKoin { modules(appModule) }`, com um guard (`koinStarted`) contra
chamada duplicada — necessário porque não há garantia de quantas vezes o
ponto de entrada do iOS (`iOSApp.init`) executa.

**`androidApp` não depende de Koin nem de `auth-kt` diretamente.** Em vez de
elevar essas dependências para `api` no `shared` (o que exporia essas
bibliotecas para qualquer consumidor do módulo), criei um wrapper Android-only
em `shared/src/androidMain/.../di/AuthDeeplinks.android.kt`:
`handleAuthDeeplink(intent: Intent)`. `MainActivity` só importa essa função —
um `Intent` (tipo do próprio Android) entra, nada sai. Isso mantém a
superfície pública do `shared` pequena.

`MainActivity` ganhou `android:launchMode="singleTask"`: sem isso, se a
Activity já estiver na memória quando o navegador devolver o controle após o
login, o Android poderia não invocar `onCreate` de novo — o retorno do OAuth
se perderia. Com `singleTask`, o retorno chega em `onNewIntent`, que também
foi implementado.

## Consequências

- Todo código de domínio que precisar do `SupabaseClient` no futuro (ex: pra
  ler a sessão atual) injeta via Koin, nunca instancia um novo — é o único
  lugar que sabe a URL/key do projeto.
- O lado iOS (`KoinKt.initKoin()` em `iOSApp.swift`) segue a mesma convenção
  de nomenclatura já usada em `MainViewControllerKt.MainViewController()`
  (`ContentView.swift`), mas **não foi compilado nem testado** — esta máquina
  não tem Xcode. Precisa de confirmação na primeira vez que alguém builda o
  `iosApp` de verdade.
- `koin-compose`/`koin-compose-viewmodel` (que dão `koinInject()`/
  `koinViewModel()`) ainda **não foram adicionados** — só entram quando a
  tela de login (próximo passo da seção 7.3) precisar deles. Adicionar agora
  seria dependência sem consumidor.
