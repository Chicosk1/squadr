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
- **Consequências:** índice em `Usuario.regiao`.

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