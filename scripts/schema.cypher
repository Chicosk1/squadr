CREATE CONSTRAINT usuario_id IF NOT EXISTS FOR (u:Usuario) REQUIRE u.id IS UNIQUE;
CREATE CONSTRAINT jogo_id    IF NOT EXISTS FOR (j:Jogo)    REQUIRE j.id IS UNIQUE;
CREATE CONSTRAINT horario_id IF NOT EXISTS FOR (h:Horario) REQUIRE h.id IS UNIQUE;
CREATE CONSTRAINT funcao_id  IF NOT EXISTS FOR (f:Funcao)  REQUIRE f.id IS UNIQUE;
// usuario_regiao foi removido (D-032): o planejador nunca o usou e, forçado, piorou a consulta
DROP INDEX usuario_regiao IF EXISTS;