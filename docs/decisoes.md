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
