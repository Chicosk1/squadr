CREATE TABLE swipes (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  swiper_id  uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  swiped_id  uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  direction  swipe_direction NOT NULL,
  created_at timestamptz DEFAULT now(),
  UNIQUE(swiper_id, swiped_id),
  CHECK(swiper_id != swiped_id)
);

CREATE TABLE matches (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  player_a   uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  player_b   uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  created_at timestamptz DEFAULT now(),
  UNIQUE(player_a, player_b),
  CHECK(player_a != player_b)
);

ALTER TABLE swipes ENABLE ROW LEVEL SECURITY;

-- Usuário só vê seus próprios swipes
CREATE POLICY "swipes_select" ON swipes
  FOR SELECT TO authenticated
  USING (auth.uid() = swiper_id);

-- Usuário só insere swipes como ele mesmo
CREATE POLICY "swipes_insert" ON swipes
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = swiper_id);

ALTER TABLE matches ENABLE ROW LEVEL SECURITY;

-- Usuário só vê matches em que está envolvido
CREATE POLICY "matches_select" ON matches
  FOR SELECT TO authenticated
  USING (
    auth.uid() = player_a OR
    auth.uid() = player_b
  );

-- Matches só são criados pelo trigger — nenhum INSERT direto