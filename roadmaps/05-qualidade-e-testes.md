# Fase 5 — Qualidade e Testes

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [04 — Desenvolvimento do MVP](./04-desenvolvimento-mvp.md)

## O que essa fase entrega

Confiança de que o app funciona antes de gastar dinheiro/tempo publicando
nas lojas. Para um projeto solo/MVP, o foco é **teste manual estruturado +
alguns testes automatizados nos pontos de maior risco** — não é necessário
(nem recomendado, nesse estágio) cobertura de testes automatizados de 100%.

---

## 1. Teste manual estruturado

Use a checklist abaixo como roteiro fixo antes de cada nova versão interna.
Marque, para cada item, se passou (✅) ou falhou (❌) — não deixe em branco.

### Onboarding
- [ ] Login com Discord funciona numa conta nova (nunca usada no app)
- [ ] Login com Discord funciona numa conta que já tem perfil
- [ ] Seleção de jogos leva direto ao feed
- [ ] Banner de perfil incompleto aparece no momento certo (30–60s ou 1º like) e some após completar
- [ ] Usuário sem convite válido não consegue completar o cadastro (regra da fase inicial — ver Fase 3, seção 6)

### Feed e Match
- [ ] Feed só mostra jogadores com jogo em comum
- [ ] Like unilateral não gera match
- [ ] Like mútuo gera match e notificação nos dois dispositivos
- [ ] Chat abre corretamente após match
- [ ] Botão "copiar convite" funciona e grava o timestamp

### Squads
- [ ] Criar squad com vagas define corretamente `slots_total`
- [ ] Pedido para entrar aparece para o criador
- [ ] Aprovar pedido atualiza `slots_filled` e adiciona o membro
- [ ] Squad lotado não aceita novos pedidos

### Reputação e Moderação
- [ ] Elogio só pode ser dado dentro do limite diário
- [ ] Reportar usuário grava o motivo corretamente
- [ ] Bloquear usuário remove ele do feed de quem bloqueou
- [ ] Usuário sem nenhum sinal positivo não tem a visibilidade reduzida (regra de neutralidade)

### Casos de borda (o que quebra apps de matchmaking com frequência)
- [ ] Dois usuários dão like simultaneamente um no outro (checar se não duplica o match)
- [ ] Usuário fecha o app no meio do onboarding e reabre (estado não deve corromper)
- [ ] Sem internet: o app avisa, não trava
- [ ] Usuário bloqueado tenta reabrir chat antigo com quem o bloqueou

---

## 2. Testes automatizados (onde vale o esforço no MVP)

Não é necessário testar toda tela. Priorize lógica que, se quebrar, corrompe
dados ou vaza informação:

- [ ] Lógica de detecção de match mútuo (unit test na função, não na tela)
- [ ] Cálculo de compatibilidade de horário/rank no feed
- [ ] Regras de RLS do Supabase — testar diretamente com queries usando
      diferentes usuários simulados, confirmando que um usuário não lê/edita
      dados de outro fora do que a política permite
- [ ] Lógica de agregação dos sinais de sessão (Bloco 6 da Fase 4)

Ferramenta sugerida: **Jest** (já compatível com projetos Expo/React Native).

```bash
npm install --save-dev jest @testing-library/react-native
```

---

## 3. Teste com pessoas reais antes do lançamento público

Como o crescimento inicial é por convite (Fase 3, seção 6), aproveite isso
para testar com um grupo pequeno e real antes do lançamento:

- [ ] Selecionar 5–10 pessoas de uma comunidade Discord já existente (ver
      também Fase 6 — estratégia de lançamento)
- [ ] Pedir que completem o fluxo de onboarding sem ajuda e observar onde travam
- [ ] Confirmar que pelo menos alguns matches reais geram sessão de jogo de fato (validação do sistema de sinais, não só da tela)
- [ ] Coletar feedback qualitativo simples (uma pergunta: "o que foi confuso ou irritante?")

---

## 4. Build de teste interno (antes do build de loja)

- [ ] Gerar um build interno via EAS (`eas build --profile preview`) para
      testar em dispositivos físicos reais, fora do Expo Go — o Expo Go não
      reflete 100% do comportamento de um build final
- [ ] Instalar esse build em pelo menos um Android físico e, se possível, um iPhone físico
- [ ] Repetir a checklist da seção 1 nesse build (não só no Expo Go)

---

## Antes de avançar

- [ ] Checklist da seção 1 passou 100%, sem item marcado como ❌ sem correção
- [ ] Testes automatizados dos pontos críticos (seção 2) passando
- [ ] Teste com grupo real (seção 3) feito e feedback incorporado
- [ ] Build interno (seção 4) testado em dispositivo físico, não só no Expo Go

➡️ Próxima fase: [`06-deploy-e-lancamento.md`](./06-deploy-e-lancamento.md)