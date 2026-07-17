CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id            uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash text NOT NULL,
    email         VARCHAR(300) UNIQUE NOT NULL,
    date_created  TIMESTAMP DEFAULT now(),
    date_updated  TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS urls (
    id           serial PRIMARY KEY,
    original_url text NOT NULL,
    short_url    text NOT NULL UNIQUE,
    short_access text DEFAULT '',
    is_custom    boolean DEFAULT false,
    access_count int DEFAULT 0,
    expire_at    TIMESTAMP DEFAULT now(),
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS numbers (
    id           serial PRIMARY KEY,
    number       int NOT NULL UNIQUE,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id           serial PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    key_hash     text NOT NULL UNIQUE,
    name         text DEFAULT '',
    revoked      boolean DEFAULT false,
    expire_at    TIMESTAMP,
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP DEFAULT now()
);
