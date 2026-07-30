# Roadmap — Squadr

> App de matchmaking para gamers brasileiros. Este é o roadmap principal do projeto: a visão geral de todas as fases, da ideia até o app publicado nas lojas. Cada fase tem um **sub-roadmap** próprio, com o passo a passo detalhado.
>
> **Como usar este documento:** marque as caixas (`- [ ]` → `- [x]`) conforme for avançando. Ele foi feito para viver dentro do repositório do projeto e ser atualizado ao longo do desenvolvimento.
>
> ⚠️ **Atualização de 29/07/2026 — troca de stack.** O projeto migrou de React Native + Expo (com Supabase como backend completo) para **Kotlin Multiplatform + Compose Multiplatform no mobile e backend próprio em Go**. Ver [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md). Todas as fases abaixo foram reescritas para a stack nova, e **as Fases 1 e 2 voltaram a ter itens pendentes** — o ambiente que estava pronto era o do Expo.

---

## 📍 Onde este arquivo mora no projeto

```
squadr/
├── docs/
│   ├── roadmaps/
│   │   ├── README.md                    ← você está aqui
│   │   ├── 00-ideacao-e-planejamento.md
│   │   ├── 01-fundacao-ambiente-e-ferramentas.md
│   │   ├── 02-estrutura-do-projeto.md
│   │   ├── 03-backend-banco-e-autenticacao.md
│   │   ├── 04-desenvolvimento-mvp.md
│   │   ├── 05-qualidade-e-testes.md
│   │   └── 06-deploy-e-lancamento.md
│   ├── context/                         ← produto, stack, arquitetura, banco
│   └── decisions/                       ← ADRs
├── backend/     mobile/     supabase/     contracts/
```

> 💡 Toda decisão técnica importante deve virar um arquivo novo em
> [`docs/decisions/`](../decisions/), numerado em sequência. Isso já está definido
> no ADR-001 e vale para o projeto todo, não só banco de dados.

---

## 🗺️ Visão geral das fases

| # | Fase | Status | Sub-roadmap |
|---|---|---|---|
| 0 | Ideação e Planejamento | ✅ Concluída (stack revisada em 29/07/2026) | [`00-ideacao-e-planejamento.md`](./00-ideacao-e-planejamento.md) |
| 1 | Fundação: Ambiente e Ferramentas | 🟨 Em andamento — reaberta pela troca de stack | [`01-fundacao-ambiente-e-ferramentas.md`](./01-fundacao-ambiente-e-ferramentas.md) |
| 2 | Estrutura do Projeto | 🟨 Em andamento — pastas prontas, projetos ainda não gerados | [`02-estrutura-do-projeto.md`](./02-estrutura-do-projeto.md) |
| 3 | Backend: Banco, API e Autenticação | 🟨 Em andamento — Supabase criado, schema pendente | [`03-backend-banco-e-autenticacao.md`](./03-backend-banco-e-autenticacao.md) |
| 4 | Desenvolvimento do MVP | ⬜ A fazer | [`04-desenvolvimento-mvp.md`](./04-desenvolvimento-mvp.md) |
| 5 | Qualidade e Testes | ⬜ A fazer | [`05-qualidade-e-testes.md`](./05-qualidade-e-testes.md) |
| 6 | Deploy e Lançamento | ⬜ A fazer | [`06-deploy-e-lancamento.md`](./06-deploy-e-lancamento.md) |

> Atualize a coluna **Status** deste arquivo conforme as fases avançam:
> `⬜ A fazer` → `🟨 Em andamento` → `✅ Concluída`.

---

## 🧭 A lógica das fases (por que essa ordem)

```
IDEIA                    FUNDAÇÃO                  CONSTRUÇÃO                 LANÇAMENTO
  │                         │                          │                          │
  ▼                         ▼                          ▼                          ▼
0. O que estamos  →  1. Ambiente e        →  3. Banco + API +   →  5. Testar   →  6. Publicar
   construindo e        ferramentas             login (Discord)      tudo antes     (lojas +
   para quem        →  2. Organização      →  4. Construir as       de lançar      Fly.io)
                          do monorepo            telas e endpoints
```

A ideia central: **você não escreve tela** antes de ter um endpoint que responda,
e não escreve endpoint antes de saber onde os dados moram e onde cada arquivo do
projeto vai viver.

O que a stack nova acrescenta a essa lógica: agora existem **duas frentes** em
cada funcionalidade (backend Go e app Kotlin). A ordem dentro de cada bloco da
Fase 4 é sempre a mesma — **contrato → backend → app** — porque testar um
endpoint com `curl` é muito mais barato que descobrir o erro através da UI.

---

## ✅ Checklist rápido (visão de altíssimo nível)

- [x] Definir visão de produto, personas, funcionalidades do MVP e stack técnica
- [x] Revisar a stack e registrar a migração como ADR (ADR-002/ADR-003)
- [x] Criar a estrutura de pastas do monorepo
- [ ] Preparar ambiente de desenvolvimento da stack nova (Go, Android Studio/KMP, Firebase, Fly.io, Codemagic)
- [ ] Gerar os projetos (módulo Go e projeto Kotlin Multiplatform) dentro da estrutura
- [ ] Configurar banco de dados, autorização e autenticação (Discord OAuth + validação de JWT no Go)
- [ ] Construir os endpoints e as telas do MVP
- [ ] Testar tudo (manual + automatizado nos pontos de risco)
- [ ] Publicar backend no Fly.io e app nas lojas (Google Play e Apple App Store)
- [ ] Coletar feedback e planejar v1.1

---

## 📚 Referências do projeto

- Visão de produto: [`docs/context/produto.md`](../context/produto.md)
- Stack e custos: [`docs/context/stack.md`](../context/stack.md)
- Arquitetura: [`docs/context/arquitetura.md`](../context/arquitetura.md)
- Banco de dados: [`docs/context/banco-de-dados.md`](../context/banco-de-dados.md)
- Decisões técnicas (ADRs): [`docs/decisions/`](../decisions/)
- Contrato da API: [`contracts/`](../../contracts/)
