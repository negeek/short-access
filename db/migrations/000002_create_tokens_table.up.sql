CREATE TABLE IF NOT EXISTS tokens(
  id serial PRIMARY KEY,
  token text NOT NULL UNIQUE,
  user_id uuid NOT NULL,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);