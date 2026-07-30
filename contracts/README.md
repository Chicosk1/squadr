# contracts — o contrato da API

[`openapi.yaml`](./openapi.yaml) é a **fonte única de verdade** da API REST do
Squadr.

## Por que isso existe

Cliente (Kotlin) e servidor (Go) são linguagens diferentes, compiladas
separadamente. Sem um contrato explícito no meio, a divergência entre os dois só
aparece em runtime — no dispositivo do usuário. Este arquivo é o árbitro.

Na stack antiga isso não era necessário: o app falava direto com o Supabase e o
"contrato" era o schema do banco. Ver
[ADR-003](../docs/decisions/003-migracao-kmp-e-backend-go.md).

## A ordem de trabalho

```
1. alterar contracts/openapi.yaml
2. implementar/ajustar o handler em backend/internal/<domínio>
3. ajustar o cliente em mobile/shared/src/commonMain/.../data/remote
```

**Nunca o inverso.** Endpoint que existe no Go e não está aqui é bug de processo,
não atalho.

## Regras

- Um endpoint não documentado aqui não existe.
- Formato de **erro** é o mesmo em todos os endpoints (`components/schemas/Error`).
- Toda rota exige `Authorization: Bearer <JWT do Supabase>`, exceto `/healthz`.
- Listagens são paginadas — não existe endpoint que devolva "tudo".
- Campos que o cliente não pode alterar (`created_at`, `slots_filled`,
  `profile_complete`) são `readOnly`.

## Estado

**Esqueleto.** As rotas e os schemas principais estão desenhados, mas os detalhes
(paginação, campos exatos, códigos de erro por rota) são fechados na Fase 3 do
roadmap, antes da implementação. Ver
[`docs/roadmaps/03-backend-banco-e-autenticacao.md`](../docs/roadmaps/03-backend-banco-e-autenticacao.md),
seção 6.
