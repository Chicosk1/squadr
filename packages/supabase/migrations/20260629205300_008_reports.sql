CREATE TABLE reports (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  reported_id uuid NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  reason      text NOT NULL,
  created_at  timestamptz DEFAULT now(),
  CHECK(reporter_id != reported_id)
);

ALTER TABLE reports ENABLE ROW LEVEL SECURITY;
CREATE POLICY "reports_insert" ON reports
  FOR INSERT TO authenticated
  WITH CHECK (auth.uid() = reporter_id);