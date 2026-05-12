-- Migration: Update call status values from Scheduled/Completed to Open/In Progress/Completed
-- Created: 2025-11-12
-- Description: Updates call status enum values to support Open, In Progress, and Completed

-- Update existing status values
UPDATE ticket_calls 
SET status = CASE 
    WHEN status = 'Scheduled' THEN 'Open'
    WHEN status = 'SCHEDULED' THEN 'Open'
    WHEN status = 'Completed' THEN 'Completed'
    WHEN status = 'COMPLETED' THEN 'Completed'
    ELSE 'Open'
END;

-- Update status column constraint to new values
ALTER TABLE ticket_calls 
DROP CONSTRAINT IF EXISTS ticket_calls_status_check;

ALTER TABLE ticket_calls 
ADD CONSTRAINT ticket_calls_status_check CHECK (status IN ('Open', 'In Progress', 'Completed'));

-- Update default value
ALTER TABLE ticket_calls 
ALTER COLUMN status SET DEFAULT 'Open';

-- Add index for new status values
DROP INDEX IF EXISTS idx_ticket_calls_status;
CREATE INDEX idx_ticket_calls_status ON ticket_calls(status);

-- Comments for documentation
COMMENT ON COLUMN ticket_calls.status IS 'Current state of the call: Open, In Progress, or Completed';
