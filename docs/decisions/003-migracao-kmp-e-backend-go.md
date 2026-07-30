# ADR-003: Migração para Kotlin Multiplatform no mobile e backend próprio em Go

**Status:** Aceito
**Data:** 29/07/2026
**Supersede:** [ADR-002 — Stack inicial: React Native + Expo com Supabase como backend completo](./002-stack-mobile-react-native-expo.md)
**Emenda:** [ADR-001 — Migrations e RLS no Supabase](./001-migrations-e-rls-supabase.md) (continua válido; o papel do RLS muda — ver seção "Consequências")

---

## Contexto

O projeto foi iniciado com React Native + Expo e Supabase como backend completo
(ADR-002). Nesse arranjo **não havia servidor de aplicação**: o app falava direto
com o Postgres pela SDK JS do Supabase, o chat usava o Supabase Realtime, o push
usava Expo Notifications e o build saía do Expo EAS.

O estado do projeto no momento desta decisão:

- Fases 0, 1 e 2 do roadmap concluídas (planejamento, ambiente, estrutura de pastas).
- Fase 3 em andamento: projeto Supabase criado e vinculado via CLI, mas
  `supabase/migrations/` **vazia** — nenhuma tabela criada.
- `app/` e `src/` existiam com a estrutura de arquivos planejada, mas **todos os
  arquivos estavam com 0 byte**. Nenhuma tela, nenhum serviço, nenhum hook
  implementado.

Ou seja: a migração acontece no último momento em que ela é barata.

## Decisão

Trocar a stack conforme abaixo. As duas mudanças estruturais são: **(a)** passar
a ter um backend de aplicação próprio e **(b)** trocar a tecnologia de UI
multiplataforma.

| Camada | Antes (ADR-002) | Agora |
|---|---|---|
| Mobile | React Native + Expo | **Kotlin Multiplatform + Compose Multiplatform** (UI compartilhada Android/iOS, desenvolvimento primário no Android Studio) |
| Backend de aplicação | *não existia* — app falava direto com o Supabase | **Go**, em dois serviços: `cmd/api` (REST, regras de negócio) e `cmd/ws` (WebSocket, chat), compartilhando os mesmos pacotes de `internal/` |
| Banco de dados | Supabase via SDK JS, do app | **Supabase como Postgres gerenciado + Storage**; o Go conecta direto via **`pgx`** |
| Autenticação | Discord OAuth via Supabase Auth, consumido pelo app | **Supabase Auth (Discord OAuth) mantido**; quem valida o JWT em cada request e em cada abertura de WebSocket agora é o **Go**, contra o **JWKS** do Supabase |
| Push | Expo Notifications | **Firebase Cloud Messaging** — recepção no app via **KMPNotifier**, disparo pelo backend via **Firebase Admin SDK** |
| Analytics | PostHog (SDK React Native) | **PostHog mantido** — SDK oficial de **Kotlin Multiplatform** no app + SDK oficial de **Go** no backend, mesmo projeto |
| Deploy mobile | Expo EAS Build | **Codemagic** (build, assinatura e publicação nas duas lojas) |
| Hospedagem do backend | Vercel/Railway (cogitado, nunca implementado) | **Fly.io** |
| Contrato da API | não se aplicava | **OpenAPI em `contracts/`** como fonte única de verdade |

Detalhamento da arquitetura resultante: [`docs/context/arquitetura.md`](../context/arquitetura.md).

## Por quê

1. **Regra de negócio no cliente é indefensável.** No arranjo antigo, matching,
   bloqueio e cálculo de reputação rodavam no dispositivo do usuário, e a única
   proteção dos dados era o RLS. Um app é código que está na mão de quem quiser
   inspecioná-lo. Com o Go no meio, a regra passa a rodar onde o usuário não
   alcança — e operações que precisam ser atômicas (like mútuo criando match,
   aprovar pedido de squad e incrementar `slots_filled`) passam a rodar em
   transação de verdade.

2. **O chat precisa ser nosso.** O Supabase Realtime entregava mensagens sem
   trabalho, mas sem controle: reconexão, presença, ordenação, backpressure e
   limites de conexão eram o que o serviço decidisse. Chat em tempo real é a
   funcionalidade central do produto (é o que acontece depois do match) — ficar
   sem controle sobre ela é risco no lugar errado. Daí o serviço `cmd/ws`
   separado, com hub próprio.

3. **Alinhamento de linguagem.** O dev tem background em Java/Kotlin. Kotlin
   Multiplatform + Compose deixa a UI e o domínio do app na linguagem que ele
   domina, e Go é uma linguagem pequena o suficiente para ser aprendida em
   paralelo sem virar um projeto à parte. O React era conhecido, mas mantê-lo
   significava manter um ecossistema (JS/TS, npm, Metro) só por causa do app.

4. **Separação de API e WebSocket como decisão de arquitetura, não de
   escala.** CRUD e conexões persistentes têm perfis de carga opostos. Se
   estivessem no mesmo processo, um deploy da API derrubaria todas as conversas
   abertas. Como os dois binários compartilham `internal/`, o custo da separação
   é próximo de zero e a refatoração dolorosa nunca precisa acontecer.

5. **Fly.io por causa do WebSocket.** Vercel e Railway eram os candidatos do
   arranjo antigo (nunca implementados). O modelo serverless/edge da Vercel é
   hostil a conexões persistentes de longa duração; nenhum dos dois dá o controle
   de região e de processo sempre-ligado que o chat exige. O Fly.io sustenta
   conexão persistente e permite rodar em GRU (São Paulo), o que também reduz a
   latência para o público-alvo brasileiro.

6. **FCM e Codemagic são consequência, não escolha independente.** Expo
   Notifications e EAS Build só existem dentro do Expo. Saindo dele, o caminho
   nativo para push nas duas plataformas é FCM, e o CI que entende projeto Kotlin
   Multiplatform e publica nas duas lojas (com máquinas macOS) é o Codemagic.

> **Nota de honestidade sobre esta seção:** a decisão de trocar a stack foi
> tomada pelo dono do projeto e chegou aqui já fechada. As justificativas acima
> descrevem o que o novo arranjo objetivamente resolve em relação ao antigo. Se
> houve motivações adicionais — limites concretos encontrados no Expo durante o
> uso, por exemplo — elas devem ser acrescentadas aqui por quem tomou a decisão.

## Consequências

### Estruturais
- O repositório passa a ser um **monorepo** com `backend/`, `mobile/`,
  `contracts/`, `supabase/`, `docs/` e `.github/`.
- `backend/internal/` é organizado **por domínio** (`auth`, `user`, `matching`,
  `chat`, `group`, `notification`, `analytics`, `platform`), não por camada
  técnica.
- Surge a necessidade de um contrato explícito entre cliente e servidor
  (`contracts/openapi.yaml`), porque agora são duas linguagens compiladas
  separadamente. Ordem de trabalho obrigatória: contrato → Go → Kotlin.

### Sobre o ADR-001 (migrations e RLS)
O ADR-001 **continua integralmente válido** no que diz respeito a migrations:
arquivo aplicado não se edita, RLS na mesma migration que cria a tabela,
`DROP POLICY IF EXISTS` antes de `CREATE POLICY`, confirmação em `pg_tables`
antes do commit.

O que muda é o **peso** do RLS: ele deixa de ser o mecanismo de autorização do
produto e passa a ser defesa em profundidade, porque a autorização agora é feita
pelo Go. **Decisão pendente obrigatória (Fase 3, antes de criar as tabelas):** com
qual papel do Postgres o Go conecta e se o RLS é mantido ativo com propagação das
claims do JWT ou tratado apenas como proteção do Storage e de acessos diretos.
Ver [`docs/context/banco-de-dados.md`](../context/banco-de-dados.md), seção 3.

### Custo
O MVP **deixa de custar R$ 0/mês**. O Fly.io não tem free tier para contas novas:
dois apps (`api` e `ws`) em `shared-cpu-1x` custam da ordem de US$ 4–7/mês, mais
egress e volumes. É o preço de tirar a regra de negócio do cliente. Números
conferidos em 29/07/2026 — ver [`docs/context/stack.md`](../context/stack.md),
seção 4.

### O que foi jogado fora
- Setup de Node.js, Expo CLI, Expo Go e contas Expo/EAS (Fases 1 e 2 do roadmap
  precisaram ser reescritas e voltaram a "a fazer").
- `app.json`, `eas.json`, `package.json`, `package-lock.json`, `tsconfig.json`,
  `babel.config.js` e as pastas `app/` e `src/`.
- **Nenhuma linha de código de produto foi perdida** — todos os arquivos de
  `app/` e `src/` estavam vazios. O custo real da migração foi de *setup e
  documentação*, não de implementação.

## Dificuldades esperadas na transição

Em ordem aproximada de risco:

1. **Duas linguagens novas em produção ao mesmo tempo, em projeto solo.** Go e
   Kotlin Multiplatform simultaneamente é o maior risco do plano. É o mesmo tipo
   de risco que fez o ADR-002 descartar Flutter ("aprender linguagem nova +
   mobile ao mesmo tempo") — só que agora aceito conscientemente, com a
   atenuante de que Kotlin não é novo para o dev. **Mitigação sugerida:** fechar
   o backend de um domínio ponta a ponta (ex: perfil) antes de abrir frente no
   mobile, para não estar aprendendo as duas coisas no mesmo dia.

2. **O risco do iOS que o ADR-002 apontou não desapareceu.** O descarte do
   Kotlin Multiplatform em 29/06/2026 citava "ecossistema iOS ainda com
   fricção". Isso segue verdade em algum grau: bibliotecas KMP têm suporte a iOS
   desigual, e Compose Multiplatform no iOS é mais novo que no Android. O risco
   foi reavaliado e aceito, não refutado. **Consequência prática:** validar
   suporte a iOS de *cada* dependência antes de adotá-la, e testar no iOS cedo —
   não na véspera do lançamento.

3. **Desenvolver para iOS exige macOS.** O Codemagic resolve build, assinatura e
   publicação, mas **não** resolve o ciclo de desenvolvimento: rodar no
   simulador, debugar, ver a UI. Sem um Mac, o desenvolvimento iOS fica limitado
   a "compila no CI e testa no dispositivo". Isso precisa entrar no planejamento
   das Fases 4 e 5 como restrição real.

4. **Operação passa a ser nossa.** Deploy, logs, métricas, alerta de queda,
   rotação de credencial, atualização de dependência: tudo isso o Supabase
   absorvia. Agora há dois serviços para manter no ar. Backend próprio é
   trabalho recorrente, não trabalho de uma vez.

5. **O chat tem que ser construído, não configurado.** Reconexão com backoff,
   ordenação de mensagens, entrega de quem estava offline, presença,
   heartbeat/timeout, limite de mensagem: tudo que o Realtime dava pronto agora é
   código nosso. É a maior massa de trabalho novo da Fase 4.

6. **O hub de WebSocket em memória limita a uma instância de `cmd/ws`.** Com duas
   instâncias e sem fan-out, a mensagem só chega a quem estiver conectado na
   mesma instância. Aceitável no MVP, **bloqueante para escalar**: antes da
   segunda instância, escolher `LISTEN/NOTIFY` do Postgres ou Redis pub/sub e
   registrar como ADR.

7. **Fim do OTA update.** O Expo permitia corrigir a UI sem passar pela loja.
   Agora qualquer correção de interface exige build novo e revisão — o que muda
   o custo de errar e recomenda mais cuidado no que é liberado.

8. **Tipos duplicados em Go e Kotlin.** O mesmo modelo existe nas duas
   linguagens; divergência silenciosa é questão de tempo. O `contracts/` existe
   para conter isso, mas só funciona se a disciplina "contrato primeiro" for
   respeitada.

9. **`pgx` + pooler do Supabase.** Se for usado o pooler em modo transação, o
   cache de prepared statements do `pgx` precisa ser desligado, senão surgem
   erros intermitentes de "prepared statement already exists". Fácil de resolver,
   difícil de diagnosticar sem saber.

10. **Firebase é uma superfície nova.** Projeto Firebase, `google-services.json`
    (Android), `GoogleService-Info.plist` (iOS), service account no backend,
    certificado APNs para iOS. Nada difícil, mas é uma fila de configuração que
    não existia com Expo Notifications.

11. **Autorização precisa de teste explícito.** Com RLS como mecanismo principal,
    um erro de política aparecia como "não vejo o dado". Com autorização no Go,
    um `if` esquecido em um handler vaza dado silenciosamente. Os testes da Fase 5
    precisam cobrir autorização por endpoint, não só o caminho felizmente
    autenticado.

## Alternativas consideradas para o backend

| Alternativa | Motivo do descarte |
|---|---|
| Manter Supabase como backend completo (status quo) | Não resolve regra de negócio no cliente nem dá controle sobre o chat |
| Supabase Edge Functions | Resolveria parte da regra de negócio, mas não sustenta WebSocket persistente próprio, e amarraria a lógica ao runtime do Supabase |
| Node.js / Java+Spring no backend | Go entrega binário único, startup quase instantâneo e pouca memória — o que importa pagando por instância e mantendo muitas conexões abertas |
| API e WebSocket no mesmo processo | Junta perfis de carga opostos; deploy da API derrubaria as conversas |
| Vercel / Railway | Descartados especificamente por conexões WebSocket persistentes e falta de controle de região |
