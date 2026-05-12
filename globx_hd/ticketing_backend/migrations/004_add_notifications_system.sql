-- Migration: Add comprehensive notifications system
-- Created: 2025-10-25
-- Description: Adds notifications table for real-time user notifications across the ticketing system

-- Table for system notifications
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    
    -- Recipient Information
    recipient_id INTEGER NOT NULL,
    recipient_type VARCHAR(20) NOT NULL CHECK (recipient_type IN ('user', 'contact')),
    
    -- Notification Content
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    notification_type VARCHAR(100) NOT NULL, -- e.g., 'ticket.assigned', 'ticket.status_changed'
    
    -- Reference Data (what this notification is about)
    related_id INTEGER,                      -- ticket_id, task_id, etc.
    related_type VARCHAR(50),                -- 'ticket', 'task', 'call', 'comment'
    related_sub_id INTEGER,                  -- comment_id, call_id for nested items
    
    -- Actor Information (who performed the action)
    actor_id INTEGER,                        -- user_id or contact_id who triggered this notification
    actor_type VARCHAR(20) CHECK (actor_type IN ('user', 'contact', 'system')),
    
    -- Metadata
    is_read BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    category VARCHAR(50) DEFAULT 'general', -- 'ticket', 'task', 'communication', 'system'
    
    -- Additional Data (JSON for flexible metadata)
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP WITH TIME ZONE NULL
    
    -- Note: Foreign key constraints will be handled by application logic
    -- since we have polymorphic relationships (recipient can be user OR contact)
);

-- Performance Indexes
CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications(recipient_id, recipient_type, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_related ON notifications(related_type, related_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(notification_type);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(recipient_id, recipient_type) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_notifications_category ON notifications(category);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);

-- Partial index for unread notifications only (without time constraint to avoid IMMUTABLE function issues)
CREATE INDEX IF NOT EXISTS idx_notifications_active ON notifications(recipient_id, recipient_type, created_at DESC) 
    WHERE is_read = FALSE;

-- Composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_category_read ON notifications(recipient_id, recipient_type, category, is_read, created_at DESC);

-- Table comments for documentation
COMMENT ON TABLE notifications IS 'System-wide notifications for users and contacts';
COMMENT ON COLUMN notifications.recipient_id IS 'ID of the user or contact receiving the notification';
COMMENT ON COLUMN notifications.recipient_type IS 'Type of recipient: user (staff) or contact (customer)';
COMMENT ON COLUMN notifications.notification_type IS 'Specific type of notification for categorization and filtering';
COMMENT ON COLUMN notifications.related_id IS 'ID of the primary related entity (ticket, task, etc.)';
COMMENT ON COLUMN notifications.related_type IS 'Type of the primary related entity';
COMMENT ON COLUMN notifications.related_sub_id IS 'ID of secondary related entity (comment, call, etc.)';
COMMENT ON COLUMN notifications.actor_id IS 'ID of user/contact who triggered this notification';
COMMENT ON COLUMN notifications.actor_type IS 'Type of actor: user, contact, or system';
COMMENT ON COLUMN notifications.priority IS 'Notification priority level for UI treatment';
COMMENT ON COLUMN notifications.category IS 'High-level category for filtering and organization';
COMMENT ON COLUMN notifications.metadata IS 'Additional flexible data as JSON';

-- Function to automatically mark notification as read and set read_at timestamp
CREATE OR REPLACE FUNCTION mark_notification_read()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_read = TRUE AND OLD.is_read = FALSE THEN
        NEW.read_at = CURRENT_TIMESTAMP;
    ELSIF NEW.is_read = FALSE THEN
        NEW.read_at = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically set read_at when is_read changes to true
CREATE TRIGGER trigger_mark_notification_read
    BEFORE UPDATE ON notifications
    FOR EACH ROW
    EXECUTE PROCEDURE mark_notification_read();

-- Function for cleanup of old notifications (to be called by background job)
CREATE OR REPLACE FUNCTION cleanup_old_notifications(
    read_retention_days INTEGER DEFAULT 30,
    unread_retention_days INTEGER DEFAULT 90
) RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER := 0;
    temp_count INTEGER;
BEGIN
    -- Delete old read notifications
    DELETE FROM notifications 
    WHERE is_read = TRUE 
      AND read_at < (CURRENT_TIMESTAMP - (read_retention_days || ' days')::INTERVAL);
    
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;
    
    -- Archive or delete very old unread notifications
    DELETE FROM notifications 
    WHERE is_read = FALSE 
      AND created_at < (CURRENT_TIMESTAMP - (unread_retention_days || ' days')::INTERVAL);
    
    GET DIAGNOSTICS temp_count = ROW_COUNT;
    deleted_count := deleted_count + temp_count;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Create notification template constants table for consistent messaging
CREATE TABLE IF NOT EXISTS notification_templates (
    id SERIAL PRIMARY KEY,
    notification_type VARCHAR(100) NOT NULL UNIQUE,
    title_template VARCHAR(255) NOT NULL,
    message_template TEXT NOT NULL,
    default_priority VARCHAR(20) DEFAULT 'normal' CHECK (default_priority IN ('low', 'normal', 'high', 'urgent')),
    category VARCHAR(50) DEFAULT 'general',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default notification templates
INSERT INTO notification_templates (notification_type, title_template, message_template, default_priority, category) VALUES
-- Customer Ticket Notifications
('ticket.created_confirmation', 'Ticket Created Successfully', 'Your ticket #{ticket_id} has been created and is being reviewed.', 'normal', 'ticket'),
('ticket.engineer_assigned', 'Engineer Assigned to Your Ticket', 'Engineer {engineer_name} has been assigned to your ticket #{ticket_id}.', 'normal', 'ticket'),
('ticket.status_changed', 'Ticket Status Updated', 'Your ticket #{ticket_id} status changed from {old_status} to {new_status}.', 'normal', 'ticket'),
('ticket.comment_added', 'New Comment on Your Ticket', 'A new comment has been added to your ticket #{ticket_id}.', 'normal', 'communication'),
('ticket.resolution_added', 'Resolution Provided', 'A resolution has been provided for your ticket #{ticket_id}.', 'high', 'communication'),
('ticket.call_logged', 'Call Logged for Your Ticket', 'A {call_direction} call has been logged for your ticket #{ticket_id}.', 'normal', 'communication'),
('ticket.call_completed', 'Call Completed', 'The scheduled call for your ticket #{ticket_id} has been completed.', 'normal', 'communication'),
('ticket.properties_updated', 'Ticket Updated', 'Your ticket #{ticket_id} has been updated.', 'normal', 'ticket'),

-- Manager Notifications
('ticket.created_by_customer', 'New Ticket Created', 'Customer {customer_name} created ticket #{ticket_id} for {product_name}.', 'high', 'ticket'),
('ticket.engineer_status_update', 'Engineer Updated Ticket Status', 'Engineer {engineer_name} changed ticket #{ticket_id} status to {new_status}.', 'normal', 'ticket'),
('task.status_updated_by_engineer', 'Task Status Updated', 'Engineer {engineer_name} updated task "{task_subject}" status to {new_status}.', 'normal', 'task'),
('task.comment_added', 'New Task Comment', 'A comment has been added to task "{task_subject}".', 'normal', 'task'),

-- Engineer Notifications
('ticket.assigned_to_you', 'New Ticket Assigned', 'Ticket #{ticket_id} has been assigned to you by {manager_name}.', 'high', 'ticket'),
('ticket.unassigned_from_you', 'Ticket Unassigned', 'Ticket #{ticket_id} has been unassigned from you.', 'normal', 'ticket'),
('ticket.customer_comment_added', 'Customer Added Comment', 'Customer added a comment to your assigned ticket #{ticket_id}.', 'normal', 'communication'),
('ticket.manager_comment_added', 'Manager Added Comment', 'Manager added a comment to your assigned ticket #{ticket_id}.', 'normal', 'communication'),
('task.assigned_to_you', 'New Task Assigned', 'Task "{task_subject}" has been assigned to you.', 'high', 'task'),
('task.unassigned_from_you', 'Task Unassigned', 'Task "{task_subject}" has been unassigned from you.', 'normal', 'task'),
('task.manager_comment_added', 'Manager Added Comment', 'Manager added a comment to your assigned task "{task_subject}".', 'normal', 'communication'),

-- System Notifications
('system.maintenance', 'System Maintenance', 'System maintenance is scheduled for {maintenance_time}.', 'normal', 'system'),
('system.update', 'System Update', 'The ticketing system has been updated with new features.', 'low', 'system');

-- Apply updated_at trigger to notification_templates (check if function exists first)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_notification_templates_updated_at 
            BEFORE UPDATE ON notification_templates 
            FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    ELSE
        -- Create the function if it doesn't exist
        CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $func$
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $func$ LANGUAGE plpgsql;
        
        CREATE TRIGGER update_notification_templates_updated_at 
            BEFORE UPDATE ON notification_templates 
            FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
    END IF;
END
$$;

-- Create indexes for notification_templates
CREATE INDEX IF NOT EXISTS idx_notification_templates_type ON notification_templates(notification_type);
CREATE INDEX IF NOT EXISTS idx_notification_templates_active ON notification_templates(is_active);
CREATE INDEX IF NOT EXISTS idx_notification_templates_category ON notification_templates(category);

COMMENT ON TABLE notification_templates IS 'Templates for consistent notification messaging';
COMMENT ON COLUMN notification_templates.title_template IS 'Template for notification title with placeholders like {variable_name}';
COMMENT ON COLUMN notification_templates.message_template IS 'Template for notification message with placeholders';
