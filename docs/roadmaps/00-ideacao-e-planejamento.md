# Fase 0 — Ideação e Planejamento

**Status: ✅ Concluída** — decisões de produto de 29/06/2026; **stack revisada em 29/07/2026**

> Voltar para o [Roadmap principal](./README.md)

## O que essa fase entrega

Antes de escrever qualquer linha de código, essa fase responde três
perguntas:
- **O quê** iremos construir;
- **Para quem** estamos construindo;
- **Com qual tecnologia** vamos construir;

Essa fase já está feita — o conteúdo completo está em
[`docs/context/`](../context/). Este arquivo é só o checklist de conferência do
que foi decidido.

> ⚠️ **A terceira pergunta foi respondida duas vezes.** A stack escolhida em
> 29/06/2026 (React Native + Expo, Supabase como backend completo) foi
> substituída em 29/07/2026. As decisões de **produto** não mudaram nada nisso —
> escopo do MVP, personas e fluxos seguem idênticos. Ver
> [ADR-002](../decisions/002-stack-mobile-react-native-expo.md) (supersedido) e
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md) (vigente).

## Checklist do que foi definido

### Produto
- [ ] Posicionamento e problema central definidos
- [ ] Pilares do produto (Velocidade, Brasileiro, Confiança)
- [ ] Personas principais mapeadas (competitivo solo, casual conectado, sem amigos no jogo, organizador)
- [ ] Dados de mercado que justificam o produto (103M jogadores no Brasil, 75,3% jogam digitalmente, 36,5% Gen Z)

### Escopo
- [ ] Lista de funcionalidades do **MVP** fechada
- [ ] Lista de funcionalidades da **v1.1** fechada
- [ ] Lista de funcionalidades da **v2.0** fechada
- [ ] Lista explícita do que **não** entra no MVP (evita scope creep)

### Fluxos
- [ ] Fluxo de onboarding desenhado (Discord OAuth → seleção de jogos → feed → perfil completo opcional)
- [ ] Fluxo de match desenhado (like → match mútuo → chat → sessão)
- [ ] Jogos suportados no MVP definidos (Valorant, League of Legends, CS2)

### Stack técnica (revisada em 29/07/2026)
- [ ] Mobile: **Kotlin Multiplatform + Compose Multiplatform** (UI compartilhada, Android Studio)
- [ ] Backend: **Go**, dois serviços — `cmd/api` (REST) e `cmd/ws` (WebSocket) — compartilhando `internal/`
- [ ] Banco: **Supabase como Postgres gerenciado + Storage**, acessado pelo Go via `pgx`
- [ ] Autenticação: **Discord OAuth via Supabase Auth**, com o **JWT validado pelo Go** (JWKS)
- [ ] Notificações: **Firebase Cloud Messaging** (KMPNotifier no app, Firebase Admin SDK no backend)
- [ ] Analytics: **PostHog** (SDK Kotlin Multiplatform + SDK Go, mesmo projeto)
- [ ] Deploy mobile: **Codemagic**
- [ ] Hospedagem do backend: **Fly.io** (escolhido por WebSocket persistente e região GRU)
- [ ] Contrato da API: **OpenAPI** em `contracts/`, fonte única de verdade
- [ ] Alternativas descartadas e motivo de cada descarte registrados ([`stack.md`](../context/stack.md), seção 2)
- [ ] Migração de stack registrada como ADR, com dificuldades esperadas mapeadas

### Arquitetura (nova nesta revisão)
- [ ] Decisão de que o app **nunca** fala direto com o banco — tudo passa pelo Go
- [ ] Decisão de separar API REST e WebSocket em dois serviços, compartilhando pacotes de domínio
- [ ] Organização do backend **por domínio**, não por camada técnica
- [ ] Limitação conhecida e aceita: hub de WebSocket em memória ⇒ uma instância de `cmd/ws` no MVP

### Modelagem de dados (nível conceitual)
- [ ] Tabelas principais mapeadas: `profiles`, `games`, `player_games`,
      `availability`, `swipes`, `matches`, `messages`, `reports`, `squads`,
      `squad_members`, `squad_requests`, `commendations`, `blocks`,
      `game_identities`, `device_tokens`, `invites`
- [ ] Custo de infraestrutura estimado — **atenção: deixou de ser R$ 0/mês** com backend próprio (~US$ 4–7/mês no Fly.io)

### Diferenciais e riscos
- [ ] Diferenciais competitivos frente a Discord/WhatsApp mapeados
- [ ] Riscos de produto identificados (ovo e galinha, toxicidade no chat, Discord com LFG nativo) com mitigação para cada um
- [ ] Riscos técnicos da troca de stack identificados ([ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), "Dificuldades esperadas na transição")
- [ ] Sistema de confiança/reputação desenhado (vínculo com comunidade Discord via escopo `guilds`, crescimento por convite na fase inicial, sinais de sessão orgânicos, elogio opcional)

➡️ Próxima fase: [`01-fundacao-ambiente-e-ferramentas.md`](./01-fundacao-ambiente-e-ferramentas.md)
