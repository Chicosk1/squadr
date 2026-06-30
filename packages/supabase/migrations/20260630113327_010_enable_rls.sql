-- Migration 010: habilita RLS e cria as policies em todas as tabelas.
-- As migrations 002 a 008 já continham esses comandos no arquivo,
-- mas como as tabelas já existiam no banco remoto quando o RLS foi
-- adicionado ao conteúdo, o push não reaplicou. Esta migration nova
-- garante que tudo seja executado de fato.

-- ============================================================
-- PROFILES
-- ============================================================
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "profiles_select" ON profiles;
CREATE POLICY "profiles_select" ON profiles
  FOR SELECT TO authenticated
  USING (true);

DROP POLICY IF EXISTS "profiles_update" ON profiles;
CREATE POLICY "profiles_update" ON profiles
  FOR UPDATE TO authenticated
  USING (auth.uid() = id);

-- ============================================================
-- GAMES
-- ============================================================
ALTER TABLE games ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "games_select" ON games;
CREATE POLICY "games_select" ON games
  FOR SELECT TO authenticated USING (true);

-- ============================================================
-- PLAYER_GAMES
-- ============================================================
ALTER TABLE player_games ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "player_games_select" ON player_games;
CREATE POLICY "player_games_select" ON player_games
  FOR SELECT TO authenticated USING (true);

DROP POLICY IF EXISTS "player_games_insert" ON player_games;
CREATE POLICY "player_games_insert" ON player_games
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = player_id);

DROP POLICY IF EXISTS "player_games_delete" ON player_games;
CREATE POLICY "player_games_delete" ON player_games
  FOR DELETE TO authenticated
  USING (auth.uid() = player_id);

-- ============================================================
-- AVAILABILITY
-- ============================================================
ALTER TABLE availability ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "availability_select" ON availability;
CREATE POLICY "availability_select" ON availability
  FOR SELECT TO authenticated USING (true);

DROP POLICY IF EXISTS "availability_insert" ON availability;
CREATE POLICY "availability_insert" ON availability
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = player_id);

DROP POLICY IF EXISTS "availability_delete" ON availability;
CREATE POLICY "availability_delete" ON availability
  FOR DELETE TO authenticated
  USING (auth.uid() = player_id);

-- ============================================================
-- SWIPES
-- ============================================================
ALTER TABLE swipes ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "swipes_select" ON swipes;
CREATE POLICY "swipes_select" ON swipes
  FOR SELECT TO authenticated
  USING (auth.uid() = swiper_id);

DROP POLICY IF EXISTS "swipes_insert" ON swipes;
CREATE POLICY "swipes_insert" ON swipes
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = swiper_id);

-- ============================================================
-- MATCHES
-- ============================================================
ALTER TABLE matches ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "matches_select" ON matches;
CREATE POLICY "matches_select" ON matches
  FOR SELECT TO authenticated
  USING (
    auth.uid() = player_a OR
    auth.uid() = player_b
  );

-- ============================================================
-- MESSAGES
-- ============================================================
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "messages_select" ON messages;
CREATE POLICY "messages_select" ON messages
  FOR SELECT TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM matches
      WHERE id = match_id
        AND (player_a = auth.uid() OR player_b = auth.uid())
    )
  );

DROP POLICY IF EXISTS "messages_insert" ON messages;
CREATE POLICY "messages_insert" ON messages
  FOR INSERT TO authenticated
  WITH CHECK (
    auth.uid() = sender_id
    AND EXISTS (
      SELECT 1 FROM matches
      WHERE id = match_id
        AND (player_a = auth.uid() OR player_b = auth.uid())
    )
  );

-- ============================================================
-- REPORTS
-- ============================================================
ALTER TABLE reports ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "reports_insert" ON reports;
CREATE POLICY "reports_insert" ON reports
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = reporter_id);