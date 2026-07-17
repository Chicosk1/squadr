# Fase 4 — Desenvolvimento do MVP

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [03 — Backend: Banco de Dados & Autenticação](./03-backend-banco-e-autenticacao.md)

## O que essa fase entrega

Todas as telas e funcionalidades listadas como MVP em
`docs/contexto-produto.md` (seção 3), funcionando de ponta a ponta. Esta
fase é a maior do roadmap — por isso está dividida em blocos que podem ser
construídos e testados um por vez, na ordem sugerida (cada bloco depende do
anterior).

> 💡 Sugestão de fluxo de trabalho por bloco: construir → testar manualmente
> no Expo Go → commitar → só então seguir para o próximo bloco. Blocos
> pequenos e testados evitam ter que debugar três funcionalidades ao mesmo
> tempo.

---

## Bloco 1 — Onboarding

Referência: `docs/contexto-produto.md`, seção 4.

- [ ] Tela de boas-vindas com proposta de valor em 1 linha e botão "Entrar com Discord"
- [ ] Fluxo OAuth do Discord (já com `src/services/discord-oauth.ts` da Fase 3) integrado à tela
- [ ] Importação automática de nick, avatar e e-mail após login
- [ ] Tela de seleção de jogos (grid com cards visuais grandes — Valorant, LoL, CS2)
- [ ] Redirecionamento direto para o feed após seleção de jogos (sem exigir mais nada)
- [ ] Banner não bloqueante no feed, disparado após 30–60s de navegação OU na tentativa do primeiro like, oferecendo completar o perfil
- [ ] Tela de "completar perfil" (rank por jogo via lista/seleção, nunca digitação livre; estilo casual/competitivo; horários disponíveis por dia da semana)
- [ ] Lógica de visibilidade: perfil completo → alta visibilidade no feed; incompleto → visibilidade reduzida, nunca excluído

✅ Critério de aceite do bloco: um usuário novo sai do zero e chega ao feed em menos de 60 segundos, sem ser bloqueado por falta de dado.

---

## Bloco 2 — Feed de Discovery

- [ ] Query no Supabase que filtra jogadores por: jogo em comum, compatibilidade de rank, horário de disponibilidade
- [ ] Componente `PlayerCard` (`src/components/feature/PlayerCard.tsx`) com nick, avatar, jogos, rank, estilo, horário
- [ ] Lista/stack de cards navegável (swipe ou botões de like/pular)
- [ ] Estado vazio bem tratado (o que mostrar quando não há mais jogadores compatíveis)

---

## Bloco 3 — Like, Match e Chat

Referência: `docs/contexto-produto.md`, seção 5.

- [ ] Registrar swipe na tabela `swipes` (like/pular)
- [ ] Detectar like mútuo e criar registro em `matches`
- [ ] Animação de match (`src/components/feature/MatchAnimation.tsx`)
- [ ] Notificação push nos dois dispositivos no momento do match (Expo Notifications)
- [ ] Tela de chat em tempo real usando Supabase Realtime, ligada à tabela `messages`
- [ ] Botão "Copiar convite para o jogo" dentro do chat, que grava `invite_copied_at` em `matches` (esse clique é também o sinal de sessão mais forte — ver Bloco 6)

---

## Bloco 4 — Squads (lobby aberto)

- [ ] Tela de criação de squad: jogo, vagas totais, rank mínimo, função, horário
- [ ] Listagem de squads abertos, filtrável por jogo
- [ ] Fluxo de pedido para entrar (`squad_requests`) e aprovação pelo criador
- [ ] Atualização de `slots_filled` ao aceitar um pedido
- [ ] Tela de detalhe do squad com membros atuais e vagas restantes

---

## Bloco 5 — Perfil e Status

- [ ] Tela de perfil próprio, editável (bio, rank, estilo, horários, jogos)
- [ ] Status "Disponível agora" / "Jogando" / "Offline", refletido em `profiles`/campos de presença
- [ ] Exibição do vínculo com comunidade Discord (selo, se o usuário consentiu o escopo `guilds`)

---

## Bloco 6 — Sinais de Sessão e Elogio

Referência: `docs/contexto-produto.md`, seção 9.2 e 9.3. Estes sinais **não
exigem nenhuma ação nova do usuário** — nascem de eventos que os blocos
anteriores já geram.

- [ ] Sinal "clique em copiar convite" já existe a partir do Bloco 3 — só precisa ser lido/agregado
- [ ] Sinal "presença simultânea pós-match", a partir dos campos de status do Bloco 5
- [ ] Sinal "reciprocidade de mensagens" (mensagens em ambas as direções na mesma conversa, não uma isolada), calculado a partir de `messages`
- [ ] Tela/modal de elogio pós-sessão no chat, com categorias fixas (ex: "Bom comunicador", "Apareceu no horário combinado", "Jogaria de novo") e limite diário
- [ ] Grava elogio em `commendations`
- [ ] **Regra a implementar explicitamente no código, não só documentar:** ausência de qualquer sinal é neutra — nunca reduzir a visibilidade de um perfil só por falta de sinal positivo

---

## Bloco 7 — Moderação

- [ ] Botão de reportar usuário (grava em `reports`, com motivo)
- [ ] Botão de bloquear usuário (grava em `blocks`)
- [ ] Efeito do bloqueio: usuário bloqueado desaparece do feed e não consegue iniciar novo chat com quem bloqueou
- [ ] (Se houver tempo no MVP) painel simples de revisão de reports, mesmo que manual/interno

---

## Bloco 8 — Notificações, Analytics e Polimento

- [ ] Configurar Expo Notifications de ponta a ponta (permissão do usuário, token salvo, envio no evento de match/mensagem)
- [ ] Instrumentar eventos principais no PostHog: cadastro concluído, perfil completo, like dado, match criado, convite copiado, squad criado, elogio enviado
- [ ] Estados de carregamento e erro tratados em todas as telas principais (nada de tela branca sem feedback)
- [ ] Revisão de textos em PT-BR (consistência de tom, sem termos em inglês desnecessários)

---

## Antes de avançar

- [ ] Todos os 8 blocos com critério de aceite cumprido, testados manualmente no Expo Go
- [ ] Fluxo completo testado com duas contas reais diferentes (não só uma conta sozinha) — match, chat e squad dependem de duas pessoas
- [ ] Nenhuma funcionalidade da lista "o que não entra no MVP" (`docs/contexto-produto.md`, fim da seção 3) foi implementada por engano

➡️ Próxima fase: [`05-qualidade-e-testes.md`](./05-qualidade-e-testes.md)