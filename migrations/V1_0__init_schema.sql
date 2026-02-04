CREATE TABLE IF NOT EXISTS chapter (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  book_id     INTEGER,
  number      TEXT,
  image_urls  TEXT,
  upload_date DATETIME,
  completed   BOOLEAN NOT NULL DEFAULT FALSE,
  downloaded  BOOLEAN NOT NULL DEFAULT FALSE,
  other_id    TEXT,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS chapter_number_unique ON chapter(book_id, number);

CREATE TABLE IF NOT EXISTS book (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  type        TEXT NOT NULL CHECK(type IN ('manga', 'novel')),
  image_url   TEXT,
  synopsis    TEXT,
  slug        TEXT,
  other_id    TEXT,
  source_id   INTEGER,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS book_url_unique ON book(source_id, url);

CREATE TABLE IF NOT EXISTS source (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT UNIQUE NOT NULL,
  domain        TEXT UNIQUE NOT NULL,
  icon_url      TEXT,
  updated_at    DATETIME NOT NULL,
  deleted_at    DATETIME
);

CREATE TABLE IF NOT EXISTS favorite (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER,
  book_id     INTEGER,
  progress    TEXT,
  categories  TEXT,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME
);

CREATE TABLE IF NOT EXISTS category (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME
);

CREATE TABLE IF NOT EXISTS user (
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  name                TEXT NOT NULL,
  email               TEXT UNIQUE,
  email_verified      BOOLEAN NOT NULL DEFAULT FALSE,
  mobile              TEXT UNIQUE,
  mobile_verified     BOOLEAN NOT NULL DEFAULT FALSE,
  password            TEXT NOT NULL,
  role                INTEGER NOT NULL DEFAULT 1,
  google_id           TEXT,
  google_avatar       TEXT,
  updated_at          DATETIME NOT NULL,
  deleted_at          DATETIME
);
CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON user(deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);

CREATE TABLE IF NOT EXISTS verify (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  target          TEXT NOT NULL UNIQUE,
  token           TEXT NOT NULL UNIQUE,
  verified        BOOLEAN NOT NULL DEFAULT FALSE,
  expires_at      DATETIME NOT NULL,
  updated_at      DATETIME NOT NULL,
  deleted_at      DATETIME
);
