package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// WebSocketHub interface for broadcasting notifications
type WebSocketHub interface {
	BroadcastToUser(userID uint, userType, messageType string, data interface{})
}

// GlobalWebSocketHub is the global WebSocket hub instance
// Set by main.go during initialization
var GlobalWebSocketHub WebSocketHub

// NotificationService handles notification creation and management
type NotificationService struct {
	db  *gorm.DB
	hub WebSocketHub // WebSocket hub for real-time notifications
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db:  db,
		hub: GlobalWebSocketHub, // Use global hub if available
	}
}

// SetWebSocketHub sets the WebSocket hub for real-time broadcasting
func (ns *NotificationService) SetWebSocketHub(hub WebSocketHub) {
	ns.hub = hub
}

// NotificationData holds data for creating notifications
type NotificationData struct {
	RecipientID      uint
	RecipientType    string // "user" or "contact"
	NotificationType string
	RelatedID        *uint
	RelatedType      string
	RelatedSubID     *uint
	ActorID          *uint
	ActorType        string
	Variables        map[string]string      // For template replacement
	Priority         string                 // Optional, will use template default if empty
	Metadata         map[string]interface{} // Additional data
}

// CreateNotification creates a new notification using templates
func (ns *NotificationService) CreateNotification(data NotificationData) error {
	// Get template for this notification type
	var template models.NotificationTemplate
	err := ns.db.Where("notification_type = ? AND is_active = ?", data.NotificationType, true).
		First(&template).Error
	if err != nil {
		return fmt.Errorf("notification template not found for type %s: %w", data.NotificationType, err)
	}

	// Replace variables in title and message templates
	title := ns.replaceVariables(template.TitleTemplate, data.Variables)
	message := ns.replaceVariables(template.MessageTemplate, data.Variables)

	// Use provided priority or template default
	priority := data.Priority
	if priority == "" {
		priority = template.DefaultPriority
	}

	// Convert metadata to JSON string
	metadataJSON := "{}"
	if data.Metadata != nil {
		if metadataBytes, err := json.Marshal(data.Metadata); err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	// Create notification
	notification := models.Notification{
		RecipientID:      data.RecipientID,
		RecipientType:    data.RecipientType,
		Title:            title,
		Message:          message,
		NotificationType: data.NotificationType,
		RelatedID:        data.RelatedID,
		RelatedType:      data.RelatedType,
		RelatedSubID:     data.RelatedSubID,
		ActorID:          data.ActorID,
		ActorType:        data.ActorType,
		Priority:         priority,
		Category:         template.Category,
		Metadata:         metadataJSON,
	}

	if err := ns.db.Create(&notification).Error; err != nil {
		return err
	}

	// Broadcast notification via WebSocket if hub is available
	if ns.hub != nil {
		ns.broadcastNotification(&notification)
	}

	return nil
}

// replaceVariables replaces {variable_name} placeholders in templates
func (ns *NotificationService) replaceVariables(template string, variables map[string]string) string {
	result := template
	for key, value := range variables {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// GetUserNotifications retrieves notifications for a user with pagination
func (ns *NotificationService) GetUserNotifications(userID uint, userType string, limit, offset int, filters NotificationFilters) ([]models.Notification, error) {
	var notifications []models.Notification

	query := ns.db.Where("recipient_id = ? AND recipient_type = ?", userID, userType)

	// Apply filters
	if filters.Category != "" && filters.Category != "all" {
		query = query.Where("category = ?", filters.Category)
	}
	if filters.Priority != "" && filters.Priority != "all" {
		query = query.Where("priority = ?", filters.Priority)
	}
	if filters.IsRead != "" && filters.IsRead != "all" {
		if filters.IsRead == "read" {
			query = query.Where("is_read = ?", true)
		} else if filters.IsRead == "unread" {
			query = query.Where("is_read = ?", false)
		}
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, err
}

// NotificationFilters holds filtering options
type NotificationFilters struct {
	Category string // "all", "ticket", "task", "communication", "system"
	Priority string // "all", "low", "normal", "high", "urgent"
	IsRead   string // "all", "read", "unread"
}

// GetUnreadCount returns count of unread notifications for a user
func (ns *NotificationService) GetUnreadCount(userID uint, userType string) (int64, error) {
	var count int64
	err := ns.db.Model(&models.Notification{}).
		Where("recipient_id = ? AND recipient_type = ? AND is_read = ?", userID, userType, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(notificationID, userID uint, userType string) error {
	now := time.Now()
	err := ns.db.Model(&models.Notification{}).
		Where("id = ? AND recipient_id = ? AND recipient_type = ?", notificationID, userID, userType).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error

	if err != nil {
		return err
	}

	// Broadcast read status update via WebSocket
	if ns.hub != nil {
		ns.hub.BroadcastToUser(userID, userType, "notification.read", map[string]interface{}{
			"notification_id": notificationID,
		})

		// Update unread count
		ns.broadcastUnreadCount(userID, userType)
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (ns *NotificationService) MarkAllAsRead(userID uint, userType string) error {
	now := time.Now()
	err := ns.db.Model(&models.Notification{}).
		Where("recipient_id = ? AND recipient_type = ? AND is_read = ?", userID, userType, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error

	if err != nil {
		return err
	}

	// Broadcast mark all read via WebSocket
	if ns.hub != nil {
		ns.hub.BroadcastToUser(userID, userType, "notification.all_read", map[string]interface{}{
			"message": "All notifications marked as read",
		})

		// Update unread count to 0
		ns.hub.BroadcastToUser(userID, userType, "count.update", map[string]interface{}{
			"unread_count": 0,
		})
	}

	return nil
}

// DeleteNotification deletes a notification (soft delete)
func (ns *NotificationService) DeleteNotification(notificationID, userID uint, userType string) error {
	// Check if notification was unread before deletion
	var notification models.Notification
	err := ns.db.Where("id = ? AND recipient_id = ? AND recipient_type = ?", notificationID, userID, userType).
		First(&notification).Error
	if err != nil {
		return err
	}

	wasUnread := !notification.IsRead

	// Delete notification
	err = ns.db.Where("id = ? AND recipient_id = ? AND recipient_type = ?", notificationID, userID, userType).
		Delete(&models.Notification{}).Error

	if err != nil {
		return err
	}

	// Broadcast deletion via WebSocket
	if ns.hub != nil {
		ns.hub.BroadcastToUser(userID, userType, "notification.deleted", map[string]interface{}{
			"notification_id": notificationID,
		})

		// Update unread count if deleted notification was unread
		if wasUnread {
			ns.broadcastUnreadCount(userID, userType)
		}
	}

	return nil
}

// Helper methods for creating specific notification types

// NotifyTicketCreated notifies about ticket creation
func (ns *NotificationService) NotifyTicketCreated(ticket models.Ticket, actorID *uint, actorType string) error {
	variables := map[string]string{
		"ticket_id": ticket.TicketID,
	}

	metadata := map[string]interface{}{
		"ticket_id":      ticket.ID,
		"ticket_subject": ticket.Subject,
	}

	// Notify customer (confirmation)
	customerData := NotificationData{
		RecipientID:      ticket.ContactID,
		RecipientType:    "contact",
		NotificationType: models.NotificationTicketCreatedConfirmation,
		RelatedID:        &ticket.ID,
		RelatedType:      "ticket",
		ActorID:          actorID,
		ActorType:        actorType,
		Variables:        variables,
		Metadata:         metadata,
	}

	if err := ns.CreateNotification(customerData); err != nil {
		return err
	}

	// Notify managers (new ticket alert) - get all managers
	var managers []models.User
	if err := ns.db.Where("role_id = ?", 2).Find(&managers).Error; err != nil {
		return err
	}

	// Get customer name for manager notification
	var contact models.Contact
	if err := ns.db.First(&contact, ticket.ContactID).Error; err == nil {
		variables["customer_name"] = contact.FirstName + " " + contact.LastName
	}

	// Get product name
	var product models.MasterProduct
	if err := ns.db.First(&product, ticket.ProductID).Error; err == nil {
		variables["product_name"] = product.ProductName
	}

	for _, manager := range managers {
		managerData := NotificationData{
			RecipientID:      manager.ID,
			RecipientType:    "user",
			NotificationType: models.NotificationTicketCreatedByCustomer,
			RelatedID:        &ticket.ID,
			RelatedType:      "ticket",
			ActorID:          &ticket.ContactID,
			ActorType:        "contact",
			Variables:        variables,
			Metadata:         metadata,
		}

		if err := ns.CreateNotification(managerData); err != nil {
			return err
		}
	}

	return nil
}

// NotifyTicketAssigned notifies about ticket assignment
func (ns *NotificationService) NotifyTicketAssigned(ticket models.Ticket, engineerID uint, actorID *uint, actorType string) error {
	variables := map[string]string{
		"ticket_id": ticket.TicketID,
	}

	metadata := map[string]interface{}{
		"ticket_id":   ticket.ID,
		"engineer_id": engineerID,
	}

	// Get engineer name
	var engineer models.User
	if err := ns.db.First(&engineer, engineerID).Error; err == nil {
		variables["engineer_name"] = engineer.FirstName + " " + engineer.LastName
	}

	// Get manager name for engineer notification
	if actorID != nil && actorType == "user" {
		var manager models.User
		if err := ns.db.First(&manager, *actorID).Error; err == nil {
			variables["manager_name"] = manager.FirstName + " " + manager.LastName
		}
	}

	// Notify customer
	customerData := NotificationData{
		RecipientID:      ticket.ContactID,
		RecipientType:    "contact",
		NotificationType: models.NotificationTicketEngineerAssigned,
		RelatedID:        &ticket.ID,
		RelatedType:      "ticket",
		ActorID:          actorID,
		ActorType:        actorType,
		Variables:        variables,
		Metadata:         metadata,
	}

	if err := ns.CreateNotification(customerData); err != nil {
		return err
	}

	// Notify engineer
	engineerData := NotificationData{
		RecipientID:      engineerID,
		RecipientType:    "user",
		NotificationType: models.NotificationTicketAssignedToYou,
		RelatedID:        &ticket.ID,
		RelatedType:      "ticket",
		ActorID:          actorID,
		ActorType:        actorType,
		Variables:        variables,
		Metadata:         metadata,
		Priority:         "high",
	}

	return ns.CreateNotification(engineerData)
}

// NotifyStatusChanged notifies about status changes
func (ns *NotificationService) NotifyStatusChanged(ticket models.Ticket, oldStatus, newStatus string, actorID *uint, actorType string) error {
	variables := map[string]string{
		"ticket_id":  ticket.TicketID,
		"old_status": oldStatus,
		"new_status": newStatus,
	}

	metadata := map[string]interface{}{
		"ticket_id":  ticket.ID,
		"old_status": oldStatus,
		"new_status": newStatus,
	}

	// Get actor name
	var actorName string
	if actorID != nil {
		if actorType == "user" {
			var user models.User
			if err := ns.db.First(&user, *actorID).Error; err == nil {
				actorName = user.FirstName + " " + user.LastName
				variables["engineer_name"] = user.FirstName + " " + user.LastName
			}
		} else if actorType == "contact" {
			var contact models.Contact
			if err := ns.db.First(&contact, *actorID).Error; err == nil {
				actorName = contact.FirstName + " " + contact.LastName
				variables["customer_name"] = contact.FirstName + " " + contact.LastName
			}
		}
	}

	// Notify customer
	customerData := NotificationData{
		RecipientID:      ticket.ContactID,
		RecipientType:    "contact",
		NotificationType: models.NotificationTicketStatusChanged,
		RelatedID:        &ticket.ID,
		RelatedType:      "ticket",
		ActorID:          actorID,
		ActorType:        actorType,
		Variables:        variables,
		Metadata:         metadata,
	}

	if err := ns.CreateNotification(customerData); err != nil {
		return err
	}

	// If engineer updated status, notify managers
	if actorType == "user" && actorName != "" {
		var managers []models.User
		if err := ns.db.Where("role_id = ?", 2).Find(&managers).Error; err == nil {
			for _, manager := range managers {
				managerData := NotificationData{
					RecipientID:      manager.ID,
					RecipientType:    "user",
					NotificationType: models.NotificationTicketEngineerStatusUpdate,
					RelatedID:        &ticket.ID,
					RelatedType:      "ticket",
					ActorID:          actorID,
					ActorType:        actorType,
					Variables:        variables,
					Metadata:         metadata,
				}

				if err := ns.CreateNotification(managerData); err != nil {
					return err
				}
			}
		}
	}

	// If manager updated status and ticket is assigned, notify engineer
	if actorType == "user" && ticket.AssignedEngineer != nil {
		engineerData := NotificationData{
			RecipientID:      *ticket.AssignedEngineer,
			RecipientType:    "user",
			NotificationType: models.NotificationTicketStatusChanged,
			RelatedID:        &ticket.ID,
			RelatedType:      "ticket",
			ActorID:          actorID,
			ActorType:        actorType,
			Variables:        variables,
			Metadata:         metadata,
		}

		if err := ns.CreateNotification(engineerData); err != nil {
			return err
		}
	}

	return nil
}

// NotifyCommentAdded notifies about new comments
func (ns *NotificationService) NotifyCommentAdded(ticket models.Ticket, comment models.TicketComment, actorID *uint, actorType string) error {
	// Skip internal comments for customer notifications
	if comment.IsInternal && actorType == "user" {
		return nil
	}

	variables := map[string]string{
		"ticket_id": ticket.TicketID,
	}

	metadata := map[string]interface{}{
		"ticket_id":    ticket.ID,
		"comment_id":   comment.ID,
		"comment_type": comment.Type,
	}

	notificationType := models.NotificationTicketCommentAdded
	if comment.Type == "resolution" {
		notificationType = models.NotificationTicketResolutionAdded
	}

	// Notify customer (if not internal comment)
	if !comment.IsInternal {
		customerData := NotificationData{
			RecipientID:      ticket.ContactID,
			RecipientType:    "contact",
			NotificationType: notificationType,
			RelatedID:        &ticket.ID,
			RelatedType:      "comment",
			RelatedSubID:     &comment.ID,
			ActorID:          actorID,
			ActorType:        actorType,
			Variables:        variables,
			Metadata:         metadata,
		}

		if err := ns.CreateNotification(customerData); err != nil {
			return err
		}
	}

	// Notify assigned engineer (if comment not from engineer)
	if ticket.AssignedEngineer != nil && (actorID == nil || *actorID != *ticket.AssignedEngineer) {
		engineerNotificationType := models.NotificationTicketCommentAdded
		if actorType == "contact" {
			engineerNotificationType = models.NotificationTicketCustomerCommentAdded
		} else if actorType == "user" {
			engineerNotificationType = models.NotificationTicketManagerCommentAdded
		}

		engineerData := NotificationData{
			RecipientID:      *ticket.AssignedEngineer,
			RecipientType:    "user",
			NotificationType: engineerNotificationType,
			RelatedID:        &ticket.ID,
			RelatedType:      "comment",
			RelatedSubID:     &comment.ID,
			ActorID:          actorID,
			ActorType:        actorType,
			Variables:        variables,
			Metadata:         metadata,
		}

		if err := ns.CreateNotification(engineerData); err != nil {
			return err
		}
	}

	return nil
}



// broadcastNotification sends a new notification via WebSocket
func (ns *NotificationService) broadcastNotification(notification *models.Notification) {
	if ns.hub == nil {
		return
	}

	// Send new notification
	ns.hub.BroadcastToUser(
		notification.RecipientID,
		notification.RecipientType,
		"notification.new",
		map[string]interface{}{
			"id":                notification.ID,
			"title":             notification.Title,
			"message":           notification.Message,
			"notification_type": notification.NotificationType,
			"priority":          notification.Priority,
			"category":          notification.Category,
			"related_id":        notification.RelatedID,
			"related_type":      notification.RelatedType,
			"is_read":           notification.IsRead,
			"created_at":        notification.CreatedAt,
		},
	)

	// Update unread count
	ns.broadcastUnreadCount(notification.RecipientID, notification.RecipientType)
}

// broadcastUnreadCount sends updated unread count via WebSocket
func (ns *NotificationService) broadcastUnreadCount(userID uint, userType string) {
	if ns.hub == nil {
		return
	}

	var count int64
	ns.db.Model(&models.Notification{}).
		Where("recipient_id = ? AND recipient_type = ? AND is_read = ?", userID, userType, false).
		Count(&count)

	ns.hub.BroadcastToUser(userID, userType, "count.update", map[string]interface{}{
		"unread_count": count,
	})
}
