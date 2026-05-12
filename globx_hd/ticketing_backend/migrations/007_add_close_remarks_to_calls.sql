-- Migration: Add close_remarks to ticket_calls
-- Created: 2025-12-08
-- Description: Add close_remarks field to support closing calls with remarks

-- Step 1: Add close_remarks column to ticket_calls table
ALTER TABLE ticket_calls
ADD COLUMN IF NOT EXISTS close_remarks TEXT;

-- Step 2: Add comment for documentation
COMMENT ON COLUMN ticket_calls.close_remarks IS 'Remarks provided when closing a call (similar to ticket status change remarks)';
