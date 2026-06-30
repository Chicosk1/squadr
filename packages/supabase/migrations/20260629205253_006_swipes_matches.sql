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