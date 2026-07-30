# Decisões de Arquitetura (ADRs)

Cada decisão técnica relevante do Squadr é um arquivo numerado nesta pasta.

## Índice

| # | Decisão | Status | Data |
|---|---|---|---|
| [001](./001-migrations-e-rls-supabase.md) | Processo de migrations e RLS no Supabase | ✅ Aceito (emendado pelo 003) | 30/06/2026 |
| [002](./002-stack-mobile-react-native-expo.md) | Stack inicial: React Native + Expo com Supabase como backend completo | ⛔ Supersedido pelo 003 | 29/06/2026 |
| [003](./003-migracao-kmp-e-backend-go.md) | Migração para Kotlin Multiplatform no mobile e backend próprio em Go | ✅ Aceito — **stack vigente** | 29/07/2026 |

## Convenções

- **Nome do arquivo:** `NNN-titulo-curto.md`, número sequencial de 3 dígitos.
  O número é atribuído na ordem em que o *registro* é criado, o que pode não
  coincidir com a ordem das datas de decisão (é o caso do 002, registrado
  retroativamente).
- **Status possíveis:** `Proposto`, `Aceito`, `Supersedido por ADR-NNN`,
  `Emendado por ADR-NNN`.
- **Decisão antiga nunca é apagada.** Quando é substituída, ganha um bloco de
  supersedência no topo apontando para quem a substituiu, e o arquivo continua no
  repositório. O histórico de *por que chegamos aqui* vale tanto quanto o estado
  atual — foi o que permitiu, no ADR-003, reconhecer que o risco do iOS apontado
  no ADR-002 não desapareceu, apenas foi aceito.
- **Emenda ≠ supersedência.** Se uma decisão continua válida mas parte dela muda
  de peso ou de escopo, ela é *emendada* (bloco no topo + link), não substituída.
  Ver o ADR-001.
- **Um ADR por decisão.** Se ao escrever surgir um "e além disso decidimos que…"
  sem relação com o título, são dois ADRs.

## Decisões pendentes conhecidas

Registradas aqui para não se perderem. Cada uma vira um ADR quando resolvida:

| Pendência | Prazo | Onde está descrita |
|---|---|---|
| Papel do Postgres usado pelo Go e estratégia de RLS (defesa em profundidade ou só Storage) | **Antes** de criar as tabelas — Fase 3 | [ADR-001 (emenda)](./001-migrations-e-rls-supabase.md), [`context/banco-de-dados.md`](../context/banco-de-dados.md) §3 |
| Mecanismo de fan-out do chat (`LISTEN/NOTIFY` vs. Redis pub/sub) | Antes de subir a 2ª instância de `cmd/ws` | [`context/arquitetura.md`](../context/arquitetura.md) §5.3 |
| Como o app fala com o Supabase Auth (biblioteca KMP de comunidade vs. REST direto) | Fase 3 | [`context/arquitetura.md`](../context/arquitetura.md) §5.1 |
| Bibliotecas do mobile (cliente HTTP, DI, cache local) e do backend (roteador, WebSocket, sqlc) | Fases 2 e 3 | [`context/stack.md`](../context/stack.md) §3 |
| Critério numérico para abrir o cadastro ao público (fim do crescimento por convite) | Antes do lançamento | [`roadmaps/03-backend-banco-e-autenticacao.md`](../roadmaps/03-backend-banco-e-autenticacao.md) |
