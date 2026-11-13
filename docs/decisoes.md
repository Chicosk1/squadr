# Registro de Decisões

## Objetivo

Registrar cada decisão importante do projeto com contexto, opções consideradas, escolha e motivo. Este documento serve de base para o capítulo de desenvolvimento da monografia.

## Modelo de entrada

```
## D-XXX — Título
- **Data:**
- **Status:** Proposta | Aceita | Substituída por D-YYY
- **Contexto:** qual problema motivou a decisão
- **Opções consideradas:**
- **Decisão:**
- **Motivo:**
- **Consequências:** o que muda no projeto a partir daqui
```

---

## Fase 1 — Modelagem

## D-001 — Escopo de jogos: ARK, Rust e League of Legends

- **Status:** Proposta
- **Contexto:** o protótipo precisa de um escopo limitado de jogos para ser viável no prazo do TCC.
- **Opções consideradas:** apenas ARK; ARK e Rust; ARK, Rust e LoL; mais de três jogos.
- **Decisão:** ARK, Rust e League of Legends.
- **Motivo:** ARK é o caso que motiva o trabalho (sistema de formação de grupos fraco). Rust traz a dinâmica de wipes e limites de grupo. LoL permite demonstrar complementaridade de funções e terá integração futura com a API da Riot.
- **Consequências:** três seções específicas no questionário e três matrizes de complementaridade.

## D-002 — Recomendação contextualizada por jogo

- **Status:** Proposta
- **Contexto:** o mesmo jogador pode ser competitivo em um jogo e casual em outro. Comparar perfis "no geral" mistura contextos diferentes.
- **Opções consideradas:** recomendação geral por usuário; recomendação por jogo.
- **Decisão:** toda recomendação responde "parceiros para *este* jogo". A API recebe o jogo como parâmetro (`GET /usuarios/{id}/recomendacoes?jogo=ark`).
- **Motivo:** reflete como os jogadores realmente procuram grupo e evita recomendações incoerentes.
- **Consequências:** likes e recusas também passam a ser registrados por jogo.

## D-003 — Atributos dependentes do jogo na relação `JOGA`

- **Status:** Proposta
- **Contexto:** consequência da D-002. Intensidade, modo, servidor e nível variam de jogo para jogo.
- **Opções consideradas:** propriedades no nó `Usuario`; propriedades na relação `JOGA`; um nó intermediário `PerfilJogo`.
- **Decisão:** propriedades na relação `(Usuario)-[:JOGA]->(Jogo)`.
- **Motivo:** mantém o modelo simples e os dados são lidos durante a própria travessia, sem saltos extras. Um nó intermediário só se justificaria se o perfil por jogo tivesse relações próprias.
- **Consequências:** a relação `JOGA` tem propriedades diferentes para cada jogo; o banco não valida essa estrutura, então o contrato fica documentado em `modelo-grafo.md` e o Go é responsável por respeitá-lo.

## D-004 — Público 18+ e idade coletada por faixa

- **Status:** Proposta
- **Contexto:** o app conecta desconhecidos. Incluir menores exigiria uma camada de segurança fora do escopo.
- **Opções consideradas:** aberto a todas as idades; separação entre menores e adultos; somente 18+.
- **Decisão:** somente 18+, com idade coletada por faixa (18–24, 25–34, 35+) e preferência de faixas aceitas.
- **Motivo:** segurança, princípio de minimização de dados da LGPD (coletar a faixa, e não a data de nascimento) e Lei Felca.
- **Consequências:** filtro recíproco de faixa etária na busca.

## D-005 — Disponibilidade em 7 dias × 4 períodos, no horário de Brasília

- **Status:** Proposta
- **Contexto:** a sobreposição de horários é um dos critérios mais importantes, mas granularidade alta torna o formulário cansativo.
- **Opções consideradas:** hora a hora (168 faixas); 7 dias × 4 períodos (28 faixas); apenas dias da semana.
- **Decisão:** 28 faixas (madrugada, manhã, tarde, noite), no horário de Brasília.
- **Motivo:** granularidade viável para um formulário e suficiente para distinguir perfis. 28 faixas cabem em um bitmask de 32 bits, permitindo comparação por operações de bits.
- **Consequências:** usuários em outros fusos ficam registrados como limitação do protótipo.

## D-006 — Funções como nós e complementaridade armazenada no grafo

- **Status:** Proposta
- **Contexto:** a compatibilidade entre funções não é similaridade (dois suportes no LoL não formam uma boa dupla), e sim complementaridade.
- **Opções consideradas:** matriz fixa no código Go; funções como propriedade do usuário; funções como nós com relação `COMPLEMENTA`.
- **Decisão:** funções como nós, ligadas por `(Funcao)-[:COMPLEMENTA {valor}]-(Funcao)`.
- **Motivo:** a matriz fica visível, editável e explicável no próprio grafo, e pode ser usada diretamente na geração de candidatos. O Go carrega a matriz em memória ao iniciar, para o scoring não depender de consultas ao banco.
- **Consequências:** script de dados de referência (`scripts/referencia.cypher`) com 16 funções e suas relações.

## D-007 — Região, idioma e faixa etária como propriedades

- **Status:** Proposta
- **Contexto:** é preciso decidir o que vira nó e o que vira propriedade no grafo.
- **Opções consideradas:** todos como nós; todos como propriedades.
- **Decisão:** região, idioma e faixa etária como propriedades do `Usuario`; jogos, horários e funções como nós.
- **Motivo:** região, idioma e faixa etária são usados apenas como filtro; jogos, horários e funções são caminhos pelos quais os usuários se conectam.
- **Consequências:** índice em `Usuario.regiao` (removido na D-032: o planejador nunca o usou).

## D-008 — Pesos opcionais com valor padrão 3

- **Status:** Proposta
- **Contexto:** o formulário já é longo (cerca de 25 perguntas para quem joga um jogo).
- **Opções consideradas:** pesos obrigatórios; pesos opcionais; pesos fixos definidos pelo sistema.
- **Decisão:** pesos opcionais, de 1 a 5, com padrão 3.
- **Motivo:** reduz o esforço do usuário sem perder a personalização, que é o que torna o score assimétrico.
- **Consequências:** no mock, os pesos variam por persona para que a assimetria apareça nos testes.

## D-009 — ASE e ASA como versões incompatíveis

- **Status:** Proposta
- **Contexto:** ARK: Survival Evolved e ARK: Survival Ascended possuem bases de jogadores e servidores separados.
- **Opções consideradas:** tratar como o mesmo jogo; tratar como jogos diferentes; tratar como versões do mesmo jogo com filtro rígido.
- **Decisão:** mesmo nó `Jogo`, com a versão na relação `JOGA` e filtro rígido de versão.
- **Motivo:** mantém um único jogo no modelo sem gerar matches impossíveis.
- **Consequências:** filtro de versão na geração de candidatos do ARK.

## D-010 — Valores de complementaridade como premissas heurísticas

- **Status:** Proposta
- **Contexto:** não há dados reais para calcular a complementaridade entre funções.
- **Opções consideradas:** valores definidos pelo autor; valores obtidos por questionário com jogadores; valores aprendidos a partir de likes.
- **Decisão:** valores definidos pelo autor, declarados como premissa.
- **Motivo:** viabilidade no prazo do TCC.
- **Consequências:** a monografia deve declarar a limitação; se houver tempo, validar com um questionário aplicado a jogadores dos três jogos. Aprender os valores a partir de likes fica como trabalho futuro.

---

## Decisões de arquitetura (definidas no planejamento)

## D-011 — Foco exclusivo no motor de recomendação

- **Status:** Aceita
- **Contexto:** o prazo do TCC não comporta o desenvolvimento de um aplicativo completo.
- **Decisão:** desenvolver apenas o motor de recomendação, com dados mockados e um formulário representativo.
- **Motivo:** concentra o esforço na contribuição principal do trabalho.
- **Consequências:** chat, autenticação, interface completa e integrações reais com Steam e Riot ficam como trabalhos futuros.

## D-012 — Go como linguagem do backend

- **Status:** Aceita
- **Contexto:** escolha da linguagem do motor.
- **Decisão:** Go.
- **Motivo:** desempenho, simplicidade, binário único, ferramentas nativas de teste e benchmark (`go test -bench`).
- **Consequências:** scoring implementado em Go puro, isolado do banco, para facilitar testes.

## D-013 — Arquitetura híbrida: Neo4j para candidatos, Go para scoring

- **Status:** Aceita
- **Contexto:** um professor recomendou banco de grafos para ligar os usuários. Grafos são fortes em relações e travessias, mas fracos em comparações numéricas finas.
- **Opções consideradas:** apenas PostgreSQL; apenas Neo4j; Neo4j para geração de candidatos e Go para scoring.
- **Decisão:** Neo4j gera os candidatos (filtros, horários e funções em comum, sinal social); Go calcula o score fino, a reciprocidade e o ranking final.
- **Motivo:** cada ferramenta faz o que faz melhor; o grafo também permite filtragem colaborativa por travessia e explicabilidade dos matches.
- **Consequências:** o pacote de matching depende de uma interface de fonte de candidatos, e não diretamente do Neo4j.

## D-014 — API REST com `net/http` em vez de GraphQL

- **Status:** Aceita
- **Contexto:** GraphQL foi considerado inicialmente para a camada de API.
- **Opções consideradas:** GraphQL (`gqlgen`); REST com framework; REST com a biblioteca padrão.
- **Decisão:** REST com a biblioteca padrão `net/http`.
- **Motivo:** GraphQL acrescenta complexidade (schema, resolvers, problema N+1) sem melhorar o motor, que é o objetivo. Desde o Go 1.22 a biblioteca padrão suporta rotas com método e parâmetros, dispensando frameworks.
- **Consequências:** menos dependências; a API é apenas a camada de acesso ao motor.

## D-015 — Reciprocidade por média harmônica

- **Status:** Aceita
- **Contexto:** em um app no estilo Tinder, o match só é bom se for bom para os dois lados.
- **Opções consideradas:** média aritmética; mínimo; média harmônica.
- **Decisão:** `S(A,B) = 2 · s(A→B) · s(B→A) / (s(A→B) + s(B→A))`.
- **Motivo:** penaliza matches desequilibrados sem ser tão rígida quanto o mínimo. Base teórica: sistemas de recomendação recíproca (Pizzato et al., 2010, RECON).
- **Consequências:** o score unilateral precisa ser assimétrico (pesos e preferências individuais) para que a reciprocidade tenha efeito.

---

## Fase 2 — Infraestrutura

## D-016 — Neo4j Community em Docker Compose

- **Status:** Aceita
- **Contexto:** é preciso um ambiente de banco reproduzível, sem instalação manual.
- **Opções consideradas:** instalação local; Neo4j Aura (nuvem); Docker Compose.
- **Decisão:** imagem `neo4j:5.26-community` via Docker Compose, com limites de memória explícitos.
- **Motivo:** reprodutibilidade (a banca pode subir o ambiente com um comando), sem custo e sem dependência de internet em tempo de execução.
- **Consequências:** a edição Community suporta apenas um banco de dados (`neo4j`); não afeta o escopo do projeto.

## D-017 — Configuração por variáveis de ambiente com arquivo `.env`

- **Status:** Aceita
- **Contexto:** credenciais e portas não devem ficar no código nem no Git.
- **Decisão:** variáveis de ambiente, com `.env` local (ignorado pelo Git) e `.env.example` versionado; leitura com `godotenv`.
- **Motivo:** prática padrão (12-factor); `godotenv` evita configurar variáveis manualmente no Windows.
- **Consequências:** uma dependência externa pequena.

## D-018 — Scripts Cypher versionados e idempotentes

- **Status:** Aceita
- **Contexto:** o grafo de referência (jogos, horários, funções, complementaridade) precisa ser recriável a qualquer momento.
- **Decisão:** scripts em `scripts/*.cypher` usando `MERGE`, executados por `cmd/migrar`, que ao final verifica as contagens esperadas.
- **Motivo:** o ambiente pode ser recriado do zero com um comando, e erros de carga são detectados imediatamente.
- **Consequências:** se o modelo mudar, os scripts são a única fonte de verdade a atualizar.

## D-019 — Driver único compartilhado

- **Status:** Aceita
- **Contexto:** o driver do Neo4j mantém um pool de conexões.
- **Decisão:** criar um único `Cliente` na inicialização e compartilhá-lo; nunca criar um por requisição.
- **Motivo:** evita o custo de abrir conexões repetidamente, que distorceria os benchmarks da Fase 7.

## D-020 — Graph Data Science adiado

- **Status:** Aceita
- **Contexto:** o plugin GDS só será usado na Fase 8 (opcional).
- **Decisão:** não instalar APOC nem GDS agora.
- **Motivo:** o download na inicialização pode falhar e atrasa o startup; será adicionado se a fase opcional for executada.

---

## Fase 3 — Domínio e dados mock

## D-021 — Dados mock em Go, reproduzíveis por semente

- **Status:** Aceita
- **Contexto:** a avaliação exige que os experimentos possam ser repetidos.
- **Opções consideradas:** biblioteca de dados falsos (`gofakeit`); gerador próprio com `math/rand/v2`.
- **Decisão:** gerador próprio, com PCG e semente explícita (padrão 42). Apelidos a partir de listas próprias.
- **Motivo:** mesma semente gera exatamente os mesmos dados; os apelidos não influenciam o motor, então uma dependência externa não se justifica.
- **Consequências:** testes automatizados garantem a determinismo; nunca iterar sobre `map` no gerador.

## D-022 — Personas como dados, não código

- **Status:** Aceita
- **Decisão:** cada persona é uma struct (`Persona` + `ModeloJogo`) em `internal/seed/personas.go`; o gerador é genérico.
- **Motivo:** ajustar uma persona (ex.: após feedback do orientador) exige alterar uma linha, sem tocar na lógica.

## D-023 — Ruído e outliers apenas em atributos flexíveis

- **Status:** Aceita
- **Contexto:** outliers em atributos dos filtros rígidos gerariam perfis impossíveis de combinar.
- **Decisão:** outliers só alteram intensidade, horários, tolerância, nível desejado, compromisso e pesos.
- **Consequências:** o gabarito permanece válido para todos os usuários.

## D-024 — Filtros rígidos reimplementados em Go como oráculo

- **Status:** Aceita
- **Contexto:** a avaliação precisa saber quais candidatos são válidos sem depender da query sob teste.
- **Decisão:** `PassaFiltros` e `Grau` em `internal/seed/gabarito.go`, independentes da query Cypher da Fase 4.
- **Consequências:** na Fase 4, as duas implementações são comparadas; divergência indica bug em uma delas. Crossplay e diferença de elo no duo são premissas documentadas.

## D-025 — `personaJogo` na relação JOGA (somente mock)

- **Status:** Aceita
- **Contexto:** no multi-jogo, a persona efetiva muda de jogo para jogo.
- **Decisão:** guardar `personaJogo` na relação `JOGA`, além de `persona` no usuário.
- **Consequências:** a avaliação pode reconstruir o gabarito a partir do grafo. O motor nunca lê esses campos.

## D-026 — Likes por amostragem por rejeição, com probabilidade por grau

- **Status:** Aceita
- **Decisão:** de 20 a 50 candidatos por usuário e jogo, sorteados entre os que passam nos filtros; probabilidade de like de 70% (grau 2), 45% (grau 1) e 15% (grau 0).
- **Motivo:** a amostragem evita comparar cada usuário com todos os outros (custo quadrático).
- **Consequências (ameaça à validade):** as curtidas derivam do gabarito usado na avaliação, o que favorece o sinal social. A ablação da Fase 7 inclui uma versão sem esse sinal, e a limitação é declarada na monografia.

## D-027 — Carga em lote com UNWIND

- **Status:** Aceita
- **Decisão:** lotes de 2.000 usuários e 5.000 interações; limpeza com `CALL ... IN TRANSACTIONS` em sessão auto-commit.
- **Motivo:** um comando por registro seria ordens de grandeza mais lento; apagar tudo em uma transação estoura memória.

---

## Fase 4 — Geração de candidatos

## D-028 — Mínimo de horários em comum como filtro configurável

- **Status:** Aceita
- **Contexto:** o questionário classificava o horário como "Grafo + Score". Quem não compartilha nenhuma faixa de horário nunca poderá jogar junto.
- **Opções consideradas:** apenas score; filtro fixo; filtro configurável.
- **Decisão:** parâmetro `MinHorarios` (padrão 1 nas opções padrão). A verificação contra o oráculo usa 0.
- **Motivo:** descarta candidatos inúteis cedo e reduz o trabalho do scoring.
- **Consequências:** o horário passa a ser filtro e score; o oráculo `PassaFiltros` não o inclui, e por isso a verificação usa `MinHorarios` 0.

## D-029 — Candidatos em duas consultas

- **Status:** Aceita
- **Contexto:** o sinal social (filtragem colaborativa por travessia) tem custo e papel diferentes dos filtros.
- **Decisão:** (1) filtros rígidos, horários e complementaridade, com corte em 300 candidatos ordenados por horários em comum e complemento; (2) sinal social calculado apenas para os candidatos escolhidos.
- **Motivo:** cada consulta pode ser medida separadamente, a segunda custa menos e a ablação "sem sinal social" é desligar uma opção.
- **Consequências:** o sinal social não resgata candidatos cortados na primeira etapa. A Fase 7 mede quanto o corte prejudica.

## D-030 — Filtros por jogo como fragmentos Cypher fixos

- **Status:** Aceita
- **Decisão:** um fragmento de Cypher por jogo, escolhido por um mapa no código Go. A diferença máxima de elo (`MaxDiffEloDuo`) vive em `domain` e é a mesma no oráculo e na consulta.
- **Motivo:** cada jogo tem seu próprio plano de consulta e nada vindo do usuário entra no texto da query (sem injeção).

## D-031 — Verificação diferencial contra o oráculo

- **Status:** Aceita
- **Contexto:** a consulta Cypher e o oráculo Go implementam as mesmas regras de forma independente.
- **Decisão:** `cmd/verificar` regenera os dados pela semente e compara: ida e volta de perfis, conjunto de candidatos, horários em comum, sinal social e ordenação.
- **Consequências:** divergência indica bug em uma das implementações. A complementaridade ainda não é verificada: será quando a matriz existir em Go (Fase 5).

## D-032 — PROFILE da consulta de candidatos e índice `usuario_regiao`

- **Data:** 09/11/2026
- **Status:** Aceita
- **Contexto:** a Fase 4 previa medir a consulta de candidatos com PROFILE e decidir se o índice `usuario_regiao` (consequência da D-007) se justifica. Os experimentos e a tabela planejados originalmente (experimentos A e B, "Passo 7.3") não foram registrados e se perderam. Os experimentos A, B e C abaixo foram definidos em 09/11/2026 a partir do PROFILE da consulta original. Esse PROFILE mostrou dois trechos caros, ambos testados candidato a candidato com `Expand(Into)` a partir de `eu`: a contagem de horários em comum e a exclusão de quem `eu` já avaliou.
- **Opções consideradas:**
  - **A (índice):** remover `usuario_regiao`; mantê-lo; forçá-lo com `USING INDEX outro:Usuario(regiao)`.
  - **B (horários em comum):** o padrão `(eu)-[:DISPONIVEL]->(h)<-[:DISPONIVEL]-(outro)` para cada candidato (original); ou coletar uma vez os horários de `eu` e contar os de `outro` com `h IN meusHorarios`.
  - **C (já avaliados):** `NOT EXISTS { MATCH (eu)-[x:CURTIU|RECUSOU]->(outro) ... }` para cada candidato (original); ou coletar uma vez a lista `vistos` e testar `NOT outro IN vistos`.
- **Decisão:** remover o índice `usuario_regiao` e adotar B e C juntos na consulta de candidatos.
- **Motivo:**
  - **O índice nunca foi usado.** `SHOW INDEXES` mostrava `readCount` 0 desde a criação (14/10/2026), apesar de milhares de consultas do seed e do verificador. O planejador parte de `eu` (constraint de `id`) e expande pelo nó `Jogo`. A região é só um filtro sobre os jogadores do jogo.
  - **Forçar o índice piorou a consulta**, inclusive na América do Norte e na Europa, onde ele é mais seletivo: +36% de db hits e +24% na latência p50.
  - **B e C devolvem exatamente o mesmo resultado, na mesma ordem.** Juntos, cortam 64% dos db hits e 44% da latência p50 (52% da p95).
- **Consequências:**
  - `scripts/schema.cypher` passa a ter `DROP INDEX usuario_regiao IF EXISTS`. Assim, bancos criados antes desta decisão perdem o índice ao rodar o `cmd/migrar`. A consequência "índice em `Usuario.regiao`" da D-007 deixa de valer.
  - A consulta de candidatos (`internal/graph/candidatos.go`) coleta `meusHorarios` e `vistos` uma vez por consulta. O `cmd/verificar` terminou com 0 divergências depois da mudança.
  - A complementaridade passou a ser o trecho mais caro: de 64% a 74% dos db hits da consulta. É o próximo ponto a otimizar, se a Fase 7 mostrar necessidade.
  - Todas as medições usam 10 mil usuários. A Fase 7 repete com 1, 50 e 100 mil.

### Medições da D-032

**Ambiente e método:**
- **Base de dados:** 10.000 usuários (semente 42), em Neo4j 5.26 Community com runtime SLOTTED, rodando em Docker Desktop no Windows.
- **Amostra:** a mesma do `cmd/verificar` (1 a cada 50 usuários): 214 pares (usuário, jogo), 159 da América do Sul e 55 da América do Norte e da Europa.
- **Opções:** as padrão (`Limite` 300, `MinHorarios` 1).
- **db hits:** soma de todos os operadores do PROFILE.
- **Latência:** só a consulta de candidatos, sem PROFILE. Cada variante teve aquecimento e 3 rodadas intercaladas, em ordem rotativa, somando 642 medições.
- **Equivalência:** cada variante foi comparada com a original em 428 execuções (opções padrão; e `MinHorarios` 0 sem limite), com 0 diferenças de ids, horários, complemento e ordem.

**PROFILE por trecho da consulta** (db hits; um usuário por jogo; original → adotada):

| Trecho da consulta | ARK (`u000004`) | Rust (`u000003`) | LoL (`u000001`) |
|---|---|---|---|
| Horários em comum | 220.871 → 20.498 | 341.774 → 26.043 | 115.652 → 8.391 |
| Exclusão de quem já foi avaliado | 65.668 → 74 | 67.665 → 58 | 30.489 → 94 |
| Complementaridade (não mudou) | 110.766 | 134.038 | 42.368 |
| **Total da consulta** | **423.760 → 157.521** | **563.567 → 180.073** | **204.111 → 66.147** |

**db hits médios por consulta** (214 pares):

| Variante | Geral | ARK | Rust | LoL | América do Sul | América do Norte e Europa |
|---|---|---|---|---|---|---|
| Original | 246.417 | 217.417 | 341.373 | 203.164 | 315.262 | 47.393 |
| A: índice forçado | 336.059 (+36%) | 306.682 | 431.048 | 293.253 | 404.667 | 137.720 |
| B: horários coletados uma vez | 125.443 (−49%) | 106.019 | 174.680 | 108.544 | 159.501 | 26.981 |
| C: avaliados coletados uma vez | 206.163 (−16%) | 182.690 | 285.066 | 169.436 | 263.047 | 41.718 |
| **B + C (adotada)** | **89.760 (−64%)** | 78.684 | 119.492 | 78.734 | 113.587 | 20.877 |

**Latência da consulta de candidatos** (642 medições por variante):

| Variante | p50 | p95 | Máx. | p50 América do Sul | p50 América do Norte e Europa |
|---|---|---|---|---|---|
| Original | 53,5 ms | 104,0 ms | 167,7 ms | 63,5 ms | 15,2 ms |
| A: índice forçado | 66,5 ms (+24%) | 117,1 ms | 159,5 ms | 76,8 ms | 27,3 ms |
| B | 35,3 ms (−34%) | 71,0 ms | 96,6 ms | 44,3 ms | 12,5 ms |
| C | 42,0 ms (−21%) | 84,9 ms | 123,9 ms | 49,9 ms | 13,5 ms |
| **B + C (adotada)** | **30,1 ms (−44%)** | **49,9 ms (−52%)** | 78,8 ms | 34,4 ms | 12,3 ms |

**`BuscarCandidatos` completo** (candidatos e sinal social), medido pelo `cmd/verificar` em uma execução de 214 consultas antes e outra depois:

| | p50 | p95 | Máx. |
|---|---|---|---|
| Antes | 64,5 ms | 101,5 ms | 126,6 ms |
| Depois | 46,5 ms | 73,4 ms | 92,7 ms |

O db hit é uma unidade abstrata de trabalho do banco; a métrica que vale para o usuário é a latência. A diferença entre as regiões vem do tamanho do conjunto filtrado: 80% dos usuários são da América do Sul.

## D-033 — Ids de funções no vocabulário da comunidade (inglês)

- **Status:** Aceita
- **Contexto:** os ids de funções foram definidos em português na Fase 1, mas os jogadores usam termos em inglês (*breeder*, *builder*, *gatherer*). Durante a Fase 4, o gerador passou a produzir ids em inglês enquanto o grafo tinha os ids em português; a relação EXERCE deixou de ser criada sem nenhum erro, porque o MATCH da carga não encontrou o nó.
- **Opções consideradas:** ids em português (alinhar o Go ao grafo); ids em inglês (alinhar o grafo ao Go).
- **Decisão:** ids em inglês (`ark:breeder`, `rust:electrician`, `lol:support`...), que é o vocabulário de quem joga. Rótulos do formulário e nomes de persona permanecem em português.
- **Consequências:** `scripts/referencia.cypher` e `modelo-grafo.md` atualizados. `Cliente.ValidarFuncoes` faz o `seed` falhar antes de apagar a base quando uma função não existe no grafo. O erro foi detectado pela verificação por ida e volta do `cmd/verificar`.

---

## Fase 5 — Scoring e reciprocidade

## D-034 — Tipos de candidato no domínio e interface de fonte

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** a D-013 exige que `internal/matching` dependa de uma interface de fonte de candidatos, e não do Neo4j. Os tipos `Candidato` e `Opcoes` ficavam no pacote `graph`, que importa o driver: uma interface que os usasse faria o `matching` importar o `graph`.
- **Opções consideradas:** mover os tipos para `domain`; tipos próprios em `matching` com um adaptador que converte os de `graph`.
- **Decisão:** `domain.Candidato`, `domain.OpcoesCandidatos` e `domain.OpcoesPadrao()`. A interface `matching.Fonte` tem `BuscarCandidatos`, `BuscarPerfis` e `BuscarComplementaridade`.
- **Motivo:** com os tipos no domínio, o `*graph.Cliente` satisfaz `matching.Fonte` diretamente, porque em Go a implementação de uma interface é implícita (tipagem estrutural). Não há adaptador nem tipos duplicados. Os testes usam uma fonte em memória.
- **Consequências:** `graph.Opcoes`, `graph.OpcoesPadrao` e `graph.Candidato` deixam de existir. O `cmd/avaliar` contém `var _ matching.Fonte = (*graph.Cliente)(nil)`, que quebra a compilação se a interface e o cliente divergirem. A matriz de complementaridade é lida uma vez, em `matching.Novo`.

## D-035 — Fórmulas do score unilateral

- **Data:** 10/11/2026
- **Status:** Aceita (os valores numéricos são premissas heurísticas, como na D-010)
- **Contexto:** o plano definia os seis componentes do score s(A→B) e os pesos P01 a P06, mas não as fórmulas de cada componente.
- **Opções consideradas:**
  - **Função:** máximo entre os arranjos, com desconto para a secundária; média ponderada dos quatro arranjos.
  - **Toxicidade:** similaridade simétrica `1 − |tA − tB|/4`; penalidade assimétrica.
- **Decisão:** s(A→B) = Σ pᵢ·cᵢ / Σ pᵢ, com os pesos de A (peso ausente ou ≤ 0 vale 3, D-008). Assim, s fica em [0, 1]. Todos os componentes cᵢ estão em [0, 1]:

  | Componente | Fórmula |
  |---|---|
  | Horário (G06) | Jaccard das grades de 28 bits; duas grades vazias valem 0 |
  | Intensidade (A07/R07/L05, G10) | média de `1 − |iA − iB|/4` e do compromisso (iguais ou algum "tanto_faz" = 1; fixo × eventual = 0) |
  | Modo e servidor | ARK: média do servidor (iguais ou "tanto_faz" 1; dois não oficiais 0,5; oficial × não oficial 0) e do tamanho de tribo (escala ordinal). Rust: média do servidor (iguais 1; oficial × comunitário 0,5; resto 0), do limite de grupo e do wipe (escalas ordinais). LoL: Jaccard das filas, contando `solo_duo` em comum só se o duo é permitido pelo elo |
  | Comunicação (G07–G09) | fator de voz × tolerância. Voz: se A exige, B "sim" 1, "às vezes" 0,5, "não" 0; se A não exige, 1. Tolerância: `1 − max(0, tB − tA)/4` |
  | Função | maior valor de COMPLEMENTA entre os arranjos (principal ou secundária de cada um); cada jogador fora da principal multiplica por 0,75 |
  | Nível | d = nível de B − nível de A. "parecido": `max(0, 1 − 2|d|)`; "mais_experiente": `clamp(0,5 + 2d, 0, 1)`; "tanto_faz": 1 |

  Escala ordinal: `1 − distância/(n − 1)`; "tanto_faz" vale 1; valor fora da escala vale 0,5 (neutro).
- **Motivo:**
  - **Função pelo melhor arranjo:** na partida cada jogador assume uma função só, então o que importa é o melhor par possível, e não a média de todos. A média diluiria um par ótimo (ADC × Suporte) com arranjos que nunca vão acontecer.
  - **Toxicidade assimétrica:** A só se incomoda se B tolera mais trash talk que A. O desconforto de B com A já é medido em s(B→A), e a reciprocidade junta os dois lados. Com a similaridade simétrica, o mesmo desconforto entraria duas vezes.
  - **Nível:** um degrau de horas vale 0,25. Com "parecido", um degrau vale 0,5 e dois zeram o componente. Com "mais_experiente", um degrau acima já vale 1.
  - Horário e função são simétricos por construção; intensidade, comunicação e nível podem diferir entre A→B e B→A.
- **Consequências:** `internal/matching/score.go`. Os testes cobrem cada fórmula em tabela; as propriedades (componentes e scores em [0, 1], recíproco simétrico e entre os dois unilaterais); e o caso de controle. No caso de controle (2.000 usuários, semente 42), s(aprendiz→grinder) = 0,643, s(grinder→aprendiz) = 0,443, recíproco = 0,524, contra 0,704 de aprendiz × aprendiz. Os valores de 0,5 e 0,75 são heurísticos e podem ser calibrados na Fase 8.

## D-036 — Normalização do sinal social e combinação final

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** o sinal social (D-029) é um inteiro sem limite e precisa ficar em [0, 1] antes de ser combinado com o score recíproco.
- **Opções consideradas:** linear sobre o máximo (`s/max`); logarítmica sobre o máximo (`ln(1+s)/ln(1+max)`).
- **Decisão:** `social = ln(1+s)/ln(1+max)`, com o máximo tomado entre os candidatos da mesma consulta (0 se ninguém tem sinal). Score final = `0,8·recíproco + 0,2·social`. O empate é desfeito pelo id.
- **Motivo:** o sinal tem cauda longa. Na escala linear, um único candidato com sinal muito alto achata todos os outros perto de zero; o log amortece esse efeito. Sem sinal (usuário novo), o social é 0 para todos e o ranking passa a depender só do recíproco.
- **Consequências:** o social é relativo à consulta: o mesmo valor bruto vale mais ou menos dependendo dos outros candidatos. A ablação "sem sinal social" da Fase 7 mede o efeito (ameaça à validade da D-026). O motivo do match são os dois componentes com o maior valor mínimo entre os dois lados, para não explicar o par com algo que só agrada a um deles.

## D-037 — Precision@K simples e saturação do gabarito

- **Data:** 10/11/2026
- **Status:** Proposta (afeta o desenho da avaliação; levar ao orientador)
- **Contexto:** a Fase 5 fecha com uma versão simples de Precision@10. O `personas.md` já define como relevante o candidato que passa em todos os filtros rígidos e tem grau ≥ 1.
- **Opções consideradas:** relevância binária (grau ≥ 1); relevância graduada (NDCG com graus 1 e 2).
- **Decisão:** `cmd/avaliar` mede o Precision@K binário do motor completo e o da ordem do grafo (linha de base). O NDCG graduado e a ablação completa ficam para a Fase 7.
- **Motivo:** é o mínimo que o plano pede para a Fase 5, e a linha de base mostra se o score acrescenta algo além da ordem do grafo.
- **Consequências:** um experimento offline, com 3.000 usuários e sem o banco (os candidatos são os que passam no oráculo), mostrou o **Precision@10 saturado**:

  | Jogo | Só horários em comum | Score unilateral | Motor (recíproco) |
  |---|---|---|---|
  | ARK | 0,974 | 0,993 | 0,993 |
  | Rust | 0,881 | 0,984 | 0,956 |
  | LoL | 0,994 | 0,994 | 0,994 |

  Os filtros rígidos e as grades de horário de cada persona já separam as personas quase sozinhos. Assim, o Precision@10 binário não distingue as versões do motor e não serve como métrica principal da ablação. Na Fase 7: usar NDCG@10 graduado (grau 2 > grau 1) como métrica principal; olhar o resultado por persona; e considerar aumentar o ruído ou a sobreposição de horários entre personas. Isso é uma ameaça à validade e precisa ir para a monografia. O Rust recíproco abaixo do unilateral também precisa ser investigado na Fase 7.
- **Medição com o banco** (10/11/2026; 10 mil usuários, semente 42, amostra de 214 consultas do `cmd/avaliar`, opções padrão):

  | Versão | Geral | ARK | Rust | LoL |
  |---|---|---|---|---|
  | Ordem do grafo (horários, complemento, id) | 0,942 | 0,959 | 0,848 | 1,000 |
  | Motor completo | 0,981 | 0,998 | 0,933 | 1,000 |

  A saturação se confirma: o LoL fica em 1,000 nas duas versões, e só o Rust deixa algum espaço para medir diferença. Latência média de `Recomendar` (candidatos, perfis e scoring): 71,7 ms. No mesmo dia, o `cmd/verificar` terminou com 0 divergências, inclusive na complementaridade, e mediu `BuscarCandidatos` em p50 44,3 ms e p95 61,8 ms. Os outros ~27 ms vêm principalmente da leitura dos até 300 perfis; a Fase 7 separa os tempos de banco e de scoring.

---

## Fase 6 — API e formulário

## D-038 — Vocabulário do formulário e cálculo do nível no domínio

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** os valores válidos de cada resposta (`rates_altas`, `solo_duo`, `as_vezes`...) existiam só como strings soltas nas personas do `seed` e no `questionario.md`. A Fase 6 precisa deles em três lugares: validação da API, opções do formulário e escalas ordinais do score. O cálculo do nível (`calcularNivel`) ficava no `seed`, mas a API também precisa calculá-lo para quem se cadastra.
- **Opções consideradas:** repetir as listas em cada pacote; uma fonte única no `domain`.
- **Decisão:** `internal/domain/vocabulario.go` guarda cada pergunta como uma lista ordenada de `Opcao{Valor, Rotulo}`. `Perfil.Validar()` (`validacao.go`) aplica as regras do questionário. `PerfilJogo.CalcularNivel()` substitui a função do `seed`. As escalas ordinais do `matching` (tamanho de tribo, limite de grupo, wipe) usam `domain.Valores(...)`, que respeita a ordem da lista.
- **Motivo:** com uma lista só, o formulário não oferece um valor que a validação rejeita, e o score não usa uma ordem diferente da do formulário. O motor trata valor desconhecido como neutro (0,5), então um erro de digitação numa persona passaria despercebido. O novo `TestPerfisValidos` (`seed`) exige que todos os perfis gerados passem em `Validar()`.
- **Consequências:** a ordem das listas passa a ter significado (é a escala ordinal); reordenar uma lista muda o score. Os rótulos das funções continuam no grafo (D-018, D-033); os das demais perguntas ficam no `domain`.

## D-039 — Contrato da API REST

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** o plano previa `GET /usuarios/{id}/recomendacoes`, `POST /usuarios` e `POST /usuarios/{id}/curtidas`, sem detalhar corpo, erros, recusas e ids.
- **Opções consideradas:**
  - **Recusa:** campo `curtiu: bool` em `/curtidas`; ou uma rota própria `/recusas`.
  - **Id do usuário:** enviado pelo cliente; ou gerado pelo servidor.
  - **Avaliação repetida:** sobrescrever; ou ignorar em silêncio; ou recusar com 409.
- **Decisão:**

  | Rota | Sucesso | Erros |
  |---|---|---|
  | `GET /` | 200, formulário HTML | — |
  | `GET /saude` | 200 | 503 sem banco |
  | `GET /usuarios/{id}/recomendacoes?jogo=&k=` | 200, lista com nome, score, os dois scores unilaterais, social e motivos | 400 jogo ou `k` inválido (`k` de 1 a 50, padrão 10); 404 usuário; 422 não joga o jogo |
  | `POST /usuarios` | 201, `{"id"}` e `Location` | 400 com a lista de todos os problemas |
  | `POST /usuarios/{id}/curtidas` e `/recusas` (`{"para", "jogo"}`) | 201, `{"match": bool}` | 400; 404 se algum dos dois não existe ou não joga o jogo; 409 se já avaliou neste jogo |

  Os erros são sempre `{"erro", "detalhes"?}`. O id é gerado pelo servidor: `w` + 8 hexadecimais de `crypto/rand`, o que separa os usuários do formulário dos do seed (`u000001`). O JSON rejeita campos desconhecidos (`persona` não entra de fora) e corpos acima de 64 KiB. Erro interno vira 500 com mensagem genérica; a causa vai só para o log.
- **Motivo:** uma rota por ação é mais clara do que um booleano no corpo. O 409 torna explícito que o candidato some da lista depois de avaliado (o filtro "vistos" da D-029). O match mútuo é o que dá sentido ao "estilo match" e custa um `EXISTS` na mesma consulta.
- **Consequências:** usuários criados pelo formulário entram no grafo, e o `cmd/verificar` e o `cmd/avaliar` passam a falhar ou a medir errado até um novo `go run ./cmd/seed`, que apaga todos os `Usuario`. Duas curtidas simultâneas do mesmo par podem criar duas relações (não há trava); para um protótipo de usuário único, o risco foi aceito.

## D-040 — Criação de usuário em uma transação com conferência

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** a carga em lote (D-027) grava cada usuário em quatro consultas separadas e depende do MATCH dos nós de referência, que falha em silêncio (o bug que já aconteceu duas vezes). No seed, `ValidarFuncoes` e o `cmd/verificar` protegem a carga; um cadastro isolado pela API não tem essa proteção.
- **Opções consideradas:** reaproveitar `CarregarUsuarios` como está; uma consulta nova que crie tudo com MERGE; as mesmas consultas da carga dentro de uma transação, com conferência antes do commit.
- **Decisão:** `graph.CriarUsuario` roda as quatro consultas da carga em uma única transação de escrita (`ExecuteWrite`). Antes do commit, conta as relações `JOGA`, `EXERCE` e `DISPONIVEL` do novo usuário e compara com o esperado. Se faltar alguma, a transação é desfeita. A função secundária vazia deixa de gerar um item na carga.
- **Motivo:** a mesma consulta serve à carga e ao cadastro, então não há duas formas de gravar um perfil. A transação garante tudo ou nada, e a conferência transforma o MATCH silencioso em erro.
- **Consequências:** a validação da API já rejeita função inexistente no domínio. A conferência pega o caso em que o domínio e o grafo divergem (por exemplo, sem `cmd/migrar`). As curtidas usam o mesmo princípio: a consulta devolve uma linha só quando os dois usuários existem e jogam o jogo, e zero linhas viram 404.

## D-041 — Formulário como página única servida pelo Go

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** a D-011 pede um formulário representativo, não uma interface completa. A API só pode usar `net/http` (D-014).
- **Opções consideradas:** um front-end separado (React, Vue); uma página `html/template` com JavaScript puro, embutida no binário com `embed`.
- **Decisão:** `internal/api/pagina.html`, embutida com `//go:embed` e renderizada com `html/template` a partir do vocabulário do `domain` e dos rótulos das funções no grafo. O JavaScript monta o JSON do cadastro, chama a API, lista as recomendações e permite curtir e recusar, avisando quando há match. Os textos vindos de usuários só entram na página por `textContent`.
- **Motivo:** sem dependências nem build de front-end, e o binário continua sendo um só. O formulário é gerado das mesmas listas que a validação usa (D-038), então não há opção que a API rejeite. A página consome a mesma API pública, o que também a demonstra.
- **Consequências:** a validação fica no servidor; o formulário não marca campos obrigatórios e mostra a lista de erros que a API devolve. `api.Novo` falha se alguma função do domínio não tiver rótulo no grafo.

## D-042 — Perfil público e explicação do score

- **Data:** 10/11/2026
- **Status:** Aceita
- **Contexto:** a tela de recomendações mostrava só o nome, o score e dois motivos. Para decidir entre curtir e recusar, quem usa precisa ver quem é o candidato. Para a avaliação e a apresentação, também é útil ver por que o motor chegou àquele número.
- **Opções consideradas:**
  - **Perfil:** devolver o `domain.Perfil` inteiro; ou um tipo público próprio, com só os campos escolhidos.
  - **Explicação:** só os dois motivos; ou os seis componentes nos dois sentidos.
- **Decisão:**
  - `GET /usuarios/{id}` devolve um `perfilPublico`, com os valores já como rótulos: faixa etária, região, idiomas, voz, tolerância, compromisso, horários e, por jogo, funções, nível desejado e as respostas específicas. **Ficam de fora:** `persona` e `personaJogo` (regra inviolável), os pesos e as faixas etárias aceitas (preferências privadas, que só servem ao motor) e o nível normalizado.
  - A resposta de recomendações ganha `componentes`: os seis critérios com `paraVoce` (c(A→B)) e `paraCandidato` (c(B→A)).
  - Na página, "Ver perfil" abre o perfil, uma mini-grade com os horários em comum e as barras dos componentes.
- **Motivo:** um tipo público separado faz com que um campo novo no domínio não vaze sem alguém decidir que ele é público. Mostrar os componentes nos dois sentidos torna visível a assimetria (D-035) e a reciprocidade (D-015). É o que justifica um par com score alto para um lado e baixo para o outro.
- **Consequências:** qualquer pessoa com o id vê o perfil público; sem autenticação (trabalho futuro), isso é aceitável só no protótipo. O teste `TestPerfilPublico` falha se `persona`, `pesos` ou `faixasAceitas` aparecerem na resposta.

## D-043 — Exploração da base: lista de usuários, comparação de pares e filtros com motivo

- **Data:** 12/11/2026
- **Status:** Proposta
- **Contexto:** a página só mostrava as recomendações de um id digitado. Para a reunião com o orientador, para a banca e para o estudo com jogadores, é preciso navegar pela base e explicar um par qualquer, inclusive um que o motor descarta ("por que esses dois não deram match?"). O oráculo dos filtros (`seed.PassaFiltros`) só respondia sim ou não e morava no pacote do gerador de dados.
- **Opções consideradas:**
  - **Motivo dos filtros:** uma terceira implementação só para a tela; ou o próprio oráculo devolvendo cada filtro com o motivo.
  - **Explicação do par:** montar no handler; ou um método do `matching.Recomendador`, testável sem banco.
  - **Paginação da lista:** `SKIP`/`LIMIT`; ou por chave (`id > depois`).
- **Decisão:**
  - `domain.Filtros(a, b, jogo)` devolve todos os filtros rígidos, na ordem, com `Passou` e `Motivo` (por exemplo, "Solo/Duo, Flex × ARAM"). `domain.PassaFiltros` passa a ser o oráculo do `cmd/verificar`, do `cmd/avaliar` e do gerador de interações. Os testes do oráculo foram para `domain/filtros_test.go`.
  - `matching.Recomendador.Explicar(a, b, jogo)` junta os filtros, o mínimo de horários da consulta, o score nos dois sentidos (calculado mesmo quando o par é barrado, se B joga o jogo), os motivos e a posição de B no ranking de A.
  - Rotas novas:

    | Rota | Sucesso | Erros |
    |---|---|---|
    | `GET /usuarios?jogo=&regiao=&q=&origem=&depois=&limite=` | 200, `{"usuarios", "total", "proximo"?}` com perfis públicos, em ordem de id | 400 parâmetro inválido (`limite` de 1 a 100, padrão 50; `origem` = `formulario` ou `seed`; busca até 40 caracteres) |
    | `GET /usuarios/{id}/par/{outro}?jogo=` | 200, perfis públicos dos dois, filtros, horários em comum, scores, seis componentes, motivos, posição e tamanho do ranking | 400 jogo inválido ou ids iguais; 404 usuário; 422 A não joga o jogo |

  - A busca procura um trecho do nome ou o começo do id, sem diferenciar maiúsculas. A origem usa o prefixo do id (D-039).
- **Motivo:** com o oráculo explicando cada filtro, a tela e a verificação diferencial usam as mesmas regras em Go, e o Cypher continua independente (uma implementação nunca chama a outra). A explicação no `matching` reaproveita `Unilateral`, `Reciproco` e `Pontuar`, então os números da comparação são os mesmos do ranking (`TestExplicar`). A paginação por chave usa o índice da constraint `usuario_id` e não fica mais lenta nas últimas páginas.
- **Consequências:** a comparação roda o ranking inteiro de A (cerca de 70 ms, D-037). Quando B passa nos filtros mas não está no ranking, a tela não sabe se foi porque A já avaliou B ou por causa do corte de 300 candidatos (D-029), e diz as duas possibilidades. A lista expõe só perfis públicos (D-042); `TestListarUsuarios` e `TestPar` falham se aparecer persona, pesos ou faixas aceitas.

## D-044 — Página de exploração e cadastro em passos

- **Data:** 12/11/2026
- **Status:** Proposta
- **Contexto:** a página da D-041 era um formulário longo seguido de uma busca por id. Ela será a interface avaliada pelo SUS e a usada na demonstração, então precisa ser rápida de entender e de usar.
- **Opções consideradas:** manter a página e só melhorar o visual; reorganizar em uma mesa de trabalho com lista, jogador escolhido e comparação; adotar um framework de front-end.
- **Decisão:** a página continua única, embutida e sem dependências (D-041), agora em três colunas:
  - **Jogadores:** lista paginada do jogo escolhido, com busca, região e origem. Cada linha mostra a semana do jogador em miniatura; com um jogador escolhido, os horários em comum com ele aparecem na cor do jogo.
  - **Jogador A:** perfil público, grade de horários e os 20 primeiros do ranking, com os dois motivos de cada um.
  - **Comparação A × B:** veredito, score nos dois sentidos e recíproco, grade com os horários sobrepostos (a única animação da página), filtros com motivo, os seis critérios em barras espelhadas e as respostas lado a lado, com as diferentes em destaque. Dá para curtir ou recusar como A.
  - **Cadastro** em painel lateral com quatro passos (você, como joga, jogos, prioridades), valores padrão e grade de horários que se marca arrastando. Ao salvar, o novo jogador vira o A da tela.
  - O jogo, a busca e os jogadores escolhidos ficam na URL (o botão voltar funciona). Abaixo de 1280 px a comparação vira um painel lateral; abaixo de 760 px, uma coluna por vez.
  - Visual: base neutra e uma cor por jogo, que vira o acento da tela; fontes do sistema (Bahnschrift no Windows, em duas larguras), sem arquivos de fonte; tema claro e escuro.
- **Motivo:** a comparação de pares é o que explica o motor na apresentação e responde "por que não deu match". O cadastro em passos com valores padrão encurta o formulário sem mudar as perguntas (o vocabulário continua vindo do `domain`, D-038). Fontes do sistema mantêm o binário único e a página funcionando sem internet.
- **Consequências:** curtir e recusar pela tela de comparação mudam a base, como na D-039. A página depende de recursos de navegadores atuais (`:has`, `color-mix`, `dialog`). O teste `TestPagina` confere a grade de 28 horários do cadastro e da comparação, o vocabulário e o corte de candidatos levado ao script.
