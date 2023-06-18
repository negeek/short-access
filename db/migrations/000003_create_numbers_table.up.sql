CREATE TABLE IF NOT EXISTS numbers (
  id serial PRIMARY KEY,
  number INT NOT NULL UNIQUE,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now()
);