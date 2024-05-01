CREATE TABLE IF NOT EXISTS users
(
    id        INTEGER PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_username ON users (email);