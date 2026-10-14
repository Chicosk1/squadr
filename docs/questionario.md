# Questionário de Perfil

> **Fase:** 1 — Modelagem
> **Status:** Proposta (aguardando validação do orientador)

## Objetivo

Definir cada pergunta do formulário de cadastro, como ela é usada pelo motor de recomendação e de onde o dado poderá vir no futuro, quando houver integração com as APIs da Steam e da Riot.

## Legenda de uso

| Uso | Significado |
|---|---|
| **Filtro** | Elimina candidatos incompatíveis antes de qualquer cálculo |
| **Filtro recíproco** | Elimina candidatos se a condição falhar para qualquer um dos dois lados |
| **Grafo** | Vira nó ou relação percorrida na busca de candidatos |
| **Score** | Entra no cálculo de compatibilidade |
| **Score (assimétrico)** | Entra no cálculo, mas o valor de A→B pode ser diferente de B→A |
| **Peso** | Ajusta a importância de cada critério no score |

---

## Seção G — Perfil geral

Respondida por todos os usuários.

| ID | Pergunta | Resposta | Uso | Fonte futura |
|---|---|---|---|---|
| G01 | Você tem 18 anos ou mais? | Sim / Não | Acesso ao app | Formulário |
| G02 | Qual a sua faixa etária? | 18–24, 25–34, 35+ | Filtro | Formulário |
| G03 | Com quais faixas etárias você aceita jogar? | Múltipla escolha (18–24, 25–34, 35+) | Filtro recíproco | Formulário |
| G04 | Qual a sua região? | América do Sul, América do Norte, Europa | Filtro | Formulário / Riot |
| G05 | Quais idiomas você fala? | Português, Inglês, Espanhol (múltipla escolha) | Filtro (pelo menos 1 em comum) | Formulário |
| G06 | Quando você costuma jogar? | Grade 7 dias × 4 períodos | Grafo + Score | Formulário / histórico de partidas da Riot |
| G07 | Você usa chat de voz? | Sim / Às vezes / Não | Score | Formulário |
| G08 | Você exige que o parceiro use voz? | Sim / Não | Score (assimétrico) | Formulário |
| G09 | Quanto você tolera trash talk? | 1 (nada) a 5 (normal pra mim) | Score | Formulário |
| G10 | O que você procura? | Parceiro/grupo fixo, Partidas eventuais, Tanto faz | Score | Formulário |
| G11 | Quais jogos você joga? | ARK, Rust, League of Legends (múltipla escolha) | Grafo (abre as seções específicas) | Steam / Riot |

### Períodos da grade de horários (G06)

| Período | Faixa |
|---|---|
| Madrugada | 00h – 06h |
| Manhã | 06h – 12h |
| Tarde | 12h – 18h |
| Noite | 18h – 24h |

Todos os horários são registrados no horário de Brasília. Usuários em outros fusos ficam registrados como limitação do protótipo.

---

## Seção A — ARK

Respondida apenas por quem marcou ARK na G11.

| ID | Pergunta | Resposta | Uso | Fonte futura |
|---|---|---|---|---|
| A01 | Qual versão você joga? | Survival Evolved (ASE), Survival Ascended (ASA) | Filtro | Steam |
| A02 | Em qual plataforma? | PC Steam, PC Epic, Xbox, PlayStation | Filtro (tabela de crossplay) | Steam |
| A03 | Qual modo? | PvE, PvP, Ambos | Filtro (PvE e PvP são incompatíveis) | Formulário |
| A04 | Qual tipo de servidor? | Oficial, Não oficial rates baixas, Não oficial rates altas, Tanto faz | Score | Formulário |
| A05 | Quantas horas você tem no jogo? | <100, 100–500, 500–2000, 2000+ | Score (nível) | Steam |
| A06 | Quais são suas funções? (até 2, indicando a principal) | Breeder/domador, Builder, Coletor, PvP, Explorador | Grafo (complementaridade) | Formulário |
| A07 | Qual a sua intensidade? | 1 (casual) a 5 (competitivo) | Score | Formulário |
| A08 | Qual tamanho de tribo você procura? | 2–3, 4–6, 7+, Tanto faz | Score | Formulário |
| A09 | Qual nível você prefere no parceiro? | Parecido, Mais experiente, Tanto faz | Score (assimétrico) | Formulário |

---

## Seção R — Rust

Respondida apenas por quem marcou Rust na G11.

| ID | Pergunta | Resposta | Uso | Fonte futura |
|---|---|---|---|---|
| R01 | Qual versão você joga? | PC, Console Edition | Filtro | Steam |
| R02 | Qual tipo de servidor? | Oficial, Comunitário vanilla, Modificado, PvE | Score | Formulário |
| R03 | Qual limite de grupo você prefere? | Solo/Duo, Trio, Quad, Sem limite | Score | Formulário |
| R04 | Você joga no dia do wipe? | Sempre, Às vezes, Raramente | Score | Formulário |
| R05 | Quantas horas você tem no jogo? | <100, 100–500, 500–2000, 2000+ | Score (nível) | Steam |
| R06 | Quais são suas funções? (até 2, indicando a principal) | Builder, Coletor, PvP/roamer, Raider, Eletricista | Grafo (complementaridade) | Formulário |
| R07 | Qual a sua intensidade? | 1 (casual) a 5 (competitivo) | Score | Formulário |
| R08 | Qual nível você prefere no parceiro? | Parecido, Mais experiente, Tanto faz | Score (assimétrico) | Formulário |

---

## Seção L — League of Legends

Respondida apenas por quem marcou League of Legends na G11.

| ID | Pergunta | Resposta | Uso | Fonte futura |
|---|---|---|---|---|
| L01 | Qual o seu elo atual? | Sem ranked, Ferro, Bronze, Prata, Ouro, Platina, Esmeralda, Diamante, Mestre, Grão-Mestre, Desafiante | Score (nível) + filtro na fila Solo/Duo | Riot |
| L02 | Qual a sua função principal? | Top, Jungle, Mid, ADC, Suporte, Preencher | Grafo (complementaridade) | Riot (histórico de partidas) |
| L03 | Qual a sua função secundária? | Top, Jungle, Mid, ADC, Suporte, Preencher | Grafo (complementaridade) | Riot (histórico de partidas) |
| L04 | Quais filas você joga? | Solo/Duo, Flex, Normal, ARAM, Clash (múltipla escolha) | Score + filtro (nenhuma fila em comum elimina) | Riot |
| L05 | Qual a sua intensidade? | 1 (casual) a 5 (competitivo) | Score | Formulário |
| L06 | Qual nível você prefere no parceiro? | Parecido, Mais experiente, Tanto faz | Score (assimétrico) | Formulário |

### Conversão de elo para nível

| Elo | Valor | Elo | Valor |
|---|---|---|---|
| Sem ranked | — | Esmeralda | 5 |
| Ferro | 0 | Diamante | 6 |
| Bronze | 1 | Mestre | 7 |
| Prata | 2 | Grão-Mestre | 8 |
| Ouro | 3 | Desafiante | 9 |
| Platina | 4 | | |

O valor é normalizado para o intervalo 0–1 antes de entrar no score. O filtro de elo na fila Solo/Duo usa uma diferença máxima configurável; as regras de duo da Riot mudam com frequência e devem ser confirmadas antes de definir o valor.

---

## Seção P — Pesos

Opcional. Se o usuário não responder, todos os pesos assumem o valor padrão **3**.

> "Quanto cada item importa para você?" (1 = pouco, 5 = muito)

| ID | Critério |
|---|---|
| P01 | Horários em comum |
| P02 | Intensidade e compromisso |
| P03 | Modo e tipo de servidor |
| P04 | Comunicação (voz e tolerância a toxicidade) |
| P05 | Função complementar |
| P06 | Nível do parceiro |

---

## Mapeamento pergunta → componente do score

| Componente | Perguntas | Peso |
|---|---|---|
| Horário | G06 | P01 |
| Intensidade e compromisso | A07 / R07 / L05, G10 | P02 |
| Modo e servidor | A04, A08 / R02, R03, R04 / L04 | P03 |
| Comunicação | G07, G08, G09 | P04 |
| Função | A06 / R06 / L02, L03 | P05 |
| Nível | A05, A09 / R05, R08 / L01, L06 | P06 |

---

## Observações e limitações

- **Tamanho do formulário:** quem joga um único jogo responde cerca de 25 perguntas, o que é muito para um app no estilo Tinder. Por isso os pesos são opcionais, e as perguntas com fonte futura Steam/Riot sairão do formulário quando as APIs forem integradas.
- **Fuso horário:** todos os horários são considerados no horário de Brasília.
- **Steam e Epic:** jogadores de ARK na Epic não terão horas disponíveis pela API da Steam no futuro.
- **Perfis privados:** perfis privados na Steam não expõem jogos nem horas; o formulário continua sendo necessário como alternativa.