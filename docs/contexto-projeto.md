# Contexto do Squadr — App de Matchmaking para Gamers Brasileiros

> Documento gerado em 29/06/2026. Serve como fonte de verdade para decisões de produto, stack e arquitetura tomadas antes do início do desenvolvimento.
>
> **Atualização 16/07/2026:** incorporadas decisões de produto sobre squads/lobby, sistema de confiança por comunidade e sinais de sessão orgânicos, definidas após análise crítica da ideia frente a concorrentes (Noobly, Sherwa, GameTree).

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
| **Brasileiro** | PT-BR, servidores BR/SA, jogos populares localmente |
| **Confiança** | Reputação, moderação e comunidade saudável desde o primeiro dia |

---

## 2. Público-alvo

### Personas principais

**O competitivo solo** — Joga Valorant ou LoL, quer subir de rank, cansa de jogar com randoms. Precisa de squad com comunicação, nível e papel compatível.

**O casual conectado** — Joga à noite depois do trabalho ou faculdade, quer companhia mais do que resultado. Free Fire, Minecraft, Fortnite.

**O sem amigos no jogo** — Migrou de plataforma ou de jogo e perdeu a base de amigos. Precisa construir nova rede do zero.

**O organizador** — Monta time para torneio, precisa de jogadores específicos por posição ou função. Usa o app como ferramenta de recrutamento.

### Dados de mercado
- Brasil tem **103 milhões de jogadores** — 5º maior mercado consumidor de games do mundo
- **75,3%** dos brasileiros jogam jogos digitais (PGB 2026)
- **36,5%** dos gamers brasileiros são Gen Z — público primário do app

---

## 3. Funcionalidades

### MVP

| Funcionalidade | Descrição |
|---|---|
| **Perfil do jogador** | Nick, avatar (Discord), bio curta, jogos favoritos, rank, estilo (casual/competitivo) e horários disponíveis |
| **Feed de discovery** | Lista de jogadores filtrada por jogo, rank compatível e horário de disponibilidade |
| **Sistema de like e match** | Curtir perfil; match mútuo gera notificação e possibilita o chat |
| **Chat de match** | Conversa em tempo real entre os dois jogadores após o match; botão "copiar convite para o jogo" |
| **Squads (lobby aberto)** | Criar ou entrar em lobbies com vagas definidas (jogo, rank mínimo, função, horário) |
| **Sinais de sessão (orgânicos)** | Clique em "copiar convite", presença simultânea pós-match e reciprocidade de mensagens, coletados como sinal implícito de que a dupla/lobby realmente jogou junto |
| **Login via Discord** | OAuth Discord — 2 cliques para entrar, sem cadastro manual; importa nick, avatar e e-mail |
| **Notificações push** | Alertas de novo match e nova mensagem |
| **Moderação básica** | Reportar usuário e bloquear |
| **Status online** | "Disponível agora" / "Jogando" / "Offline" |

### v1.1
- Reputação visível no perfil, consolidando os sinais de sessão já coletados desde o MVP + elogios recebidos
- Importação automática de rank via Riot API (Valorant) e Steam API (CS2) — também usada para confirmar presença real na mesma partida, substituindo os sinais implícitos por dado definitivo quando disponível
- Grupos fixos/squads permanentes

### v2.0
- Plano Pro, badge de verificação e histórico completo
- Integração com plataforma de campeonatos amadores
- Recomendação por IA baseada em histórico de matches
- Perfil de streamer/criador de conteúdo

### O que **não** entra no MVP
- Feed social (posts, stories, curtidas de conteúdo)
- Ranking global de jogadores
- Streaming ou VOD
- Matchmaking automático com IA
- Grupos fixos/squads permanentes com histórico

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
       └── Estilo de jogo (casual / competitivo)
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
1. Ver card de jogador
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

**Feito para o Brasil** — Jogos populares localmente, servidores BR/SA, PT-BR nativo, fuso de Brasília.

**Moderação desde o MVP** — Report e bloqueio desde o dia 1. Comunidades gamer abandonam apps rápido quando há toxicidade sem moderação.

**Discord como aliado, não concorrente** — OAuth Discord como identidade + distribuição via servidores de Discord BR. O app é um complemento ao ecossistema que os gamers já usam.

**Confiança ancorada em comunidade, não em autodeclaração** — rank, bio e idade continuam autodeclarados, mas o vínculo com uma comunidade Discord verificável (via escopo `guilds` do OAuth) e o crescimento por convite na fase inicial dão um tipo de confiança que apps globais anônimos (GameTree, Noobly) não conseguem replicar em escala.

---

## 8. Riscos identificados

| Risco | Nível | Mitigação |
|---|---|---|
| Problema do ovo e galinha (sem usuários = sem matches) | Alto | Lançar numa comunidade já existente |
| Comportamento tóxico no chat | Alto | Report + moderação desde o MVP, não depois |
| Discord já oferece LFG nativo | Médio | Diferencial é velocidade e foco; encontrar alguém em menos de 2 minutos |

---

## 9. Sistema de Confiança e Reputação

> Contexto: decisão tomada após análise crítica da ideia frente a concorrentes diretos (Noobly, Sherwa, GameTree), cujas avaliações públicas mostram os mesmos três problemas recorrentes: contas falsas/bots, matches que nunca viram sessão real, e confiança baseada só em autodeclaração.

### 9.1 Confiança por comunidade, não por autodeclaração

Rank, bio e idade continuam autodeclarados pelo usuário — isso sozinho é o que já falha em todo concorrente analisado. Dois mecanismos adicionais:

- **Vínculo com comunidade verificada**: usando o escopo `guilds` do Discord OAuth, o app exibe (com consentimento do usuário) um selo indicando que o jogador é membro de um servidor Discord conhecido e ativo. Um perfil com vínculo verificável é muito mais difícil de forjar em massa do que uma conta anônima nova — é o que difere um bot de uma pessoa real com histórico social.
- **Crescimento por convite nos primeiros meses**: no lançamento, novos usuários entram via convite de alguém que já está na plataforma, não por cadastro aberto. Reduz a superfície de bots/perfis falsos justamente na fase em que o app é mais vulnerável a isso. Migra pra cadastro aberto assim que a base tiver massa crítica suficiente pra se autorregular via reports e reputação.

### 9.2 Sinais de sessão orgânicos

Em vez de perguntar diretamente "vocês jogaram junto?", o sistema infere isso a partir de comportamento que já existe:

| Sinal | Fonte | Força |
|---|---|---|
| Clique em "copiar convite para o jogo" | Evento já existente no chat (MVP) | Forte |
| Presença simultânea pós-match | Campos `is_online`/`last_seen_at` já no schema | Média |
| Reciprocidade de mensagens (não uma msg isolada) | Tabela `messages` já existente | Média |
| Confirmação via Riot/Steam API (mesma partida) | Integração v1.1+ | Definitiva |

Nenhum desses sinais exige ação nova do usuário — todos vêm de dados/eventos que o produto já gera.

### 9.3 Elogio opcional

Complementarmente aos sinais implícitos, o chat oferece uma opção com limite diário, sempre visível de elogiar o outro jogador após uma sessão, em categorias fixas (ex: "Bom comunicador", "Apareceu no horário combinado", "Jogaria de novo").

**Regra central: ausência de sinal é neutra, nunca penalidade.** Um usuário que não gera nenhum sinal (não clica no convite, não fica online junto, não é elogiado) não sofre penalização de visibilidade no feed — só não tem dado positivo ainda. Isso evita que o sistema se torne coercitivo ou dependente de resposta forçada.