# ADR-002: Stack inicial — React Native + Expo com Supabase como backend completo

> ## ⛔ SUPERSEDIDO
>
> **Substituído pelo [ADR-003 — Migração para Kotlin Multiplatform e backend Go](./003-migracao-kmp-e-backend-go.md), em 29/07/2026.**
>
> Este registro é mantido **de propósito**. Ele documenta o racional que levou à
> escolha de React Native + Expo e do Supabase como backend completo, incluindo
> as alternativas que foram descartadas na época — entre elas, ironicamente, o
> Kotlin Multiplatform. Apagar isso esconderia por que o projeto começou como
> começou e por que a troca custou o que custou.
>
> **Nada aqui deve ser usado como orientação para escrever código hoje.** Para o
> estado atual, ver [`docs/context/stack.md`](../context/stack.md).

---

**Status:** Supersedido pelo ADR-003
**Data da decisão:** 29/06/2026
**Data deste registro:** 29/07/2026

> Nota sobre a numeração: a decisão original foi tomada em 29/06/2026 e vivia no
> documento de contexto do projeto (`docs/contexto-projeto.md`, seção 6), não
> como ADR. Ela foi formalizada como ADR em 29/07/2026, durante a migração de
> stack, justamente para poder ser marcada como supersedida sem que o racional
> fosse perdido na reescrita do contexto. O número **002** reflete a ordem em que
> o registro foi criado, não a data da decisão.

---

## Contexto (em 29/06/2026)

Projeto solo, primeiro app do dev, com meta de MVP nas lojas rapidamente e
orçamento próximo de zero. O dev tinha experiência com React e background em
Java/Kotlin. O produto (ver [`docs/context/produto.md`](../context/produto.md))
precisava de perfil, feed, match, chat em tempo real, push e moderação.

A pergunta que se colocava era: **como uma pessoa entrega isso nas duas lojas
sem time de backend?**

## Decisão (original)

| Camada | Tecnologia | Justificativa registrada na época |
|---|---|---|
| Mobile | React Native + Expo | Dev tem experiência com React; um código para iOS e Android; Expo elimina configuração de Xcode/Android Studio no início |
| Backend / Banco | Supabase | Postgres gerenciado + auth + realtime + storage numa única SDK JS; SQL é familiar para dev com background Java/Kotlin; grátis até escala relevante |
| Autenticação | Discord OAuth via Supabase Auth | Gamers já têm Discord; login em 2 cliques sem cadastro; Supabase Auth suporta nativamente |
| Push Notifications | Expo Notifications | Já incluso no Expo; gerencia APNs (iOS) e FCM (Android) automaticamente |
| Analytics | PostHog | Open source; SDK React Native; free até 1M eventos/mês; instalado desde o dia 1 |
| Deploy mobile | Expo EAS Build | Compila `.ipa` e `.apk` para as lojas sem precisar de Mac; OTA updates sem esperar aprovação da loja |
| Hospedagem extra | Vercel ou Railway (se necessário) | Para Edge Functions além do que o Supabase oferece; free tier suficiente no MVP |

Ponto central do arranjo: **não havia backend próprio.** O app falava direto com
o Supabase pela SDK JS — inclusive as regras de matching, bloqueio e reputação
rodavam no dispositivo do usuário, e a proteção dos dados dependia inteiramente
do RLS do Postgres (daí a relevância do [ADR-001](./001-migrations-e-rls-supabase.md)).

### Alternativas descartadas na época

| Alternativa | Motivo do descarte registrado em 29/06/2026 |
|---|---|
| Flutter | Dart é linguagem nova para o dev; aprender Dart + mobile simultaneamente é risco de atraso |
| Next.js PWA | Usuário quer app nas lojas; PWA no iOS tem limitações sérias (sem push nativo) |
| **Kotlin Multiplatform** | **Ecossistema iOS ainda com fricção; risco alto para projeto solo no MVP** |
| Firebase / Firestore | NoSQL é armadilha para dados relacionais; queries de matchmaking são naturalmente SQL |
| API própria (Node ou Django) | Custo alto de setup no MVP; Supabase já resolve auth, WebSocket e storage |

As duas linhas em destaque são exatamente as que o ADR-003 reverte. O risco
apontado no descarte do Kotlin Multiplatform (fricção no iOS) **não foi
refutado** — ele foi reavaliado e aceito. Ver a seção "Dificuldades esperadas na
transição" do ADR-003.

## Consequências que se materializaram

- Ambiente de desenvolvimento montado sobre Node.js, Expo CLI e Expo Go, e
  contas criadas no Expo/EAS (Fases 1 e 2 do roadmap, concluídas).
- Estrutura de pastas `app/` (Expo Router) + `src/` criada, com arquivos
  placeholder.
- Projeto Supabase criado e vinculado via CLI; o incidente de RLS que gerou o
  [ADR-001](./001-migrations-e-rls-supabase.md) aconteceu neste período.
- **Nenhuma tela ou serviço chegou a ser implementado.** Todos os arquivos de
  `app/` e `src/` estavam com 0 byte quando a migração começou — o que reduziu
  drasticamente o custo da troca. Ver ADR-003.

## Por que foi supersedido

Resumo: o arranjo otimizava para *velocidade de entrega inicial* e pagava com
*regra de negócio rodando no cliente* e *chat sem controle próprio*. Quando a
prioridade mudou, o arranjo deixou de servir. O racional completo está no
[ADR-003](./003-migracao-kmp-e-backend-go.md).
