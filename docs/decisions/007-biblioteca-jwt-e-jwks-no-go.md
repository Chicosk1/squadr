# 007 — Biblioteca de JWT e JWKS no backend Go

**Status:** Aceito
**Data:** 04/08/2026

## Contexto

A seção 7.2 do roadmap ([`03-backend-banco-e-autenticacao.md`](../roadmaps/03-backend-banco-e-autenticacao.md))
exige que o Go valide o JWT emitido pelo Supabase Auth. O requisito não é só
"conferir assinatura":

- buscar o JWKS do Supabase e **manter em cache** — sem refazer a busca a
  cada requisição;
- **não quebrar quando a chave rotacionar** — `kid` desconhecido tem que
  disparar uma recarga;
- validar assinatura, `iss`, `aud` e `exp`;
- extrair o `sub`.

O [ADR-004](./004-bibliotecas-do-mobile-e-backend.md) escolheu roteador HTTP e
biblioteca de WebSocket, mas não tratou de JWT.

Verificação feita no projeto antes de decidir: o JWKS em
`<SUPABASE_URL>/auth/v1/.well-known/jwks.json` devolve chaves **ES256**
(ECC P-256), e já traz **duas chaves** simultâneas — rotação não é hipótese
distante, é o estado atual do projeto.

## Decisão

**`github.com/golang-jwt/jwt/v5`** (v5.3.1) para parse e validação, com
**`github.com/MicahParks/keyfunc/v3`** (v3.8.1) para o JWKS.

A alternativa considerada foi `lestrrat-go/jwx/v3` (v3.2.0), implementação
JOSE completa que também resolve cache e rotação. A escolha pelo par
golang-jwt + keyfunc foi por **superfície de API menor e mais fácil de
aprender** — quem mantém o backend está aprendendo Go durante o projeto, e
golang-jwt é a biblioteca de JWT mais difundida do ecossistema, o que torna
exemplo e resposta de comunidade mais fáceis de achar.

O keyfunc cobre o requisito difícil sem código próprio: `NewDefaultCtx` sobe
uma goroutine de refresh em background (encerrada quando o contexto é
cancelado) e recarrega o JWK Set ao encontrar um `kid` desconhecido.

### Trava de algoritmo

`jwt.WithValidMethods([]string{"ES256", "RS256"})` é obrigatório, não
opcional. Sem essa lista, um invasor poderia assinar um token **HS256**
usando a chave *pública* do JWKS (que é pública por definição) como segredo
HMAC, e a verificação passaria — é o ataque clássico de confusão de
algoritmo. RS256 entra na lista junto com ES256 porque é o outro tipo de
chave assimétrica que o Supabase pode emitir numa rotação; nenhum `HS*` é
aceito. Há teste dedicado a isso em `internal/auth/verifier_test.go`.

## Consequências

- `internal/auth` depende de duas bibliotecas em vez de uma. Elas são do
  mesmo autor no caso do keyfunc/jwkset e integram-se por design com o
  golang-jwt, então não é combinação exótica.
- A validação de `iss`/`aud`/`exp` é declarada como opção de parse
  (`WithIssuer`, `WithAudience`, `WithExpirationRequired`), não escrita à
  mão — menos código próprio para errar.
- Todo erro de validação é normalizado em `auth.ErrInvalidToken`, mas o erro
  original fica embrulhado junto (dois `%w`): o cliente sempre recebe o mesmo
  401 genérico, e o log do servidor ainda identifica a causa exata. Devolver
  ao cliente qual validação falhou entregaria informação a quem está sondando
  a API.
- Se um dia o projeto precisar de operações JOSE além de verificar token
  (emitir JWT, criptografar payload, JWE), vale reabrir a comparação com o
  `jwx` — hoje o backend só **valida** token, nunca emite.
