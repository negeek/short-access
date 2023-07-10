ALTER TABLE numbers
RENAME COLUMN created_at TO date_created;
ALTER TABLE numbers
RENAME COLUMN updated_at TO date_updated;
ALTER TABLE numbers
ALTER COLUMN date_created TYPE TIMESTAMP;
ALTER TABLE numbers
ALTER COLUMN date_updated TYPE TIMESTAMP;