CREATE TABLE messages (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  match_id   uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  sender_id  uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  content    text NOT NULL,
  read_at    timestamptz,
  created_at timestamptz DEFAULT now()
);

ALTER TABLE messages ENABLE ROW LEVEL SECURITY;

-- Usuário só lê mensagens de matches em que está
CREATE POLICY "messages_select" ON messages
  FOR SELECT TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM matches
      WHERE id = match_id
        AND (player_a = auth.uid() OR player_b = auth.uid())
    )
  );

-- Usuário só envia mensagem como ele mesmo e em match que está
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