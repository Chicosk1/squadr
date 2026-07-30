# Fase 3 — Backend: Banco de Dados, API e Autenticação

**Status: 🟨 Em andamento** — projeto Supabase criado e vinculado; schema, API Go e autenticação pendentes

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [02 — Estrutura do Projeto](./02-estrutura-do-projeto.md)

## O que essa fase entrega

O banco de dados criado no Supabase, o backend Go conectado a ele, o contrato da
API definido, e o login "Entrar com Discord" funcionando de ponta a ponta — com o
**JWT sendo validado pelo Go**, não mais consumido direto pelo app.

> ⚠️ **Leia isto antes de continuar:** este projeto já teve um problema real
> aqui — documentado no [ADR-001](../decisions/001-migrations-e-rls-supabase.md).
> As tabelas foram criadas e só depois vieram as políticas de segurança (RLS),
> editando os mesmos arquivos de migration que já tinham sido aplicados. Isso não
> funcionou: o Supabase CLI não reaplica uma migration só porque o conteúdo do
> arquivo mudou — ele rastreia pelo nome/timestamp do arquivo. Resultado: 8
> tabelas ficaram sem proteção por um tempo, mesmo com o código "certo" já
> escrito. A regra da seção 3 abaixo existe por causa disso, e **continua valendo
> integralmente** na stack nova.

> 🔄 **O que a troca de stack mudou nesta fase:** o app não fala mais com o banco.
> Tudo passa pelo Go. Isso acrescenta as seções 5 (conexão `pgx`), 6 (contrato) e
> 7 (validação de JWT no Go), e transforma o RLS de "mecanismo de autorização" em
> "defesa em profundidade" — o que exige uma decisão explícita (seção 4).

---

## 1. Projeto no Supabase — ✅ já feito

- [ ] Acessar https://supabase.com/dashboard e criar um novo projeto
- [ ] Escolher uma região próxima do Brasil (ex: São Paulo)
- [ ] Guardar a senha do banco gerada
- [ ] No painel do projeto, ir em **Project Settings → API** e copiar `Project URL` e `anon public key`
- [ ] ⚠️ **Refazer o `.env` com os nomes novos das variáveis.** As variáveis
      antigas tinham prefixo `EXPO_PUBLIC_`, que não existe mais. Ver
      [`.env.example`](../../.env.example) para a lista atual (`DATABASE_URL`,
      `SUPABASE_URL`, `SUPABASE_JWKS_URL`, etc.)
- [ ] Copiar também a **connection string** do banco (Project Settings →
      Database) — é o que o Go usa, e ela **nunca** vai para o app

---

## 2. Supabase CLI vinculado ao projeto remoto — ✅ já feito

- [ ] `supabase login`
- [ ] `supabase link --project-ref SEU_PROJECT_REF`
- [ ] Confirmar quais migrations já foram aplicadas **antes de criar qualquer migration nova**:
  ```bash
  supabase migration list
  ```
  Se uma migration aparece como `Applied`, editar o arquivo dela **não terá
  efeito** — é exatamente o problema do ADR-001.

---

## 3. Regra de migrations (obrigatória — vem do ADR-001)

- [ ] **Nunca edito uma migration que já rodou no remoto.** Qualquer mudança de
      schema — incluindo "esqueci de proteger essa tabela" — é um arquivo novo,
      com número seguinte.
- [ ] **RLS é ativado na mesma migration que cria a tabela**, sempre.
- [ ] Toda `CREATE POLICY` é precedida de `DROP POLICY IF EXISTS`, para a
      migration poder ser reaplicada sem erro em outro ambiente.
- [ ] Depois de rodar `supabase db push`, confirmar visualmente:
  1. No Table Editor, a tabela não deve aparecer como `UNRESTRICTED`
  2. Rodar esta query no SQL Editor como segunda confirmação:
     ```sql
     select tablename, rowsecurity from pg_tables where schemaname = 'public';
     ```
     Todas as linhas devem ter `rowsecurity = true`.
- [ ] Só depois de confirmado, fazer o commit da migration no Git

---

## 4. 🚧 Decisão bloqueante: papel do banco e estratégia de RLS

**Isto tem que ser resolvido antes de criar a primeira tabela.** Não é
formalidade: é a repetição do erro do ADR-001 esperando acontecer, agora em outra
camada.

O contexto: com o Go no meio, a autorização é feita em Go. Mas o app **continua
tendo a `anon key`** do Supabase, porque precisa dela para falar com o Supabase
Auth (login com Discord). Se a Data API (PostgREST) do projeto estiver exposta,
essa chave também permite consultar tabelas direto, sem passar pelo Go.

- [ ] Decidir **com qual papel do Postgres o Go conecta**:
  - papel privilegiado (RLS ignorado; autorização 100% no Go), ou
  - papel sujeito a RLS, propagando as claims do JWT (defesa em profundidade)
- [ ] Decidir o que fazer com a **Data API do Supabase**:
  - desabilitar a exposição do schema `public` na Data API, ou
  - mantê-la exposta **com RLS negando tudo por padrão** para `anon` e
    `authenticated`
- [ ] Registrar a decisão como **ADR novo** em `docs/decisions/`
- [ ] Testar a decisão do jeito errado de propósito: com a `anon key` do app na
      mão, tentar ler uma tabela direto pela Data API. O resultado esperado é
      falha. Se vier dado, a decisão não está implementada

> Referência: [`banco-de-dados.md`](../context/banco-de-dados.md), seção 3, e a
> emenda no topo do [ADR-001](../decisions/001-migrations-e-rls-supabase.md).

---

## 5. Conectando o Go ao Postgres

- [ ] Adicionar o `pgx` ao módulo:
  ```bash
  cd backend && go get github.com/jackc/pgx/v5
  ```
- [ ] Implementar `internal/platform/config` — leitura e **validação** das
      variáveis de ambiente (falhar no start se faltar alguma, em vez de descobrir
      em produção)
- [ ] Implementar `internal/platform/database` — pool de conexões e helper de
      transação
- [ ] ⚠️ Decidir entre **conexão direta** (porta 5432) e **pooler** (Supavisor).
      Se for pooler em **modo transação**, o cache de prepared statements do `pgx`
      precisa ser desligado, senão aparecem erros intermitentes de "prepared
      statement already exists". Registrar a escolha
- [ ] Implementar `internal/platform/logger` (log estruturado) e
      `internal/platform/httpserver` (servidor HTTP com graceful shutdown)
- [ ] Confirmar a conexão com um endpoint de saúde (`GET /healthz`) que faça um
      `select 1` no banco

---

## 6. Contrato da API primeiro

A ordem de trabalho no projeto é **contrato → Go → Kotlin**. Sem isso, cliente e
servidor divergem e o erro só aparece no dispositivo do usuário.

- [ ] Preencher `contracts/openapi.yaml` com os endpoints do MVP (o esqueleto já
      está lá): perfil, feed, swipe, matches, mensagens, squads, report/bloqueio,
      device token
- [ ] Definir o formato padrão de **erro** (mesmo corpo em todos os endpoints) —
      decidir uma vez aqui evita cinco formatos diferentes depois
- [ ] Definir paginação do feed e do histórico de mensagens
- [ ] Validar que o arquivo é um OpenAPI válido antes de commitar

---

## 7. Autenticação — Discord OAuth + validação de JWT no Go

### 7.1 No Discord e no Supabase (igual à stack antiga)
- [ ] No [Discord Developer Portal](https://discord.com/developers/applications), confirmar a aplicação criada (ex: "Squadr")
- [ ] Em **OAuth2 → Redirect**, adicionar a URL de callback que o Supabase fornece (Painel Supabase → Authentication → Providers → Discord)
- [ ] Copiar `Client ID` e `Client Secret` da aplicação Discord
- [ ] No painel do Supabase, em **Authentication → Providers → Discord**, ativar e colar `Client ID`/`Client Secret`
- [ ] Definir os **escopos** solicitados: `identify`, `email` e `guilds` (o escopo
      `guilds` é o que permite, com consentimento do usuário, mostrar o selo de
      vínculo com comunidade — parte do sistema de confiança do produto)

### 7.2 No backend Go (novo)
- [ ] Implementar `internal/auth`:
  - [ ] buscar o **JWKS** do Supabase (`<SUPABASE_URL>/auth/v1/.well-known/jwks.json`)
  - [ ] manter em **cache**, com renovação — mas sem refazer a busca a cada
        requisição, e sem quebrar quando a chave rotacionar (`kid` desconhecido →
        recarrega uma vez)
  - [ ] validar assinatura, `iss`, `aud` e `exp`
  - [ ] extrair o `sub` e colocá-lo no contexto da requisição
- [ ] Middleware HTTP que rejeita requisição sem `Authorization: Bearer` válido
      (401), aplicado a **todas** as rotas exceto `/healthz`
- [ ] Mesma validação no **handshake do WebSocket**, antes do upgrade — conexão
      não autenticada não deve subir
- [ ] Preferir o token no **header** do handshake, não em query string (query
      string vaza token em log de servidor e de proxy)
- [ ] Testar com token inválido, expirado e assinado por outra chave — os três
      devem falhar

### 7.3 No app (Kotlin)
- [ ] Decidir como o app fala com o Supabase Auth: biblioteca de comunidade de
      Supabase para Kotlin Multiplatform (**confirmar o suporte a iOS na versão
      corrente**) ou chamar os endpoints REST de auth direto pelo cliente HTTP.
      Registrar como ADR
- [ ] Abrir o fluxo OAuth no navegador do sistema (Custom Tabs no Android,
      `ASWebAuthenticationSession` no iOS) — implementação `actual` em
      `shared/src/androidMain` e `iosMain`, sob `platform/`
- [ ] Guardar `access_token` e `refresh_token` no armazenamento seguro da
      plataforma (EncryptedSharedPreferences / Keychain), também via
      `expect`/`actual`
- [ ] Renovar o token expirado direto com o Supabase Auth (o Go não emite token,
      só valida)
- [ ] Testar o login de ponta a ponta com o app rodando no emulador contra o
      backend local

---

## 8. Schema do banco (tabelas do MVP)

Para cada tabela: criar a migration já **com RLS habilitado desde o início**
(seção 3) e conforme a decisão da seção 4 — nunca como um passo separado depois.

- [ ] `profiles` (id, discord_id, username, avatar, bio, rank, profile_complete, is_online, last_seen_at)
- [ ] `games` (id, name, slug, platform)
- [ ] `player_games` (player_id, game_id, rank, style)
- [ ] `availability` (player_id, weekday, period)
- [ ] `swipes` (swiper_id, swiped_id, direction, created_at)
- [ ] `matches` (id, player_a, player_b, created_at, invite_copied_at)
- [ ] `messages` (id, match_id, sender_id, content, created_at)
- [ ] `reports` (id, reporter_id, reported_id, reason, created_at)
- [ ] `squads` (id, creator_id, game_id, slots_total, slots_filled, min_rank, created_at)
- [ ] `squad_members` (squad_id, player_id, joined_at)
- [ ] `squad_requests` (squad_id, player_id, status, created_at)
- [ ] `commendations` (from_id, to_id, match_id, tag, created_at)
- [ ] `blocks` (blocker_id, blocked_id, created_at)
- [ ] `game_identities` (player_id, game_id, identifier)
- [ ] `device_tokens` (player_id, token, platform, updated_at) — **nova**, exigida pelo FCM
- [ ] `invites` (code, created_by, used_by, created_at, used_at)

### Checklist de autorização — pensar por tabela

Antes, essas perguntas eram respondidas por políticas RLS. Agora elas são
respondidas **no Go** (e, se a seção 4 decidir manter RLS, nos dois lugares).
Registre a resposta junto ao código do domínio, não só na cabeça:

- [ ] Quem pode **ler** essa linha? (o próprio usuário, os dois lados de um match, todo mundo, etc.)
- [ ] Quem pode **criar** uma linha nova?
- [ ] Quem pode **atualizar**? Existe algum campo que nem o dono deveria poder alterar (ex: `created_at`, `slots_filled`)?
- [ ] Quem pode **apagar**, se alguém puder?
- [ ] Existe alguma leitura que precisa **excluir usuários bloqueados**? (feed, chat, squads — todos precisam)

---

## 9. Estratégia de crescimento por convite (fase inicial)

Conforme decidido no documento de produto ([`produto.md`](../context/produto.md),
seção 8.1), o cadastro no lançamento **não é aberto** — só entra quem foi
convidado por alguém já na plataforma.

- [ ] Implementar o mecanismo de convite **no Go** (tabela `invites` + validação
      no fluxo de primeiro login). Como a validação agora é server-side, ela é de
      fato inescapável — antes dependia do app cooperar
- [ ] Definir e documentar (em `docs/decisions/`) o critério de quando migrar
      para cadastro aberto — o texto de produto diz "massa crítica suficiente para
      se autorregular via reports e reputação", mas o número concreto (ex: X
      usuários ativos, Y reports por semana) deve ser decidido e registrado antes
      do lançamento

---

## Antes de avançar

- [ ] Decisão da seção 4 (papel do banco + Data API + RLS) tomada, registrada como ADR e **testada**
- [ ] Todas as tabelas do MVP criadas, com RLS confirmado via `pg_tables` (não só visualmente no painel)
- [ ] `GET /healthz` responde e confirma conexão com o banco
- [ ] `contracts/openapi.yaml` cobrindo os endpoints do MVP, com formato de erro padronizado
- [ ] Login via Discord funcionando do app até o backend, com escopo `guilds` autorizado
- [ ] Token inválido/expirado rejeitado com 401, tanto no HTTP quanto no handshake do WebSocket
- [ ] Mecanismo de convite implementado no backend antes de qualquer teste com usuários reais
- [ ] Variáveis de ambiente (`.env`) migradas para os nomes novos e **confirmadas fora do Git**

➡️ Próxima fase: [`04-desenvolvimento-mvp.md`](./04-desenvolvimento-mvp.md)
