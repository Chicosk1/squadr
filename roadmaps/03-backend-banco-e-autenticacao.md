# Fase 3 — Backend: Banco de Dados & Autenticação

**Status: 🟨 Em andamento** — schema e RLS das 8 tabelas iniciais ainda pendentes (ver ADR-001)

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [02 — Estrutura do Projeto](./02-estrutura-do-projeto.md)

## O que essa fase entrega

O banco de dados criado no Supabase, protegido corretamente (cada usuário só
acessa o que deveria), e o login "Entrar com Discord" funcionando de ponta a
ponta.

> ⚠️ **Leia isto antes de continuar:** este projeto já teve um problema real
> aqui — documentado no [ADR-001](../docs/decisions/001-migrations-e-rls-supabase.md).
> As tabelas foram criadas e só depois vieram as políticas de segurança
> (RLS), editando os mesmos arquivos de migration que já tinham sido
> aplicados. Isso não funcionou: o Supabase CLI não reaplica uma migration
> só porque o conteúdo do arquivo mudou — ele rastreia pelo nome/timestamp do
> arquivo, não pelo conteúdo. Resultado: 8 tabelas ficaram sem proteção por
> um tempo, mesmo com o código "certo" já escrito. A regra daqui em diante
> (seção 3 abaixo) existe por causa disso.

---

## 1. Criando o projeto no Supabase

- [x] Acessar https://supabase.com/dashboard e criar um novo projeto
- [x] Escolher uma região próxima do Brasil (ex: São Paulo, se disponível, ou a mais próxima)
- [X] Guardar a senha do banco gerada
- [x] No painel do projeto, ir em **Project Settings → API** e copiar:
  - `Project URL`
  - `anon public key`
- [x] Preencher o `.env` (nunca o `.env.example`) na raiz do projeto:
  ```
  EXPO_PUBLIC_SUPABASE_URL=https://SEU-PROJETO.supabase.co
  EXPO_PUBLIC_SUPABASE_ANON_KEY=sua-chave-anon-aqui
  ```
- [x] Espelhar os nomes das variáveis (sem os valores reais) em `.env.example`, para quem clonar o projeto saber o que precisa configurar

---

## 2. Conectando o Supabase CLI ao projeto remoto

- [x] Fazer login:
  ```bash
  supabase login
  ```
- [x] Dentro da pasta do projeto, vincular ao projeto remoto (o ID fica na URL do painel do Supabase):
  ```bash
  supabase link --project-ref SEU_PROJECT_REF
  ```
- [x] Confirmar quais migrations já foram aplicadas **antes de criar qualquer migration nova**:
  ```bash
  supabase migration list
  ```
  Isso mostra o que já está `Applied` remotamente. Se uma migration aparece
  como aplicada, editar o arquivo dela **não terá efeito** — é exatamente o
  problema descrito no ADR-001.

---

## 3. Regra de migrations (obrigatória — vem do ADR-001)

- [x] **Nunca edito uma migration que já rodou no remoto.** Qualquer mudança de
      schema — incluindo "esqueci de proteger essa tabela" — é um arquivo
      novo, com número seguinte.
- [x] **RLS é ativado na mesma migration que cria a tabela**, sempre. Nunca
      "criar tabela" numa migration e "proteger tabela" numa migration
      posterior.
- [x] Toda `CREATE POLICY` é precedida de `DROP POLICY IF EXISTS`, para a
      migration poder ser reaplicada sem erro em outro ambiente.
- [x] Depois de rodar `supabase db push`, confirmar visualmente:
  1. No Table Editor do painel Supabase, a tabela não deve aparecer como `UNRESTRICTED`
  2. Rodar esta query no SQL Editor como segunda confirmação:
     ```sql
     select tablename, rowsecurity from pg_tables where schemaname = 'public';
     ```
     Todas as linhas devem ter `rowsecurity = true`.
- [x] Só depois de confirmado (passo anterior), fazer o commit da migration no Git

---

## 4. Schema do banco (tabelas do MVP)

Ainda faltam :

- [ ] `profiles` (id, discord_id, username, avatar, bio, rank, profile_complete)
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

Para cada tabela desta lista (as 8 do MVP e as 6 da Fase 4): criar a migration
já **com RLS habilitado desde o início** (seção 3), não como um passo
separado depois.

### Checklist de política de segurança (RLS) — pensar por tabela

Para cada tabela, perguntar e registrar a resposta na própria migration:

- [ ] Quem pode **ler** essa linha? (o próprio usuário, os dois lados de um match, todo mundo, etc.)
- [ ] Quem pode **criar** uma linha nova?
- [ ] Quem pode **atualizar**? Existe algum campo que nem o dono deveria poder alterar (ex: `created_at`)?
- [ ] Quem pode **apagar**, se alguém puder?

---

## 5. Autenticação — Discord OAuth

- [ ] No [Discord Developer Portal](https://discord.com/developers/applications), criar uma aplicação nova (ex: "Squadr")
- [ ] Em **OAuth2 → Redirect**, adicionar a URL de callback que o Supabase fornece (Painel Supabase → Authentication → Providers → Discord)
- [ ] Copiar `Client ID` e `Client Secret` da aplicação Discord
- [ ] No painel do Supabase, ir em **Authentication → Providers → Discord**, ativar e colar `Client ID`/`Client Secret`
- [ ] Definir os **escopos** solicitados no OAuth: `identify`, `email` e
      `guilds` (o escopo `guilds` é o que permite, com consentimento do
      usuário, mostrar o selo de vínculo com comunidade — parte do sistema
      de confiança do produto)
- [ ] Implementar `src/services/discord-oauth.ts`, usando o cliente do
      Supabase Auth para iniciar o fluxo e tratar o retorno
- [ ] Testar o login de ponta a ponta no app rodando localmente (Expo Go)

---

## 6. Estratégia de crescimento por convite (fase inicial)

Conforme decidido no documento de produto (seção 9.1), o cadastro no
lançamento **não é aberto** — só entra quem foi convidado por alguém já na
plataforma. Isso tem implicação técnica nesta fase:

- [ ] Adicionar um mecanismo simples de convite antes de abrir o cadastro
      publicamente (pode ser tão simples quanto um código de convite
      validado numa tabela `invites`, ou um allowlist de `discord_id`)
- [ ] Definir e documentar (em `docs/decisions/`) o critério de quando migrar
      para cadastro aberto — o texto de produto diz "massa crítica
      suficiente para se autorregular via reports e reputação", mas o número
      concreto (ex: X usuários ativos, Y reports por semana) deve ser
      decidido e registrado antes do lançamento

---

## Antes de avançar

- [ ] Todas as tabelas do MVP criadas, com RLS confirmado via `pg_tables` (não só visualmente no painel)
- [ ] Login via Discord funcionando localmente, com escopo `guilds` autorizado
- [ ] Mecanismo de convite implementado (mesmo que simples) antes de qualquer teste com usuários reais
- [ ] Variáveis de ambiente (`.env`) configuradas e **confirmadas fora do Git** (`.env` está no `.gitignore`)

➡️ Próxima fase: [`04-desenvolvimento-mvp.md`](./04-desenvolvimento-mvp.md)