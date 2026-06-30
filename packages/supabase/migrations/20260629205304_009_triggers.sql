-- Quando um like é inserido, verifica se existe o like inverso.
-- Se sim, cria o match automaticamente.
CREATE OR REPLACE FUNCTION handle_new_swipe()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.direction = 'like' THEN
    IF EXISTS (
      SELECT 1 FROM swipes
      WHERE swiper_id = NEW.swiped_id
        AND swiped_id = NEW.swiper_id
        AND direction = 'like'
    ) THEN
      INSERT INTO matches (player_a, player_b)
      VALUES (
        LEAST(NEW.swiper_id, NEW.swiped_id),
        GREATEST(NEW.swiper_id, NEW.swiped_id)
      )
      ON CONFLICT DO NOTHING;
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_swipe_inserted
  AFTER INSERT ON swipes
  FOR EACH ROW EXECUTE FUNCTION handle_new_swipe();