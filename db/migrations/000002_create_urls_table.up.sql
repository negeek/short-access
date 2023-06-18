CREATE TABLE IF NOT EXISTS urls (
  id serial PRIMARY KEY,
  original_url text NOT NULL,
  short_url text NOT NULL UNIQUE,
  user_id uuid NOT NULL,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);