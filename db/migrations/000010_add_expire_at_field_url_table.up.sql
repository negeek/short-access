ALTER TABLE urls
ADD COLUMN expire_at TIMESTAMP DEFAULT now();