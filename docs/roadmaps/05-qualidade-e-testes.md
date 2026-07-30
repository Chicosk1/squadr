# Fase 5 — Qualidade e Testes

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [04 — Desenvolvimento do MVP](./04-desenvolvimento-mvp.md)

## O que essa fase entrega

Confiança de que o app funciona antes de gastar dinheiro/tempo publicando nas
lojas. Para um projeto solo/MVP, o foco é **teste manual estruturado + testes
automatizados nos pontos de maior risco** — não é necessário (nem recomendado,
nesse estágio) cobertura de 100%.

> 🔄 **O que mudou com a troca de stack.** Duas coisas importantes:
> 1. **Autorização virou risco de teste, não de configuração.** Antes, um erro de
>    RLS aparecia como "não vejo o dado". Agora um `if` esquecido num handler Go
>    **vaza dado silenciosamente**. Por isso a seção 3 é nova e não é opcional.
> 2. **O chat é código nosso.** Reconexão, ordenação e entrega offline precisam ser
>    testadas de propósito — não são mais garantia de um serviço de terceiro.

---

## 1. Teste manual estruturado

Use a checklist abaixo como roteiro fixo antes de cada nova versão interna.
Marque, para cada item, se passou (✅) ou falhou (❌) — não deixe em branco.

### Onboarding
- [ ] Login com Discord funciona numa conta nova (nunca usada no app)
- [ ] Login com Discord funciona numa conta que já tem perfil
- [ ] Seleção de jogos leva direto ao feed
- [ ] Banner de perfil incompleto aparece no momento certo (30–60s ou 1º like) e some após completar
- [ ] Usuário sem convite válido não consegue completar o cadastro (regra da fase inicial — ver Fase 3, seção 9)
- [ ] Token expirado é renovado sem o usuário perceber (deixar o app aberto além da validade do token)

### Feed e Match
- [ ] Feed só mostra jogadores com jogo em comum
- [ ] Like unilateral não gera match
- [ ] Like mútuo gera match e notificação nos dois dispositivos
- [ ] Chat abre corretamente após match
- [ ] Botão "copiar convite" funciona e grava o timestamp

### Chat (atenção redobrada — é código novo)
- [ ] Mensagem chega no outro dispositivo em tempo real
- [ ] Histórico carrega ao abrir a conversa, sem duplicar as mensagens que chegaram pelo WebSocket
- [ ] App em background e retorno: a conversa reconecta e não perde mensagem
- [ ] Troca de Wi-Fi para dados móveis no meio da conversa: reconecta
- [ ] Destinatário offline recebe push e, ao abrir, encontra a mensagem no histórico
- [ ] Ordem das mensagens está correta nos dois lados

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
- [ ] Dois usuários pedem a última vaga de um squad ao mesmo tempo (só um deve entrar)
- [ ] Usuário fecha o app no meio do onboarding e reabre (estado não deve corromper)
- [ ] Sem internet: o app avisa, não trava
- [ ] Usuário bloqueado tenta reabrir chat antigo com quem o bloqueou
- [ ] Backend reiniciado no meio de uma conversa (deploy simulado): o app reconecta

---

## 2. Testes automatizados — backend (Go)

Ferramenta: a biblioteca padrão (`testing`) + `net/http/httptest`. Não precisa de
framework extra.

```bash
cd backend && go test ./...
```

Priorize lógica que, se quebrar, corrompe dados ou vaza informação:

- [ ] Detecção de like mútuo e criação de match, incluindo o caso simultâneo (teste de concorrência com transação)
- [ ] Cálculo de compatibilidade de horário/rank do feed
- [ ] Aprovação de pedido de squad com `slots_filled` (não deve passar do total, mesmo com dois pedidos ao mesmo tempo)
- [ ] Agregação dos sinais de sessão (Bloco 6 da Fase 4), incluindo a regra "ausência de sinal é neutra"
- [ ] Limite diário de elogios
- [ ] Validação de JWT em `internal/auth`: token válido, expirado, assinatura errada, `iss`/`aud` errados, `kid` desconhecido (rotação de chave)
- [ ] Efeito do bloqueio em cada leitura (feed, perfil, squads, entrega no WebSocket)

Para os testes que precisam de banco: usar o Postgres local do `supabase start`,
aplicando as migrations do repositório — assim o teste roda contra o mesmo schema
de produção.

- [ ] Configurar o CI (`.github/workflows/`) para rodar `go vet`, `go test ./...` e o build a cada PR

---

## 3. Testes de autorização (novo — não opcional)

Com a autorização em Go, cada endpoint precisa ser testado com **três atores**:

| Ator | Resultado esperado |
|---|---|
| Sem token / token inválido | 401 |
| Token válido, mas **outro** usuário (sem direito ao recurso) | 403 ou 404 — nunca o dado |
| Token válido do dono do recurso | 200 |

- [ ] Cobrir **todo** endpoint que recebe um ID na URL (`/players/{id}`, `/matches/{id}/messages`, `/squads/{id}`, `/squad-requests/{id}/accept`…)
- [ ] Confirmar que ninguém lê mensagem de um match do qual não participa
- [ ] Confirmar que só o criador aprova pedido de squad
- [ ] Confirmar que campos que o usuário não deveria alterar (`created_at`, `slots_filled`, `profile_complete`) são ignorados quando enviados no corpo
- [ ] Testar o handshake do WebSocket com token de outro usuário tentando abrir conversa alheia
- [ ] **Teste do "jeito errado" da Fase 3, seção 4:** com a `anon key` do app, tentar ler tabelas direto pela Data API do Supabase. Deve falhar

---

## 4. Testes automatizados — mobile (Kotlin)

Ferramenta: `kotlin.test`, em `shared/src/commonTest` — roda para as duas
plataformas de uma vez.

```bash
cd mobile && ./gradlew :shared:allTests
```

- [ ] Mapeamento de DTO ↔ modelo de domínio (é onde a divergência com o contrato aparece primeiro)
- [ ] Lógica de `domain/usecase` (validação de formulário, regras de exibição)
- [ ] Deduplicação e ordenação de mensagens no chat
- [ ] Comportamento do cliente HTTP em 401 (deve renovar o token e repetir a chamada uma vez, não entrar em loop)

---

## 5. Teste com pessoas reais antes do lançamento público

Como o crescimento inicial é por convite (Fase 3, seção 9), aproveite isso para
testar com um grupo pequeno e real antes do lançamento:

- [ ] Selecionar 5–10 pessoas de uma comunidade Discord já existente (ver também Fase 6 — estratégia de lançamento)
- [ ] Pedir que completem o fluxo de onboarding sem ajuda e observar onde travam
- [ ] Confirmar que pelo menos alguns matches reais geram sessão de jogo de fato (validação do sistema de sinais, não só da tela)
- [ ] Coletar feedback qualitativo simples (uma pergunta: "o que foi confuso ou irritante?")

---

## 6. Build de teste interno (antes do build de loja)

- [ ] Backend em ambiente real: fazer o primeiro deploy no Fly.io (Fase 6, seção 2) e apontar o app de teste para ele — testar contra `localhost` esconde problemas de TLS, latência e timeout de WebSocket
- [ ] Gerar um build interno pelo **Codemagic** (workflow de `debug`/`internal`) e instalar em dispositivo físico
- [ ] Instalar em pelo menos um Android físico e, se possível, um iPhone físico
- [ ] Repetir a checklist da seção 1 nesse build, contra o backend hospedado
- [ ] Confirmar push notification funcionando no dispositivo físico das duas plataformas

> ⚠️ Sem macOS, o iPhone físico é o **único** jeito de ver o app iOS rodando (não
> há simulador disponível). Planeje conseguir um emprestado, ou aceite lançar
> primeiro só no Android — ver
> [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), dificuldade #3.

---

## Antes de avançar

- [ ] Checklist da seção 1 passou 100%, sem item marcado como ❌ sem correção
- [ ] `go test ./...` e `./gradlew :shared:allTests` passando, e rodando no CI
- [ ] Testes de autorização (seção 3) cobrindo todos os endpoints com ID na URL
- [ ] Teste com grupo real (seção 5) feito e feedback incorporado
- [ ] Build interno testado em dispositivo físico, contra o backend hospedado

➡️ Próxima fase: [`06-deploy-e-lancamento.md`](./06-deploy-e-lancamento.md)
