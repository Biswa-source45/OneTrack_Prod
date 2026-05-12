package services

import (
	"fmt"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// ActivityService handles ticket activity logging
type ActivityService struct {
	db *gorm.DB
}

// NewActivityService creates a new activity service
func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

// LogActivity logs a new activity for a ticket
func (s *ActivityService) LogActivity(ticketID uint, userID *uint, activityType, description string) error {
	activity := models.TicketActivity{
		TicketID:     ticketID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
	}
	return s.db.Create(&activity).Error
}

// LogActivityWithRemarks logs a new activity for a ticket with remarks
func (s *ActivityService) LogActivityWithRemarks(ticketID uint, userID *uint, activityType, description, remarks string) error {
	activity := models.TicketActivity{
		TicketID:     ticketID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		Remarks:      remarks,
	}
	return s.db.Create(&activity).Error
}

// LogFieldChange logs a field change activity
func (s *ActivityService) LogFieldChange(ticketID uint, userID *uint, activityType, fieldName, oldValue, newValue string) error {
	description := fmt.Sprintf("%s changed from '%s' to '%s'", fieldName, oldValue, newValue)
	activity := models.TicketActivity{
		TicketID:     ticketID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		OldValue:     oldValue,
		NewValue:     newValue,
	}
	return s.db.Create(&activity).Error
}

// LogFieldChangeWithRemarks logs a field change activity with remarks
func (s *ActivityService) LogFieldChangeWithRemarks(ticketID uint, userID *uint, activityType, fieldName, oldValue, newValue, remarks string) error {
	description := fmt.Sprintf("%s changed from '%s' to '%s'", fieldName, oldValue, newValue)
	activity := models.TicketActivity{
		TicketID:     ticketID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		OldValue:     oldValue,
		NewValue:     newValue,
		Remarks:      remarks,
	}
	return s.db.Create(&activity).Error
}

// LogTicketCreation logs ticket creation activity
func (s *ActivityService) LogTicketCreation(ticketID uint, userID *uint, ticketIDStr string) error {
	description := fmt.Sprintf("Ticket %s created", ticketIDStr)
	return s.LogActivity(ticketID, userID, models.ActivityTicketCreated, description)
}

// LogStatusChange logs ticket status change
func (s *ActivityService) LogStatusChange(ticketID uint, userID *uint, oldStatus, newStatus string) error {
	return s.LogFieldChange(ticketID, userID, models.ActivityStatusChanged, "Status", oldStatus, newStatus)
}

// LogStatusChangeWithRemarks logs ticket status change with remarks
func (s *ActivityService) LogStatusChangeWithRemarks(ticketID uint, userID *uint, oldStatus, newStatus, remarks string) error {
	return s.LogFieldChangeWithRemarks(ticketID, userID, models.ActivityStatusChanged, "Status", oldStatus, newStatus, remarks)
}

// LogApprovalRequested logs approval request activity
func (s *ActivityService) LogApprovalRequested(ticketID uint, requesterID *uint, approverName, subject string) error {
	description := fmt.Sprintf("Approval requested from %s: %s", approverName, subject)
	return s.LogActivity(ticketID, requesterID, models.ActivityApprovalRequested, description)
}

// LogApprovalApproved logs approval approved activity
func (s *ActivityService) LogApprovalApproved(ticketID uint, approverID *uint, subject, remarks string) error {
	description := fmt.Sprintf("Approval approved: %s", subject)
	return s.LogActivityWithRemarks(ticketID, approverID, models.ActivityApprovalApproved, description, remarks)
}

// LogApprovalRejected logs approval rejected activity
func (s *ActivityService) LogApprovalRejected(ticketID uint, approverID *uint, subject, remarks string) error {
	description := fmt.Sprintf("Approval rejected: %s", subject)
	return s.LogActivityWithRemarks(ticketID, approverID, models.ActivityApprovalRejected, description, remarks)
}

// LogAssignment logs ticket assignment
func (s *ActivityService) LogAssignment(ticketID uint, userID *uint, engineerName string) error {
	description := fmt.Sprintf("Ticket assigned to %s", engineerName)
	return s.LogActivity(ticketID, userID, models.ActivityAssigned, description)
}

// LogUnassignment logs ticket unassignment
func (s *ActivityService) LogUnassignment(ticketID uint, userID *uint, engineerName string) error {
	description := fmt.Sprintf("Ticket unassigned from %s", engineerName)
	return s.LogActivity(ticketID, userID, models.ActivityUnassigned, description)
}

// LogPriorityChange logs priority change
func (s *ActivityService) LogPriorityChange(ticketID uint, userID *uint, oldPriority, newPriority string) error {
	return s.LogFieldChange(ticketID, userID, models.ActivityPriorityChanged, "Priority", oldPriority, newPriority)
}

// LogComment logs comment addition
func (s *ActivityService) LogComment(ticketID uint, userID *uint, commentType string) error {
	description := fmt.Sprintf("%s added", commentType)
	activityType := models.ActivityCommentAdded
	if commentType == "resolution" {
		activityType = models.ActivityResolutionAdded
	}
	return s.LogActivity(ticketID, userID, activityType, description)
}

// LogCallScheduled logs call scheduling
func (s *ActivityService) LogCallScheduled(ticketID uint, userID *uint, partyName string, scheduledAt time.Time) error {
	description := fmt.Sprintf("Call scheduled with %s for %s", partyName, scheduledAt.Format("2006-01-02 15:04"))
	return s.LogActivity(ticketID, userID, models.ActivityCallScheduled, description)
}

// LogCallCompleted logs call completion
func (s *ActivityService) LogCallCompleted(ticketID uint, userID *uint, partyName string) error {
	description := fmt.Sprintf("Call with %s completed", partyName)
	return s.LogActivity(ticketID, userID, models.ActivityCallCompleted, description)
}

// LogCallCancelled logs call cancellation
func (s *ActivityService) LogCallCancelled(ticketID uint, userID *uint, partyName string) error {
	description := fmt.Sprintf("Call with %s cancelled", partyName)
	return s.LogActivity(ticketID, userID, models.ActivityCallCancelled, description)
}

// LogProductChange logs product change
func (s *ActivityService) LogProductChange(ticketID uint, userID *uint, oldProduct, newProduct string) error {
	return s.LogFieldChange(ticketID, userID, models.ActivityProductChanged, "Product", oldProduct, newProduct)
}

// LogSubjectChange logs subject change
func (s *ActivityService) LogSubjectChange(ticketID uint, userID *uint, oldSubject, newSubject string) error {
	return s.LogFieldChange(ticketID, userID, models.ActivitySubjectChanged, "Subject", oldSubject, newSubject)
}

// GetTicketActivities retrieves activities for a ticket with pagination
func (s *ActivityService) GetTicketActivities(ticketID uint, limit, offset int) ([]models.TicketActivity, error) {
	var activities []models.TicketActivity
	err := s.db.Where("ticket_id = ?", ticketID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	return activities, err
}

// GetTicketActivityCount gets total count of activities for a ticket
func (s *ActivityService) GetTicketActivityCount(ticketID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.TicketActivity{}).Where("ticket_id = ?", ticketID).Count(&count).Error
	return count, err
}
