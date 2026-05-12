-- Add missing approval notification templates
-- These templates are required for the approval system notifications to work

INSERT INTO notification_templates (notification_type, title_template, message_template, default_priority, category) VALUES
-- Approval Notifications
('ticket.approval_requested', 'Approval Request', '{requester_name} has requested your approval for ticket #{ticket_id}: {subject}', 'high', 'ticket'),
('ticket.approval_approved', 'Approval Granted', '{approver_name} has approved your request: {subject}. Remarks: {remarks}', 'normal', 'ticket'),
('ticket.approval_rejected', 'Approval Rejected', '{approver_name} has rejected your request: {subject}. Remarks: {remarks}', 'normal', 'ticket')
ON CONFLICT (notification_type) DO UPDATE SET
    title_template = EXCLUDED.title_template,
    message_template = EXCLUDED.message_template,
    default_priority = EXCLUDED.default_priority,
    category = EXCLUDED.category,
    updated_at = CURRENT_TIMESTAMP;
