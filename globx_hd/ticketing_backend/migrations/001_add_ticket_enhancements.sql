-- Migration: Add ticket enhancement tables
-- Created: 2025-10-15
-- Description: Adds tables for ticket comments, calls, and activity tracking

-- Table for ticket comments and resolutions
CREATE TABLE IF NOT EXISTS ticket_comments (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('comment', 'resolution')),
    content TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_ticket_comments_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_comments_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Table for ticket call logs
CREATE TABLE IF NOT EXISTS ticket_calls (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    scheduled_by INTEGER NOT NULL,
    call_type VARCHAR(50) NOT NULL,
    party_name VARCHAR(255) NOT NULL,
    party_contact VARCHAR(255),
    purpose TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'SCHEDULED' CHECK (status IN ('SCHEDULED', 'COMPLETED', 'CANCELLED')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_ticket_calls_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_calls_scheduled_by FOREIGN KEY (scheduled_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Table for ticket activity history
CREATE TABLE IF NOT EXISTS ticket_activities (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL,
    user_id INTEGER NULL, -- Nullable for system activities
    activity_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_ticket_activities_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_activities_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_ticket_comments_ticket_id ON ticket_comments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_comments_created_at ON ticket_comments(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ticket_comments_type ON ticket_comments(type);

CREATE INDEX IF NOT EXISTS idx_ticket_calls_ticket_id ON ticket_calls(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_calls_scheduled_at ON ticket_calls(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_ticket_calls_status ON ticket_calls(status);

CREATE INDEX IF NOT EXISTS idx_ticket_activities_ticket_id ON ticket_activities(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_activities_created_at ON ticket_activities(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ticket_activities_activity_type ON ticket_activities(activity_type);

-- Update trigger for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to tables with updated_at columns
CREATE TRIGGER update_ticket_comments_updated_at 
    BEFORE UPDATE ON ticket_comments 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ticket_calls_updated_at 
    BEFORE UPDATE ON ticket_calls 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE ticket_comments IS 'Stores comments and resolutions for tickets';
COMMENT ON TABLE ticket_calls IS 'Stores scheduled calls and call logs for tickets';
COMMENT ON TABLE ticket_activities IS 'Stores activity history and audit trail for tickets';

COMMENT ON COLUMN ticket_comments.type IS 'Type of entry: comment or resolution';
COMMENT ON COLUMN ticket_comments.is_internal IS 'Whether comment is internal (not visible to customers)';
COMMENT ON COLUMN ticket_calls.call_type IS 'Type of call: OEM, Customer, Internal, etc.';
COMMENT ON COLUMN ticket_activities.activity_type IS 'Type of activity for categorization';
