-- Estilo de jogo do jogador
CREATE TYPE play_style AS ENUM ('casual', 'competitive', 'chill');

-- Direção do swipe
CREATE TYPE swipe_direction AS ENUM ('like', 'dislike');

-- Período de disponibilidade
CREATE TYPE day_period AS ENUM ('morning', 'afternoon', 'night');