CREATE TABLE IF NOT EXISTS source (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT UNIQUE NOT NULL,
  domain        TEXT UNIQUE NOT NULL,
  icon_url      TEXT,
  updated_at  TEXT
);

ALTER TABLE manga
  ADD COLUMN source_id INTEGER;
