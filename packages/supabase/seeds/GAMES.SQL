INSERT INTO games (name, slug, platform) VALUES
  ('Valorant',          'valorant',   'pc'),
  ('League of Legends', 'lol',        'pc'),
  ('CS2',               'cs2',        'pc')
ON CONFLICT (slug) DO NOTHING;