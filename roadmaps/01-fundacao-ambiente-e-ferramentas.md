# Fase 1 — Fundação: Ambiente e Ferramentas

**Status: 🟨 Em andamento**

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [00 — Ideação e Planejamento](./00-ideacao-e-planejamento.md)

## O que essa fase entrega

No fim dessa fase, seu computador tem **tudo instalado e configurado** para começar a escrever código do Squadr — mesmo que você nunca tenha programado antes. Isso inclui os programas necessários e as contas (gratuitas) nos serviços que o projeto usa.

> 💡 **Se você não tem experiência técnica:** siga a ordem exata desta lista.
> Cada item explica o "porquê", não só o "como" — isso ajuda a entender o que está fazendo mesmo sem background técnico.

---

## 1. Ferramentas a instalar no computador

### 1.1 Git
**O que é:** o programa que guarda o histórico de todas as versões do código do projeto — cada alteração salva vira um "commit", e é possível voltar no tempo se algo quebrar.

- [x] Baixar em: https://git-scm.com/downloads
- [x] Instalar com as opções padrão (não precisa mudar nada na instalação)
- [x] Confirmar que instalou corretamente, abrindo o terminal e rodando:
  ```bash
  git --version
  ```

### 1.2 Node.js (versão LTS)
**O que é:** o motor que executa o código JavaScript/TypeScript fora do navegador — é o que permite rodar o Expo, instalar pacotes e, depois, compilar o app.

- [x] Baixar a versão **LTS** (não a "Current") em: https://nodejs.org
- [x] Instalar com as opções padrão
- [x] Confirmar a instalação:
  ```bash
  node --version
  npm --version
  ```

### 1.3 Editor de código — VS Code
**O que é:** o programa onde você vai escrever e ler o código.

- [x] Baixar em: https://code.visualstudio.com
- [x] Instalar com as opções padrão
- [x] Extensões recomendadas para instalar dentro do VS Code (ícone de quadrados na barra lateral esquerda → buscar pelo nome → Instalar):
  - [x] **ES7+ React/Redux/React-Native snippets**
  - [x] **Prettier - Code formatter** (formata o código automaticamente)
  - [x] **ESLint** (avisa sobre erros comuns antes de você rodar o código)

### 1.4 Expo CLI
**O que é:** a ferramenta de linha de comando do Expo — cria o projeto, inicia o app em modo de desenvolvimento e gera o instalador final para as lojas.

- [x] Não precisa instalar globalmente; será usado via `npx` (explicado na Fase 2), que já vem com o Node.

### 1.5 App Expo Go
**O que é:** um app que você instala no seu celular, que permite testar o Squadr em desenvolvimento sem precisar compilar nada — basta escanear um QR code.

- [ ] Instalar **Expo Go** na Play Store (Android) ou App Store (iPhone)

### 1.6 Supabase CLI
**O que é:** a ferramenta de linha de comando que aplica as migrations (mudanças no banco de dados) no projeto Supabase.

- [x] Instalar via npm:
  ```bash
  npm install -g supabase
  ```
- [x] Confirmar:
  ```bash
  supabase --version
  ```

---

## 2. Contas a criar (todas com plano gratuito suficiente para o MVP)

| Conta | Para quê serve | Link |
|---|---|---|
| **GitHub** | Guardar o código do projeto na nuvem e colaborar | https://github.com/signup |
| **Expo (EAS)** | Compilar o app e gerenciar builds | https://expo.dev/signup |
| **Supabase** | Banco de dados, autenticação e realtime | https://supabase.com/dashboard |
| **Discord Developer Portal** | Criar a aplicação OAuth para o login "Entrar com Discord" | https://discord.com/developers/applications |
| **PostHog** | Analytics (saber como as pessoas usam o app) | https://posthog.com/signup |

- [x] Criar conta no GitHub
- [x] Criar conta no Expo
- [x] Criar conta no Supabase
- [ ] Criar uma aplicação no Discord Developer Portal
- [x] Criar conta no PostHog

> As contas de **Apple Developer Program** (pago, ~R$ 580/ano) e **Google Play Developer** (taxa única de ~R$ 130) só são necessárias na Fase 6 (Deploy). Não precisa criar agora.

---

## 3. Criando o projeto Expo

- [ ] Escolher/criar a pasta onde o projeto vai viver no computador (ex:`Documentos/projetos/`)
- [ ] Abrir o terminal nessa pasta e rodar:
  ```bash
  npx create-expo-app@latest squadr --template
  ```
  Quando perguntar o template, escolher **"Blank (TypeScript)"**.
- [ ] Entrar na pasta criada:
  ```bash
  cd squadr
  ```
- [ ] Testar que está tudo funcionando:
  ```bash
  npx expo start
  ```
  Vai aparecer um QR code no terminal. Escaneie com o app **Expo Go** no celular (Android: dentro do próprio app; iPhone: pela câmera nativa) — o app deve abrir mostrando a tela padrão do Expo.

✅ **Se o QR code abriu o app no celular, a fundação está funcionando.**

---

## 4. Conectando o projeto ao GitHub

- [ ] Criar um repositório novo no GitHub (vazio, sem README) chamado `squadr`
- [ ] Dentro da pasta do projeto:
  ```bash
  git init
  git add .
  git commit -m "chore: setup inicial do projeto Expo"
  git branch -M main
  git remote add origin https://github.com/SEU_USUARIO/squadr.git
  git push -u origin main
  ```
- [ ] Confirmar que o código apareceu no GitHub, atualizando a página do repositório no navegador

---

## Antes de avançar

- [ ] `git --version`, `node --version`, `npm --version` e `supabase --version` funcionam sem erro
- [ ] Expo Go instalado no celular e testado com sucesso
- [ ] Contas criadas: GitHub, Expo, Supabase, Discord Developer Portal, PostHog
- [ ] Projeto Expo criado, rodando localmente e já no GitHub

➡️ Próxima fase: [`02-estrutura-do-projeto.md`](./02-estrutura-do-projeto.md)