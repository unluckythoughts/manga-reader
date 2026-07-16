CREATE TABLE IF NOT EXISTS chapters (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  book_id     INTEGER,
  number      TEXT,
  content     TEXT,
  upload_date DATETIME,
  completed   BOOLEAN NOT NULL DEFAULT FALSE,
  downloaded  BOOLEAN NOT NULL DEFAULT FALSE,
  other_id    TEXT,
  created_at  DATETIME NOT NULL,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME,
  FOREIGN KEY (book_id) REFERENCES books(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS chapter_number_unique ON chapters(book_id, number);

CREATE TABLE IF NOT EXISTS books (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  url         TEXT NOT NULL,
  title       TEXT NOT NULL,
  type        TEXT NOT NULL CHECK(type IN ('manga', 'novel')),
  image_url   TEXT,
  synopsis    TEXT,
  slug        TEXT,
  other_id    TEXT,
  source_id   INTEGER,
  created_at  DATETIME NOT NULL,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME,
  FOREIGN KEY (source_id) REFERENCES sources(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS book_url_unique ON books(source_id, url);

CREATE TABLE IF NOT EXISTS sources (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          TEXT UNIQUE NOT NULL,
  domain        TEXT UNIQUE NOT NULL,
  icon_url      TEXT,
  created_at    DATETIME NOT NULL,
  updated_at    DATETIME NOT NULL,
  deleted_at    DATETIME
);

CREATE TABLE IF NOT EXISTS favorites (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER,
  book_id     INTEGER,
  progress    TEXT,
  categories  TEXT,
  created_at  DATETIME NOT NULL,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (book_id) REFERENCES books(id)
);

CREATE TABLE IF NOT EXISTS categories (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT,
  created_at  DATETIME NOT NULL,
  updated_at  DATETIME NOT NULL,
  deleted_at  DATETIME
);

CREATE TABLE IF NOT EXISTS users (
  id                  		INTEGER PRIMARY KEY AUTOINCREMENT,
  name                		TEXT NOT NULL,
  email               		TEXT,
  email_verified      		BOOLEAN NOT NULL DEFAULT FALSE,
  mobile              		TEXT,
  mobile_verified     		BOOLEAN NOT NULL DEFAULT FALSE,
  password            		TEXT NOT NULL,
  role                		INTEGER NOT NULL DEFAULT 1,
  google_id           		TEXT,
  google_avatar       		TEXT,
	google_refresh_token 		TEXT,
	google_token_expires_at DATETIME,
  created_at          		DATETIME NOT NULL,
  updated_at          		DATETIME NOT NULL,
  deleted_at          		DATETIME
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
-- Unique constraint only when email/mobile is not NULL and not empty
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(email) WHERE email IS NOT NULL AND email != '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_mobile_unique ON users(mobile) WHERE mobile IS NOT NULL AND mobile != '';

CREATE TABLE IF NOT EXISTS verify (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  target          TEXT NOT NULL UNIQUE,
  token           TEXT NOT NULL UNIQUE,
  verified        BOOLEAN NOT NULL DEFAULT FALSE,
  expires_at      DATETIME NOT NULL,
  created_at      DATETIME NOT NULL,
  updated_at      DATETIME NOT NULL,
  deleted_at      DATETIME
);
