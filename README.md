# Squadr

> Conecte-se com gamers brasileiros e monte seu squad em menos de 3 minutos.

O Squadr é um app mobile que resolve um problema simples e real: encontrar alguém para jogar online no Brasil é lento e frustrante. Grupos de WhatsApp, canais LFG no Discord, posts no Reddit — tudo isso exige 15 a 30 minutos de ida e volta antes de confirmar alguém. O Squadr entrega isso em uma tela, com filtro automático por jogo, rank e horário.

---

## Status do projeto

🟡 Em desenvolvimento — configuração do ambiente e MVP inicial

---

## Stack

| Camada | Tecnologia |
|---|---|
| Mobile | React Native + Expo |
| Backend / Banco | Supabase (Postgres + Realtime + Auth) |
| Autenticação | Discord OAuth via Supabase Auth |
| Push Notifications | Expo Notifications |
| Analytics | PostHog |
| Gerenciamento de estado | Zustand |
| Monorepo | Turborepo |

---

## Estrutura do repositório

```
squadr/
├── apps/
│   └── mobile/               # App React Native + Expo
│       └── src/
│           ├── components/
│           ├── screens/
│           ├── hooks/
│           ├── services/      # Integração com Supabase
│           ├── store/         # Estado global (Zustand)
│           └── utils/
│
├── packages/
│   ├── supabase/              # Migrations, seed e tipos do banco
│   │   ├── migrations/
│   │   └── seed/
│   └── shared/                # Tipos e constantes compartilhados
│       ├── types/
│       └── constants/
│
├── docs/
│   ├── PROJECT_CONTEXT.md     # Fonte de verdade do projeto
│   └── ADR/                   # Decisões de arquitetura
│
└── .github/
    └── workflows/             # CI/CD (em breve)
```

---

## Pré-requisitos

- [Node.js](https://nodejs.org) v18 ou superior
- [Expo Go](https://expo.dev/go) instalado no celular
- Conta no [Supabase](https://supabase.com)
- Conta no [GitHub](https://github.com)

---

## Como rodar localmente

### 1. Clone o repositório

```bash
git clone https://github.com/Chicosk1/squadr.git
cd squadr
```

### 2. Instale as dependências

```bash
npm install
```

### 3. Configure as variáveis de ambiente

Copie o arquivo de exemplo e preencha com suas chaves do Supabase:

```bash
cp .env.example .env
cp apps/mobile/.env.example apps/mobile/.env
```

Edite o arquivo `apps/mobile/.env` com os valores do seu projeto Supabase:

```env
EXPO_PUBLIC_SUPABASE_URL=https://SEU_PROJETO.supabase.co
EXPO_PUBLIC_SUPABASE_ANON_KEY=sua_chave_anon_aqui
```

> ⚠️ Nunca suba arquivos `.env` para o GitHub. As chaves reais ficam apenas na sua máquina.

### 4. Rode o app

```bash
cd apps/mobile
npx expo start
```

Abra o **Expo Go** no celular e escaneie o QR code que aparecer no terminal. Seu celular e computador precisam estar na mesma rede Wi-Fi.

---

## Documentação

- [Contexto do projeto](./docs/squadr_context.md) — visão, funcionalidades, stack e decisões tomadas
- [ADR](./docs/ADR/) — registro de decisões de arquitetura

---

## Jogos suportados no MVP

- Valorant
- League of Legends
- CS2

---

## Roadmap

- [x] Estrutura do monorepo
- [x] Configuração do Supabase
- [ ] Tela de login com Discord OAuth
- [ ] Perfil do jogador
- [ ] Feed de discovery com filtros
- [ ] Sistema de like e match
- [ ] Chat em tempo real
- [ ] Notificações push
- [ ] Moderação básica (report e bloqueio)
- [ ] Publicação nas lojas (App Store + Play Store)

---

## Convenção de commits

Este projeto segue o padrão [Conventional Commits](https://www.conventionalcommits.org):

| Prefixo | Uso |
|---|---|
| `feat:` | Nova funcionalidade |
| `fix:` | Correção de bug |
| `chore:` | Configuração, dependências, setup |
| `docs:` | Documentação |
| `refactor:` | Refatoração sem mudança de comportamento |
| `test:` | Testes |

---

## Licença

MIT