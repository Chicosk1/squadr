# 006 — Conexão do backend Go ao Postgres via pooler (Supavisor) em modo sessão

**Status:** Aceito
**Data:** 04/08/2026

## Contexto

Item da seção 5 do roadmap ([`03-backend-banco-e-autenticacao.md`](../roadmaps/03-backend-banco-e-autenticacao.md)):
decidir entre conexão direta (porta 5432) e pooler (Supavisor) para o `pgx` do
backend Go.

Teste de conectividade a partir do ambiente de desenvolvimento local:

- `nslookup -type=A db.eliqtpvclmyljsvsgeyu.supabase.co` não retornou endereço.
- `nslookup -type=AAAA` retornou endereço.
- `Test-NetConnection -ComputerName db.eliqtpvclmyljsvsgeyu.supabase.co -Port 5432`
  falhou na resolução de nome (`Name resolution ... failed`).

Conclusão do teste: o host de conexão direta do Supabase é IPv6-only por
padrão (IPv4 só via add-on pago), e a rede de desenvolvimento local não tem
rota IPv6 até ele. Conexão direta não é viável nesse ambiente.

## Decisão

O `DATABASE_URL` do backend Go usa o **pooler (Supavisor) em modo sessão**,
em todos os ambientes (dev local e produção no Fly.io) — não só como
contorno do problema de rede local, mas como escolha padrão do projeto.

Motivos:

1. **Resolve o bloqueio de IPv6.** O pooler é acessível via IPv4, então
   funciona no ambiente de desenvolvimento local sem depender de rota IPv6.
2. **Modo sessão, não modo transação.** Modo sessão preserva a mesma
   ergonomia da conexão direta para o `pgx` — não exige desligar o cache de
   prepared statements. Modo transação exigiria essa configuração extra e
   resolve um problema (muitas conexões efêmeras, ex: funções serverless)
   que não é o cenário aqui: `cmd/api` e `cmd/ws` são processos Go de vida
   longa, mantendo um pool pequeno e estável de conexões — o caso de uso
   clássico de conexão direta ou pooler em modo sessão, não de modo
   transação.
3. **Mesma string em todos os ambientes.** Mesmo que o Fly.io em produção
   tivesse saída IPv6 e conseguisse usar conexão direta, manter dev e
   produção na mesma estratégia de conexão evita comportamento divergente
   entre ambientes por causa de uma diferença de infraestrutura de rede.

## Consequências

- O usuário na connection string passa a ser `postgres.<PROJECT_REF>` (formato
  exigido pelo Supavisor para rotear a conexão), diferente do `postgres` usado
  na conexão direta.
- Caracteres especiais na senha do banco precisam ser **percent-encoded** na
  connection string (ex.: `/` → `%2F`) — sem isso, o parser de URL usado pelo
  `pgx` interpreta mal a authority da URI e a conexão falha antes mesmo de
  tentar autenticar. Isso já valia para a conexão direta, mas foi só ao montar
  esta string que o problema foi notado; `.env` e `.env.example` foram
  corrigidos.
- Se, no futuro, algum ambiente precisar de conexão direta (ex.: uma instância
  rodando em rede com suporte IPv6 completo e um motivo concreto para evitar o
  pooler), esta decisão precisa ser reaberta, não alterada silenciosamente.
