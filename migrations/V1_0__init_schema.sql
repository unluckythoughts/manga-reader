CREATE TABLE IF NOT EXISTS chapter (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT UNIQUE NOT NULL,
  title       TEXT NOT NULL,
  manga_id    INTEGER,
  number      TEXT,
  image_urls  TEXT,
  upload_date TEXT,
  completed   BOOLEAN NOT NULL DEFAULT FALSE,
  downloaded  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS manga (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT UNIQUE NOT NULL,
  title       TEXT NOT NULL,
  image_url   TEXT,
  synopsis    TEXT
);

CREATE TABLE IF NOT EXISTS favorite (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER,
  manga_id    INTEGER,
  progress    TEXT,
  categories  TEXT
);

CREATE TABLE IF NOT EXISTS category (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT
);

CREATE TABLE IF NOT EXISTS user (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT
);