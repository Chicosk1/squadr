# Banco de Dados

> **Vigente desde 29/07/2026.** O banco **não mudou** na migração de stack — o
> Supabase continua sendo o Postgres do projeto. O que mudou é **quem** fala com
> ele: antes o app, via SDK JS; agora só o backend Go, via `pgx`.
> Ver [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md).

---

## 1. Tabelas principais

```sql
profiles        (id, discord_id, username, avatar, bio, rank, profile_complete, is_online, last_seen_at)
games           (id, name, slug, platform)
player_games    (player_id, game_id, rank, style)
availability    (player_id, weekday, period)
swipes          (swiper_id, swiped_id, direction, created_at)
matches         (id, player_a, player_b, created_at, invite_copied_at)
messages        (id, match_id, sender_id, content, created_at)
reports         (id, reporter_id, reported_id, reason, created_at)
squads          (id, creator_id, game_id, slots_total, slots_filled, min_rank, created_at)
squad_members   (squad_id, player_id, joined_at)
squad_requests  (squad_id, player_id, status, created_at)
commendations   (from_id, to_id, match_id, tag, created_at)
blocks          (blocker_id, blocked_id, created_at)
game_identities (player_id, game_id, identifier)
device_tokens   (player_id, token, platform, updated_at)
invites         (code, created_by, used_by, created_at, used_at)
```

Duas anotações sobre essa lista, para não passarem batidas:

- **`device_tokens` é nova e existe por causa da troca de stack.** Com Expo
  Notifications, o registro do dispositivo era responsabilidade do Expo. Com FCM
  + backend próprio, o token de cada dispositivo tem que ser guardado por nós —
  é o que `internal/notification` usa para disparar o push.
- **`is_online` / `last_seen_at` estão explicitados em `profiles`.** O sistema de
  sinais de sessão (ver [`produto.md`](./produto.md), seção 8.2) sempre dependeu
  desses campos, mas eles não apareciam na listagem do contexto antigo. Agora
  aparecem.
- **`invites`** cobre a estratégia de crescimento por convite da fase inicial
  (`produto.md`, seção 8.1). Formato final a definir na Fase 3.

> Nenhuma migration foi escrita ainda — `supabase/migrations/` está vazia. Esta
> lista é o modelo **conceitual**; o schema real nasce na Fase 3 do roadmap.

---

## 2. Como o Go acessa

- **Driver:** `pgx`, com pool de conexões em `internal/platform/database`. Sem
  ORM e sem SDK do Supabase no caminho.
- **Transações:** operações que tocam mais de uma tabela e precisam ser atômicas
  (aceitar pedido de squad = inserir em `squad_members` + incrementar
  `slots_filled`; like mútuo = inserir `swipes` + criar `matches`) rodam em
  transação explícita. Esse tipo de garantia era difícil de obter com o app
  falando direto com o banco — é um dos ganhos concretos da arquitetura nova.
- **Conexão:** o Supabase oferece conexão direta (5432) e o pooler (Supavisor).
  ⚠️ Se for usado o pooler em **modo transação**, o cache de prepared statements
  do `pgx` precisa ser desligado, senão aparecem erros intermitentes de
  "prepared statement already exists". Conexão direta ou pooler em **modo
  sessão** não têm esse problema. Decidir e registrar na Fase 3.
- **Credencial:** a connection string vive só no ambiente do backend
  (`DATABASE_URL`). **Nunca** no app.

---

## 3. Onde a autorização mora agora (ponto de atenção)

Na stack antiga, o RLS do Postgres **era** o mecanismo de autorização: o app
tinha a `anon key` e o banco decidia, linha por linha, o que aquele usuário podia
ver. Daí a importância do [ADR-001](../decisions/001-migrations-e-rls-supabase.md).

Com o backend Go no meio, isso muda de lugar: **a autorização passa a ser
responsabilidade do Go**, que valida o JWT e aplica a regra antes de qualquer
query. O que não muda:

| Continua valendo | Muda |
|---|---|
| Migrations aplicadas nunca são editadas (ADR-001) | O RLS deixa de ser a **única** linha de defesa |
| Migrations versionadas em `supabase/migrations/`, aplicadas via Supabase CLI | O app deixa de usar a Data API (PostgREST) para ler e escrever dados |
| RLS continua relevante para o **Storage** e para qualquer acesso que não passe pelo Go | Se o Go conectar com um papel privilegiado, o RLS é **ignorado** nessas queries |

⚠️ **Cuidado com uma conclusão fácil e errada:** "o app não fala mais com o
banco, então o RLS não importa". O app **continua tendo a `anon key`**, porque
precisa dela para falar com o Supabase Auth (login com Discord). Se a Data API
(PostgREST) do projeto continuar expondo o schema `public`, essa mesma chave
permite consultar tabelas direto, **sem passar pelo Go** — e aí toda a
autorização escrita em Go é irrelevante. As duas saídas aceitáveis são
desabilitar a exposição do schema na Data API ou manter o RLS negando tudo por
padrão para `anon` e `authenticated`.

> ⚠️ **Decisão obrigatória antes de criar as tabelas (Fase 3):** com qual papel
> do Postgres o Go conecta, e se o RLS é mantido como defesa em profundidade
> (conectando com papel sujeito a RLS e propagando as claims do JWT) ou tratado
> apenas como proteção do Storage e de acessos diretos. As duas saídas são
> defensáveis; **decidir por omissão é o que não pode acontecer** — foi
> exatamente esse tipo de descuido que gerou o ADR-001. A decisão vira um ADR
> próprio.

---

## 4. Storage

Avatares e imagens continuam no Supabase Storage. Como o app não tem mais
credencial do Supabase para escrever direto, o caminho passa a ser: o app pede
ao Go, o Go gera uma URL assinada de upload e devolve. Assim o backend continua
sendo o único ponto que decide quem pode escrever o quê.

---

## 5. Migrations

Processo e regras: [ADR-001](../decisions/001-migrations-e-rls-supabase.md) —
**integralmente válido**, é a única decisão anterior que a troca de stack não
afetou. Em resumo:

1. Migration já aplicada **nunca** é editada; mudança vira arquivo novo.
2. RLS habilitado na **mesma** migration que cria a tabela.
3. `DROP POLICY IF EXISTS` antes de `CREATE POLICY` (idempotência).
4. Confirmar com `select tablename, rowsecurity from pg_tables where schemaname = 'public';`
   antes de commitar.
