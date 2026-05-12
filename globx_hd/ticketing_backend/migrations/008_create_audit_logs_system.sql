-- ============================================================================
-- Migration: 008_create_audit_logs_system.sql
-- Description: Creates comprehensive audit logging system for all activities
-- Author: System
-- Date: 2025-01-12
-- ============================================================================

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    
    -- Actor Information (who performed the action)
    actor_id INTEGER,
    actor_type VARCHAR(20) CHECK (actor_type IN ('user', 'contact', 'system')),
    actor_name VARCHAR(255),
    actor_email VARCHAR(255),
    actor_ip_address VARCHAR(45),
    
    -- Action Information (what was done)
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER,
    entity_name VARCHAR(255),
    
    -- Change Details (what changed)
    description TEXT NOT NULL,
    old_values JSONB,
    new_values JSONB,
    changes_summary TEXT,
    
    -- Context Information (how it was done)
    http_method VARCHAR(10),
    endpoint VARCHAR(255),
    user_agent TEXT,
    request_id VARCHAR(100),
    
    -- Metadata
    severity VARCHAR(20) DEFAULT 'info' CHECK (severity IN ('info', 'warning', 'critical')),
    status VARCHAR(20) DEFAULT 'success' CHECK (status IN ('success', 'failure', 'error')),
    error_message TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_audit_actor ON audit_logs(actor_id, actor_type);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_severity ON audit_logs(severity);
CREATE INDEX IF NOT EXISTS idx_audit_status ON audit_logs(status);
CREATE INDEX IF NOT EXISTS idx_audit_actor_created ON audit_logs(actor_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_entity_created ON audit_logs(entity_type, entity_id, created_at DESC);

-- Create GIN index for JSONB columns for faster JSON queries
CREATE INDEX IF NOT EXISTS idx_audit_old_values_gin ON audit_logs USING GIN (old_values);
CREATE INDEX IF NOT EXISTS idx_audit_new_values_gin ON audit_logs USING GIN (new_values);
CREATE INDEX IF NOT EXISTS idx_audit_metadata_gin ON audit_logs USING GIN (metadata);

-- Create composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_audit_actor_action_created ON audit_logs(actor_id, action, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_entity_action_created ON audit_logs(entity_type, entity_id, action, created_at DESC);

-- Add comment to table
COMMENT ON TABLE audit_logs IS 'Comprehensive audit log for all system activities including authentication, CRUD operations, and system events';

-- Add comments to important columns
COMMENT ON COLUMN audit_logs.actor_id IS 'ID of the user/contact who performed the action';
COMMENT ON COLUMN audit_logs.actor_type IS 'Type of actor: user, contact, or system';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (e.g., USER_CREATED, TICKET_UPDATED)';
COMMENT ON COLUMN audit_logs.entity_type IS 'Type of entity affected (e.g., user, ticket, account)';
COMMENT ON COLUMN audit_logs.entity_id IS 'ID of the affected entity';
COMMENT ON COLUMN audit_logs.old_values IS 'JSON representation of values before the change';
COMMENT ON COLUMN audit_logs.new_values IS 'JSON representation of values after the change';
COMMENT ON COLUMN audit_logs.severity IS 'Severity level: info, warning, or critical';
COMMENT ON COLUMN audit_logs.status IS 'Operation status: success, failure, or error';

-- Create function to automatically partition audit logs by month (optional, for future scalability)
CREATE OR REPLACE FUNCTION create_audit_log_partition()
RETURNS TRIGGER AS $$
DECLARE
    partition_date TEXT;
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    partition_date := TO_CHAR(NEW.created_at, 'YYYY_MM');
    partition_name := 'audit_logs_' || partition_date;
    start_date := TO_CHAR(DATE_TRUNC('month', NEW.created_at), 'YYYY-MM-DD');
    end_date := TO_CHAR(DATE_TRUNC('month', NEW.created_at) + INTERVAL '1 month', 'YYYY-MM-DD');
    
    -- Check if partition exists, if not create it
    IF NOT EXISTS (
        SELECT 1 FROM pg_class WHERE relname = partition_name
    ) THEN
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF audit_logs FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Note: Partitioning is commented out by default for simplicity
-- Uncomment below to enable automatic monthly partitioning for better performance at scale
-- ALTER TABLE audit_logs RENAME TO audit_logs_old;
-- CREATE TABLE audit_logs (LIKE audit_logs_old INCLUDING ALL) PARTITION BY RANGE (created_at);
-- CREATE TRIGGER audit_log_partition_trigger BEFORE INSERT ON audit_logs FOR EACH ROW EXECUTE FUNCTION create_audit_log_partition();

-- Create view for recent audit logs (last 30 days) for faster queries
CREATE OR REPLACE VIEW recent_audit_logs AS
SELECT * FROM audit_logs
WHERE created_at >= CURRENT_TIMESTAMP - INTERVAL '30 days'
ORDER BY created_at DESC;

-- Create view for critical audit logs
CREATE OR REPLACE VIEW critical_audit_logs AS
SELECT * FROM audit_logs
WHERE severity = 'critical'
ORDER BY created_at DESC;

-- Create view for failed operations
CREATE OR REPLACE VIEW failed_audit_logs AS
SELECT * FROM audit_logs
WHERE status IN ('failure', 'error')
ORDER BY created_at DESC;

-- Create view for authentication events
CREATE OR REPLACE VIEW authentication_audit_logs AS
SELECT * FROM audit_logs
WHERE action IN (
    'USER_LOGIN_SUCCESS', 'USER_LOGIN_FAILURE',
    'CONTACT_LOGIN_SUCCESS', 'CONTACT_LOGIN_FAILURE',
    'LOGOUT', 'PASSWORD_RESET', 'TOKEN_REFRESH'
)
ORDER BY created_at DESC;

-- Create view for user management events
CREATE OR REPLACE VIEW user_management_audit_logs AS
SELECT * FROM audit_logs
WHERE entity_type IN ('user', 'contact')
AND action IN (
    'USER_CREATED', 'USER_UPDATED', 'USER_DELETED',
    'CONTACT_CREATED', 'CONTACT_UPDATED', 'CONTACT_DELETED'
)
ORDER BY created_at DESC;

-- Create view for master data changes
CREATE OR REPLACE VIEW master_data_audit_logs AS
SELECT * FROM audit_logs
WHERE entity_type IN ('product', 'role', 'designation', 'account', 'product_issue')
ORDER BY created_at DESC;

-- Grant permissions (adjust as needed for your security model)
-- GRANT SELECT ON audit_logs TO manager_role;
-- GRANT SELECT ON recent_audit_logs TO manager_role;
-- GRANT SELECT ON critical_audit_logs TO manager_role;

-- Insert initial system audit log
INSERT INTO audit_logs (
    actor_type,
    actor_name,
    action,
    entity_type,
    entity_name,
    description,
    severity,
    status
) VALUES (
    'system',
    'Database Migration',
    'AUDIT_SYSTEM_INITIALIZED',
    'system',
    'Audit Logs',
    'Audit logging system initialized and ready for use',
    'info',
    'success'
);

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- The audit_logs table is now ready to capture all system activities.
-- Next steps:
-- 1. Implement AuditService in Go backend
-- 2. Integrate audit logging into all handlers
-- 3. Create manager UI for viewing audit logs
-- ============================================================================
