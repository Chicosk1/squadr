# Produto — Squadr, app de matchmaking para gamers brasileiros

> Documento gerado em 29/06/2026. Serve como fonte de verdade para decisões de produto.
>
> **Atualização 16/07/2026:** incorporadas decisões de produto sobre squads/lobby, sistema de confiança por comunidade e sinais de sessão orgânicos, definidas após análise crítica da ideia frente a concorrentes (Noobly, Sherwa, GameTree).
>
> **Atualização 29/07/2026:** este arquivo era o contexto único do projeto (`docs/contexto-projeto.md`) e passou a conter **só produto**. A stack técnica saiu para [`stack.md`](./stack.md), a arquitetura para [`arquitetura.md`](./arquitetura.md) e o schema para [`banco-de-dados.md`](./banco-de-dados.md), na migração de stack descrita no [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md). **Nenhuma decisão de produto mudou por causa da troca de stack** — o escopo do MVP é exatamente o mesmo.

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
   └── Conversa em tempo real (serviço WebSocket do backend Go)
       └── Botão "Copiar convite para o jogo"

5. Sessão acontece
   └── Jogadores se adicionam no jogo e jogam juntos
```

> O detalhamento técnico de cada passo (qual endpoint, quem valida o token, como
> a mensagem chega no outro dispositivo) está em
> [`arquitetura.md`](./arquitetura.md), seção 5.

---

## 6. Diferenciais Competitivos

**Velocidade** — Feed já filtrado por jogo e rank. Match em 1 clique. Chat imediato. Menos de 3 minutos do download ao primeiro match — vs. 15–30 minutos no fluxo atual via Discord/WhatsApp.

**Feito para o Brasil** — Jogos populares localmente, servidores BR/SA, PT-BR nativo, fuso de Brasília.

**Moderação desde o MVP** — Report e bloqueio desde o dia 1. Comunidades gamer abandonam apps rápido quando há toxicidade sem moderação.

**Discord como aliado, não concorrente** — OAuth Discord como identidade + distribuição via servidores de Discord BR. O app é um complemento ao ecossistema que os gamers já usam.

**Confiança ancorada em comunidade, não em autodeclaração** — rank, bio e idade continuam autodeclarados, mas o vínculo com uma comunidade Discord verificável (via escopo `guilds` do OAuth) e o crescimento por convite na fase inicial dão um tipo de confiança que apps globais anônimos (GameTree, Noobly) não conseguem replicar em escala.

---

## 7. Riscos identificados

| Risco | Nível | Mitigação |
|---|---|---|
| Problema do ovo e galinha (sem usuários = sem matches) | Alto | Lançar numa comunidade já existente |
| Comportamento tóxico no chat | Alto | Report + moderação desde o MVP, não depois |
| Discord já oferece LFG nativo | Médio | Diferencial é velocidade e foco; encontrar alguém em menos de 2 minutos |

> Riscos **de execução técnica** introduzidos pela troca de stack (curva de
> aprendizado de Go e Kotlin Multiplatform, build iOS exigindo macOS, custo de
> infra deixando de ser zero) estão no
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), seção "Dificuldades
> esperadas na transição" — não aqui, porque não são riscos de produto.

---

## 8. Sistema de Confiança e Reputação

> Contexto: decisão tomada após análise crítica da ideia frente a concorrentes diretos (Noobly, Sherwa, GameTree), cujas avaliações públicas mostram os mesmos três problemas recorrentes: contas falsas/bots, matches que nunca viram sessão real, e confiança baseada só em autodeclaração.

### 8.1 Confiança por comunidade, não por autodeclaração

Rank, bio e idade continuam autodeclarados pelo usuário — isso sozinho é o que já falha em todo concorrente analisado. Dois mecanismos adicionais:

- **Vínculo com comunidade verificada**: usando o escopo `guilds` do Discord OAuth, o app exibe (com consentimento do usuário) um selo indicando que o jogador é membro de um servidor Discord conhecido e ativo. Um perfil com vínculo verificável é muito mais difícil de forjar em massa do que uma conta anônima nova — é o que difere um bot de uma pessoa real com histórico social.
- **Crescimento por convite nos primeiros meses**: no lançamento, novos usuários entram via convite de alguém que já está na plataforma, não por cadastro aberto. Reduz a superfície de bots/perfis falsos justamente na fase em que o app é mais vulnerável a isso. Migra pra cadastro aberto assim que a base tiver massa crítica suficiente pra se autorregular via reports e reputação.

### 8.2 Sinais de sessão orgânicos

Em vez de perguntar diretamente "vocês jogaram junto?", o sistema infere isso a partir de comportamento que já existe:

| Sinal | Fonte | Força |
|---|---|---|
| Clique em "copiar convite para o jogo" | Evento já existente no chat (MVP) | Forte |
| Presença simultânea pós-match | Campos `is_online`/`last_seen_at` já no schema | Média |
| Reciprocidade de mensagens (não uma msg isolada) | Tabela `messages` já existente | Média |
| Confirmação via Riot/Steam API (mesma partida) | Integração v1.1+ | Definitiva |

Nenhum desses sinais exige ação nova do usuário — todos vêm de dados/eventos que o produto já gera.

> **O que a stack nova muda aqui (para melhor):** a agregação desses sinais
> passa a ser responsabilidade do backend Go (pacote `internal/user`), não do
> app. O cálculo fica server-side e o cliente não consegue influenciá-lo — o que
> a arquitetura antiga (app falando direto com o banco via SDK) não garantia.

### 8.3 Elogio opcional

Complementarmente aos sinais implícitos, o chat oferece uma opção com limite diário, sempre visível de elogiar o outro jogador após uma sessão, em categorias fixas (ex: "Bom comunicador", "Apareceu no horário combinado", "Jogaria de novo").

**Regra central: ausência de sinal é neutra, nunca penalidade.** Um usuário que não gera nenhum sinal (não clica no convite, não fica online junto, não é elogiado) não sofre penalização de visibilidade no feed — só não tem dado positivo ainda. Isso evita que o sistema se torne coercitivo ou dependente de resposta forçada.
