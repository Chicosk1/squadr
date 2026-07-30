# Fase 4 — Desenvolvimento do MVP

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [03 — Backend: Banco de Dados, API e Autenticação](./03-backend-banco-e-autenticacao.md)

## O que essa fase entrega

Todas as telas e funcionalidades listadas como MVP em
[`docs/context/produto.md`](../context/produto.md) (seção 3), funcionando de ponta
a ponta. Esta fase é a maior do roadmap — por isso está dividida em blocos que
podem ser construídos e testados um por vez, na ordem sugerida.

> 💡 **Fluxo de trabalho dentro de cada bloco — sempre nesta ordem:**
> **contrato → backend Go → app Kotlin.** Testar um endpoint com `curl` custa
> segundos; descobrir o mesmo erro através da UI custa uma tarde. Ao fim de cada
> bloco: testar → commitar → só então seguir.

> 🔄 **O que mudou nesta fase com a troca de stack.** Antes, cada bloco era só
> "tela + query no Supabase". Agora cada bloco tem **duas frentes** (endpoint Go +
> tela Compose) e um contrato no meio. Também apareceu um **Bloco 0** que não
> existia: navegação, cliente HTTP e sessão eram fornecidos de graça pelo Expo
> Router e pela SDK do Supabase — agora são código nosso.

---

## Bloco 0 — Fundação de código (novo)

Sem isso, todos os outros blocos travam.

### Backend (Go)
- [ ] Estrutura de resposta e de **erro** padronizada, conforme definido no contrato (Fase 3, seção 6)
- [ ] Middleware de log de requisição, `recover` de panic e CORS (se necessário)
- [ ] Registro de rotas por domínio, em `cmd/api`, montando os handlers de `internal/*`

### App (Compose Multiplatform)
- [ ] Cliente HTTP em `data/remote/`, com o `Authorization: Bearer` injetado automaticamente e renovação de token quando der 401
- [ ] Navegação entre telas em Compose Multiplatform (login → onboarding → abas → chat)
- [ ] Tema em `ui/theme/` (cores, tipografia, formas) — equivalente ao antigo `constants/theme.ts`
- [ ] Composição de dependências em `di/`
- [ ] Armazenamento seguro do token e leitura da sessão no start (`platform/`, `expect`/`actual`)
- [ ] Estado de carregamento e de erro reaproveitável — para não reinventar em cada tela

✅ Critério de aceite: o app abre, descobre se há sessão salva e chama `GET /healthz` com token, mostrando sucesso ou erro na tela.

---

## Bloco 1 — Onboarding

Referência: [`produto.md`](../context/produto.md), seção 4.

### Backend (Go) — `internal/user`
- [ ] `POST /auth/session` (ou equivalente): primeiro login cria o `profiles` a partir das claims do JWT (nick, avatar, e-mail vindos do Discord)
- [ ] Validação do convite no primeiro login (tabela `invites`, Fase 3 seção 9)
- [ ] `GET /me` e `PATCH /me` (bio, rank, estilo, horários)
- [ ] `PUT /me/games` — jogos selecionados
- [ ] Cálculo de `profile_complete` **no servidor**, não no app

### App (Compose)
- [ ] Tela de boas-vindas com proposta de valor em 1 linha e botão "Entrar com Discord"
- [ ] Fluxo OAuth do Discord (Fase 3, seção 7.3) integrado à tela
- [ ] Exibição de nick e avatar importados após login
- [ ] Tela de seleção de jogos (grid com cards visuais grandes — Valorant, LoL, CS2)
- [ ] Redirecionamento direto para o feed após seleção de jogos (sem exigir mais nada)
- [ ] Banner não bloqueante no feed, disparado após 30–60s de navegação OU na tentativa do primeiro like, oferecendo completar o perfil
- [ ] Tela de "completar perfil" (rank por jogo via lista/seleção, nunca digitação livre; estilo casual/competitivo; horários disponíveis por dia da semana)

### Regra de visibilidade
- [ ] Perfil completo → alta visibilidade no feed; incompleto → visibilidade reduzida, nunca excluído. **Implementar na query do feed, no Go**

✅ Critério de aceite do bloco: um usuário novo sai do zero e chega ao feed em menos de 60 segundos, sem ser bloqueado por falta de dado.

---

## Bloco 2 — Feed de Discovery

### Backend (Go) — `internal/matching`
- [ ] `GET /feed` com paginação: filtra por jogo em comum, compatibilidade de rank e horário de disponibilidade
- [ ] Excluir do feed: o próprio usuário, quem já recebeu swipe, e **quem bloqueou ou foi bloqueado**
- [ ] Ordenação considerando `profile_complete` (Bloco 1)
- [ ] Cuidar do desempenho da query desde já (índices nas colunas de filtro) — é a query mais chamada do produto

### App (Compose)
- [ ] Componente `PlayerCard` em `ui/matching/` com nick, avatar, jogos, rank, estilo, horário
- [ ] Lista/stack de cards navegável (swipe ou botões de like/pular)
- [ ] Estado vazio bem tratado (o que mostrar quando não há mais jogadores compatíveis)
- [ ] Estado de carregamento e de erro de rede

---

## Bloco 3 — Like, Match e Chat

Referência: [`produto.md`](../context/produto.md), seção 5.

### Backend — `internal/matching`
- [ ] `POST /swipes` registra like/pular em `swipes`
- [ ] Detectar like mútuo e criar `matches` — **na mesma transação** do swipe, para não criar match duplicado quando duas pessoas curtem ao mesmo tempo
- [ ] `GET /matches` — lista de matches do usuário
- [ ] `POST /matches/{id}/invite-copied` grava `invite_copied_at` (sinal de sessão mais forte — ver Bloco 6)
- [ ] Disparar push do match via `internal/notification` (Bloco 8)

### Backend — `internal/chat` (serviço `cmd/ws`)
- [ ] Hub de conexões: `user_id` → conexões ativas, com registro/remoção seguros para uso concorrente
- [ ] Autenticação no handshake (Fase 3, seção 7.2)
- [ ] Receber mensagem → validar que existe match e que **não há bloqueio** → persistir em `messages` → entregar ao destinatário conectado
- [ ] Destinatário offline → push via `internal/notification`
- [ ] Heartbeat/ping-pong e timeout de conexão morta
- [ ] Limite de tamanho de mensagem e proteção contra flood
- [ ] `GET /matches/{id}/messages` (na API REST) para o histórico paginado

### App (Compose)
- [ ] Animação de match em `ui/matching/`
- [ ] Tela de chat em `ui/chat/`: histórico via REST + mensagens novas via WebSocket
- [ ] Reconexão automática com backoff quando a conexão cai (trocar de rede, app em background)
- [ ] Ordenação e deduplicação de mensagens (a que chegou pelo WS pode já estar no histórico)
- [ ] Botão "Copiar convite para o jogo" dentro do chat, chamando o endpoint correspondente

> ⚠️ Este é o bloco mais pesado da fase, e o mais afetado pela troca de stack: o
> Supabase Realtime entregava boa parte disso pronto. Ver
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), dificuldade #5.
> Lembrar da limitação de **uma instância** de `cmd/ws` no MVP (dificuldade #6).

---

## Bloco 4 — Squads (lobby aberto)

### Backend (Go) — `internal/group`
- [ ] `POST /squads` — jogo, vagas totais, rank mínimo, função, horário
- [ ] `GET /squads` — listagem filtrável por jogo, com paginação
- [ ] `GET /squads/{id}` — detalhe com membros e vagas restantes
- [ ] `POST /squads/{id}/requests` — pedido para entrar
- [ ] `POST /squad-requests/{id}/accept` — aprovação pelo criador: inserir em `squad_members` **e** incrementar `slots_filled` na **mesma transação**
- [ ] Recusar pedido quando o squad estiver lotado (checagem dentro da transação, não antes)
- [ ] Só o criador pode aprovar — validar autorização, não confiar no app

### App (Compose) — `ui/groups/`
- [ ] Tela de criação de squad
- [ ] Listagem de squads abertos, filtrável por jogo
- [ ] Tela de detalhe com membros atuais e vagas restantes
- [ ] Fluxo de pedido para entrar e visualização do status do pedido

---

## Bloco 5 — Perfil e Status

### Backend (Go) — `internal/user`
- [ ] `GET /players/{id}` — perfil público (respeitando bloqueio)
- [ ] `PATCH /me/status` — "Disponível agora" / "Jogando" / "Offline", refletido em `is_online`/`last_seen_at`
- [ ] Atualizar presença também na conexão/desconexão do WebSocket — é o sinal mais confiável de "está online"
- [ ] Expor o vínculo com comunidade Discord (selo), se o usuário consentiu o escopo `guilds`

### App (Compose) — `ui/profile/`
- [ ] Tela de perfil próprio, editável (bio, rank, estilo, horários, jogos)
- [ ] Seletor de status
- [ ] Exibição do selo de comunidade

---

## Bloco 6 — Sinais de Sessão e Elogio

Referência: [`produto.md`](../context/produto.md), seções 8.2 e 8.3. Estes sinais
**não exigem nenhuma ação nova do usuário** — nascem de eventos que os blocos
anteriores já geram.

### Backend (Go) — `internal/user`
- [ ] Sinal "clique em copiar convite" já existe a partir do Bloco 3 — só precisa ser lido/agregado
- [ ] Sinal "presença simultânea pós-match", a partir dos campos de status do Bloco 5
- [ ] Sinal "reciprocidade de mensagens" (mensagens em ambas as direções na mesma conversa, não uma isolada), calculado a partir de `messages`
- [ ] `POST /commendations` grava elogio, com **limite diário aplicado no servidor**
- [ ] **Regra a implementar explicitamente no código, não só documentar:** ausência de qualquer sinal é neutra — nunca reduzir a visibilidade de um perfil só por falta de sinal positivo
- [ ] Toda a agregação roda no backend. O app apenas exibe o resultado — assim o cálculo não é manipulável pelo cliente

### App (Compose)
- [ ] Modal de elogio pós-sessão no chat, com categorias fixas (ex: "Bom comunicador", "Apareceu no horário combinado", "Jogaria de novo")
- [ ] Feedback claro quando o limite diário já foi atingido

---

## Bloco 7 — Moderação

### Backend (Go) — `internal/user`
- [ ] `POST /reports` grava em `reports`, com motivo
- [ ] `POST /blocks` grava em `blocks`
- [ ] Efeito do bloqueio aplicado em **todas** as leituras: feed, perfil, squads, entrega de mensagem no WebSocket
- [ ] (Se houver tempo no MVP) consulta simples de revisão de reports, mesmo que só via SQL interno

### App (Compose)
- [ ] Botão de reportar usuário (com seleção de motivo)
- [ ] Botão de bloquear usuário, com confirmação
- [ ] Report e bloqueio **visíveis e acessíveis** já na primeira sessão — a Apple rejeita app social sem isso (ver Fase 6)

---

## Bloco 8 — Notificações, Analytics e Polimento

### Push (FCM)
- [ ] Backend: `internal/notification` com Firebase Admin SDK, credencial via variável de ambiente (**nunca** commitada)
- [ ] Backend: `POST /me/device-tokens` para registrar o token do dispositivo em `device_tokens`
- [ ] Backend: disparar push em novo match e em nova mensagem para destinatário offline
- [ ] App: integrar **KMPNotifier**, pedir permissão de notificação, registrar o token no login e ao renovar
- [ ] App: adicionar `google-services.json` (Android) e `GoogleService-Info.plist` (iOS)
- [ ] Testar em dispositivo real nas duas plataformas — push não funciona em todo emulador

### Analytics (PostHog)
- [ ] App: SDK de Kotlin Multiplatform, instrumentando eventos de interface (tela vista, like dado, convite copiado)
- [ ] Backend: SDK de Go, instrumentando eventos de servidor (cadastro concluído, perfil completo, match criado, squad criado, elogio enviado, push enviado)
- [ ] Usar o **mesmo `distinct_id`** (o `sub` do JWT) nos dois SDKs, para o funil cruzar app e servidor
- [ ] Confirmar que os eventos das duas origens chegam no **mesmo projeto** PostHog

### Polimento
- [ ] Estados de carregamento e erro tratados em todas as telas principais (nada de tela branca sem feedback)
- [ ] Comportamento sem internet: avisar, não travar; e reconectar o chat quando a rede voltar
- [ ] Revisão de textos em PT-BR (consistência de tom, sem termos em inglês desnecessários)

---

## Antes de avançar

- [ ] Todos os blocos com critério de aceite cumprido, testados no emulador Android e (se possível) em iOS
- [ ] Fluxo completo testado com duas contas reais diferentes (não só uma conta sozinha) — match, chat e squad dependem de duas pessoas
- [ ] Chat testado em condição ruim: app em background, troca de Wi-Fi para dados, destinatário offline
- [ ] Nenhuma funcionalidade da lista "o que não entra no MVP" ([`produto.md`](../context/produto.md), fim da seção 3) foi implementada por engano
- [ ] `contracts/openapi.yaml` refletindo exatamente o que o backend implementa — sem endpoint fantasma nem endpoint não documentado

➡️ Próxima fase: [`05-qualidade-e-testes.md`](./05-qualidade-e-testes.md)
