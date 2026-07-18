package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditService handles comprehensive audit logging
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

// Log creates a basic audit log entry
func (s *AuditService) Log(log models.AuditLog) error {
	return s.db.Create(&log).Error
}

// LogWithContext creates an audit log with request context from Gin
func (s *AuditService) LogWithContext(
	c *gin.Context,
	action string,
	entityType string,
	entityID *uint,
	entityName string,
	description string,
	oldValues interface{},
	newValues interface{},
) error {
	log.Printf("[AUDIT] LogWithContext - Action: %s, Entity: %s, Name: %s", action, entityType, entityName)

	auditLog := models.AuditLog{
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		EntityName:  entityName,
		Description: description,
		Severity:    models.SeverityInfo,
		Status:      models.StatusSuccess,
	}

	// Extract actor information from context
	s.extractActorInfo(c, &auditLog)
	log.Printf("[AUDIT] Actor extracted - Type: %s, Name: %s", auditLog.ActorType, auditLog.ActorName)

	// Extract request context
	s.extractRequestContext(c, &auditLog)

	// Serialize old and new values (use pointers to avoid empty string issues with JSONB)
	if oldValues != nil {
		if oldJSON, err := json.Marshal(oldValues); err == nil {
			jsonStr := string(oldJSON)
			auditLog.OldValues = &jsonStr
		}
	}
	if newValues != nil {
		if newJSON, err := json.Marshal(newValues); err == nil {
			jsonStr := string(newJSON)
			auditLog.NewValues = &jsonStr
		}
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		log.Printf("[AUDIT ERROR] Failed to create audit log: %v", err)
		return err
	}

	log.Printf("[AUDIT SUCCESS] Created audit log ID: %d", auditLog.ID)
	return nil
}

// LogAuthentication logs authentication events (login, logout, password reset)
func (s *AuditService) LogAuthentication(
	actorType string,
	actorID *uint,
	actorName string,
	actorEmail string,
	action string,
	success bool,
	ipAddress string,
	userAgent string,
	errorMessage string,
) error {
	log.Printf("[AUDIT] Logging authentication - Actor: %s, Action: %s, Success: %v", actorName, action, success)

	auditLog := models.AuditLog{
		ActorID:        actorID,
		ActorType:      actorType,
		ActorName:      actorName,
		ActorEmail:     actorEmail,
		ActorIPAddress: ipAddress,
		Action:         action,
		EntityType:     "authentication",
		Description:    fmt.Sprintf("%s: %s", action, actorEmail),
		UserAgent:      userAgent,
		Severity:       models.SeverityInfo,
	}

	if success {
		auditLog.Status = models.StatusSuccess
	} else {
		auditLog.Status = models.StatusFailure
		auditLog.ErrorMessage = errorMessage
		auditLog.Severity = models.SeverityWarning
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		log.Printf("[AUDIT ERROR] Failed to create audit log: %v", err)
		return err
	}

	log.Printf("[AUDIT SUCCESS] Created audit log ID: %d", auditLog.ID)
	return nil
}

// LogCRUD logs Create, Read, Update, Delete operations
func (s *AuditService) LogCRUD(
	c *gin.Context,
	action string,
	entityType string,
	entityID *uint,
	entityName string,
	description string,
	oldValues interface{},
	newValues interface{},
) error {
	return s.LogWithContext(c, action, entityType, entityID, entityName, description, oldValues, newValues)
}

// LogError logs an error/failure event
func (s *AuditService) LogError(
	c *gin.Context,
	action string,
	entityType string,
	description string,
	errorMessage string,
) error {
	log := models.AuditLog{
		Action:       action,
		EntityType:   entityType,
		Description:  description,
		ErrorMessage: errorMessage,
		Severity:     models.SeverityCritical,
		Status:       models.StatusError,
	}

	s.extractActorInfo(c, &log)
	s.extractRequestContext(c, &log)

	return s.db.Create(&log).Error
}

// LogSystemEvent logs system-level events
func (s *AuditService) LogSystemEvent(
	action string,
	entityType string,
	description string,
	metadata map[string]interface{},
) error {
	log := models.AuditLog{
		ActorType:   models.ActorTypeSystem,
		ActorName:   "System",
		Action:      action,
		EntityType:  entityType,
		Description: description,
		Severity:    models.SeverityInfo,
		Status:      models.StatusSuccess,
	}

	if metadata != nil {
		if metaJSON, err := json.Marshal(metadata); err == nil {
			jsonStr := string(metaJSON)
			log.Metadata = &jsonStr
		}
	}

	return s.db.Create(&log).Error
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (s *AuditService) GetAuditLogs(filters AuditLogFilters) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{})

	// Apply filters
	if filters.ActorID != nil {
		query = query.Where("actor_id = ?", *filters.ActorID)
	}
	if filters.ActorType != "" {
		query = query.Where("actor_type = ?", filters.ActorType)
	}
	if filters.Action != "" {
		query = query.Where("action = ?", filters.Action)
	}
	if filters.EntityType != "" {
		query = query.Where("entity_type = ?", filters.EntityType)
	}
	if filters.EntityID != nil {
		query = query.Where("entity_id = ?", *filters.EntityID)
	}
	if filters.Severity != "" {
		query = query.Where("severity = ?", filters.Severity)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where(
			"description ILIKE ? OR actor_name ILIKE ? OR actor_email ILIKE ? OR entity_name ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (filters.Page - 1) * filters.Limit
	if err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID retrieves a single audit log by ID
func (s *AuditService) GetAuditLogByID(id uint64) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// GetRecentAuditLogs retrieves recent audit logs (last 30 days)
func (s *AuditService) GetRecentAuditLogs(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	err := s.db.Where("created_at >= ?", thirtyDaysAgo).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// GetCriticalAuditLogs retrieves critical severity audit logs
func (s *AuditService) GetCriticalAuditLogs(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := s.db.Where("severity = ?", models.SeverityCritical).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// GetFailedAuditLogs retrieves failed operation audit logs
func (s *AuditService) GetFailedAuditLogs(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := s.db.Where("status IN (?)", []string{models.StatusFailure, models.StatusError}).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// GetAuditLogsByEntity retrieves audit logs for a specific entity
func (s *AuditService) GetAuditLogsByEntity(entityType string, entityID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := s.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// GetAuditLogsByActor retrieves audit logs for a specific actor
func (s *AuditService) GetAuditLogsByActor(actorType string, actorID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := s.db.Where("actor_type = ? AND actor_id = ?", actorType, actorID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}

// Helper function to extract actor information from Gin context
func (s *AuditService) extractActorInfo(c *gin.Context, log *models.AuditLog) {
	// Try to get user information
	if userVal, exists := c.Get("user"); exists {
		if user, ok := userVal.(models.User); ok {
			log.ActorID = &user.ID
			log.ActorType = models.ActorTypeUser
			log.ActorName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
			log.ActorEmail = user.Email
		}
	}

	// Try to get contact information
	if contactVal, exists := c.Get("contact"); exists {
		if contact, ok := contactVal.(models.Contact); ok {
			log.ActorID = &contact.ID
			log.ActorType = models.ActorTypeContact
			log.ActorName = fmt.Sprintf("%s %s", contact.FirstName, contact.LastName)
			if contact.Email != nil {
				log.ActorEmail = *contact.Email
			}
		}
	}

	// If no user or contact, check for contact_id directly
	if log.ActorID == nil {
		if contactIDVal, exists := c.Get("contact_id"); exists {
			if contactID, ok := contactIDVal.(uint); ok {
				log.ActorID = &contactID
				log.ActorType = models.ActorTypeContact
				// Fetch contact details if needed
				var contact models.Contact
				if err := s.db.First(&contact, contactID).Error; err == nil {
					log.ActorName = fmt.Sprintf("%s %s", contact.FirstName, contact.LastName)
					if contact.Email != nil {
						log.ActorEmail = *contact.Email
					}
				}
			}
		}
	}

	// Get IP address
	log.ActorIPAddress = c.ClientIP()
}

// Helper function to extract request context from Gin context
func (s *AuditService) extractRequestContext(c *gin.Context, log *models.AuditLog) {
	log.HTTPMethod = c.Request.Method
	log.Endpoint = c.Request.URL.Path
	log.UserAgent = c.Request.UserAgent()

	// Get request ID if available
	if requestID, exists := c.Get("request_id"); exists {
		if reqID, ok := requestID.(string); ok {
			log.RequestID = reqID
		}
	}
}

// AuditLogFilters defines filters for querying audit logs
type AuditLogFilters struct {
	ActorID    *uint
	ActorType  string
	Action     string
	EntityType string
	EntityID   *uint
	Severity   string
	Status     string
	StartDate  time.Time
	EndDate    time.Time
	Search     string
	Page       int
	Limit      int
}
