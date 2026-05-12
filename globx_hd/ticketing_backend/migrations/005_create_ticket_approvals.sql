-- Create ticket_approvals table for approval requests
CREATE TABLE ticket_approvals (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    requester_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approver_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    remarks TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_ticket_approvals_ticket_id ON ticket_approvals(ticket_id);
CREATE INDEX idx_ticket_approvals_requester_id ON ticket_approvals(requester_id);
CREATE INDEX idx_ticket_approvals_approver_id ON ticket_approvals(approver_id);
CREATE INDEX idx_ticket_approvals_status ON ticket_approvals(status);
CREATE INDEX idx_ticket_approvals_created_at ON ticket_approvals(created_at);
