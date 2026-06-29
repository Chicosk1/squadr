# Contexto do Squadr — App de Matchmaking para Gamers Brasileiros

> Documento gerado em 29/06/2026. Serve como fonte de verdade para decisões de produto, stack e arquitetura tomadas antes do início do desenvolvimento.

---

## 1. Visão do Produto

### Posicionamento
**"O app que conecta gamers brasileiros para jogar juntos — sem depender de grupo de WhatsApp, sem spam de Discord e sem esperar o amigo estar online."**

### Problema central
Encontrar alguém para jogar no Brasil hoje é lento e frustrante. O processo atual exige entrar num servidor de Discord, postar num canal LFG, esperar resposta, checar compatibilidade de rank e combinar horário — tudo manualmente, com ruído e sem garantia. O app resolve esse problema em menos de 3 minutos.

### Pilares do produto
| Pilar | Descrição |
|---|---|
| **Velocidade** | Do download ao primeiro match em menos de 3 minutos |
| **Brasileiro** | PT-BR, servidores BR/SA, jogos populares localmente — não é tradução de app gringo |
| **Confiança** | Reputação, moderação e comunidade saudável desde o dia 1 |

---

## 2. Público-alvo

### Personas principais

**O competitivo solo** — Joga Valorant ou LoL, quer subir de rank, cansa de jogar com randoms. Precisa de squad com comunicação e nível compatível.

**O casual conectado** — Joga à noite depois do trabalho ou faculdade, quer companhia mais do que resultado. Free Fire, Minecraft, Fortnite.

**O sem amigos no jogo** — Migrou de plataforma ou de jogo e perdeu a base de amigos. Precisa construir nova rede do zero.

**O organizador** — Monta time para torneio, precisa de jogadores específicos por posição ou função. Usa o app como ferramenta de recrutamento.

### Dados de mercado
- Brasil tem **103 milhões de jogadores** — 5º maior mercado consumidor de games do mundo
- **75,3%** dos brasileiros jogam jogos digitais (PGB 2026)
- **36,5%** dos gamers brasileiros são Gen Z — público primário do app
- Jogos mais populares no PC: Valorant, League of Legends, CS2
- Jogos mais populares no mobile: Free Fire, Roblox

---

## 3. Funcionalidades

### MVP (lançamento inicial — semanas 1 a 12)

| Funcionalidade | Descrição |
|---|---|
| **Perfil do jogador** | Nick, avatar (Discord), bio curta, jogos favoritos, rank, estilo (casual/competitivo/chill) e horários disponíveis |
| **Feed de discovery** | Lista de jogadores filtrada por jogo, rank compatível e horário de disponibilidade |
| **Sistema de like e match** | Curtir perfil; match mútuo gera notificação e abre chat |
| **Chat de match** | Conversa em tempo real entre os dois jogadores após o match; botão "copiar convite para o jogo" |
| **Login via Discord** | OAuth Discord — 2 cliques para entrar, sem cadastro manual; importa nick, avatar e e-mail |
| **Notificações push** | Alertas de novo match e nova mensagem — loop de retenção principal |
| **Moderação básica** | Reportar usuário e bloquear |
| **Status online** | "Disponível agora" / "Jogando" / "Offline" — filtrável no discovery |

### v1.1 (pós-lançamento, primeiras semanas)
- Sistema de reputação pós-sessão (comunicação, pontualidade, fair play)
- Importação automática de rank via Riot API (Valorant) e Steam API (CS2)
- Squads e grupos fixos de até 5 jogadores
- Notificação "jogador compatível ficou online"

### v2.0 (quando tiver tração validada)
- Plano Pro com filtros avançados, badge de verificação e histórico completo
- Integração com plataforma de campeonatos amadores
- Recomendação por IA baseada em histórico de matches
- Perfil de streamer/criador de conteúdo

### O que **não** entra no MVP
- Feed social (posts, stories, curtidas de conteúdo)
- Ranking global de jogadores
- Streaming ou VOD
- Matchmaking automático com IA

---

## 4. Fluxo de Onboarding (decisão de produto)

### Princípio adotado: Progressive Onboarding
O usuário nunca é bloqueado por falta de informação. A primeira sessão é fluida. Dados adicionais (rank e horários) são incentivados, não obrigatórios.

### Fluxo detalhado

```
1. Baixar o app
   └── Tela de boas-vindas com proposta de valor em 1 linha
       └── Botão "Entrar com Discord"

2. Discord OAuth (2 cliques)
   └── Redireciona para Discord → usuário aceita → volta autenticado
       └── Nick, avatar e e-mail importados automaticamente

3. Seleção de jogos
   └── Grid com jogos disponíveis (cards visuais grandes)
       └── Seleciona os que joga → entra direto no feed
           └── [ Total até aqui: menos de 60 segundos ]

4. Feed de discovery
   └── Após 30–60 segundos de navegação OU ao tentar o primeiro like:
       └── Banner não bloqueante:
           "Seu perfil está incompleto. Adicionar rank e horários
            aumenta suas chances de aparecer para jogadores compatíveis."
           [ Completar agora ]  [ Depois ]

5. Perfil completo (opcional, incentivado)
   └── Rank por jogo (seleção de lista, não digitação)
       └── Estilo de jogo (casual / competitivo / chill)
           └── Horários disponíveis (manhã/tarde/noite por dia da semana)

6. Algoritmo de distribuição
   └── Perfil completo → alta visibilidade no feed de outros usuários
       └── Perfil incompleto → visibilidade reduzida, mas não excluído
```

### Jogos suportados no MVP
- Valorant
- League of Legends
- CS2

---

## 5. Fluxo de Match

```
1. Ver card de jogador no feed
   └── Nick, avatar, jogos, rank, estilo, horário

2. Dar like
   └── Registra o swipe
       └── Se o outro ainda não curtiu → espera silenciosa

3. Match mútuo
   └── Animação de match
       └── Notificação push nos dois dispositivos

4. Chat abre
   └── Conversa em tempo real (Supabase Realtime)
       └── Botão "Copiar convite para o jogo"

5. Sessão acontece
   └── Jogadores se adicionam no jogo e jogam juntos
```

---

## 6. Stack Técnica

### Decisões e justificativas

| Camada | Tecnologia | Justificativa |
|---|---|---|
| **Mobile** | React Native + Expo | Dev tem experiência com React; um código para iOS e Android; Expo elimina configuração de Xcode/Android Studio no início |
| **Backend / Banco** | Supabase | Postgres gerenciado + auth + realtime + storage numa única SDK JS; SQL é familiar para dev com background Java/Kotlin; grátis até escala relevante |
| **Autenticação** | Discord OAuth via Supabase Auth | Gamers já têm Discord; login em 2 cliques sem cadastro; Supabase Auth suporta nativamente |
| **Push Notifications** | Expo Notifications | Já incluso no Expo; gerencia APNs (iOS) e FCM (Android) automaticamente |
| **Analytics** | PostHog | Open source; SDK React Native; free até 1M eventos/mês; instalado desde o dia 1 |
| **Deploy mobile** | Expo EAS Build | Compila .ipa e .apk para as lojas sem precisar de Mac; OTA updates sem esperar aprovação da loja |
| **Hospedagem extra** | Vercel ou Railway (se necessário) | Para Edge Functions além do que o Supabase oferece; free tier suficiente no MVP |

### Por que não outras opções

| Alternativa | Motivo do descarte |
|---|---|
| Flutter | Dart é linguagem nova para o dev; aprender Dart + mobile simultaneamente é risco de atraso |
| Next.js PWA | Usuário quer app nas lojas; PWA no iOS tem limitações sérias (sem push nativo) |
| Kotlin Multiplatform | Ecossistema iOS ainda com fricção; risco alto para projeto solo no MVP |
| Firebase / Firestore | NoSQL é armadilha para dados relacionais; queries de matchmaking são naturalmente SQL |
| API própria (Node ou Django) | Custo alto de setup no MVP; Supabase já resolve auth, WebSocket e storage |

### Estrutura do banco de dados (tabelas principais)

```sql
profiles        (id, discord_id, username, avatar, bio, rank, profile_complete)
games           (id, name, slug, platform)
player_games    (player_id, game_id, rank, style)
availability    (player_id, weekday, period)  -- period: morning/afternoon/night
swipes          (swiper_id, swiped_id, direction, created_at)
matches         (id, player_a, player_b, created_at)
messages        (id, match_id, sender_id, content, created_at)
reports         (id, reporter_id, reported_id, reason, created_at)
```

### Custo de infra estimado no MVP

| Serviço | Custo |
|---|---|
| Supabase (free tier) | R$ 0/mês |
| Expo EAS Build (free tier) | R$ 0/mês |
| PostHog (free tier) | R$ 0/mês |
| Google Play Developer | R$ 130 (taxa única) |
| Apple Developer Program | R$ 580/ano |
| **Total mensal até ~2k usuários** | **R$ 0/mês** |

---

## 7. Diferenciais Competitivos

**Velocidade** — Feed já filtrado por jogo e rank. Match em 1 clique. Chat imediato. Menos de 3 minutos do download ao primeiro match — vs. 15–30 minutos no fluxo atual via Discord/WhatsApp.

**Feito para o Brasil** — Jogos populares localmente, servidores BR/SA, PT-BR nativo, fuso de Brasília. Apps americanos como GamerLink existem mas não entendem o mercado local.

**Moderação desde o MVP** — Report e bloqueio desde o dia 1. Comunidades gamer abandonam apps rápido quando há toxicidade sem moderação.

**Discord como aliado, não concorrente** — OAuth Discord como identidade + distribuição via servidores de Discord BR. O app é um complemento ao ecossistema que os gamers já usam.

---

## 8. Riscos identificados

| Risco | Nível | Mitigação |
|---|---|---|
| Problema do ovo e galinha (sem usuários = sem matches) | Alto | Lançar numa comunidade já existente (Discord do dev), não em geral |
| Comportamento tóxico no chat | Alto | Report + moderação desde o MVP, não depois |
| Discord já oferece LFG nativo | Médio | Diferencial é velocidade e foco; encontrar alguém em menos de 2 minutos |
| APIs da Riot/Valve com limitações | Médio | No MVP, rank declarado pelo usuário; integração de API na v1.1 |