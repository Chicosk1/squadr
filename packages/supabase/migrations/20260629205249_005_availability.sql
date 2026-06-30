CREATE TABLE availability (
  id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  weekday   int NOT NULL CHECK (weekday BETWEEN 0 AND 6),
  period    day_period NOT NULL,
  UNIQUE(player_id, weekday, period)
);

ALTER TABLE availability ENABLE ROW LEVEL SECURITY;
CREATE POLICY "availability_select" ON availability
  FOR SELECT TO authenticated USING (true);
CREATE POLICY "availability_insert" ON availability
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = player_id);
CREATE POLICY "availability_delete" ON availability
  FOR DELETE TO authenticated
  USING (auth.uid() = player_id);