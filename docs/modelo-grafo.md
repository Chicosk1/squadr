# Modelo do Grafo

> **Fases:** 1 — Modelagem; 4 — Geração de candidatos (filtros por jogo e sinais)
> **Status:** Proposta (aguardando validação do orientador)
> **Banco:** Neo4j 5.26 Community

## Objetivo

Descrever os nós, relações, propriedades, dados de referência, índices e as perguntas que o modelo de grafo precisa responder.

## Princípios do modelo

1. **A recomendação é contextualizada por jogo.** O mesmo jogador pode ter perfis diferentes em jogos diferentes, então a busca sempre responde "parceiros para *este* jogo".
2. **Nó ou propriedade?** Se vários usuários compartilham aquilo e a busca passa por ele para encontrar pessoas, é **nó**. Se serve apenas para filtrar ou comparar, é **propriedade**.
3. **O que depende do jogo fica na relação `JOGA`**, não no usuário.

## Visão geral

```
(:Usuario)-[:JOGA]->(:Jogo)
(:Usuario)-[:EXERCE]->(:Funcao)
(:Usuario)-[:DISPONIVEL]->(:Horario)
(:Funcao)-[:COMPLEMENTA]-(:Funcao)
(:Usuario)-[:CURTIU]->(:Usuario)
(:Usuario)-[:RECUSOU]->(:Usuario)
```

---

## Nós

### `Usuario`

| Propriedade | Tipo | Origem | Observação |
|---|---|---|---|
| `id` | string | Sistema | Único |
| `nome` | string | Formulário | |
| `faixaEtaria` | string | G02 | `18-24`, `25-34`, `35+` |
| `faixasAceitas` | lista de string | G03 | |
| `regiao` | string | G04 | `AS`, `AN`, `EU` |
| `idiomas` | lista de string | G05 | `pt`, `en`, `es` |
| `usaVoz` | string | G07 | `sim`, `as_vezes`, `nao` |
| `exigeVoz` | boolean | G08 | |
| `toleranciaToxicidade` | int | G09 | 1 a 5 |
| `compromisso` | string | G10 | `fixo`, `eventual`, `tanto_faz` |
| `pesoHorario` | int | P01 | 1 a 5, padrão 3 |
| `pesoIntensidade` | int | P02 | 1 a 5, padrão 3 |
| `pesoModo` | int | P03 | 1 a 5, padrão 3 |
| `pesoComunicacao` | int | P04 | 1 a 5, padrão 3 |
| `pesoFuncao` | int | P05 | 1 a 5, padrão 3 |
| `pesoNivel` | int | P06 | 1 a 5, padrão 3 |
| `persona` | string | Gerador | **Somente dados mock.** Nunca usado pelo motor, apenas pela avaliação |

Os pesos ficam "achatados" em propriedades separadas porque o Neo4j não aceita objetos aninhados como propriedade, apenas valores simples e listas.

### `Jogo`

| Propriedade | Tipo | Valores |
|---|---|---|
| `id` | string | `ark`, `rust`, `lol` |
| `nome` | string | `ARK`, `Rust`, `League of Legends` |

### `Horario`

28 nós fixos: 7 dias × 4 períodos.

| Propriedade | Tipo | Exemplo | Observação |
|---|---|---|---|
| `id` | string | `sex-noite` | Único |
| `dia` | string | `sex` | `seg` a `dom` |
| `periodo` | string | `noite` | `madrugada`, `manha`, `tarde`, `noite` |
| `bit` | int | `19` | `indiceDia × 4 + indicePeriodo`, de 0 a 27 |

O `bit` permite montar, no Go, um bitmask de 28 bits com a disponibilidade do usuário. Comparar horários vira uma operação de bits.

### `Funcao`

16 nós fixos.

| Propriedade | Tipo | Exemplo |
|---|---|---|
| `id` | string | `ark:breeder`, `lol:support` |
| `jogo` | string | `ark` |
| `nome` | string | `Breeder/domador` |

| Jogo | Funções |
|---|---|
| ARK | `ark:breeder`, `ark:builder`, `ark:gatherer`, `ark:pvp`, `ark:explorer` |
| Rust | `rust:builder`, `rust:gatherer`, `rust:pvp`, `rust:raider`, `rust:electrician` |
| LoL | `lol:top`, `lol:jungle`, `lol:mid`, `lol:adc`, `lol:support`, `lol:fill` |

---

## Relações

### Normalização do nível

| Jogo | Regra |
|---|---|
| ARK e Rust | `<100` = 0,125; `100-500` = 0,375; `500-2000` = 0,625; `2000+` = 0,875 |
| LoL | Sem ranked = 0; demais = (elo + 1) / 10, de Ferro (0,1) a Desafiante (1,0) |

### `(Usuario)-[:JOGA]->(Jogo)`

Guarda tudo que depende do jogo. O Neo4j aceita relações do mesmo tipo com propriedades diferentes, mas **o banco não valida isso**, então a estrutura abaixo é o contrato.

**Propriedades comuns a todos os jogos**

| Propriedade | Tipo | Origem | Observação |
|---|---|---|---|
| `intensidade` | int | A07 / R07 / L05 | 1 a 5 |
| `nivel` | float | A05 / R05 / L01 | Normalizado de 0 a 1 |
| `nivelDesejado` | string | A09 / R08 / L06 | `parecido`, `mais_experiente`, `tanto_faz` |
| `personaJogo` | string | Gerador | **Somente dados mock.** Persona do usuário naquele jogo (no multi-jogo, a persona copiada). Nunca usada pelo motor |

**Propriedades do ARK**

| Propriedade | Tipo | Origem | Valores |
|---|---|---|---|
| `horas` | string | A05 | `<100`, `100-500`, `500-2000`, `2000+` |
| `versao` | string | A01 | `ase`, `asa` |
| `plataforma` | string | A02 | `pc_steam`, `pc_epic`, `xbox`, `playstation` |
| `modo` | string | A03 | `pve`, `pvp`, `ambos` |
| `servidor` | string | A04 | `oficial`, `rates_baixas`, `rates_altas`, `tanto_faz` |
| `tamanhoTribo` | string | A08 | `2-3`, `4-6`, `7+`, `tanto_faz` |

**Propriedades do Rust**

| Propriedade | Tipo | Origem | Valores |
|---|---|---|---|
| `horas` | string | R05 | `<100`, `100-500`, `500-2000`, `2000+` |
| `versao` | string | R01 | `pc`, `console` |
| `servidor` | string | R02 | `oficial`, `comunitario`, `modificado`, `pve` |
| `limiteGrupo` | string | R03 | `solo_duo`, `trio`, `quad`, `sem_limite` |
| `wipe` | string | R04 | `sempre`, `as_vezes`, `raramente` |

**Propriedades do LoL**

| Propriedade | Tipo | Origem | Valores |
|---|---|---|---|
| `elo` | int | L01 | 0 (Ferro) a 9 (Desafiante); ausente se sem ranked |
| `filas` | lista de string | L04 | `solo_duo`, `flex`, `normal`, `aram`, `clash` |

### `(Usuario)-[:EXERCE]->(Funcao)`

| Propriedade | Tipo | Valores |
|---|---|---|
| `prioridade` | int | `1` (principal), `2` (secundária) |

### `(Usuario)-[:DISPONIVEL]->(Horario)`

Sem propriedades.

### `(Funcao)-[:COMPLEMENTA]-(Funcao)`

| Propriedade | Tipo | Valores |
|---|---|---|
| `valor` | float | 0 (péssima dupla) a 1 (dupla ideal) |

Cada par é gravado uma única vez e consultado **sem direção**. Uma função pode complementar a si mesma (relação de uma função com ela mesma).

### `(Usuario)-[:CURTIU]->(Usuario)` e `(Usuario)-[:RECUSOU]->(Usuario)`

| Propriedade | Tipo | Observação |
|---|---|---|
| `jogo` | string | Likes são contextualizados por jogo |
| `em` | datetime | Momento da ação |

---

## Justificativa: nó ou propriedade

| Elemento | Escolha | Motivo |
|---|---|---|
| Jogo | Nó | Compartilhado; a busca parte dele para encontrar pessoas |
| Horário | Nó | Compartilhado; contar horários em comum é uma travessia |
| Função | Nó | Compartilhada; permite guardar a complementaridade como relação |
| Região | Propriedade | Usada só como filtro |
| Idioma | Propriedade (lista) | Usado só como filtro |
| Faixa etária | Propriedade | Usada só como filtro |
| Intensidade, modo, nível | Propriedade da relação `JOGA` | Dependem do jogo |

A matriz de complementaridade fica no grafo porque assim é visível, editável e explicável. O Go carrega essa matriz em memória ao iniciar, para que o scoring não precise consultar o banco.

---

## Matrizes de complementaridade

Os valores são **premissas heurísticas** e devem ser declarados como tal na monografia. Se possível, validar com um questionário rápido aplicado a jogadores dos três jogos. As matrizes são simétricas; apenas a metade superior é mostrada.

### ARK

| | Breeder | Builder | Coletor | PvP | Explorador |
|---|---|---|---|---|---|
| **Breeder** | 0,4 | 0,9 | 0,8 | 0,8 | 0,7 |
| **Builder** | | 0,4 | 0,9 | 0,8 | 0,6 |
| **Coletor** | | | 0,6 | 0,7 | 0,6 |
| **PvP** | | | | 0,8 | 0,7 |
| **Explorador** | | | | | 0,7 |

### Rust

| | Builder | Coletor | PvP | Raider | Eletricista |
|---|---|---|---|---|---|
| **Builder** | 0,4 | 0,9 | 0,8 | 0,7 | 0,8 |
| **Coletor** | | 0,6 | 0,8 | 0,7 | 0,7 |
| **PvP** | | | 0,8 | 0,9 | 0,6 |
| **Raider** | | | | 0,6 | 0,7 |
| **Eletricista** | | | | | 0,3 |

### League of Legends

| | Top | Jungle | Mid | ADC | Suporte | Preencher |
|---|---|---|---|---|---|---|
| **Top** | 0,1 | 0,8 | 0,5 | 0,5 | 0,5 | 0,7 |
| **Jungle** | | 0,1 | 0,9 | 0,7 | 0,8 | 0,7 |
| **Mid** | | | 0,1 | 0,5 | 0,6 | 0,7 |
| **ADC** | | | | 0,1 | 1,0 | 0,7 |
| **Suporte** | | | | | 0,1 | 0,7 |
| **Preencher** | | | | | | 0,5 |

No LoL, a mesma função vale 0,1 porque duas pessoas não podem ocupar a mesma rota no mesmo time.

---

## Dados de referência

Ficam no arquivo `scripts/referencia.cypher` e são carregados uma única vez.

### Jogos

```cypher
UNWIND [
  {id: 'ark',  nome: 'ARK'},
  {id: 'rust', nome: 'Rust'},
  {id: 'lol',  nome: 'League of Legends'}
] AS jogo
MERGE (j:Jogo {id: jogo.id})
SET j.nome = jogo.nome;
```

### Horários

```cypher
UNWIND range(0, 6) AS d
UNWIND range(0, 3) AS p
WITH d, p,
     ['seg','ter','qua','qui','sex','sab','dom'][d] AS dia,
     ['madrugada','manha','tarde','noite'][p] AS periodo
MERGE (h:Horario {id: dia + '-' + periodo})
SET h.dia = dia, h.periodo = periodo, h.bit = d * 4 + p;
```

### Funções

```cypher
UNWIND [
  {id: 'ark:breeder',      jogo: 'ark',  nome: 'Breeder/domador'},
  {id: 'ark:builder',   jogo: 'ark',  nome: 'Builder'},
  {id: 'ark:gatherer',      jogo: 'ark',  nome: 'Coletor'},
  {id: 'ark:pvp',          jogo: 'ark',  nome: 'PvP'},
  {id: 'ark:explorer',   jogo: 'ark',  nome: 'Explorador'},
  {id: 'rust:builder',  jogo: 'rust', nome: 'Builder'},
  {id: 'rust:gatherer',     jogo: 'rust', nome: 'Coletor'},
  {id: 'rust:pvp',         jogo: 'rust', nome: 'PvP/roamer'},
  {id: 'rust:raider',      jogo: 'rust', nome: 'Raider'},
  {id: 'rust:electrician', jogo: 'rust', nome: 'Eletricista'},
  {id: 'lol:top',          jogo: 'lol',  nome: 'Top'},
  {id: 'lol:jungle',       jogo: 'lol',  nome: 'Jungle'},
  {id: 'lol:mid',          jogo: 'lol',  nome: 'Mid'},
  {id: 'lol:adc',          jogo: 'lol',  nome: 'ADC'},
  {id: 'lol:support',      jogo: 'lol',  nome: 'Suporte'},
  {id: 'lol:fill',    jogo: 'lol',  nome: 'Preencher'}
] AS funcao
MERGE (f:Funcao {id: funcao.id})
SET f.jogo = funcao.jogo, f.nome = funcao.nome;
```

### Complementaridade

```cypher
UNWIND [
  // ARK
  {a: 'ark:breeder',    b: 'ark:breeder',    valor: 0.4},
  {a: 'ark:breeder',    b: 'ark:builder', valor: 0.9},
  {a: 'ark:breeder',    b: 'ark:gatherer',    valor: 0.8},
  {a: 'ark:breeder',    b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:breeder',    b: 'ark:explorer', valor: 0.7},
  {a: 'ark:builder', b: 'ark:builder', valor: 0.4},
  {a: 'ark:builder', b: 'ark:gatherer',    valor: 0.9},
  {a: 'ark:builder', b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:builder', b: 'ark:explorer', valor: 0.6},
  {a: 'ark:gatherer',    b: 'ark:gatherer',    valor: 0.6},
  {a: 'ark:gatherer',    b: 'ark:pvp',        valor: 0.7},
  {a: 'ark:gatherer',    b: 'ark:explorer', valor: 0.6},
  {a: 'ark:pvp',        b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:pvp',        b: 'ark:explorer', valor: 0.7},
  {a: 'ark:explorer', b: 'ark:explorer', valor: 0.7},
  // Rust
  {a: 'rust:builder',  b: 'rust:builder',  valor: 0.4},
  {a: 'rust:builder',  b: 'rust:gatherer',     valor: 0.9},
  {a: 'rust:builder',  b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:builder',  b: 'rust:raider',      valor: 0.7},
  {a: 'rust:builder',  b: 'rust:electrician', valor: 0.8},
  {a: 'rust:gatherer',     b: 'rust:gatherer',     valor: 0.6},
  {a: 'rust:gatherer',     b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:gatherer',     b: 'rust:raider',      valor: 0.7},
  {a: 'rust:gatherer',     b: 'rust:electrician', valor: 0.7},
  {a: 'rust:pvp',         b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:pvp',         b: 'rust:raider',      valor: 0.9},
  {a: 'rust:pvp',         b: 'rust:electrician', valor: 0.6},
  {a: 'rust:raider',      b: 'rust:raider',      valor: 0.6},
  {a: 'rust:raider',      b: 'rust:electrician', valor: 0.7},
  {a: 'rust:electrician', b: 'rust:electrician', valor: 0.3},
  // LoL
  {a: 'lol:top',       b: 'lol:top',       valor: 0.1},
  {a: 'lol:top',       b: 'lol:jungle',    valor: 0.8},
  {a: 'lol:top',       b: 'lol:mid',       valor: 0.5},
  {a: 'lol:top',       b: 'lol:adc',       valor: 0.5},
  {a: 'lol:top',       b: 'lol:support',   valor: 0.5},
  {a: 'lol:top',       b: 'lol:fill', valor: 0.7},
  {a: 'lol:jungle',    b: 'lol:jungle',    valor: 0.1},
  {a: 'lol:jungle',    b: 'lol:mid',       valor: 0.9},
  {a: 'lol:jungle',    b: 'lol:adc',       valor: 0.7},
  {a: 'lol:jungle',    b: 'lol:support',   valor: 0.8},
  {a: 'lol:jungle',    b: 'lol:fill', valor: 0.7},
  {a: 'lol:mid',       b: 'lol:mid',       valor: 0.1},
  {a: 'lol:mid',       b: 'lol:adc',       valor: 0.5},
  {a: 'lol:mid',       b: 'lol:support',   valor: 0.6},
  {a: 'lol:mid',       b: 'lol:fill', valor: 0.7},
  {a: 'lol:adc',       b: 'lol:adc',       valor: 0.1},
  {a: 'lol:adc',       b: 'lol:support',   valor: 1.0},
  {a: 'lol:adc',       b: 'lol:fill', valor: 0.7},
  {a: 'lol:support',   b: 'lol:support',   valor: 0.1},
  {a: 'lol:support',   b: 'lol:fill', valor: 0.7},
  {a: 'lol:fill', b: 'lol:fill', valor: 0.5}
] AS par
MATCH (a:Funcao {id: par.a}), (b:Funcao {id: par.b})
MERGE (a)-[c:COMPLEMENTA]->(b)
SET c.valor = par.valor;
```

---

## Constraints e índices

Executar antes de qualquer carga de dados.

```cypher
CREATE CONSTRAINT usuario_id IF NOT EXISTS FOR (u:Usuario) REQUIRE u.id IS UNIQUE;
CREATE CONSTRAINT jogo_id    IF NOT EXISTS FOR (j:Jogo)    REQUIRE j.id IS UNIQUE;
CREATE CONSTRAINT horario_id IF NOT EXISTS FOR (h:Horario) REQUIRE h.id IS UNIQUE;
CREATE CONSTRAINT funcao_id  IF NOT EXISTS FOR (f:Funcao)  REQUIRE f.id IS UNIQUE;
// usuario_regiao foi removido (D-032): o planejador nunca o usou e, forçado, piorou a consulta
DROP INDEX usuario_regiao IF EXISTS;
```

Sem a constraint em `Usuario.id`, toda busca por usuário percorre o banco inteiro.

Não há índice em `Usuario.regiao`. A consulta de candidatos parte do usuário (pela constraint de `id`) e expande pelo nó `Jogo`. A região é só um filtro sobre os jogadores do jogo, e o índice nunca foi usado pelo planejador (D-032).

---

## Perguntas que o modelo precisa responder

Um modelo de grafo é validado pelas perguntas que consegue responder. Cada uma deve ser testada assim que houver dados.

1. Quem joga o mesmo jogo, na mesma região, com faixas etárias mutuamente aceitas e pelo menos um idioma em comum?
2. Quantos horários essas pessoas têm em comum comigo?
3. Qual é a melhor complementaridade entre as minhas funções e as delas?
4. Quem curtiu as mesmas pessoas que eu curti neste jogo?
5. Por que esta pessoa foi recomendada (quais caminhos nos ligam)?

### Consulta de validação (perguntas 1 e 3)

```cypher
MATCH (eu:Usuario {id: $id})-[:JOGA]->(j:Jogo {id: $jogo})<-[:JOGA]-(outro:Usuario)
WHERE outro <> eu
  AND outro.regiao = eu.regiao
  AND outro.faixaEtaria IN eu.faixasAceitas
  AND eu.faixaEtaria IN outro.faixasAceitas
  AND any(i IN outro.idiomas WHERE i IN eu.idiomas)
  AND NOT EXISTS { (eu)-[x:CURTIU|RECUSOU]->(outro) WHERE x.jogo = $jogo }
OPTIONAL MATCH (eu)-[:EXERCE]->(f1:Funcao {jogo: $jogo})-[c:COMPLEMENTA]-(f2:Funcao)<-[:EXERCE]-(outro)
RETURN outro.id AS id, max(c.valor) AS melhorComplemento
ORDER BY melhorComplemento DESC
LIMIT 20
```

Os filtros específicos de cada jogo foram implementados na Fase 4 (`internal/graph/candidatos.go`):

| Jogo | Filtros específicos |
|---|---|
| ARK | Mesma versão; mesmo grupo de plataforma (PC Steam e Epic jogam juntos, cada console é isolado); modo igual ou algum dos dois em "ambos" |
| Rust | Mesma versão |
| LoL | Ao menos uma fila em comum; Solo/Duo só conta se ambos são ranqueados e a diferença de elo é no máximo 2 |

Sinais calculados no grafo:

| Sinal | Como é calculado | Consulta |
|---|---|---|
| Horários em comum | Contagem de nós `Horario` compartilhados | Candidatos |
| Complementaridade | Maior `valor` entre as relações `COMPLEMENTA` ligando uma função minha a uma função do candidato | Candidatos |
| Sinal social | Soma, sobre os usuários que curtiram alguém que eu também curti, do número de curtidas em comum, para cada candidato que eles curtiram | Social |

---

## Pontos de atenção

- **Supernós:** o nó de cada `Jogo` estará ligado a milhares de usuários. Os filtros rígidos devem ser aplicados o mais cedo possível na consulta para cortar caminhos.
- **Propriedades heterogêneas em `JOGA`:** o banco não valida a estrutura; o código Go é responsável por gravar e ler cada jogo corretamente.
- **Complementaridade:** valores heurísticos, sujeitos a revisão.