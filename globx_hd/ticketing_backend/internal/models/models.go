package models

import "time"

// Master tables
type MasterProduct struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	ProductName        string    `gorm:"not null" json:"product_name"`
	ProductDescription string    `json:"product_description"`
	CreatedAt          time.Time `json:"created_at"`
}

type MasterUserDesignation struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DesignationName string    `gorm:"unique;not null" json:"designation_name"`
	CreatedAt       time.Time `json:"created_at"`
}

type MasterContactDesignation struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DesignationName string    `gorm:"unique;not null" json:"designation_name"`
	CreatedAt       time.Time `json:"created_at"`
}

type MasterRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleName  string    `gorm:"unique;not null" json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}

type MasterProductIssue struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	ProductID uint          `gorm:"not null" json:"product_id"`
	Product   MasterProduct `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	IssueName string        `gorm:"not null" json:"issue_name"`
	CreatedAt time.Time     `json:"created_at"`
}

// Core tables
type Account struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AccountName  string    `gorm:"not null" json:"account_name"`
	AccountOwner string    `json:"account_owner"`
	CustomerCode string    `gorm:"size:3;unique;not null" json:"customer_code"` // "001" - "999"
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Contacts     []Contact `gorm:"constraint:OnDelete:CASCADE" json:"contacts,omitempty"`
}

type Contact struct {
	ID            uint                     `gorm:"primaryKey" json:"id"`
	AccountID     *uint                    `json:"account_id"` // Made optional for Individual contacts
	Account       Account                  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	DesignationID uint                     `gorm:"not null" json:"designation_id"`
	Designation   MasterContactDesignation `gorm:"foreignKey:DesignationID" json:"designation,omitempty"`
	ContactType   string                   `gorm:"not null;size:20;check:contact_type IN ('Govt.','Private','Individual')" json:"contact_type"`
	Department    string                   `gorm:"size:100" json:"department"`
	Location      string                   `gorm:"size:100" json:"location"`
	FirstName     string                   `gorm:"not null" json:"first_name"`
	LastName      string                   `json:"last_name"`
	Email         string                   `gorm:"unique;not null" json:"email"`
	Mobile        string                   `json:"mobile"`
	PasswordHash  string                   `gorm:"not null" json:"-"`
	FirstLogin    bool                     `gorm:"default:true" json:"first_login"`
	CustomerCode  string                   `json:"customer_code"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

type User struct {
	ID            uint                  `gorm:"primaryKey" json:"id"`
	EmployeeID    string                `gorm:"unique;not null" json:"employee_id"`
	Username      string                `gorm:"unique;not null" json:"username"`
	PasswordHash  string                `gorm:"not null" json:"-"`
	FirstName     string                `gorm:"not null" json:"first_name"`
	LastName      string                `json:"last_name"`
	Email         string                `gorm:"unique;not null" json:"email"`
	Phone         string                `json:"phone"`
	DesignationID uint                  `gorm:"not null" json:"designation_id"`
	Designation   MasterUserDesignation `gorm:"foreignKey:DesignationID" json:"designation,omitempty"`
	RoleID        uint                  `gorm:"not null" json:"role_id"`
	Role          MasterRole            `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	FirstLogin    bool                  `gorm:"default:true" json:"first_login"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// Ticket model
type Ticket struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	TicketID         string        `gorm:"unique;not null;size:50" json:"ticket_id"`
	AccountID        *uint         `json:"account_id"` // Made optional for Individual contacts
	Account          Account       `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ContactID        uint          `gorm:"not null" json:"contact_id"`
	Contact          Contact       `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
	ProductID        uint          `gorm:"not null" json:"product_id"`
	Product          MasterProduct `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Subject          string        `gorm:"type:text;not null" json:"subject"`
	TicketDetails    string        `gorm:"type:text;not null" json:"ticket_details"`
	TicketStatus     string        `gorm:"not null;default:OPEN;size:50" json:"ticket_status"`
	Priority         string        `gorm:"not null;default:Medium;size:10" json:"priority"`
	Channel          string        `gorm:"not null;default:Phone;size:10;check:channel IN ('Phone','Mail')" json:"channel"`
	AssignedEngineer *uint         `gorm:"default:null" json:"assigned_engineer"`
	Engineer         User          `gorm:"foreignKey:AssignedEngineer" json:"engineer,omitempty"`
	CreatedAt        time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	// New relationships for enhanced ticket management
	Comments    []TicketComment    `gorm:"foreignKey:TicketID" json:"comments,omitempty"`
	Calls       []TicketCall       `gorm:"foreignKey:TicketID" json:"calls,omitempty"`
	Activities  []TicketActivity   `gorm:"foreignKey:TicketID" json:"activities,omitempty"`
	Approvals   []TicketApproval   `gorm:"foreignKey:TicketID" json:"approvals,omitempty"`
	Attachments []TicketAttachment `gorm:"foreignKey:TicketID;references:TicketID" json:"attachments,omitempty"`
}

// TicketAttachment model
type TicketAttachment struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TicketID         string    `gorm:"size:50;not null;index" json:"ticket_id"`
	OriginalFilename string    `gorm:"size:255;not null" json:"original_filename"`
	StoredFilename   string    `gorm:"size:255;not null" json:"stored_filename"`
	FilePath         string    `gorm:"size:500;not null" json:"file_path"`
	FileSize         int       `gorm:"not null" json:"file_size"`
	MimeType         string    `gorm:"size:100;not null" json:"mime_type"`
	UploadedBy       uint      `gorm:"not null" json:"uploaded_by"`
	Contact          Contact   `gorm:"foreignKey:UploadedBy" json:"contact,omitempty"`
	UploadedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"uploaded_at"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TicketComment model for comments and resolutions
type TicketComment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"not null;index" json:"ticket_id"`
	Ticket     Ticket    `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type       string    `gorm:"not null;size:20;check:type IN ('comment','resolution')" json:"type"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsInternal bool      `gorm:"default:false" json:"is_internal"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TicketCall model for call scheduling and logs
type TicketCall struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TicketID    uint   `gorm:"not null;index" json:"ticket_id"`
	Ticket      Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	ScheduledBy uint   `gorm:"not null" json:"scheduled_by"`
	User        User   `gorm:"foreignKey:ScheduledBy" json:"user,omitempty"`

	// Enhanced fields for advanced call logging
	Subject     string     `gorm:"size:255" json:"subject"`
	Direction   string     `gorm:"size:20;check:direction IN ('Inbound','Outbound')" json:"direction"`
	Status      string     `gorm:"not null;default:Open;size:20;check:status IN ('Open','In Progress','Completed')" json:"status"`
	StartTime   *time.Time `gorm:"index" json:"start_time"`
	Description string     `gorm:"type:text" json:"description"`
	CallType    string     `gorm:"size:50" json:"call_type"`

	// New fields for OEM tracking and mail content
	OEMTicketID string     `gorm:"size:255;index" json:"oem_ticket_id"`
	DueDate     *time.Time `gorm:"index" json:"due_date"`
	MailContent string     `gorm:"type:text" json:"mail_content"`

	// Close remarks (for closing calls with context)
	CloseRemarks string `gorm:"type:text" json:"close_remarks"`

	// Relationships
	Attachments []TicketCallAttachment `gorm:"foreignKey:CallID" json:"attachments,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TicketCallAttachment model for call log attachments
type TicketCallAttachment struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CallID           uint      `gorm:"not null;index" json:"call_id"`
	TicketID         string    `gorm:"size:50;not null;index" json:"ticket_id"`
	OriginalFilename string    `gorm:"size:255;not null" json:"original_filename"`
	StoredFilename   string    `gorm:"size:255;not null" json:"stored_filename"`
	FilePath         string    `gorm:"size:500;not null" json:"file_path"`
	FileSize         int64     `gorm:"not null" json:"file_size"`
	MimeType         string    `gorm:"size:100" json:"mime_type"`
	UploadedBy       uint      `gorm:"not null" json:"uploaded_by"`
	User             User      `gorm:"foreignKey:UploadedBy" json:"user,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TicketActivity model for activity history and audit trail
type TicketActivity struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TicketID     uint      `gorm:"not null;index" json:"ticket_id"`
	Ticket       Ticket    `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	UserID       *uint     `json:"user_id"` // Nullable for system activities
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActivityType string    `gorm:"not null;size:50;index" json:"activity_type"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	OldValue     string    `gorm:"type:text" json:"old_value"`
	NewValue     string    `gorm:"type:text" json:"new_value"`
	Remarks      string    `gorm:"type:text" json:"remarks"` // New field for status change remarks
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TicketApproval model for approval requests
type TicketApproval struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TicketID    uint      `gorm:"not null;index" json:"ticket_id"`
	Ticket      Ticket    `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	RequesterID uint      `gorm:"not null" json:"requester_id"`
	Requester   User      `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	ApproverID  uint      `gorm:"not null" json:"approver_id"`
	Approver    User      `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	Subject     string    `gorm:"not null;size:255" json:"subject"`
	Status      string    `gorm:"not null;default:PENDING;size:20;check:status IN ('PENDING','APPROVED','REJECTED')" json:"status"`
	Remarks     string    `gorm:"type:text" json:"remarks"` // Approver's remarks when approving/rejecting
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Task model
type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Subject      string     `gorm:"not null" json:"subject"`
	Description  string     `gorm:"type:text" json:"description"`
	DueDate      *time.Time `json:"due_date"`
	TaskStatus   string     `gorm:"not null;default:Not Started;size:50" json:"task_status"`
	Priority     string     `gorm:"not null;default:Medium;size:10" json:"priority"`
	AssignedTo   *uint      `gorm:"default:null" json:"assigned_to"`
	AssignedUser User       `gorm:"foreignKey:AssignedTo" json:"assigned_user,omitempty"`
	CreatedBy    uint       `gorm:"not null" json:"created_by"`
	Creator      User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TaskComment model
type TaskComment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TaskID     uint      `gorm:"not null" json:"task_id"`
	Task       Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsInternal bool      `gorm:"default:false" json:"is_internal"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TaskActivity model
type TaskActivity struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TaskID       uint      `gorm:"not null;index" json:"task_id"`
	Task         Task      `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	UserID       *uint     `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActivityType string    `gorm:"not null;size:50;index" json:"activity_type"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	OldValue     string    `gorm:"type:text" json:"old_value"`
	NewValue     string    `gorm:"type:text" json:"new_value"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// Activity type constants
const (
	ActivityTicketCreated     = "TICKET_CREATED"
	ActivityStatusChanged     = "STATUS_CHANGED"
	ActivityAssigned          = "ASSIGNED"
	ActivityUnassigned        = "UNASSIGNED"
	ActivityAssigneeChanged   = "ASSIGNEE_CHANGED"
	ActivityPriorityChanged   = "PRIORITY_CHANGED"
	ActivityCommentAdded      = "COMMENT_ADDED"
	ActivityResolutionAdded   = "RESOLUTION_ADDED"
	ActivityCallScheduled     = "CALL_SCHEDULED"
	ActivityCallCompleted     = "CALL_COMPLETED"
	ActivityCallCancelled     = "CALL_CANCELLED"
	ActivityTicketUpdated     = "TICKET_UPDATED"
	ActivityProductChanged    = "PRODUCT_CHANGED"
	ActivitySubjectChanged    = "SUBJECT_CHANGED"
	ActivityApprovalRequested = "APPROVAL_REQUESTED"
	ActivityApprovalApproved  = "APPROVAL_APPROVED"
	ActivityApprovalRejected  = "APPROVAL_REJECTED"
	ActivityContactChanged    = "CONTACT_CHANGED"
	ActivityTicketDeleted     = "TICKET_DELETED"
	// Task activity constants
	ActivityTaskCreated         = "TASK_CREATED"
	ActivityTaskUpdated         = "TASK_UPDATED"
	ActivityTaskStatusChanged   = "TASK_STATUS_CHANGED"
	ActivityTaskAssigneeChanged = "TASK_ASSIGNEE_CHANGED"
	ActivityTaskPriorityChanged = "TASK_PRIORITY_CHANGED"
	ActivityTaskCommentAdded    = "TASK_COMMENT_ADDED"
	ActivityTaskDeleted         = "TASK_DELETED"
)

// Notification model
type Notification struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	RecipientID      uint       `gorm:"not null;index" json:"recipient_id"`
	RecipientType    string     `gorm:"not null;size:20;check:recipient_type IN ('user','contact')" json:"recipient_type"`
	Title            string     `gorm:"not null;size:255" json:"title"`
	Message          string     `gorm:"type:text;not null" json:"message"`
	NotificationType string     `gorm:"not null;size:100;index" json:"notification_type"`
	RelatedID        *uint      `gorm:"index" json:"related_id"`
	RelatedType      string     `gorm:"size:50" json:"related_type"`
	RelatedSubID     *uint      `json:"related_sub_id"`
	ActorID          *uint      `json:"actor_id"`
	ActorType        string     `gorm:"size:20;check:actor_type IN ('user','contact','system')" json:"actor_type"`
	IsRead           bool       `gorm:"default:false;index" json:"is_read"`
	Priority         string     `gorm:"default:normal;size:20;check:priority IN ('low','normal','high','urgent')" json:"priority"`
	Category         string     `gorm:"default:general;size:50" json:"category"`
	Metadata         string     `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	CreatedAt        time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	ReadAt           *time.Time `json:"read_at"`

	// Optional relationships (loaded when needed)
	User    User    `gorm:"foreignKey:RecipientID" json:"user,omitempty"`
	Contact Contact `gorm:"foreignKey:RecipientID" json:"contact,omitempty"`
	Actor   User    `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

// NotificationTemplate model
type NotificationTemplate struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	NotificationType string    `gorm:"not null;unique;size:100" json:"notification_type"`
	TitleTemplate    string    `gorm:"not null;size:255" json:"title_template"`
	MessageTemplate  string    `gorm:"type:text;not null" json:"message_template"`
	DefaultPriority  string    `gorm:"default:normal;size:20;check:default_priority IN ('low','normal','high','urgent')" json:"default_priority"`
	Category         string    `gorm:"default:general;size:50" json:"category"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Notification type constants
const (
	// Customer Ticket Notifications
	NotificationTicketCreatedConfirmation = "ticket.created_confirmation"
	NotificationTicketEngineerAssigned    = "ticket.engineer_assigned"
	NotificationTicketStatusChanged       = "ticket.status_changed"
	NotificationTicketCommentAdded        = "ticket.comment_added"
	NotificationTicketResolutionAdded     = "ticket.resolution_added"
	NotificationTicketCallLogged          = "ticket.call_logged"
	NotificationTicketCallCompleted       = "ticket.call_completed"
	NotificationTicketPropertiesUpdated   = "ticket.properties_updated"
	NotificationTicketApprovalRequested   = "ticket.approval_requested"
	NotificationTicketApprovalApproved    = "ticket.approval_approved"
	NotificationTicketApprovalRejected    = "ticket.approval_rejected"

	// Manager Notifications
	NotificationTicketCreatedByCustomer     = "ticket.created_by_customer"
	NotificationTicketEngineerStatusUpdate  = "ticket.engineer_status_update"
	NotificationTaskStatusUpdatedByEngineer = "task.status_updated_by_engineer"
	NotificationTaskCommentAdded            = "task.comment_added"

	// Engineer Notifications
	NotificationTicketAssignedToYou        = "ticket.assigned_to_you"
	NotificationTicketUnassignedFromYou    = "ticket.unassigned_from_you"
	NotificationTicketCustomerCommentAdded = "ticket.customer_comment_added"
	NotificationTicketManagerCommentAdded  = "ticket.manager_comment_added"
	NotificationTaskAssignedToYou          = "task.assigned_to_you"
	NotificationTaskUnassignedFromYou      = "task.unassigned_from_you"
	NotificationTaskManagerCommentAdded    = "task.manager_comment_added"

	// System Notifications
	NotificationSystemMaintenance = "system.maintenance"
	NotificationSystemUpdate      = "system.update"
)

// DumpedQuery model for failed email-to-ticket attempts
type DumpedQuery struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	N8nID           string    `json:"n8n_id"`
	SenderEmail     string    `gorm:"not null" json:"sender_email"`
	SenderName      string    `json:"sender_name"`
	Subject         string    `gorm:"type:text" json:"subject"`
	Body            string    `gorm:"type:text" json:"body"`
	FailureReason   string    `gorm:"type:text" json:"failure_reason"`
	AIExtractedData string    `gorm:"type:jsonb" json:"ai_extracted_data"` // Stored as JSON string
	Status          string    `gorm:"default:OPEN" json:"status"`          // OPEN, RESOLVED, IGNORED
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
