-- Migration: Add OEM Ticket ID, Due Date, Mail Content, and Attachments to Call Logging
-- Created: 2025-12-08
-- Description: Enhance ticket_calls table with OEM tracking, due dates, mail content, and attachment support

-- Step 1: Add new columns to ticket_calls table
ALTER TABLE ticket_calls 
ADD COLUMN IF NOT EXISTS oem_ticket_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS due_date TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS mail_content TEXT;

-- Step 2: Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_ticket_calls_oem_ticket_id ON ticket_calls(oem_ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_calls_due_date ON ticket_calls(due_date);

-- Step 3: Create ticket_call_attachments table
CREATE TABLE IF NOT EXISTS ticket_call_attachments (
    id SERIAL PRIMARY KEY,
    call_id INTEGER NOT NULL,
    ticket_id VARCHAR(50) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    stored_filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100),
    uploaded_by INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraints
    CONSTRAINT fk_call_attachments_call_id FOREIGN KEY (call_id) REFERENCES ticket_calls(id) ON DELETE CASCADE,
    CONSTRAINT fk_call_attachments_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets(ticket_id) ON DELETE CASCADE,
    CONSTRAINT fk_call_attachments_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 4: Create indexes for ticket_call_attachments
CREATE INDEX IF NOT EXISTS idx_call_attachments_call_id ON ticket_call_attachments(call_id);
CREATE INDEX IF NOT EXISTS idx_call_attachments_ticket_id ON ticket_call_attachments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_call_attachments_uploaded_by ON ticket_call_attachments(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_call_attachments_created_at ON ticket_call_attachments(created_at DESC);

-- Step 5: Add comments for documentation
COMMENT ON COLUMN ticket_calls.oem_ticket_id IS 'External OEM ticket reference ID';
COMMENT ON COLUMN ticket_calls.due_date IS 'Follow-up deadline for the call';
COMMENT ON COLUMN ticket_calls.mail_content IS 'Content of email related to this call';
COMMENT ON TABLE ticket_call_attachments IS 'Stores file attachments for call logs';
COMMENT ON COLUMN ticket_call_attachments.call_id IS 'Reference to the call log';
COMMENT ON COLUMN ticket_call_attachments.ticket_id IS 'Reference to the parent ticket';
COMMENT ON COLUMN ticket_call_attachments.file_size IS 'File size in bytes (max 3MB = 3145728 bytes)';
