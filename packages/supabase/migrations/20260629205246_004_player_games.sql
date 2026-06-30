CREATE TABLE player_games (
  id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  game_id   uuid NOT NULL REFERENCES games(id) ON DELETE CASCADE,
  rank      text,
  created_at timestamptz DEFAULT now(),
  UNIQUE(player_id, game_id)
);

ALTER TABLE player_games ENABLE ROW LEVEL SECURITY;
CREATE POLICY "player_games_select" ON player_games
  FOR SELECT TO authenticated USING (true);
CREATE POLICY "player_games_insert" ON player_games
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = player_id);
CREATE POLICY "player_games_delete" ON player_games
  FOR DELETE TO authenticated
  USING (auth.uid() = player_id);