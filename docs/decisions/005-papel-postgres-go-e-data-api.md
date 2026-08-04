# 005 — Papel do Postgres para o backend Go e exposição da Data API

**Status:** Aceito
**Data:** 04/08/2026

## Contexto

Decisão bloqueante da Fase 3 do roadmap (seção 4), levantada pela emenda de
29/07/2026 no [ADR-001](./001-migrations-e-rls-supabase.md). Com a troca de
stack, a autorização deixou de ser feita por RLS no Postgres e passou a ser
responsabilidade do backend Go, que valida o JWT do Supabase Auth e aplica as
regras de negócio antes de qualquer query. Duas coisas precisavam ser
decididas antes de criar a primeira tabela (seção 8 do roadmap):

1. Com qual papel do Postgres o Go se conecta ao banco.
2. O que fazer com a Data API (PostgREST) do Supabase, que continua acessível
   pelo app através da `publishable key`, independentemente do backend Go.

Contexto adicional relevante para a decisão: quem mantém o backend está
aprendendo Go durante o desenvolvimento do projeto — complexidade evitável na
camada de acesso a dados tem custo maior que o normal neste momento.

## Decisão

### 1. Papel do Postgres

O Go conecta com um **papel privilegiado**, que ignora RLS. A autorização é
100% responsabilidade do código Go, validada por endpoint — sem propagação de
claims do JWT para o Postgres.

Motivo: propagar claims do JWT (defesa em profundidade via RLS) exigiria
implementar e manter autorização em dois lugares — Go e Postgres —
simultaneamente, o que não compensa para quem está aprendendo Go agora. O
risco que essa camada extra mitigaria (um bug de autorização no Go passar
direto) já é coberto pela exigência de testes explícitos por endpoint com os
três atores — sem token, token de outro usuário, token do dono — definida na
seção 3 do roadmap de qualidade
([`05-qualidade-e-testes.md`](../roadmaps/05-qualidade-e-testes.md)).

RLS continua sendo habilitado em toda tabela (regra herdada do ADR-001), mas
nesse cenário sua função é proteger o Storage e qualquer acesso que não passe
pelo Go — não servir de defesa em profundidade para as próprias queries do Go.

**Gatilhos para reavaliar esta decisão:**

- Antes de abrir cadastro público (hoje o acesso é só por convite, seção 9 do
  roadmap) — mais superfície de ataque justifica mais camadas.
- Se outro serviço, além de `cmd/api`/`cmd/ws`, passar a acessar o mesmo banco.
- Se o produto passar a armazenar dado sensível de categoria diferente da
  atual (ex: dado de pagamento).

### 2. Data API do Supabase

A exposição do schema `public` na Data API (PostgREST) é **desabilitada**, em
Project Settings → Data API.

Motivo: manter a Data API exposta — mesmo com RLS negando tudo por padrão para
`anon`/`authenticated` — deixa um caminho residual de app → banco direto, que
só permanece seguro se toda policy futura continuar restritiva. Isso viola o
espírito da regra 1 do [AGENTS.md](../../AGENTS.md) ("o app nunca fala direto
com o banco") e repete a categoria de risco que originou o ADR-001: proteção
que depende de disciplina contínua em vez de ser estrutural. Desabilitar a
exposição do schema remove essa classe de risco por completo, sem afetar o
Supabase Auth (API separada) — o login com Discord continua funcionando
normalmente.

## Consequências

- Toda autorização do MVP é escrita e testada em Go; não existe política RLS
  de autorização a manter em paralelo (RLS cobre só Storage e acessos fora do
  Go).
- A Data API do Supabase não pode ser usada por nenhum motivo enquanto esta
  decisão estiver em vigor — qualquer necessidade futura de usá-la (ex:
  Supabase Realtime) exige reabrir esta decisão.
- O teste "jeito errado" da seção 4 do roadmap (tentar ler uma tabela direto
  pela Data API com a `publishable key`) deve falhar porque a Data API está
  desabilitada — não por causa de RLS. Isso deve ser confirmado
  explicitamente, não assumido.
- Esta decisão não substitui o ADR-001: RLS continua sendo habilitado na
  mesma migration que cria cada tabela, protegendo Storage e qualquer futuro
  acesso fora do Go.
