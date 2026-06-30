CREATE TABLE games (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name       text NOT NULL,
  slug       text UNIQUE NOT NULL,
  platform   text,
  cover_url  text
);

-- games: leitura pública, sem escrita pelo usuário
ALTER TABLE games ENABLE ROW LEVEL SECURITY;
CREATE POLICY "games_select" ON games
  FOR SELECT TO authenticated USING (true);