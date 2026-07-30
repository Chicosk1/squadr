# Contexto do Squadr

Esta pasta é a **fonte de verdade** do projeto: o que estamos construindo, para
quem, com qual tecnologia e como as peças se encaixam.

| Documento | O que responde |
|---|---|
| [`produto.md`](./produto.md) | Visão, personas, escopo do MVP, fluxos de onboarding/match, diferenciais, riscos e sistema de confiança |
| [`stack.md`](./stack.md) | Quais tecnologias usamos, por quê, o que foi descartado e quanto custa |
| [`arquitetura.md`](./arquitetura.md) | Como o monorepo é organizado, como mobile e backend conversam, fluxo de autenticação e de chat |
| [`banco-de-dados.md`](./banco-de-dados.md) | Tabelas, como o Go acessa o Postgres e onde a autorização mora |

## Relação com as outras pastas de `docs/`

- **`docs/decisions/`** — ADRs. Toda decisão técnica relevante vira um arquivo
  numerado lá. Quando uma decisão muda, o ADR antigo é marcado como
  supersedido, **nunca apagado**. Este contexto reflete sempre o estado atual;
  o histórico de *por que chegamos aqui* vive nos ADRs.
- **`docs/roadmaps/`** — o passo a passo de execução, fase por fase. Se um
  documento de contexto e um roadmap divergirem, o contexto ganha — o roadmap
  é que deve ser corrigido.

> ⚠️ A stack deste projeto **mudou** em 29/07/2026 (React Native/Expo +
> Supabase-como-backend → Kotlin Multiplatform + backend Go). Todo documento
> desta pasta já reflete a stack nova. O racional da stack antiga está
> preservado no [ADR-002](../decisions/002-stack-mobile-react-native-expo.md) e
> o motivo da troca no [ADR-003](../decisions/003-migracao-kmp-e-backend-go.md).
