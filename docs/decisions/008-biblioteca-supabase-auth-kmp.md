# 008 — Biblioteca para o app falar com o Supabase Auth (Kotlin Multiplatform)

**Status:** Aceito
**Data:** 05/08/2026

## Contexto

Pendência registrada desde a Fase 2/[ADR-003](./003-migracao-kmp-e-backend-go.md) e
formalizada como item 1 da seção 7.3 do roadmap
([`03-backend-banco-e-autenticacao.md`](../roadmaps/03-backend-banco-e-autenticacao.md)):
como o app (Kotlin Multiplatform) fala com o Supabase Auth para o login "Entrar
com Discord" — biblioteca de comunidade específica para Supabase, ou os
endpoints REST de auth chamados direto pelo cliente HTTP (Ktor, já escolhido no
[ADR-004](./004-bibliotecas-do-mobile-e-backend.md)).

A dúvida que travava essa decisão era suporte a iOS — [`arquitetura.md`](../context/arquitetura.md)
§5.1 registrava isso como "ponto aberto". Antes de decidir, verifiquei
diretamente no repositório `supabase-community/supabase-kt` (não presumido):

- Versão atual: **3.7.0**.
- Módulo `auth-kt`: iOS (`iosArm64`, `iosSimulatorArm64`, `iosX64`) está listado
  como **totalmente suportado**, com OAuth nativo incluso — ao contrário de
  Windows/Linux, que não têm OAuth embutido.
- `supabase.auth.signInWith(Discord)` abre o navegador do sistema por conta
  própria e trata o retorno via `handleDeeplinks()` nas plataformas mobile —
  não é preciso escrever `expect`/`actual` de Custom Tabs / `ASWebAuthenticationSession`.
- `autoLoadFromStorage` e `alwaysAutoRefresh` são `true` por padrão — sessão
  persistida e renovação de token vêm do plugin, não de código próprio.
- O módulo é construído sobre **Ktor** internamente (confirmado no
  `build.gradle.kts` do projeto) — mesmo cliente do ADR-004, sem stack de rede
  duplicada.
- Requisito encontrado durante a verificação: **`auth-kt` exige Android API 26+**
  (ou core library desugaring abaixo disso). O projeto estava em `minSdk = 24`.

## Decisão

**`io.github.jan-tennert.supabase:auth-kt` v3.7.0.**

Motivo: a alternativa (Ktor direto) exigiria escrever à mão, em Kotlin/Native
para iOS — a parte da stack onde há menos experiência acumulada no projeto —
exatamente as três coisas que os itens 2, 3 e 4 da seção 7.3 do roadmap listam
como trabalho manual: abrir o navegador do sistema por plataforma, persistir
token em storage seguro nativo, e implementar renovação. O `auth-kt` resolve
essas três com comportamento padrão, e está construído sobre a mesma base
(Ktor) já adotada no ADR-004 — não introduz uma segunda forma de fazer HTTP no
app.

Módulos do supabase-kt **não** instalados: `postgrest-kt`, `storage-kt`,
`realtime-kt` — o [ADR-005](./005-papel-postgres-go-e-data-api.md) já decidiu
que o app não acessa o banco nem a Data API do Supabase diretamente. Só o
plugin `Auth` entra no projeto.

**`minSdk` sobe de 24 para 26** (Android 8.0), para evitar a configuração de
desugaring. Android 7.x tem participação de mercado residual no Brasil em
2026; o custo de suportá-lo não compensa a complexidade extra.

## Consequências

- `mobile/gradle/libs.versions.toml`: `android-minSdk` passa de `"24"` para
  `"26"`; adiciona a versão do `supabase-kt` e do Ktor (exigido como
  dependência direta do app, não só transitiva do `auth-kt`) e os engines por
  plataforma (`ktor-client-okhttp` no Android, `ktor-client-darwin` no iOS).
- `mobile/shared/build.gradle.kts`: `commonMain` ganha `auth-kt`; `androidMain`
  e a source set de iOS ganham o engine Ktor correspondente.
- O deep link de retorno do OAuth (esquema de URL) ainda precisa ser
  registrado no `AndroidManifest.xml` e no `Info.plist` do `iosApp` — é
  configuração declarativa, não código de plataforma, mas não é automático.
- Exato mecanismo de storage por trás de `autoLoadFromStorage` (se é
  Keychain/EncryptedSharedPreferences de fato, ou algo menos seguro por
  padrão) ainda precisa ser confirmado ao implementar o login de verdade —
  a documentação consultada não detalhou isso. Registrar como nota ao
  implementar a seção 7.3, itens 2–4.
