CREATE TABLE IF NOT EXISTS chapter (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  manga_id    INTEGER,
  number      TEXT,
  image_urls  TEXT,
  upload_date TEXT,
  completed   BOOLEAN NOT NULL DEFAULT FALSE,
  downloaded  BOOLEAN NOT NULL DEFAULT FALSE,
  other_id    TEXT,
  updated_at  TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS chapter_number_unique ON chapter(manga_id, number);

CREATE TABLE IF NOT EXISTS manga (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  image_url   TEXT,
  synopsis    TEXT,
  slug        TEXT,
  other_id    TEXT,
  source_id   INTEGER,
  updated_at  TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS manga_title_unique ON manga(source_id, title);

CREATE TABLE IF NOT EXISTS source (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT UNIQUE NOT NULL,
  domain        TEXT UNIQUE NOT NULL,
  icon_url      TEXT,
  updated_at    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS favorite (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER,
  manga_id    INTEGER,
  progress    TEXT,
  categories  TEXT,
  updated_at  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS category (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT,
  updated_at  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT,
  updated_at  TEXT NOT NULL
);