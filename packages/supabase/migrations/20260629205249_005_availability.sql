CREATE TABLE availability (
  id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  weekday   int NOT NULL CHECK (weekday BETWEEN 0 AND 6),
  period    day_period NOT NULL,
  UNIQUE(player_id, weekday, period)
);