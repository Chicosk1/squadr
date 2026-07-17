# Fase 2 — Estrutura do Projeto (pastas e arquivos)

**Status: ⬜ A fazer**

> Voltar para o [Roadmap principal](../ROADMAP.md) · Fase anterior: [01 — Fundação: Ambiente e Ferramentas](./01-fundacao-ambiente-e-ferramentas.md)

## O que essa fase entrega

Uma organização de pastas fixa, que todo o resto do projeto vai seguir.
Definir isso **antes** de escrever telas evita o problema mais comum em
projetos que crescem sem planejamento: arquivos jogados em qualquer lugar,
sem critério, até ficar difícil achar qualquer coisa.

> 💡 Você não precisa criar todas as pastas agora, vazias. Crie a estrutura
> de nível 1 (as pastas principais) agora; as subpastas nascem naturalmente
> quando o primeiro arquivo de cada tipo for criado, na Fase 4.

---

## 1. Estrutura completa do projeto

```
squadr/
├── ROADMAP.md                      # roadmap principal (já existe)
├── roadmap/                        # sub-roadmaps (já existe)
│
├── app/                             # TELAS e NAVEGAÇÃO (Expo Router)
│   ├── (auth)/                     # telas antes do login
│   │   └── login.tsx
│   ├── (onboarding)/                # seleção de jogos, perfil opcional
│   │   ├── selecionar-jogos.tsx
│   │   └── completar-perfil.tsx
│   ├── (tabs)/                      # navegação principal (abas)
│   │   ├── feed.tsx
│   │   ├── matches.tsx
│   │   ├── squads.tsx
│   │   └── perfil.tsx
│   ├── chat/
│   │   └── [matchId].tsx
│   ├── _layout.tsx
│   └── index.tsx
│
├── src/                             # TODO O CÓDIGO QUE NÃO É TELA
│   ├── components/                  # peças de UI reutilizáveis
│   │   ├── ui/                      # botão, card, avatar, badge...
│   │   └── feature/                 # componentes específicos de uma função
│   │       ├── PlayerCard.tsx
│   │       ├── MatchAnimation.tsx
│   │       └── SquadSlotList.tsx
│   ├── services/                    # comunicação com o mundo de fora
│   │   ├── supabase/
│   │   │   ├── client.ts            # inicialização do Supabase
│   │   │   ├── profiles.ts
│   │   │   ├── swipes.ts
│   │   │   ├── matches.ts
│   │   │   └── squads.ts
│   │   ├── discord-oauth.ts
│   │   └── posthog.ts
│   ├── hooks/                       # lógica reutilizável de estado/efeito
│   │   ├── useAuth.ts
│   │   ├── useFeed.ts
│   │   └── useRealtimeChat.ts
│   ├── types/                       # tipos TypeScript compartilhados
│   │   ├── profile.ts
│   │   ├── match.ts
│   │   └── squad.ts
│   ├── constants/                   # valores fixos (cores, jogos suportados)
│   │   ├── theme.ts
│   │   └── games.ts
│   └── utils/                       # funções puras de apoio
│       └── formatters.ts
│
├── supabase/                        # BANCO DE DADOS
│   ├── migrations/                  # histórico de mudanças no schema
│   │   ├── 002_create_profiles.sql
│   │   ├── ...
│   │   └── 010_enable_rls.sql
│   └── config.toml
│
├── docs/                             # DOCUMENTAÇÃO DO PROJETO
│   ├── contexto-produto.md
│   └── decisions/                   # ADRs — decisões técnicas registradas
│       └── 001-migrations-e-rls-supabase.md
│
├── assets/                           # imagens, ícones, fontes
│   ├── images/
│   └── fonts/
│
├── .env.example                      # modelo das variáveis de ambiente
├── .env                               # variáveis reais (NUNCA vai pro Git)
├── .gitignore
├── app.json                           # configuração do app Expo
├── eas.json                           # configuração de build (Fase 6)
├── package.json
├── tsconfig.json
└── babel.config.js
```

---

## 2. Por que cada pasta existe (o critério de organização)

| Pasta | Critério para algo entrar aqui |
|---|---|
| `app/` | É uma **tela** que o usuário vê e navega até ela. Se o arquivo representa uma rota, ele mora aqui. |
| `src/components/ui/` | Componente visual **genérico**, sem lógica de negócio — funcionaria em qualquer app (botão, input, avatar). |
| `src/components/feature/` | Componente visual **específico do Squadr** — não faria sentido fora deste app (card de jogador, animação de match). |
| `src/services/` | Qualquer código que **fala com algo de fora do app** — banco de dados, API externa, SDK de terceiro. |
| `src/hooks/` | Lógica de estado/comportamento que **se repete em mais de uma tela**. |
| `src/types/` | Formatos de dados usados em mais de um lugar — evita redefinir o mesmo tipo em arquivos diferentes. |
| `src/constants/` | Valores que não mudam em tempo de execução (cores do tema, lista de jogos suportados). |
| `src/utils/` | Funções pequenas, sem estado, que só transformam um dado em outro (ex: formatar data). |
| `supabase/migrations/` | Cada arquivo é uma mudança no banco, na ordem em que foi aplicada. **Nunca edite um arquivo antigo** — ver ADR-001. |
| `docs/decisions/` | Toda decisão técnica relevante ganha um arquivo numerado aqui, seguindo o padrão do ADR-001. |

## 3. Regra de nomenclatura (convenção usada no projeto)

- [x] Telas e componentes React: `PascalCase.tsx` (ex: `PlayerCard.tsx`)
- [x] Hooks: `useAlgumaCoisa.ts`, sempre começando com `use`
- [x] Serviços e utilitários: `kebab-case.ts` ou `camelCase.ts` (ex:
      `discord-oauth.ts`)
- [x] Migrations do Supabase: `NNN_descricao_curta.sql`, número sequencial de
      3 dígitos (ex: `011_create_commendations_table.sql`)
- [x] ADRs: `NNN-titulo-curto.md`, mesmo padrão de numeração

---

## 4. Criando a estrutura no computador

- [x] Dentro da pasta `squadr/`, criar as pastas de nível 1:
  ```bash
  mkdir -p src/components/ui src/components/feature src/services/supabase src/hooks src/types src/constants src/utils
  mkdir -p supabase/migrations
  mkdir -p docs/decisions
  mkdir -p assets/images assets/fonts
  ```
- [x] Mover o arquivo de contexto de produto (o markdown com a visão do
      produto) para `docs/contexto-produto.md`
- [x] Mover o ADR-001 para `docs/decisions/001-migrations-e-rls-supabase.md`
- [x] Criar o arquivo `.env.example` (vazio por enquanto, será preenchido na Fase 3):
  ```bash
  touch .env.example
  ```
- [x] Confirmar que `.gitignore` já inclui `.env`, `node_modules` e pastas de
      build do Expo (o template do `create-expo-app` já gera isso por
      padrão — só confirme abrindo o arquivo)

---

## Antes de avançar

- [x] Estrutura de pastas de nível 1 criada e commitada no Git
- [x] Documentos existentes (`contexto-produto.md`, ADR-001) movidos para
      dentro de `docs/`
- [x] Todos entendem o critério de "o que vai em cada pasta" (tabela da
      seção 2) — isso evita decisões diferentes por quem trabalhar no
      projeto depois

➡️ Próxima fase: [`03-backend-banco-e-autenticacao.md`](./03-backend-banco-e-autenticacao.md)