-- Migration to make email and password optional for contacts
ALTER TABLE contacts ALTER COLUMN email DROP NOT NULL;
ALTER TABLE contacts ALTER COLUMN password_hash DROP NOT NULL;
