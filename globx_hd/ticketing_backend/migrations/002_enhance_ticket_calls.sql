-- Migration: Enhance ticket calls table for advanced call logging
-- Created: 2025-10-15
-- Description: Adds new fields for enhanced call logging system

-- Add new columns to ticket_calls table
ALTER TABLE ticket_calls 
ADD COLUMN IF NOT EXISTS subject VARCHAR(255),
ADD COLUMN IF NOT EXISTS direction VARCHAR(20) CHECK (direction IN ('Inbound', 'Outbound')),
ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS description TEXT;

-- Rename existing columns to match new requirements
-- purpose -> subject (if subject doesn't exist and purpose does)
DO $$ 
BEGIN
    -- Check if subject column doesn't exist but purpose does
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'ticket_calls' AND column_name = 'subject') 
       AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'ticket_calls' AND column_name = 'purpose') THEN
        -- Copy data from purpose to subject
        UPDATE ticket_calls SET subject = purpose WHERE purpose IS NOT NULL;
    END IF;
END $$;

-- Update status values to match new enum (Completed, Scheduled)
UPDATE ticket_calls 
SET status = CASE 
    WHEN status = 'COMPLETED' THEN 'Completed'
    WHEN status = 'SCHEDULED' THEN 'Scheduled'
    WHEN status = 'CANCELLED' THEN 'Scheduled' -- Map cancelled to scheduled for now
    ELSE 'Scheduled'
END;

-- Update status column constraint
ALTER TABLE ticket_calls 
DROP CONSTRAINT IF EXISTS ticket_calls_status_check,
ADD CONSTRAINT ticket_calls_status_check CHECK (status IN ('Completed', 'Scheduled'));

-- Add direction constraint
ALTER TABLE ticket_calls 
ADD CONSTRAINT ticket_calls_direction_check CHECK (direction IN ('Inbound', 'Outbound'));

-- Copy scheduled_at to start_time for existing records
UPDATE ticket_calls 
SET start_time = scheduled_at 
WHERE start_time IS NULL AND scheduled_at IS NOT NULL;

-- Copy notes to description for existing records
UPDATE ticket_calls 
SET description = notes 
WHERE description IS NULL AND notes IS NOT NULL;

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_ticket_calls_direction ON ticket_calls(direction);
CREATE INDEX IF NOT EXISTS idx_ticket_calls_status ON ticket_calls(status);
CREATE INDEX IF NOT EXISTS idx_ticket_calls_start_time ON ticket_calls(start_time);

-- Comments for documentation
COMMENT ON COLUMN ticket_calls.subject IS 'What the call is about';
COMMENT ON COLUMN ticket_calls.direction IS 'Direction of the call: Inbound or Outbound';
COMMENT ON COLUMN ticket_calls.status IS 'Current state of the call: Completed or Scheduled';
COMMENT ON COLUMN ticket_calls.start_time IS 'Start date and time of the call';
COMMENT ON COLUMN ticket_calls.description IS 'Optional call details or notes';
