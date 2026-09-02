# 010 — Esquema de deep link para o retorno do OAuth

**Status:** Aceito
**Data:** 05/08/2026

## Contexto

O [ADR-008](./008-biblioteca-supabase-auth-kmp.md) decidiu usar `auth-kt`, que
trata o retorno do fluxo OAuth via deep link (`handleDeeplinks()`), em vez de
código próprio de plataforma. Falta decidir **qual** URL de retorno o app usa
— isso precisa ser o mesmo valor em quatro lugares: painel do Supabase,
`AndroidManifest.xml`, `Info.plist` do `iosApp` e a configuração do
`SupabaseClient` no código Kotlin (essa última entra num passo futuro, quando
o cliente for montado com Koin).

Consultei o guia oficial de deep linking do Supabase (aba "Android Kotlin")
antes de decidir o formato — não presumido: a URL de retorno segue
`[SCHEME]://[HOSTNAME]`, e o `scheme` precisa ser único no dispositivo do
usuário. A convenção recomendada é usar o domínio reverso do app.

## Decisão

**`com.squadr://login-callback`** — reaproveita o `applicationId` do
`androidApp` (`com.squadr`), que já é único no dispositivo por definir a
identidade do app na Play Store. Não inventei um domínio reverso novo porque
o app já tem um identificador exclusivo pronto.

Configurado em:

- `mobile/androidApp/src/main/AndroidManifest.xml` — `intent-filter` novo na
  `MainActivity`, com `android:scheme="com.squadr"` e
  `android:host="login-callback"`, além do já existente (que só trata o
  ícone do launcher).
- `mobile/iosApp/iosApp/Info.plist` — `CFBundleURLTypes` com
  `CFBundleURLSchemes: ["com.squadr"]`.

**Ação manual pendente, fora do código:** cadastrar
`com.squadr://login-callback` no painel do Supabase, em
**Authentication → URL Configuration → Redirect URLs**. Sem isso, o Supabase
recusa o redirect por não estar na allow-list — mesmo mecanismo de segurança
que vimos ao testar o fluxo manualmente pelo navegador (seção 7.1 do
roadmap).

## Consequências

- Os quatro lugares (painel, manifest, plist, código do cliente) têm que
  ficar sincronizados manualmente. Se o esquema mudar no futuro, precisa
  mudar nos quatro.
- `android:exported="true"` já estava na `MainActivity` (necessário pro
  launcher); o novo `intent-filter` de deep link não exige mudança nesse
  atributo.
- Falta ainda: chamar `supabase.handleDeeplinks(intent)` no `MainActivity`
  (Android) e a configuração equivalente no lado Kotlin/iOS — isso entra
  junto com a montagem do `SupabaseClient` (próximo passo da seção 7.3).
