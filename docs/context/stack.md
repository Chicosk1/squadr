# Stack Técnica

> **Vigente desde 29/07/2026.** Esta é a segunda stack do projeto. A primeira
> (React Native + Expo, com Supabase como backend completo) está registrada no
> [ADR-002](../decisions/002-stack-mobile-react-native-expo.md), marcado como
> supersedido, e o motivo da troca no
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md).

---

## 1. Decisões e justificativas

| Camada | Tecnologia | Justificativa |
|---|---|---|
| **Mobile** | Kotlin Multiplatform + Compose Multiplatform | UI compartilhada de verdade entre Android e iOS, escrita uma vez em Compose; Kotlin é a linguagem com que o dev tem mais afinidade (background Java/Kotlin), o que reduz a curva comparada a manter TypeScript/React só por causa do app; desenvolvimento primário no Android Studio |
| **Backend — API REST** | Go (`cmd/api`) | Regras de negócio (perfil, matching, bloqueio, avaliação, grupos) saem do cliente e passam a viver no servidor, onde não podem ser burladas; Go compila para um binário único, sobe rápido e consome pouca memória — bom encaixe com cobrança por recurso no Fly.io |
| **Backend — WebSocket** | Go (`cmd/ws`), serviço separado | Chat em tempo real tem perfil de carga totalmente diferente do CRUD (conexões longas, muitas ociosas). Serviço separado permite escalar e reiniciar um sem derrubar o outro; ambos compartilham os mesmos pacotes de domínio em `internal/`, então não há duplicação de regra |
| **Banco de dados** | Supabase (só Postgres gerenciado + Storage) | Postgres gerenciado, backup e painel sem custo de operação; as queries de matchmaking são naturalmente relacionais. O Go conecta direto via **`pgx`** — sem SDK intermediária |
| **Autenticação** | Supabase Auth (Discord OAuth), **validação do JWT no Go** | Mantém o login em 2 cliques que os gamers já esperam e evita reimplementar OAuth. A diferença é onde o token é conferido: cada request HTTP e cada abertura de conexão WebSocket é validada pelo próprio Go, contra o **JWKS** do Supabase |
| **Push notifications** | Firebase Cloud Messaging | Padrão de mercado nas duas plataformas; integrado ao Kotlin Multiplatform via **KMPNotifier** (uma API de notificação para Android e iOS). O disparo é feito pelo backend Go via **Firebase Admin SDK**, no mesmo lugar onde o evento (match, mensagem) acontece |
| **Analytics** | PostHog | Mantido da stack anterior. Agora com o **SDK oficial de Kotlin Multiplatform** no app e o **SDK oficial de Go** no backend, apontando para o **mesmo projeto PostHog** — eventos de cliente e de servidor no mesmo funil |
| **Deploy mobile** | Codemagic | Build, assinatura de código e publicação nas duas lojas a partir de um CI que entende projetos Kotlin Multiplatform; roda os builds iOS em máquinas macOS gerenciadas, o que dispensa ter um Mac para publicar |
| **Hospedagem do backend** | Fly.io | Escolhido especificamente por sustentar **conexões WebSocket persistentes** e por permitir rodar a aplicação perto do usuário (região GRU/São Paulo), reduzindo latência do chat |
| **Contrato da API** | OpenAPI em [`contracts/`](../../contracts/) | Fonte única de verdade entre Go e Kotlin. Evita o problema clássico de cliente e servidor discordando sobre o formato do JSON |

---

## 2. Por que não outras opções

| Alternativa | Motivo do descarte |
|---|---|
| **React Native + Expo** (stack anterior) | Ver [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md). Resumo: a UI compartilhada não era o gargalo — o gargalo era ter regra de negócio no cliente e chat preso ao Realtime do Supabase |
| **Supabase como backend completo** (arranjo anterior) | Regras de negócio (matching, bloqueio, reputação) ficavam no app, ou seja, no dispositivo do usuário — impossível garantir integridade. Também amarrava o chat ao Realtime do Supabase, sem controle sobre reconexão, presença e limites |
| Flutter | Dart é linguagem nova para o dev; com Kotlin Multiplatform o dev reaproveita a linguagem que já domina |
| Next.js PWA | Usuário quer app nas lojas; PWA no iOS tem limitações sérias (sem push nativo) |
| Firebase / Firestore como banco | NoSQL é armadilha para dados relacionais; queries de matchmaking são naturalmente SQL. (O Firebase entra **só** para push, via FCM) |
| Node.js ou Java/Spring no backend | Go entrega binário único, startup quase instantâneo e consumo de memória baixo — o que importa quando se paga por instância no Fly.io e se mantém milhares de conexões WebSocket abertas |
| **Vercel / Railway** para o backend | Descartados **por causa do WebSocket**: o modelo serverless/edge da Vercel é hostil a conexões persistentes de longa duração, e nenhum dos dois oferece o controle de região e de processo sempre-ligado que o serviço de chat exige |
| Expo Notifications | Deixou de existir como opção junto com o Expo; FCM é o caminho nativo nas duas plataformas |
| Expo EAS Build | Idem — Codemagic é o equivalente para Kotlin Multiplatform |
| API REST e WebSocket no mesmo processo | Tentador no começo, mas junta dois perfis de carga opostos: um deploy da API derrubaria todas as conversas abertas. Separar desde o início custa pouco (os dois compartilham `internal/`) e evita uma refatoração dolorosa depois |

---

## 3. Decisões ainda abertas

Estas **não** estão decididas e não devem ser tratadas como se estivessem. Cada
uma vira um ADR próprio quando for resolvida:

| Ponto | Candidatos | Quando decidir |
|---|---|---|
| Cliente HTTP no mobile | Ktor Client, Fuel | Fase 2 (estrutura do mobile) |
| Injeção de dependência no mobile | Koin, manual via factories | Fase 2 |
| Cache local no mobile | SQLDelight, Room Multiplatform, nenhum no MVP | Fase 4 (só se o MVP precisar de offline) |
| Roteador HTTP no Go | `net/http` puro (stdlib), chi | Fase 3 |
| Biblioteca de WebSocket no Go | `coder/websocket`, `gorilla/websocket` | Fase 3 |
| Acesso a dados no Go | `pgx` puro, `pgx` + sqlc | Fase 3 |
| Estratégia de RLS vs. autorização no Go | Ver [`banco-de-dados.md`](./banco-de-dados.md), seção 3 | Fase 3 — **antes** de criar as tabelas |

---

## 4. Custo de infra estimado no MVP

> ⚠️ **Mudança relevante em relação à stack antiga:** o MVP **deixou de ser
> R$ 0/mês**. Rodar backend próprio tem custo fixo, e o Fly.io não tem free tier
> para contas novas. É o preço de tirar a regra de negócio do cliente.

| Serviço | Custo |
|---|---|
| Supabase (free tier) | R$ 0/mês |
| **Fly.io — 2 apps (`api` + `ws`)** | **~US$ 4–7/mês.** Sem free tier para contas novas; `shared-cpu-1x` com 256 MB ≈ US$ 2,02/mês por máquina, 512 MB ≈ US$ 3,32/mês. Somar volumes (US$ 0,15/GB/mês) e egress (US$ 0,02–0,12/GB) se usados |
| Firebase Cloud Messaging | R$ 0 — o envio de mensagens FCM não é cobrado |
| PostHog (free tier) | R$ 0/mês |
| Codemagic (free tier) | R$ 0/mês — 500 min de macOS M2/mês, 1 usuário, 1 build por vez. Suficiente para um projeto solo; acima disso, pay-as-you-go (~US$ 0,095/min no M2) |
| Google Play Developer | ~R$ 130 (taxa única) |
| Apple Developer Program | ~R$ 580/ano |
| **Total mensal até ~2k usuários** | **~US$ 4–7/mês** (era R$ 0/mês na stack antiga) |

Valores de Fly.io e Codemagic conferidos em 29/07/2026 nas páginas oficiais de
preço. Reconfirme antes de assumir compromisso — preço de infra muda sem aviso,
e a conversão para real depende do câmbio do dia.
