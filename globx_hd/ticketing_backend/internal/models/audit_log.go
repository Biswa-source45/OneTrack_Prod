package models

import "time"

// AuditLog represents a comprehensive audit trail entry
type AuditLog struct {
	ID uint64 `gorm:"primaryKey;column:id" json:"id"`

	// Actor Information (who performed the action)
	ActorID        *uint  `gorm:"column:actor_id" json:"actor_id"`
	ActorType      string `gorm:"column:actor_type;size:20" json:"actor_type"`
	ActorName      string `gorm:"column:actor_name;size:255" json:"actor_name"`
	ActorEmail     string `gorm:"column:actor_email;size:255" json:"actor_email"`
	ActorIPAddress string `gorm:"column:actor_ip_address;size:45" json:"actor_ip_address"`

	// Action Information (what was done)
	Action     string `gorm:"column:action;size:100;not null" json:"action"`
	EntityType string `gorm:"column:entity_type;size:50;not null" json:"entity_type"`
	EntityID   *uint  `gorm:"column:entity_id" json:"entity_id"`
	EntityName string `gorm:"column:entity_name;size:255" json:"entity_name"`

	// Change Details (what changed)
	Description    string  `gorm:"column:description;type:text;not null" json:"description"`
	OldValues      *string `gorm:"column:old_values;type:jsonb" json:"old_values"`
	NewValues      *string `gorm:"column:new_values;type:jsonb" json:"new_values"`
	ChangesSummary *string `gorm:"column:changes_summary;type:text" json:"changes_summary"`

	// Context Information (how it was done)
	HTTPMethod string `gorm:"column:http_method;size:10" json:"http_method"`
	Endpoint   string `gorm:"column:endpoint;size:255" json:"endpoint"`
	UserAgent  string `gorm:"column:user_agent;type:text" json:"user_agent"`
	RequestID  string `gorm:"column:request_id;size:100" json:"request_id"`

	// Metadata
	Severity     string  `gorm:"column:severity;size:20;default:info" json:"severity"`
	Status       string  `gorm:"column:status;size:20;default:success" json:"status"`
	ErrorMessage string  `gorm:"column:error_message;type:text" json:"error_message"`
	Metadata     *string `gorm:"column:metadata;type:jsonb" json:"metadata"`

	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Audit Action Constants - Authentication
const (
	AuditUserLoginSuccess    = "USER_LOGIN_SUCCESS"
	AuditUserLoginFailure    = "USER_LOGIN_FAILURE"
	AuditContactLoginSuccess = "CONTACT_LOGIN_SUCCESS"
	AuditContactLoginFailure = "CONTACT_LOGIN_FAILURE"
	AuditLogout              = "LOGOUT"
	AuditTokenRefresh        = "TOKEN_REFRESH"
	AuditPasswordReset       = "PASSWORD_RESET"
	AuditFirstLoginComplete  = "FIRST_LOGIN_COMPLETE"
)

// Audit Action Constants - User Management
const (
	AuditUserCreated            = "USER_CREATED"
	AuditUserUpdated            = "USER_UPDATED"
	AuditUserDeleted            = "USER_DELETED"
	AuditUserRoleChanged        = "USER_ROLE_CHANGED"
	AuditUserDesignationChanged = "USER_DESIGNATION_CHANGED"
)

// Audit Action Constants - Contact Management
const (
	AuditContactCreated        = "CONTACT_CREATED"
	AuditContactUpdated        = "CONTACT_UPDATED"
	AuditContactDeleted        = "CONTACT_DELETED"
	AuditContactTypeChanged    = "CONTACT_TYPE_CHANGED"
	AuditContactAccountChanged = "CONTACT_ACCOUNT_CHANGED"
)

// Audit Action Constants - Account Management
const (
	AuditAccountCreated      = "ACCOUNT_CREATED"
	AuditAccountUpdated      = "ACCOUNT_UPDATED"
	AuditAccountDeleted      = "ACCOUNT_DELETED"
	AuditAccountOwnerChanged = "ACCOUNT_OWNER_CHANGED"
)

// Audit Action Constants - Master Data
const (
	AuditProductCreated      = "PRODUCT_CREATED"
	AuditProductUpdated      = "PRODUCT_UPDATED"
	AuditProductDeleted      = "PRODUCT_DELETED"
	AuditRoleCreated         = "ROLE_CREATED"
	AuditRoleUpdated         = "ROLE_UPDATED"
	AuditRoleDeleted         = "ROLE_DELETED"
	AuditDesignationCreated  = "DESIGNATION_CREATED"
	AuditDesignationUpdated  = "DESIGNATION_UPDATED"
	AuditDesignationDeleted  = "DESIGNATION_DELETED"
	AuditProductIssueCreated = "PRODUCT_ISSUE_CREATED"
	AuditProductIssueUpdated = "PRODUCT_ISSUE_UPDATED"
	AuditProductIssueDeleted = "PRODUCT_ISSUE_DELETED"
)

// Audit Action Constants - Ticket Management
const (
	AuditTicketCreated         = "TICKET_CREATED"
	AuditTicketUpdated         = "TICKET_UPDATED"
	AuditTicketDeleted         = "TICKET_DELETED"
	AuditTicketStatusChanged   = "TICKET_STATUS_CHANGED"
	AuditTicketAssigned        = "TICKET_ASSIGNED"
	AuditTicketUnassigned      = "TICKET_UNASSIGNED"
	AuditTicketPriorityChanged = "TICKET_PRIORITY_CHANGED"
)

// Audit Action Constants - Task Management
const (
	AuditTaskCreated         = "TASK_CREATED"
	AuditTaskUpdated         = "TASK_UPDATED"
	AuditTaskDeleted         = "TASK_DELETED"
	AuditTaskStatusChanged   = "TASK_STATUS_CHANGED"
	AuditTaskAssigned        = "TASK_ASSIGNED"
	AuditTaskPriorityChanged = "TASK_PRIORITY_CHANGED"
)

// Audit Action Constants - Attachments
const (
	AuditAttachmentUploaded   = "ATTACHMENT_UPLOADED"
	AuditAttachmentDownloaded = "ATTACHMENT_DOWNLOADED"
	AuditAttachmentDeleted    = "ATTACHMENT_DELETED"
)

// Audit Action Constants - Comments
const (
	AuditCommentCreated = "COMMENT_CREATED"
	AuditCommentUpdated = "COMMENT_UPDATED"
	AuditCommentDeleted = "COMMENT_DELETED"
)

// Audit Action Constants - System
const (
	AuditEmailToTicket       = "EMAIL_TO_TICKET"
	AuditDumpedQueryCreated  = "DUMPED_QUERY_CREATED"
	AuditDumpedQueryResolved = "DUMPED_QUERY_RESOLVED"
	AuditDumpedQueryDeleted  = "DUMPED_QUERY_DELETED"
	AuditSystemInitialized   = "AUDIT_SYSTEM_INITIALIZED"
)

// Entity Type Constants
const (
	EntityTypeUser         = "user"
	EntityTypeContact      = "contact"
	EntityTypeAccount      = "account"
	EntityTypeTicket       = "ticket"
	EntityTypeTask         = "task"
	EntityTypeProduct      = "product"
	EntityTypeRole         = "role"
	EntityTypeDesignation  = "designation"
	EntityTypeProductIssue = "product_issue"
	EntityTypeAttachment   = "attachment"
	EntityTypeComment      = "comment"
	EntityTypeCall         = "call"
	EntityTypeApproval     = "approval"
	EntityTypeDumpedQuery  = "dumped_query"
	EntityTypeSystem       = "system"
)

// Actor Type Constants
const (
	ActorTypeUser    = "user"
	ActorTypeContact = "contact"
	ActorTypeSystem  = "system"
)

// Severity Constants
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// Status Constants
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusError   = "error"
)
