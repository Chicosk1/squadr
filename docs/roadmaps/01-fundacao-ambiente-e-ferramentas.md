# Fase 1 — Fundação: Ambiente e Ferramentas

**Status: 🟨 Em andamento** — reaberta em 29/07/2026 pela troca de stack

> Voltar para o [Roadmap principal](./README.md) · Fase anterior: [00 — Ideação e Planejamento](./00-ideacao-e-planejamento.md)

## O que essa fase entrega

No fim dessa fase, seu computador tem **tudo instalado e configurado** para
começar a escrever código do Squadr. Isso inclui os programas necessários e as
contas nos serviços que o projeto usa.

> ⚠️ **Por que esta fase foi reaberta.** Ela estava concluída para a stack antiga
> (Node.js, Expo CLI, Expo Go, conta no Expo/EAS). Com a migração para Kotlin
> Multiplatform + Go ([ADR-003](../decisions/003-migracao-kmp-e-backend-go.md)),
> boa parte do ambiente mudou. **O que já estava instalado e continua servindo
> está marcado como feito** — o que aparece desmarcado é o que realmente falta.

### O que saiu da lista (não instale, não precisa mais)

| Removido | Por quê |
|---|---|
| Node.js / npm | Era o motor do Expo/React Native. Nenhuma parte da stack nova depende dele |
| Expo CLI | O Expo saiu do projeto |
| App **Expo Go** | Não existe equivalente: o teste passa a ser em emulador/dispositivo com build real |
| Extensões de React Native no editor | Sem React Native no projeto |
| Conta no Expo (EAS) | Substituída pelo Codemagic. Pode ser mantida, mas não é usada |

> A **conta no Supabase** e a **aplicação no Discord Developer Portal** continuam
> sendo usadas exatamente como antes. Não precisa refazer nada nelas.

---

## 1. Ferramentas a instalar no computador

> Verificação feita em 29/07/2026 nesta máquina de desenvolvimento. Em outra
> máquina, confira tudo de novo.

### 1.1 Git
**O que é:** o programa que guarda o histórico de todas as versões do código.

- [x] Instalado e funcionando (`git --version`)

### 1.2 JDK (Java Development Kit)
**O que é:** o kit que compila e roda código Kotlin/Java. O Gradle (ferramenta de
build do projeto mobile) precisa dele.

- [x] JDK instalado (`java -version`)
- [x] ⚠️ Se o Gradle reclamar de versão de JDK não suportada, instale também um
      JDK LTS anterior (ex: 21) e aponte o projeto para ele em **Android Studio →
      Settings → Build Tools → Gradle → Gradle JDK**.

### 1.3 Android Studio
**O que é:** o ambiente de desenvolvimento principal do projeto. É nele que a UI
compartilhada em Compose Multiplatform é escrita, e é ele que roda o app no
emulador Android.

- [x] Android Studio instalado (`C:\Program Files\Android\Android Studio`)
- [x] Android SDK presente (`%LOCALAPPDATA%\Android\Sdk`)
- [x] Instalar o plugin **Kotlin Multiplatform** (Settings → Plugins → buscar
      "Kotlin Multiplatform" → Install → reiniciar)
- [x] Criar um emulador Android (Device Manager → Create Device) — o app precisa
      de algum lugar para rodar
- [x] Confirmar que o emulador abre e roda

### 1.4 Go
**O que é:** a linguagem do backend. Compila os dois serviços (`cmd/api` e
`cmd/ws`) em binários únicos.

- [x] Baixar a versão estável mais recente em: https://go.dev/dl/
- [x] Instalar com as opções padrão
- [x] **Fechar e reabrir o terminal** (o instalador altera o `PATH`)
- [x] Confirmar (`go version`)
- [x] Anotar a versão que apareceu — ela vai ser usada no `go.mod` (Fase 2) e na
      variável `GO_VERSION` do `.env`

### 1.5 Editor para o backend Go
**O que é:** o Android Studio é ótimo para Kotlin, mas não para Go. Vale ter um
editor separado para a pasta `backend/`.

- [x] **VS Code** (https://code.visualstudio.com) + extensão oficial **Go**
      (publisher: Go Team at Google)

### 1.6 Docker Desktop
**O que é:** roda os serviços em contêiner. Usado pelo `docker-compose.yml` (sobe
`api` e `ws` localmente) e pelo Supabase CLI (que sobe o Postgres local).

- [x] Instalado — detectado **Docker 29.4.3** (`docker --version`)
- [x] Confirmar que o Docker Desktop está **rodando** antes de usar (o comando
      funciona, mas falha se o serviço estiver parado)

### 1.7 Supabase CLI
**O que é:** aplica as migrations (mudanças no banco) e sobe uma cópia local do
Supabase para desenvolvimento.

- [x] Instalado (`supabase --version`)

### 1.8 flyctl (CLI do Fly.io)
**O que é:** a ferramenta que faz o deploy do backend no Fly.io.

- [x] Instalar seguindo https://fly.io/docs/flyctl/install/
- [x] Confirmar (`flyctl version`)

### 1.9 Xcode — só em macOS
**O que é:** necessário para compilar, rodar no simulador e assinar o app iOS.

- [ ] **Se você tem um Mac:** instalar o Xcode pela App Store e rodar
      `xcode-select --install`
- [ ] **Se você não tem um Mac:** não há o que instalar. O Codemagic gera e
      publica o build iOS em máquinas macOS na nuvem (Fase 6). A limitação real é
      que **você não consegue rodar nem debugar o app iOS localmente** — isso
      precisa estar no seu planejamento das Fases 4 e 5. Ver
      [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md), dificuldade #3

---

## 2. Contas a criar

| Conta | Para quê serve | Link |
|---|---|---|---|
| **GitHub** | Guardar o código e rodar o CI do backend | https://github.com/signup |
| **Supabase** | Postgres gerenciado + Storage + Auth | https://supabase.com/dashboard |
| **Discord Developer Portal** | Aplicação OAuth do "Entrar com Discord" | https://discord.com/developers/applications |
| **PostHog** | Analytics de app e de servidor | https://posthog.com/signup |
| **Firebase** | Push notifications (FCM) | https://console.firebase.google.com |
| **Codemagic** | Build, assinatura e publicação do app | https://codemagic.io/signup |
| **Fly.io** | Hospedagem do backend Go | https://fly.io/app/sign-up |

- [x] Criar conta no GitHub
- [x] Criar conta no Supabase
- [x] Criar uma aplicação no Discord Developer Portal
- [x] Criar conta no PostHog
- [x] Criar projeto no **Firebase** e adicionar os dois apps (Android e iOS) — os
      arquivos `google-services.json` e `GoogleService-Info.plist` saem daqui, e
      a **service account** que o backend usa também
- [x] Criar conta no **Codemagic** e conectar ao repositório do GitHub
- [x] Criar conta no **Fly.io** — ⚠️ **exige cartão de crédito**: não há free
      tier para contas novas. O custo estimado é de ~US$ 4–7/mês para os dois
      serviços. Ver [`stack.md`](../context/stack.md), seção 4

> As contas de **Apple Developer Program** (~R$ 580/ano) e **Google Play
> Developer** (taxa única de ~R$ 130) só são necessárias na Fase 6 (Deploy). Não
> precisa criar agora.

---

## 3. Validando que a fundação funciona

Não crie os projetos aqui — isso é a Fase 2. Aqui só confirmamos que as
ferramentas respondem:

- [x] `go version` retorna uma versão
- [x] `docker ps` roda sem erro (com o Docker Desktop aberto)
- [x] `supabase --version` retorna uma versão
- [x] `flyctl version` retorna uma versão
- [x] Android Studio abre, reconhece o SDK e o emulador liga

---

## 4. Repositório no GitHub

- [x] Repositório `squadr` criado em https://github.com/Chicosk1/squadr
- [x] Projeto local vinculado ao remoto (`git remote -v`)

---

## Antes de avançar

- [x] Git, JDK, Android Studio, Docker e Supabase CLI instalados e respondendo
- [x] Go instalado, com a versão anotada
- [x] flyctl instalado
- [x] Plugin Kotlin Multiplatform instalado no Android Studio, com emulador criado
- [x] Contas criadas: Firebase (com projeto e os dois apps), Codemagic e Fly.io
- [x] Ciente da restrição de iOS sem macOS (item 1.9) e do custo mensal do Fly.io

➡️ Próxima fase: [`02-estrutura-do-projeto.md`](./02-estrutura-do-projeto.md)
