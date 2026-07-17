# Fase 6 — Deploy e Lançamento

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [05 — Qualidade e Testes](./05-qualidade-e-testes.md)

## O que essa fase entrega

O Squadr publicado e instalável de verdade — Google Play primeiro
(processo mais simples e rápido), Apple App Store em seguida — e o
lançamento inicial dentro de uma comunidade já existente, conforme a
estratégia definida em `docs/contexto-produto.md` (mitigação do "problema
do ovo e da galinha", seção 8).

---

## 1. Contas de desenvolvedor (só agora, não antes)

- [ ] Criar conta no **Google Play Console** (taxa única, ~R$ 130): https://play.google.com/console/signup
- [ ] Criar conta no **Apple Developer Program** (~R$ 580/ano): https://developer.apple.com/programs/enroll
- [ ] Reservar 1–2 dias de folga entre a criação da conta Apple e a submissão — a verificação da Apple pode levar até 48h

---

## 2. Configuração do EAS Build

- [ ] Autenticar o projeto com o EAS:
  ```bash
  eas login
  eas build:configure
  ```
  Isso cria o arquivo `eas.json` na raiz do projeto.
- [ ] Definir os perfis de build em `eas.json`:
  - `development` — para testar com Expo Dev Client
  - `preview` — build interno de teste (usado na Fase 5)
  - `production` — build final para as lojas
- [ ] Preencher em `app.json`/`app.config.ts`: nome do app, ícone, splash
      screen, identificador único (`com.seudominio.squadr`), versão

---

## 3. Ícone, splash screen e ficha da loja

- [ ] Ícone do app em alta resolução (1024×1024px), sem transparência para iOS
- [ ] Splash screen configurada
- [ ] Screenshots do app em pelo menos 2 tamanhos de tela diferentes (obrigatório nas duas lojas)
- [ ] Texto curto (para a busca) e descrição longa da ficha da loja, em PT-BR, usando o posicionamento já definido: **"O app que conecta gamers brasileiros para jogar juntos"**
- [ ] Política de privacidade publicada (obrigatória em ambas as lojas — pode ser uma página simples hospedada gratuitamente, ex: GitHub Pages)

---

## 4. Build de produção

- [ ] Gerar o build final:
  ```bash
  eas build --profile production --platform android
  eas build --profile production --platform ios
  ```
- [ ] Baixar os artefatos gerados (`.aab` para Android, `.ipa` para iOS)

---

## 5. Submissão às lojas

### Google Play
- [ ] Criar a ficha do app no Play Console
- [ ] Fazer upload do `.aab` numa faixa de teste interno primeiro (não direto em produção)
- [ ] Testar essa versão de teste interno com os mesmos usuários convidados da Fase 5
- [ ] Promover para produção depois de validado

### Apple App Store
- [ ] Criar o app no App Store Connect
- [ ] Fazer upload do `.ipa` (via EAS Submit ou Transporter)
- [ ] Enviar para revisão da Apple — atenção aos motivos mais comuns de
      rejeição em apps sociais/matchmaking: política de privacidade
      insuficiente, ausência de mecanismo de report/bloqueio visível (o
      Squadr já tem isso, do Bloco 7 da Fase 4 — só confirmar que está
      visível na primeira sessão), ausência de moderação de conteúdo
      gerado por usuário

---

## 6. Estratégia de lançamento (mitigação do risco "ovo e galinha")

Referência: `docs/contexto-produto.md`, seção 8 (riscos) e seção 9.1
(crescimento por convite).

- [ ] Escolher 1–2 comunidades Discord de jogos (Valorant, LoL ou CS2) já
      ativas, de preferência onde você (ou alguém do time) já tem alguma
      relação
- [ ] Combinar com os administradores da comunidade antes de divulgar —
      divulgação não combinada em servidor de terceiros é rota rápida para
      ban e má vontade
- [ ] Distribuir os primeiros convites (mecanismo da Fase 3, seção 6) dentro
      dessa comunidade, não abertamente
- [ ] Acompanhar de perto os primeiros dias: taxa de conclusão do
      onboarding, número de matches que geram sinal de sessão real, volume
      de reports
- [ ] Definir o gatilho concreto (número/critério, já deveria estar
      registrado desde a Fase 3) para abrir o cadastro ao público geral

---

## Antes de avançar

- [ ] App aprovado e disponível na Google Play (pelo menos)
- [ ] App aprovado e disponível na Apple App Store, ou em revisão sem
      pendência bloqueante
- [ ] Lançamento inicial feito dentro de comunidade combinada, não aberto
- [ ] Métricas básicas (PostHog) sendo observadas desde o primeiro dia

➡️ Próxima fase: [`07-pos-lancamento.md`](./07-pos-lancamento.md)