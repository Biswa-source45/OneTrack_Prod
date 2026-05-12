-- Migration: Recreate ticket_calls table with enhanced fields
-- Created: 2025-10-16
-- Description: Drop and recreate ticket_calls table with all enhanced fields from scratch

-- First, backup existing data (optional - uncomment if you want to preserve data)
-- CREATE TABLE ticket_calls_backup AS SELECT * FROM ticket_calls;

-- Drop existing table with all constraints and indexes
DROP TABLE IF EXISTS ticket_calls CASCADE;

-- Recreate ticket_calls table with enhanced structure
CREATE TABLE ticket_calls (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    scheduled_by INTEGER NOT NULL,
    
    -- Enhanced fields for advanced call logging
    subject VARCHAR(255),
    direction VARCHAR(20) CHECK (direction IN ('Inbound', 'Outbound')),
    status VARCHAR(20) DEFAULT 'Scheduled' CHECK (status IN ('Completed', 'Scheduled')),
    start_time TIMESTAMP WITH TIME ZONE,
    description TEXT,
    
    -- Legacy fields (for backward compatibility)
    call_type VARCHAR(50),
    party_name VARCHAR(255),
    party_contact VARCHAR(255),
    purpose TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_ticket_calls_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_calls_scheduled_by FOREIGN KEY (scheduled_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_ticket_calls_ticket_id ON ticket_calls(ticket_id);
CREATE INDEX idx_ticket_calls_scheduled_at ON ticket_calls(scheduled_at);
CREATE INDEX idx_ticket_calls_start_time ON ticket_calls(start_time);
CREATE INDEX idx_ticket_calls_direction ON ticket_calls(direction);
CREATE INDEX idx_ticket_calls_status ON ticket_calls(status);
CREATE INDEX idx_ticket_calls_created_at ON ticket_calls(created_at DESC);

-- Create trigger for updated_at
CREATE TRIGGER update_ticket_calls_updated_at 
    BEFORE UPDATE ON ticket_calls 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE ticket_calls IS 'Stores scheduled calls and call logs for tickets with enhanced logging capabilities';
COMMENT ON COLUMN ticket_calls.subject IS 'What the call is about';
COMMENT ON COLUMN ticket_calls.direction IS 'Direction of the call: Inbound or Outbound';
COMMENT ON COLUMN ticket_calls.status IS 'Current state of the call: Completed or Scheduled';
COMMENT ON COLUMN ticket_calls.start_time IS 'Start date and time of the call';
COMMENT ON COLUMN ticket_calls.description IS 'Optional call details or notes';
COMMENT ON COLUMN ticket_calls.call_type IS 'Legacy: Type of call: OEM, Customer, Internal, etc.';
COMMENT ON COLUMN ticket_calls.party_name IS 'Legacy: Name of the party involved in the call';
COMMENT ON COLUMN ticket_calls.party_contact IS 'Legacy: Contact information of the party';
COMMENT ON COLUMN ticket_calls.purpose IS 'Legacy: Purpose of the call';
COMMENT ON COLUMN ticket_calls.scheduled_at IS 'Legacy: When the call was scheduled';
COMMENT ON COLUMN ticket_calls.notes IS 'Legacy: Additional notes about the call';

-- If you want to restore data from backup (uncomment and modify as needed):
-- INSERT INTO ticket_calls (ticket_id, scheduled_by, call_type, party_name, party_contact, purpose, scheduled_at, notes, created_at, updated_at)
-- SELECT ticket_id, scheduled_by, call_type, party_name, party_contact, purpose, scheduled_at, notes, created_at, updated_at
-- FROM ticket_calls_backup;
