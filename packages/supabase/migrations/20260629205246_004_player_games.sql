CREATE TABLE player_games (
  id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  game_id   uuid NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  rank      text,
  created_at timestamptz DEFAULT now(),
  UNIQUE(player_id, game_id)
);