-- Migration: Add ticket attachments functionality
-- File: 005_add_ticket_attachments.sql
-- Description: Creates ticket_attachments table for file upload feature

-- Create ticket_attachments table
CREATE TABLE IF NOT EXISTS ticket_attachments (
    id SERIAL PRIMARY KEY,
    ticket_id VARCHAR(50) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    stored_filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by INTEGER NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraints
ALTER TABLE ticket_attachments 
ADD CONSTRAINT fk_ticket_attachments_ticket_id 
FOREIGN KEY (ticket_id) REFERENCES tickets(ticket_id) ON DELETE CASCADE;

ALTER TABLE ticket_attachments 
ADD CONSTRAINT fk_ticket_attachments_uploaded_by 
FOREIGN KEY (uploaded_by) REFERENCES contacts(id) ON DELETE CASCADE;

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_ticket_attachments_ticket_id ON ticket_attachments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_attachments_uploaded_by ON ticket_attachments(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_ticket_attachments_uploaded_at ON ticket_attachments(uploaded_at);

-- Add updated_at trigger (if the function exists)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_ticket_attachments_updated_at 
        BEFORE UPDATE ON ticket_attachments 
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Add comments for documentation
COMMENT ON TABLE ticket_attachments IS 'Stores file attachments for tickets';
COMMENT ON COLUMN ticket_attachments.ticket_id IS 'References the formatted ticket_id (e.g., 942-2025-003)';
COMMENT ON COLUMN ticket_attachments.original_filename IS 'Original filename as uploaded by user';
COMMENT ON COLUMN ticket_attachments.stored_filename IS 'Filename as stored on filesystem (ticket_id_filename)';
COMMENT ON COLUMN ticket_attachments.file_path IS 'Full relative path to the stored file';
COMMENT ON COLUMN ticket_attachments.file_size IS 'File size in bytes';
COMMENT ON COLUMN ticket_attachments.mime_type IS 'MIME type of the uploaded file';
COMMENT ON COLUMN ticket_attachments.uploaded_by IS 'Contact ID who uploaded the file';
