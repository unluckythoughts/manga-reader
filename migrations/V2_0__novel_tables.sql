CREATE TABLE IF NOT EXISTS novel_chapter (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  novel_id    INTEGER,
  number      TEXT,
  other_id    TEXT,
  image_urls  TEXT,
  upload_date TEXT,
  completed   BOOLEAN NOT NULL DEFAULT FALSE,
  downloaded  BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at  TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS novel_chapter_url_unique ON novel_chapter(novel_id, url);

CREATE TABLE IF NOT EXISTS novel (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  source_id   INTEGER NOT NULL,
  slug        TEXT,
  other_id    TEXT,
  image_url   TEXT,
  synopsis    TEXT,
  updated_at  TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS novel_url_unique ON novel(source_id, url);

CREATE TABLE IF NOT EXISTS novel_source (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT UNIQUE NOT NULL,
  domain        TEXT UNIQUE NOT NULL,
  icon_url      TEXT,
  updated_at    TEXT
);

CREATE TABLE IF NOT EXISTS novel_favorite (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER,
  novel_id    INTEGER,
  progress    TEXT,
  categories  TEXT,
  updated_at  TEXT NOT NULL
);
