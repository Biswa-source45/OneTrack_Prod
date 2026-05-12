package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTicketInput struct {
	ProductID     uint   `json:"product_id" binding:"required"`
	Subject       string `json:"subject" binding:"required"`
	TicketDetails string `json:"ticket_details" binding:"required"`
}

type ManagerCreateTicketInput struct {
	ContactID        uint   `json:"contact_id" binding:"required"`
	ProductID        uint   `json:"product_id" binding:"required"`
	Subject          string `json:"subject" binding:"required"`
	TicketDetails    string `json:"ticket_details" binding:"required"`
	TicketStatus     string `json:"ticket_status" binding:"required"`
	AssignedEngineer *uint  `json:"assigned_engineer"`
	Priority         string `json:"priority" binding:"required"`
	Channel          string `json:"channel" binding:"required"`
}

func CreateTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get logged-in contact_id from context (set by AuthMiddleware)
		contactIDVal, exists := c.Get("contact_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "contact_id not found in token"})
			return
		}
		contactID, ok := contactIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact_id type"})
			return
		}

		// Fetch account_id from contacts table
		var contact models.Contact
		if err := db.First(&contact, contactID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contact not found"})
			return
		}
		// Handle optional AccountID for Individual contacts
		var accountID uint
		customerCode := contact.CustomerCode // Use the contact's customer code
		if contact.AccountID != nil {
			accountID = *contact.AccountID
		} else {
			// For Individual contacts without account
			accountID = 0
		}

		// Bind input
		var input CreateTicketInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get current date
		now := time.Now()

		// Get sequence number for this account and date
		seq, err := utils.GetNextTicketSequence(db, accountID, now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate ticket sequence"})
			return
		}

		// Format ticket_id with date in DDMMYY format
		dateStr := utils.FormatDateForTicketID(now)
		ticketID := utils.FormatTicketID(customerCode, dateStr, seq)

		// Create ticket
		ticket := models.Ticket{
			TicketID:         ticketID,
			AccountID:        contact.AccountID,
			ContactID:        contactID,
			ProductID:        input.ProductID,
			Subject:          input.Subject,
			TicketDetails:    input.TicketDetails,
			TicketStatus:     "OPEN",
			Priority:         "Medium", // Default priority for customer tickets
			AssignedEngineer: nil,
		}
		if err := db.Create(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create ticket"})
			return
		}

		// Log ticket creation activity
		activityService := services.NewActivityService(db)
		activityService.LogTicketCreation(ticket.ID, &contactID, ticket.TicketID)

		// Audit log for ticket creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketCreated,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Ticket created: %s - %s", ticket.TicketID, ticket.Subject),
			nil,
			ticket,
		)

		// Send notifications for ticket creation
		notificationService := services.NewNotificationService(db)
		if err := notificationService.NotifyTicketCreated(ticket, &contactID, "contact"); err != nil {
			// Log error but don't fail the request
			// log.Printf("Failed to send ticket creation notifications: %v", err)
		}

		c.JSON(http.StatusCreated, gin.H{"ticket": ticket})
	}
}

// GetTicketsByContactHandler returns all tickets for the logged-in contact
func GetTicketsByContactHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contactIDVal, exists := c.Get("contact_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "contact_id not found in token"})
			return
		}
		contactID, ok := contactIDVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact_id type"})
			return
		}

		var tickets []models.Ticket
		if err := db.Preload("Product").Preload("Engineer").Preload("Account").Where("contact_id = ?", contactID).Order("created_at DESC").Find(&tickets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tickets"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tickets": tickets})
	}
}

// Manager: Create ticket on behalf of customer
func ManagerCreateTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// Bind input
		var input ManagerCreateTicketInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch contact and account info
		var contact models.Contact
		if err := db.Preload("Account").First(&contact, input.ContactID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contact not found"})
			return
		}

		// Handle optional AccountID for Individual contacts
		var accountID uint
		customerCode := contact.CustomerCode // Use the contact's customer code
		if contact.AccountID != nil {
			accountID = *contact.AccountID
		} else {
			// For Individual contacts without account
			accountID = 0
		}

		// Generate ticket sequence and ID
		now := time.Now()
		seq, err := utils.GetNextTicketSequence(db, accountID, now)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate ticket sequence"})
			return
		}

		// Format ticket_id with date in DDMMYY format
		dateStr := utils.FormatDateForTicketID(now)
		ticketID := utils.FormatTicketID(customerCode, dateStr, seq)

		// Create ticket
		ticket := models.Ticket{
			TicketID:         ticketID,
			AccountID:        contact.AccountID,
			ContactID:        input.ContactID,
			ProductID:        input.ProductID,
			Subject:          input.Subject,
			TicketDetails:    input.TicketDetails,
			TicketStatus:     input.TicketStatus,
			Priority:         input.Priority,
			Channel:          input.Channel,
			AssignedEngineer: input.AssignedEngineer,
		}
		if err := db.Create(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create ticket"})
			return
		}

		// Log ticket creation activity
		activityService := services.NewActivityService(db)
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok {
				activityService.LogTicketCreation(ticket.ID, &user.ID, ticket.TicketID)

				// Audit log for manager ticket creation
				auditService := services.NewAuditService(db)
				auditService.LogCRUD(
					c,
					models.AuditTicketCreated,
					models.EntityTypeTicket,
					&ticket.ID,
					ticket.TicketID,
					fmt.Sprintf("Ticket created by manager: %s - %s", ticket.TicketID, ticket.Subject),
					nil,
					ticket,
				)

				// Send notifications for manager ticket creation
				notificationService := services.NewNotificationService(db)
				if err := notificationService.NotifyTicketCreated(ticket, &user.ID, "user"); err != nil {
					// Log error but don't fail the request
					fmt.Printf("❌ Failed to send manager ticket creation notifications: %v\n", err)
				} else {
					fmt.Printf("✅ Manager ticket creation notifications sent successfully\n")
				}
			} else {
				fmt.Printf("❌ User context type assertion failed in manager create ticket\n")
			}
		} else {
			fmt.Printf("❌ No user context found in manager create ticket handler\n")
		}

		c.JSON(http.StatusCreated, gin.H{"ticket": ticket})
	}
}

// Helper to check manager role
func IsManager(c *gin.Context) bool {
	userVal, exists := c.Get("user")
	if !exists {
		return false
	}
	user, ok := userVal.(models.User)
	if !ok {
		return false
	}
	return user.RoleID == 2
}

// Manager: List all tickets
func ManagerListTicketsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var tickets []models.Ticket
		if err := db.Preload("Engineer").Preload("Product").Preload("Contact.Account").Preload("Account").Order("created_at DESC").Find(&tickets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tickets"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tickets": tickets})
	}
}

// Manager: Edit ticket
func ManagerEditTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		id := c.Param("id")

		// Bind to a map to check which fields are explicitly provided
		var rawInput map[string]interface{}
		if err := c.ShouldBindJSON(&rawInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if assigned_engineer field was explicitly provided (even if null)
		_, assignedEngineerProvided := rawInput["assigned_engineer"]

		// Parse the input properly
		var input struct {
			ContactID        *uint   `json:"contact_id"`
			TicketDetails    *string `json:"ticket_details"`
			Subject          *string `json:"subject"`
			ProductID        *uint   `json:"product_id"`
			AssignedEngineer *uint   `json:"assigned_engineer"`
			TicketStatus     *string `json:"ticket_status"`
			Priority         *string `json:"priority"`
		}

		// Convert map back to struct
		jsonBytes, _ := json.Marshal(rawInput)
		if err := json.Unmarshal(jsonBytes, &input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ticket models.Ticket
		if err := db.First(&ticket, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Get user from context for activity logging
		userVal, exists := c.Get("user")
		var userID *uint
		if exists {
			if user, ok := userVal.(models.User); ok {
				userID = &user.ID
			}
		}

		// Initialize activity service
		activityService := services.NewActivityService(db)

		// Store original values for activity logging
		originalContactID := ticket.ContactID
		originalTicketDetails := ticket.TicketDetails
		originalSubject := ticket.Subject
		originalProductID := ticket.ProductID
		originalAssignedEngineer := ticket.AssignedEngineer
		originalTicketStatus := ticket.TicketStatus
		originalPriority := ticket.Priority

		// Update fields and log activities
		if input.ContactID != nil && *input.ContactID != originalContactID {
			ticket.ContactID = *input.ContactID

			// Get contact names for activity logging
			oldVal := "Unknown Contact"
			if contactName, err := getContactName(db, originalContactID); err == nil {
				oldVal = contactName
			}
			newVal := "Unknown Contact"
			if contactName, err := getContactName(db, *input.ContactID); err == nil {
				newVal = contactName
			}
			activityService.LogFieldChange(ticket.ID, userID, models.ActivityContactChanged, "Contact", oldVal, newVal)
		}
		if input.TicketDetails != nil && *input.TicketDetails != originalTicketDetails {
			ticket.TicketDetails = *input.TicketDetails
			activityService.LogFieldChange(ticket.ID, userID, models.ActivityTicketUpdated, "Description",
				originalTicketDetails, *input.TicketDetails)
		}
		if input.Subject != nil && *input.Subject != originalSubject {
			ticket.Subject = *input.Subject
			activityService.LogFieldChange(ticket.ID, userID, models.ActivitySubjectChanged, "Subject",
				originalSubject, *input.Subject)
		}
		if input.ProductID != nil && *input.ProductID != originalProductID {
			ticket.ProductID = *input.ProductID

			// Get product names for activity logging
			oldVal := "Unknown Product"
			if productName, err := getProductName(db, originalProductID); err == nil {
				oldVal = productName
			}
			newVal := "Unknown Product"
			if productName, err := getProductName(db, *input.ProductID); err == nil {
				newVal = productName
			}
			activityService.LogFieldChange(ticket.ID, userID, models.ActivityProductChanged, "Product", oldVal, newVal)
		}
		// Handle assigned_engineer changes (including explicit null for unassignment)
		if assignedEngineerProvided {
			var originalEngineerID uint
			if originalAssignedEngineer != nil {
				originalEngineerID = *originalAssignedEngineer
			}

			var newEngineerID uint
			if input.AssignedEngineer != nil {
				newEngineerID = *input.AssignedEngineer
			}

			// Check if there's actually a change
			if (input.AssignedEngineer == nil && originalAssignedEngineer != nil) ||
				(input.AssignedEngineer != nil && originalAssignedEngineer == nil) ||
				(input.AssignedEngineer != nil && originalAssignedEngineer != nil && newEngineerID != originalEngineerID) {

				ticket.AssignedEngineer = input.AssignedEngineer

				// Get user names for activity logging
				oldVal := "Unassigned"
				if originalAssignedEngineer != nil {
					if oldUser, err := getUserName(db, *originalAssignedEngineer); err == nil {
						oldVal = oldUser
					}
				}
				newVal := "Unassigned"
				if input.AssignedEngineer != nil && *input.AssignedEngineer != 0 {
					if newUser, err := getUserName(db, *input.AssignedEngineer); err == nil {
						newVal = newUser
					}
				}
				activityService.LogFieldChange(ticket.ID, userID, models.ActivityAssigneeChanged, "Assigned Engineer", oldVal, newVal)
			}
		}
		// Note: Don't set to nil if not provided - preserve existing value
		if input.TicketStatus != nil && *input.TicketStatus != originalTicketStatus {
			ticket.TicketStatus = *input.TicketStatus
			activityService.LogFieldChange(ticket.ID, userID, models.ActivityStatusChanged, "Status",
				originalTicketStatus, *input.TicketStatus)
		}
		if input.Priority != nil && *input.Priority != originalPriority {
			ticket.Priority = *input.Priority
			activityService.LogFieldChange(ticket.ID, userID, models.ActivityPriorityChanged, "Priority",
				originalPriority, *input.Priority)
		}
		if err := db.Save(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update ticket"})
			return
		}

		// Audit log for ticket update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketUpdated,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Ticket updated: %s", ticket.TicketID),
			rawInput,
			ticket,
		)

		// Send notifications for changes made in edit
		fmt.Printf("🔥 DEBUG: ManagerEditTicketHandler - checking for notification triggers\n")
		notificationService := services.NewNotificationService(db)

		// Check for engineer assignment change
		if input.AssignedEngineer != nil && originalAssignedEngineer != input.AssignedEngineer {
			if *input.AssignedEngineer != 0 {
				// Engineer was assigned
				fmt.Printf("🔥 DEBUG: Engineer assignment detected in edit handler\n")
				if userID != nil {
					if err := notificationService.NotifyTicketAssigned(ticket, *input.AssignedEngineer, userID, "user"); err != nil {
						fmt.Printf("❌ Failed to send ticket assignment notifications from edit: %v\n", err)
					} else {
						fmt.Printf("✅ Ticket assignment notifications sent successfully from edit\n")
					}
				}
			}
		}

		// Check for status change
		if input.TicketStatus != nil && *input.TicketStatus != originalTicketStatus {
			fmt.Printf("🔥 DEBUG: Status change detected in edit handler: %s -> %s\n", originalTicketStatus, *input.TicketStatus)
			if userID != nil {
				if err := notificationService.NotifyStatusChanged(ticket, originalTicketStatus, *input.TicketStatus, userID, "user"); err != nil {
					fmt.Printf("❌ Failed to send status change notifications from edit: %v\n", err)
				} else {
					fmt.Printf("✅ Status change notifications sent successfully from edit\n")
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"ticket": ticket})
	}
}

// Manager: Delete ticket
func ManagerDeleteTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		id := c.Param("id")

		// Start transaction for data integrity
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Get ticket details first (for logging and validation)
		var ticket models.Ticket
		if err := tx.First(&ticket, id).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find ticket"})
			}
			return
		}

		fmt.Printf("🗑️ DELETING TICKET: %s (ID: %s)\n", ticket.TicketID, id)

		// Delete related records in correct order

		// 1. Delete ticket attachments (and their files)
		var attachments []models.TicketAttachment
		if err := tx.Where("ticket_id = ?", ticket.TicketID).Find(&attachments).Error; err == nil {
			for _, attachment := range attachments {
				// TODO: Delete physical files from filesystem
				fmt.Printf("📎 Would delete file: %s\n", attachment.FilePath)
			}
			// Delete attachment records
			if err := tx.Where("ticket_id = ?", ticket.TicketID).Delete(&models.TicketAttachment{}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete ticket attachments"})
				return
			}
		}

		// 2. Delete ticket comments
		if err := tx.Where("ticket_id = ?", ticket.ID).Delete(&models.TicketComment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete ticket comments"})
			return
		}

		// 3. Delete ticket calls
		if err := tx.Where("ticket_id = ?", ticket.ID).Delete(&models.TicketCall{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete ticket calls"})
			return
		}

		// 4. Delete ticket activities
		if err := tx.Where("ticket_id = ?", ticket.ID).Delete(&models.TicketActivity{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete ticket activities"})
			return
		}

		// 5. Log deletion activity before deleting the ticket
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok {
				activityService := services.NewActivityService(tx)
				description := fmt.Sprintf("Ticket %s deleted by manager", ticket.TicketID)
				if err := activityService.LogActivity(ticket.ID, &user.ID, models.ActivityTicketDeleted, description); err != nil {
					fmt.Printf("⚠️ Warning: Could not log deletion activity: %v\n", err)
					// Don't fail the deletion for logging issues
				}

				// Audit log for ticket deletion
				auditService := services.NewAuditService(db)
				auditService.LogCRUD(
					c,
					models.AuditTicketDeleted,
					models.EntityTypeTicket,
					&ticket.ID,
					ticket.TicketID,
					fmt.Sprintf("Ticket deleted: %s - %s", ticket.TicketID, ticket.Subject),
					ticket,
					nil,
				)
			}
		}

		// 6. Finally, delete the main ticket
		if err := tx.Delete(&ticket).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete ticket"})
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete ticket deletion"})
			return
		}

		fmt.Printf("✅ TICKET DELETED SUCCESSFULLY: %s\n", ticket.TicketID)
		c.JSON(http.StatusOK, gin.H{
			"message":   "Ticket deleted successfully",
			"ticket_id": ticket.TicketID,
		})
	}
}

// Manager: Change ticket status
func ManagerChangeStatusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("🔥 DEBUG: ManagerChangeStatusHandler called\n")
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		id := c.Param("id")
		var input struct {
			Status  string `json:"status" binding:"required"`
			Remarks string `json:"remarks"` // Optional remarks for status change
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ticket models.Ticket
		if err := db.First(&ticket, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		oldStatus := ticket.TicketStatus
		ticket.TicketStatus = input.Status
		if err := db.Save(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
			return
		}

		// Log status change activity with remarks
		activityService := services.NewActivityService(db)
		userVal, exists := c.Get("user")
		var userID *uint
		if exists {
			if user, ok := userVal.(models.User); ok {
				userID = &user.ID
			}
		}

		if input.Remarks != "" {
			activityService.LogStatusChangeWithRemarks(ticket.ID, userID, oldStatus, input.Status, input.Remarks)
		} else {
			activityService.LogStatusChange(ticket.ID, userID, oldStatus, input.Status)
		}

		// Audit log for status change
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketStatusChanged,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Ticket status changed: %s -> %s", oldStatus, input.Status),
			map[string]string{"status": oldStatus},
			map[string]string{"status": input.Status},
		)

		// Send notifications for status change
		notificationService := services.NewNotificationService(db)
		if exists {
			if user, ok := userVal.(models.User); ok {
				if err := notificationService.NotifyStatusChanged(ticket, oldStatus, input.Status, &user.ID, "user"); err != nil {
					// Log error but don't fail the request
					fmt.Printf("❌ Failed to send status change notifications: %v\n", err)
				} else {
					fmt.Printf("✅ Status change notifications sent successfully\n")
				}
			} else {
				fmt.Printf("❌ User context type assertion failed\n")
			}
		} else {
			fmt.Printf("❌ No user context found in status change handler\n")
		}

		c.JSON(http.StatusOK, gin.H{"ticket": ticket})
	}
}

// Manager: Assign ticket to engineer
func ManagerAssignTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("🔥 DEBUG: ManagerAssignTicketHandler called\n")
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		id := c.Param("id")
		var input struct {
			EngineerID uint `json:"engineer_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var engineer models.User
		if err := db.Where("id = ? AND role_id = ?", input.EngineerID, 3).First(&engineer).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid engineer"})
			return
		}
		var ticket models.Ticket
		if err := db.First(&ticket, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		ticket.AssignedEngineer = &input.EngineerID
		if err := db.Save(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not assign engineer"})
			return
		}

		// Audit log for ticket assignment
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketAssigned,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Ticket %s assigned to engineer: %s %s", ticket.TicketID, engineer.FirstName, engineer.LastName),
			nil,
			map[string]interface{}{"assigned_engineer_id": input.EngineerID, "engineer_name": fmt.Sprintf("%s %s", engineer.FirstName, engineer.LastName)},
		)

		// Send notifications for ticket assignment
		notificationService := services.NewNotificationService(db)
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok {
				if err := notificationService.NotifyTicketAssigned(ticket, input.EngineerID, &user.ID, "user"); err != nil {
					// Log error but don't fail the request
					fmt.Printf("❌ Failed to send ticket assignment notifications: %v\n", err)
				} else {
					fmt.Printf("✅ Ticket assignment notifications sent successfully\n")
				}
			} else {
				fmt.Printf("❌ User context type assertion failed in assignment handler\n")
			}
		} else {
			fmt.Printf("❌ No user context found in assignment handler\n")
		}

		c.JSON(http.StatusOK, gin.H{"ticket": ticket})
	}
}

// Manager: List engineers with assigned (not CLOSED) ticket count
func ManagerListEngineersWithTicketCountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var engineers []models.User
		if err := db.Where("role_id = ?", 3).Find(&engineers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch engineers"})
			return
		}
		// For each engineer, count assigned tickets (not CLOSED)
		result := make([]gin.H, 0, len(engineers))
		for _, eng := range engineers {
			var count int64
			db.Model(&models.Ticket{}).Where("assigned_engineer = ? AND ticket_status != ?", eng.ID, "CLOSED").Count(&count)
			result = append(result, gin.H{
				"id":                     eng.ID,
				"first_name":             eng.FirstName,
				"last_name":              eng.LastName,
				"phone":                  eng.Phone,
				"assigned_tickets_count": count,
			})
		}
		c.JSON(http.StatusOK, gin.H{"engineers": result})
	}
}

// Manager: List engineers (legacy, unchanged)
func ManagerListEngineersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var engineers []models.User
		if err := db.Where("role_id = ?", 3).Find(&engineers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch engineers"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"engineers": engineers})
	}
}

// Manager: List tickets assigned to engineer
func ManagerEngineerTicketsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		id := c.Param("id")
		var tickets []models.Ticket
		if err := db.Preload("Product").Preload("Contact").Preload("Account").Preload("Engineer").Where("assigned_engineer = ? AND ticket_status != ?", id, "CLOSED").Find(&tickets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tickets"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tickets": tickets})
	}
}

// Public: Get full ticket details by ID
func GetTicketDetailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var ticket models.Ticket
		// Preload related models for full details
		err := db.Preload("Product").Preload("Contact.Designation").Preload("Engineer.Role").Preload("Engineer.Designation").Preload("Account").First(&ticket, id).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ticket": ticket})
	}
}

// Helper function to get user name by ID for activity logging
func getUserName(db *gorm.DB, userID uint) (string, error) {
	var user models.User
	err := db.First(&user, userID).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName), nil
}

// Helper function to get contact name by ID for activity logging
func getContactName(db *gorm.DB, contactID uint) (string, error) {
	var contact models.Contact
	err := db.First(&contact, contactID).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", contact.FirstName, contact.LastName), nil
}

// Helper function to get product name by ID for activity logging
func getProductName(db *gorm.DB, productID uint) (string, error) {
	var product models.MasterProduct
	err := db.First(&product, productID).Error
	if err != nil {
		return "", err
	}
	return product.ProductName, nil
}
