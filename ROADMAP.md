# Roadmap — Squadr

> App de matchmaking para gamers brasileiros. Este é o roadmap principal do projeto: a visão geral de todas as fases, da ideia até o app publicado nas lojas. Cada fase tem um **sub-roadmap** próprio, com o passo a passo detalhado.
>
> **Como usar este documento:** marque as caixas (`- [ ]` → `- [x]`) conforme for avançando. Ele foi feito para viver dentro do repositório do projeto e ser atualizado ao longo do desenvolvimento.

---

## 📍 Onde este arquivo mora no projeto

```
squadr/
├── ROADMAP.md                  ← você está aqui
├── roadmap/                    ← os sub-roadmaps, um por fase
│   ├── 00-ideacao-e-planejamento.md
│   ├── 01-fundacao-ambiente-e-ferramentas.md
│   ├── 02-estrutura-do-projeto.md
│   ├── 03-backend-banco-e-autenticacao.md
│   ├── 04-desenvolvimento-mvp.md
│   ├── 05-qualidade-e-testes.md
│   ├── 06-deploy-e-lancamento.md
├── docs/
│   ├── contexto-produto.md
│   └── decisions/              ← ADRs
```

> 💡 Toda decisão técnica importante deve virar um arquivo novo em `docs/decisions/`, numerado em sequência. Isso já está definido no ADR-001 e vale para o projeto todo, não só banco de dados.

---

## 🗺️ Visão geral das fases

| # | Fase | Status | Sub-roadmap |
|---|---|---|---|
| 0 | Ideação e Planejamento | ✅ Concluída | [`00-ideacao-e-planejamento.md`](./roadmap/00-ideacao-e-planejamento.md) |
| 1 | Fundação: Ambiente e Ferramentas | ⬜ A fazer | [`01-fundacao-ambiente-e-ferramentas.md`](./roadmap/01-fundacao-ambiente-e-ferramentas.md) |
| 2 | Estrutura do Projeto | ⬜ A fazer | [`02-estrutura-do-projeto.md`](./roadmap/02-estrutura-do-projeto.md) |
| 3 | Backend: Banco de Dados & Autenticação | ⬜ A fazer | [`03-backend-banco-e-autenticacao.md`](./roadmap/03-backend-banco-e-autenticacao.md) |
| 4 | Desenvolvimento do MVP | ⬜ A fazer | [`04-desenvolvimento-mvp.md`](./roadmap/04-desenvolvimento-mvp.md) |
| 5 | Qualidade e Testes | ⬜ A fazer | [`05-qualidade-e-testes.md`](./roadmap/05-qualidade-e-testes.md) |
| 6 | Deploy e Lançamento | ⬜ A fazer | [`06-deploy-e-lancamento.md`](./roadmap/06-deploy-e-lancamento.md) |

> Atualize a coluna **Status** deste arquivo conforme as fases avançam:
> `⬜ A fazer` → `🟨 Em andamento` → `✅ Concluída`.

---

## 🧭 A lógica das fases (por que essa ordem)

```
IDEIA                         FUNDAÇÃO                      CONSTRUÇÃO                    LANÇAMENTO
  │                              │                              │                              │
  ▼                              ▼                              ▼                              ▼
0. O que estamos    →   1. Ambiente e        →   3. Banco de dados   →   5. Testar      →   6. Publicar
   construindo e           ferramentas             e login (Discord)       tudo antes          nas lojas
   para quem          →   2. Organização      →   4. Construir as          de lançar     →   7. Iterar
                            das pastas               telas e funções                          
```

A ideia central: **você não escreve código** de tela antes de ter onde guardar os dados, e não guarda dados antes de saber onde cada arquivo do projeto vai morar.

---

## ✅ Checklist rápido (visão de altíssimo nível)

- [x] Definir visão de produto, personas, funcionalidades do MVP e stack técnica
- [ ] Preparar ambiente de desenvolvimento e instalar ferramentas
- [ ] Criar a estrutura de pastas do projeto
- [ ] Configurar banco de dados, RLS e autenticação (Discord OAuth)
- [ ] Construir as telas e funcionalidades do MVP
- [ ] Testar tudo (manual + automatizado básico)
- [ ] Publicar nas lojas (Google Play e Apple App Store)
- [ ] Coletar feedback e planejar v1.1

---

## 📚 Referências do projeto

- Visão completa de produto: [`docs/contexto-produto.md`](./docs/contexto-produto.md)
- Decisões técnicas (ADRs): [`docs/decisions/`](./docs/decisions/)