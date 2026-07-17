# Fase 0 — Ideação e Planejamento

**Status: ✅ Concluída**

> Voltar para o [Roadmap principal](../ROADMAP.md)

## O que essa fase entrega

Antes de escrever qualquer linha de código, essa fase responde três
perguntas: 
- **O quê** iremos construir; 
- **Para quem** estamos construindo;
- **Com qual tecnologia** vamos construir; 

Essa fase já está feita — o conteúdo completo está em [`docs/contexto-produto.md`](../docs/contexto-produto.md). Este arquivo é só o checklist de conferência do que foi decidido, para você confirmar que nada ficou de fora antes de seguir para a Fase 1.

## Checklist do que foi definido

### Produto
- [x] Posicionamento e problema central definidos
- [x] Pilares do produto (Velocidade, Brasileiro, Confiança)
- [x] Personas principais mapeadas (competitivo solo, casual conectado, sem amigos no jogo, organizador)
- [x] Dados de mercado que justificam o produto (103M jogadores no Brasil, 75,3% jogam digitalmente, 36,5% Gen Z)

### Escopo
- [x] Lista de funcionalidades do **MVP** fechada
- [x] Lista de funcionalidades da **v1.1** fechada
- [x] Lista de funcionalidades da **v2.0** fechada
- [x] Lista explícita do que **não** entra no MVP (evita scope creep)

### Fluxos
- [x] Fluxo de onboarding desenhado (Discord OAuth → seleção de jogos → feed → perfil completo opcional)
- [x] Fluxo de match desenhado (like → match mútuo → chat → sessão)
- [x] Jogos suportados no MVP definidos (Valorant, League of Legends, CS2)

### Stack técnica
- [x] Mobile: React Native + Expo
- [x] Backend/Banco: Supabase (Postgres + Auth + Realtime + Storage)
- [x] Autenticação: Discord OAuth via Supabase Auth
- [x] Notificações: Expo Notifications
- [x] Analytics: PostHog
- [x] Deploy mobile: Expo EAS Build
- [x] Alternativas descartadas e motivo de cada descarte (Flutter, Next.js PWA, Kotlin Multiplatform, Firebase, API própria)

### Modelagem de dados (nível conceitual)
- [x] Tabelas principais mapeadas: `profiles`, `games`, `player_games`,
      `availability`, `swipes`, `matches`, `messages`, `reports`, `squads`,
      `squad_members`, `squad_requests`, `commendations`, `blocks`,
      `game_identities`
- [x] Custo de infraestrutura estimado (R$ 0/mês até ~2k usuários)

### Diferenciais e riscos
- [x] Diferenciais competitivos frente a Discord/WhatsApp mapeados
- [x] Riscos identificados (problema do ovo e da galinha, toxicidade no chat, Discord já ter LFG nativo) com mitigação para cada um
- [x] Sistema de confiança/reputação desenhado (vínculo com comunidade Discord via escopo `guilds`, crescimento por convite na fase inicial, sinais de sessão orgânicos, elogio opcional)

➡️ Próxima fase: [`01-fundacao-ambiente-e-ferramentas.md`](./01-fundacao-ambiente-e-ferramentas.md)