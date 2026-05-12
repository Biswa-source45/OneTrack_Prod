-- Add channel field to tickets table
-- This field will store the communication channel (Phone or Mail) for the ticket

ALTER TABLE tickets 
ADD COLUMN channel VARCHAR(10) NOT NULL DEFAULT 'Phone' 
CHECK (channel IN ('Phone', 'Mail'));
