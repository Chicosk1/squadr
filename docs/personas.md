# Personas

> **Fase:** 1 — Modelagem
> **Status:** Proposta

## Objetivo

Definir os perfis-tipo usados para gerar os dados mock e o gabarito de relevância usado na avaliação do motor.

As personas existem por dois motivos:

1. **Gerar dados realistas:** em vez de perfis totalmente aleatórios, cada usuário mock é sorteado em torno de uma persona, com ruído.
2. **Ter um gabarito:** como sabemos a persona de cada usuário, sabemos quem *deveria* combinar com quem, e podemos medir se o motor acerta (Precision@K e NDCG@K).

> **Regra importante:** o campo `persona` existe apenas nos dados mock. O motor **nunca** pode usá-lo para recomendar; ele é usado somente pela avaliação.

---

## Distribuição da base

| Jogo | Personas | Proporção total |
|---|---|---|
| ARK | ARK-PVP-RAIDER, ARK-PVE-BREEDER, ARK-PVE-CONSTRUTOR | 35% |
| Rust | RUST-GRINDER, RUST-CASUAL, RUST-APRENDIZ | 30% |
| League of Legends | LOL-SUPORTE, LOL-ADC, LOL-ARAM | 30% |
| Multi-jogo | MULTI-CASUAL | 5% |

## Atributos gerais (todas as personas)

| Atributo | Distribuição |
|---|---|
| Região | América do Sul 80%, América do Norte 12%, Europa 8% |
| Idiomas (América do Sul) | Português 95%; adicionalmente Inglês 30%, Espanhol 15% |
| Idiomas (América do Norte) | Inglês 100%; adicionalmente Espanhol 20% |
| Idiomas (Europa) | Inglês 100% |
| Faixa etária | 18–24: 50%, 25–34: 35%, 35+: 15% |
| Faixas aceitas | Própria e vizinhas 70%, todas 20%, só a própria 10% |

---

## Personas de ARK

| Atributo | ARK-PVP-RAIDER | ARK-PVE-BREEDER | ARK-PVE-CONSTRUTOR |
|---|---|---|---|
| **Descrição** | Raider noturno | Criador de fim de semana | Construtor dedicado |
| **Proporção** | 12% | 12% | 11% |
| **Versão** | ASA | ASA | ASA 70%, ASE 30% |
| **Plataforma** | PC Steam 80%, console 20% | PC Steam | PC Steam |
| **Modo** | PvP | PvE | PvE |
| **Servidor** | Oficial | Não oficial rates altas | Não oficial rates baixas |
| **Intensidade** | 4–5 | 1–2 | 2–3 |
| **Horas** | 2000+ | 500–2000 | 100–500 |
| **Funções** | PvP (principal) + criador | Criador (principal) + coletor | Construtor (principal) + coletor |
| **Horários** | Todos os dias, noite e madrugada | Sex noite; sáb e dom tarde e noite | Qui e sex noite; sáb e dom tarde e noite |
| **Voz** | Usa e exige | Às vezes, não exige | Usa, não exige |
| **Tolerância a toxicidade** | 4–5 | 1–2 | 2–3 |
| **Tamanho de tribo** | 7+ | 2–3 | 4–6 |
| **Compromisso** | Fixo | Eventual | Fixo |
| **Nível desejado** | Parecido | Tanto faz | Tanto faz |
| **Pesos altos** | Horário 5, modo 5, função 4 | Modo 5, comunicação 4 | Função 5, horário 4 |

**Casos de teste embutidos:**

- Os 30% de ASE no ARK-PVE-CONSTRUTOR testam o filtro de versão: esses usuários nunca podem aparecer para quem joga ASA.
- ARK-PVP-RAIDER nunca pode aparecer para as personas PvE (filtro de modo).
- ARK-PVE-BREEDER e ARK-PVE-CONSTRUTOR são complementares: funções diferentes que se completam e horários com boa sobreposição.

---

## Personas de Rust

| Atributo | RUST-GRINDER | RUST-CASUAL | RUST-APRENDIZ |
|---|---|---|---|
| **Descrição** | Grinder de wipe | Construtor casual | Aprendiz |
| **Proporção** | 10% | 10% | 10% |
| **Versão** | PC | PC | PC 85%, console 15% |
| **Servidor** | Oficial | Comunitário vanilla | Comunitário vanilla ou modificado |
| **Limite de grupo** | Quad ou sem limite | Solo/Duo | Duo/Trio |
| **Wipe** | Sempre | Às vezes | Às vezes |
| **Intensidade** | 5 | 2 | 3 |
| **Horas** | 2000+ | 500–2000 | <100 |
| **Funções** | PvP (principal) + raider | Construtor (principal) + eletricista | Coletor (principal) + construtor |
| **Horários** | Qui tarde, noite e madrugada; sex a dom noite e madrugada | Seg a sex noite | Ter a sex noite; sáb tarde |
| **Voz** | Usa e exige | Usa, não exige | Usa, não exige |
| **Tolerância a toxicidade** | 5 | 2 | 3 |
| **Compromisso** | Fixo | Fixo | Fixo |
| **Nível desejado** | Parecido | Tanto faz | Mais experiente |
| **Pesos altos** | Horário 5, nível 5 | Horário 4, comunicação 4 | Nível 5, função 4 |

**Casos de teste embutidos:**

- **Par de controle de reciprocidade:** RUST-APRENDIZ quer alguém mais experiente e daria nota alta ao RUST-GRINDER. Porém o RUST-GRINDER quer alguém de nível parecido e dá nota baixa ao RUST-APRENDIZ. O motor com reciprocidade **não** deve recomendar esse par nas primeiras posições; o motor sem reciprocidade provavelmente vai recomendar.
- Os 15% de console no RUST-APRENDIZ testam o filtro de versão.
- RUST-APRENDIZ e RUST-CASUAL são compatíveis: o casual aceita qualquer nível e é mais experiente que o aprendiz.

---

## Personas de League of Legends

| Atributo | LOL-SUPORTE | LOL-ADC | LOL-ARAM |
|---|---|---|---|
| **Descrição** | Suporte ranqueado | ADC ranqueado | Casual de ARAM |
| **Proporção** | 10% | 10% | 10% |
| **Elo** | Ouro–Platina | Ouro–Platina | Sem ranked ou Ferro–Prata |
| **Função principal** | Suporte | ADC | Preencher |
| **Função secundária** | Mid | Mid | Qualquer (sorteada) |
| **Filas** | Solo/Duo, Flex | Solo/Duo, Flex | ARAM, Normal, Clash |
| **Intensidade** | 4 | 4 | 1–2 |
| **Horários** | Seg a sex noite; dom noite | Seg a sex noite; sáb noite | Sex a dom tarde e noite |
| **Voz** | Usa, não exige | Usa e exige | Às vezes, não exige |
| **Tolerância a toxicidade** | 2 | 3–4 | 2 |
| **Compromisso** | Fixo | Fixo | Eventual |
| **Nível desejado** | Parecido | Parecido | Tanto faz |
| **Pesos altos** | Função 5, nível 4 | Função 5, horário 4 | Horário 4, modo 4 |

**Casos de teste embutidos:**

- LOL-SUPORTE e LOL-ADC formam a dupla ideal de bot lane; o motor deve preferir esse par ao par de mesma persona.
- LOL-ARAM não tem filas em comum com as personas ranqueadas e deve ser eliminado pelo filtro de filas.

---

## Persona multi-jogo

### MULTI-CASUAL (5%)

- **Jogos:** dois dos três jogos, sorteados.
- **Perfil em cada jogo:** copia o perfil da persona casual daquele jogo (ARK-PVE-BREEDER, RUST-CASUAL ou LOL-ARAM), com intensidade 1–2.
- **Horários:** espalhados; cada período tem 30% de chance de ser marcado.
- **Voz:** às vezes, não exige.
- **Tolerância a toxicidade:** 2–3.
- **Compromisso:** eventual.

**Caso de teste embutido:** valida a recomendação contextualizada por jogo. No ARK, este usuário deve aparecer para criadores; no LoL, para jogadores de ARAM; e nunca para personas competitivas.

---

## Regras de ruído

Sem ruído, todos os usuários da mesma persona seriam idênticos e qualquer algoritmo acertaria. As regras abaixo tornam a avaliação honesta.

| Atributo | Regra |
|---|---|
| Horários | Cada período da persona é invertido (marcado ↔ desmarcado) com 10% de chance |
| Intensidade | Varia ±1 com 20% de chance, limitada entre 1 e 5 |
| Função secundária | Trocada por outra função do mesmo jogo com 15% de chance |
| Pesos | Cada peso varia ±1 com 30% de chance, limitado entre 1 e 5 |
| Outliers | 5% dos usuários recebem um atributo totalmente aleatório |
| Faixas de valores | Valores definidos como intervalo (ex.: intensidade 4–5) são sorteados uniformemente dentro do intervalo |

A geração usa **semente fixa**, para que todos os experimentos sejam reproduzíveis.

---

## Gabarito de relevância

Usado para calcular Precision@K e NDCG@K. **Só conta como relevante o candidato que passa em todos os filtros rígidos** (região, idioma, faixa etária, versão, modo, filas).

| Grau | Significado | Pares |
|---|---|---|
| **2** | Ideal | Mesma persona no ARK; mesma persona no Rust; LOL-SUPORTE ↔ LOL-ADC; ARK-PVE-BREEDER ↔ ARK-PVE-CONSTRUTOR |
| **1** | Bom | LOL-SUPORTE ↔ LOL-SUPORTE; LOL-ADC ↔ LOL-ADC; LOL-ARAM ↔ LOL-ARAM; RUST-APRENDIZ ↔ RUST-CASUAL; MULTI-CASUAL ↔ persona casual de cada jogo que ele joga |
| **0** | Controle (não deve ser recomendado) | RUST-APRENDIZ ↔ RUST-GRINDER; qualquer par PvE ↔ PvP; ASE ↔ ASA; PC ↔ console no Rust |

No LoL, o par complementar vale mais que a mesma persona: dois ADCs não formam uma boa dupla, e o gabarito precisa refletir isso.

Qualquer par não listado na tabela tem grau 0.

---

## Likes mock

Para alimentar o sinal social do grafo, cada usuário mock recebe likes e recusas coerentes com o gabarito.

| Regra | Valor |
|---|---|
| Candidatos avaliados por usuário | 20 a 50, sorteados entre os que passam nos filtros do jogo |
| Probabilidade de like (grau 2) | 70% |
| Probabilidade de like (grau 1) | 45% |
| Probabilidade de like (grau 0) | 15% |
| Demais avaliações | Viram `RECUSOU` |

Todos os likes e recusas registram o jogo em que aconteceram.