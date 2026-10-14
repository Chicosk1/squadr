UNWIND [
  {id: 'ark',  nome: 'ARK'},
  {id: 'rust', nome: 'Rust'},
  {id: 'lol',  nome: 'League of Legends'}
] AS jogo
MERGE (j:Jogo {id: jogo.id})
SET j.nome = jogo.nome;

UNWIND range(0, 6) AS d
UNWIND range(0, 3) AS p
WITH d, p,
     ['seg','ter','qua','qui','sex','sab','dom'][d] AS dia,
     ['madrugada','manha','tarde','noite'][p] AS periodo
MERGE (h:Horario {id: dia + '-' + periodo})
SET h.dia = dia, h.periodo = periodo, h.bit = d * 4 + p;

UNWIND [
  {id: 'ark:breeder',      jogo: 'ark',  nome: 'Breeder/domador'},
  {id: 'ark:builder',   jogo: 'ark',  nome: 'Builder'},
  {id: 'ark:coletor',      jogo: 'ark',  nome: 'Coletor'},
  {id: 'ark:pvp',          jogo: 'ark',  nome: 'PvP'},
  {id: 'ark:explorador',   jogo: 'ark',  nome: 'Explorador'},
  {id: 'rust:builder',  jogo: 'rust', nome: 'Builder'},
  {id: 'rust:coletor',     jogo: 'rust', nome: 'Coletor'},
  {id: 'rust:pvp',         jogo: 'rust', nome: 'PvP/roamer'},
  {id: 'rust:raider',      jogo: 'rust', nome: 'Raider'},
  {id: 'rust:eletricista', jogo: 'rust', nome: 'Eletricista'},
  {id: 'lol:top',          jogo: 'lol',  nome: 'Top'},
  {id: 'lol:jungle',       jogo: 'lol',  nome: 'Jungle'},
  {id: 'lol:mid',          jogo: 'lol',  nome: 'Mid'},
  {id: 'lol:adc',          jogo: 'lol',  nome: 'ADC'},
  {id: 'lol:suporte',      jogo: 'lol',  nome: 'Suporte'},
  {id: 'lol:preencher',    jogo: 'lol',  nome: 'Preencher'}
] AS funcao
MERGE (f:Funcao {id: funcao.id})
SET f.jogo = funcao.jogo, f.nome = funcao.nome;

UNWIND [
  // ARK
  {a: 'ark:breeder',    b: 'ark:breeder',    valor: 0.4},
  {a: 'ark:breeder',    b: 'ark:builder', valor: 0.9},
  {a: 'ark:breeder',    b: 'ark:coletor',    valor: 0.8},
  {a: 'ark:breeder',    b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:breeder',    b: 'ark:explorador', valor: 0.7},
  {a: 'ark:builder', b: 'ark:builder', valor: 0.4},
  {a: 'ark:builder', b: 'ark:coletor',    valor: 0.9},
  {a: 'ark:builder', b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:builder', b: 'ark:explorador', valor: 0.6},
  {a: 'ark:coletor',    b: 'ark:coletor',    valor: 0.6},
  {a: 'ark:coletor',    b: 'ark:pvp',        valor: 0.7},
  {a: 'ark:coletor',    b: 'ark:explorador', valor: 0.6},
  {a: 'ark:pvp',        b: 'ark:pvp',        valor: 0.8},
  {a: 'ark:pvp',        b: 'ark:explorador', valor: 0.7},
  {a: 'ark:explorador', b: 'ark:explorador', valor: 0.7},
  // Rust
  {a: 'rust:builder',  b: 'rust:builder',  valor: 0.4},
  {a: 'rust:builder',  b: 'rust:coletor',     valor: 0.9},
  {a: 'rust:builder',  b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:builder',  b: 'rust:raider',      valor: 0.7},
  {a: 'rust:builder',  b: 'rust:eletricista', valor: 0.8},
  {a: 'rust:coletor',     b: 'rust:coletor',     valor: 0.6},
  {a: 'rust:coletor',     b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:coletor',     b: 'rust:raider',      valor: 0.7},
  {a: 'rust:coletor',     b: 'rust:eletricista', valor: 0.7},
  {a: 'rust:pvp',         b: 'rust:pvp',         valor: 0.8},
  {a: 'rust:pvp',         b: 'rust:raider',      valor: 0.9},
  {a: 'rust:pvp',         b: 'rust:eletricista', valor: 0.6},
  {a: 'rust:raider',      b: 'rust:raider',      valor: 0.6},
  {a: 'rust:raider',      b: 'rust:eletricista', valor: 0.7},
  {a: 'rust:eletricista', b: 'rust:eletricista', valor: 0.3},
  // LoL
  {a: 'lol:top',       b: 'lol:top',       valor: 0.1},
  {a: 'lol:top',       b: 'lol:jungle',    valor: 0.8},
  {a: 'lol:top',       b: 'lol:mid',       valor: 0.5},
  {a: 'lol:top',       b: 'lol:adc',       valor: 0.5},
  {a: 'lol:top',       b: 'lol:suporte',   valor: 0.5},
  {a: 'lol:top',       b: 'lol:preencher', valor: 0.7},
  {a: 'lol:jungle',    b: 'lol:jungle',    valor: 0.1},
  {a: 'lol:jungle',    b: 'lol:mid',       valor: 0.9},
  {a: 'lol:jungle',    b: 'lol:adc',       valor: 0.7},
  {a: 'lol:jungle',    b: 'lol:suporte',   valor: 0.8},
  {a: 'lol:jungle',    b: 'lol:preencher', valor: 0.7},
  {a: 'lol:mid',       b: 'lol:mid',       valor: 0.1},
  {a: 'lol:mid',       b: 'lol:adc',       valor: 0.5},
  {a: 'lol:mid',       b: 'lol:suporte',   valor: 0.6},
  {a: 'lol:mid',       b: 'lol:preencher', valor: 0.7},
  {a: 'lol:adc',       b: 'lol:adc',       valor: 0.1},
  {a: 'lol:adc',       b: 'lol:suporte',   valor: 1.0},
  {a: 'lol:adc',       b: 'lol:preencher', valor: 0.7},
  {a: 'lol:suporte',   b: 'lol:suporte',   valor: 0.1},
  {a: 'lol:suporte',   b: 'lol:preencher', valor: 0.7},
  {a: 'lol:preencher', b: 'lol:preencher', valor: 0.5}
] AS par
MATCH (a:Funcao {id: par.a}), (b:Funcao {id: par.b})
MERGE (a)-[c:COMPLEMENTA]->(b)
SET c.valor = par.valor;