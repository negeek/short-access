ALTER TABLE urls
RENAME COLUMN created_at TO date_created;
ALTER TABLE urls
RENAME COLUMN updated_at TO date_updated;
ALTER TABLE urls
ALTER COLUMN date_created TYPE TIMESTAMP;
ALTER TABLE urls
ALTER COLUMN date_updated TYPE TIMESTAMP;