# Fase 6 — Deploy e Lançamento

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [05 — Qualidade e Testes](./05-qualidade-e-testes.md)

## O que essa fase entrega

O Squadr publicado e instalável de verdade — Google Play primeiro (processo mais
simples e rápido), Apple App Store em seguida — **e o backend rodando no Fly.io**
— com o lançamento inicial dentro de uma comunidade já existente, conforme a
estratégia definida em [`produto.md`](../context/produto.md) (mitigação do
"problema do ovo e da galinha", seção 7).

> 🔄 **O que mudou com a troca de stack.** Esta fase agora tem **dois deploys**,
> não um: o backend (Fly.io) e o app (Codemagic → lojas). E o backend vem
> primeiro — sem ele, o app publicado não funciona. Todo o conteúdo de EAS Build
> foi removido; nada dele se aproveita.

---

## 1. Contas de desenvolvedor (só agora, não antes)

- [ ] Criar conta no **Google Play Console** (taxa única, ~R$ 130): https://play.google.com/console/signup
- [ ] Criar conta no **Apple Developer Program** (~R$ 580/ano): https://developer.apple.com/programs/enroll
- [ ] Reservar 1–2 dias de folga entre a criação da conta Apple e a submissão — a verificação da Apple pode levar até 48h

---

## 2. Deploy do backend no Fly.io (fazer **antes** do app)

- [ ] Autenticar:
  ```bash
  flyctl auth login
  ```
- [ ] Criar **dois apps separados** — a API e o WebSocket escalam e reiniciam
      independentemente (é o motivo de existirem dois binários):
  ```bash
  cd backend
  flyctl launch --name squadr-api --region gru --no-deploy
  flyctl launch --name squadr-ws  --region gru --no-deploy
  ```
  Cada um gera um arquivo de configuração próprio (ex: `fly.api.toml` e
  `fly.ws.toml`), apontando para o mesmo `Dockerfile` com um `target` diferente.
- [ ] Configurar os **secrets** em cada app (nunca no arquivo de config, nunca no Git):
  ```bash
  flyctl secrets set DATABASE_URL="..." SUPABASE_URL="..." \
    FIREBASE_SERVICE_ACCOUNT_JSON="..." POSTHOG_API_KEY="..." -a squadr-api
  ```
- [ ] Deploy dos dois:
  ```bash
  flyctl deploy --config fly.api.toml -a squadr-api
  flyctl deploy --config fly.ws.toml  -a squadr-ws
  ```
- [ ] Configurar o health check da API apontando para `/healthz`
- [ ] ⚠️ **Fixar `squadr-ws` em uma única máquina** e desligar o autoescalonamento
      dele. O hub de conexões é em memória: com duas instâncias, a mensagem só
      chega a quem estiver conectado na mesma. Ver
      [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), dificuldade #6
- [ ] Desligar o "auto stop" das máquinas do `squadr-ws` — máquina que dorme
      derruba conexão persistente
- [ ] Confirmar que o WebSocket funciona sobre **WSS** (TLS) no domínio do Fly.io
- [ ] Testar os dois endpoints públicos antes de mexer no app
- [ ] Anotar as URLs finais — elas vão para a configuração do app

---

## 3. Configuração do Codemagic

- [ ] Conectar o repositório do GitHub no Codemagic
- [ ] Criar o `codemagic.yaml` com os workflows:
  - `android-internal` — build de teste (usado na Fase 5)
  - `android-release` — `.aab` assinado para o Google Play
  - `ios-release` — `.ipa` assinado para a App Store
- [ ] Configurar a **assinatura Android**: gerar o keystore, subir no Codemagic
      como variável segura e **guardar o backup em lugar seguro** (perder o
      keystore significa nunca mais poder atualizar o app na Play Store)
- [ ] Configurar a **assinatura iOS**: chave da App Store Connect API, certificado
      de distribuição e provisioning profile (o Codemagic consegue gerenciar isso
      automaticamente com a chave da API)
- [ ] Definir as variáveis de ambiente do app no Codemagic: URL da API e do
      WebSocket (as do passo 2), `SUPABASE_URL`, `SUPABASE_PUBLISHABLE_KEY`, chave do
      PostHog
- [ ] Adicionar `google-services.json` e `GoogleService-Info.plist` como arquivos
      seguros no Codemagic (não commitados no repositório)
- [ ] Rodar um build de teste e confirmar que passa antes de tentar publicar

---

## 4. Ícone, splash screen e ficha da loja

- [ ] Ícone do app em alta resolução (1024×1024px), sem transparência para iOS
- [ ] Ícone adaptativo do Android configurado em `androidApp/`
- [ ] Splash screen configurada nas duas plataformas
- [ ] Screenshots do app em pelo menos 2 tamanhos de tela diferentes (obrigatório nas duas lojas)
- [ ] Texto curto (para a busca) e descrição longa da ficha da loja, em PT-BR, usando o posicionamento já definido: **"O app que conecta gamers brasileiros para jogar juntos"**
- [ ] Política de privacidade publicada (obrigatória em ambas as lojas — pode ser
      uma página simples hospedada gratuitamente, ex: GitHub Pages). Precisa
      mencionar: dados vindos do Discord, analytics (PostHog) e push (FCM)

---

## 5. Submissão às lojas

### Google Play
- [ ] Criar a ficha do app no Play Console
- [ ] Fazer upload do `.aab` (gerado pelo Codemagic) numa faixa de teste interno primeiro (não direto em produção)
- [ ] Testar essa versão de teste interno com os mesmos usuários convidados da Fase 5
- [ ] Preencher o questionário de **Segurança dos dados** (o Google exige declarar coleta de dados — analytics e conta Discord entram aqui)
- [ ] Promover para produção depois de validado

### Apple App Store
- [ ] Criar o app no App Store Connect
- [ ] Fazer upload do `.ipa` (o Codemagic publica direto, se configurado)
- [ ] Preencher o **App Privacy** (equivalente ao questionário do Google)
- [ ] Enviar para revisão da Apple — atenção aos motivos mais comuns de rejeição
      em apps sociais/matchmaking: política de privacidade insuficiente, ausência
      de mecanismo de report/bloqueio **visível** (o Squadr já tem isso, do Bloco 7
      da Fase 4 — só confirmar que está acessível na primeira sessão), ausência de
      moderação de conteúdo gerado por usuário
- [ ] Ter em mãos uma conta de teste com convite válido para a revisão da Apple —
      o revisor **não consegue passar do login** se o cadastro for só por convite.
      Isso é motivo garantido de rejeição se esquecido

---

## 6. Estratégia de lançamento (mitigação do risco "ovo e galinha")

Referência: [`produto.md`](../context/produto.md), seção 7 (riscos) e seção 8.1
(crescimento por convite).

- [ ] Escolher 1–2 comunidades Discord de jogos (Valorant, LoL ou CS2) já ativas, de preferência onde você (ou alguém do time) já tem alguma relação
- [ ] Combinar com os administradores da comunidade antes de divulgar — divulgação não combinada em servidor de terceiros é rota rápida para ban e má vontade
- [ ] Distribuir os primeiros convites (mecanismo da Fase 3, seção 9) dentro dessa comunidade, não abertamente
- [ ] Acompanhar de perto os primeiros dias: taxa de conclusão do onboarding, número de matches que geram sinal de sessão real, volume de reports
- [ ] Definir o gatilho concreto (número/critério, já deveria estar registrado desde a Fase 3) para abrir o cadastro ao público geral

---

## 7. Operação depois do lançamento (novo — backend é nosso agora)

Isto não existia na stack antiga, onde o Supabase cuidava de tudo:

- [ ] Acompanhar logs e métricas dos dois apps no Fly.io nos primeiros dias
- [ ] Definir como você descobre que o backend caiu (alerta, ou pelo menos uma
      checagem diária) — usuário reclamando não deve ser o primeiro aviso
- [ ] Confirmar que o número de conexões WebSocket abertas cabe na máquina
      contratada, e saber qual é o limite antes de bater nele
- [ ] Ter o procedimento de rollback claro: `flyctl releases` e voltar a versão
      anterior
- [ ] Acompanhar o custo real no painel do Fly.io no primeiro mês — a estimativa é
      de US$ 4–7/mês, mas egress e volume podem mudar isso

---

## Antes de avançar

- [ ] Backend rodando no Fly.io, com `squadr-ws` fixado em uma instância e WSS funcionando
- [ ] App aprovado e disponível na Google Play (pelo menos)
- [ ] App aprovado e disponível na Apple App Store, ou em revisão sem pendência bloqueante
- [ ] Backup do keystore Android guardado fora do repositório e fora da máquina
- [ ] Lançamento inicial feito dentro de comunidade combinada, não aberto
- [ ] Métricas básicas (PostHog) sendo observadas desde o primeiro dia, com eventos de app **e** de servidor chegando

---

## Depois desta fase

Não existe (ainda) um sub-roadmap de pós-lançamento. Quando o app estiver no ar, o
próximo passo é coletar feedback, revisar as métricas e planejar a **v1.1** —
cujo escopo já está esboçado em [`produto.md`](../context/produto.md), seção 3.
Se essa etapa ganhar corpo, ela vira o `07-pos-lancamento.md`.
