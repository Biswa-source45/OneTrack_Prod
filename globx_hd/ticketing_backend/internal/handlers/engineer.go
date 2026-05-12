package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Engineer role check
func IsEngineer(c *gin.Context) bool {
	userVal, exists := c.Get("user")
	if !exists {
		return false
	}
	user, ok := userVal.(models.User)
	if !ok {
		return false
	}
	return user.RoleID == 3
}

// Engineer: List tickets assigned to the logged-in engineer
func EngineerListTicketsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		userVal, _ := c.Get("user")
		user := userVal.(models.User)
		var tickets []models.Ticket
		if err := db.Where("assigned_engineer = ?", user.ID).Preload("Engineer").Preload("Product").Preload("Contact.Account").Preload("Account").Order("created_at DESC").Find(&tickets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tickets"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tickets": tickets})
	}
}

// Engineer: Change status of assigned ticket from OPEN to RESOLVED only
func EngineerChangeStatusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		userVal, _ := c.Get("user")
		user := userVal.(models.User)
		id := c.Param("id")
		var input struct {
			Status  string `json:"status" binding:"required"`
			Remarks string `json:"remarks"` // Optional remarks for status change
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Status == "CLOSED" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "engineer cannot set status to CLOSED"})
			return
		}
		var ticket models.Ticket
		if err := db.First(&ticket, id).Error; err != nil {
		}
		if ticket.AssignedEngineer == nil || *ticket.AssignedEngineer != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not assigned to this ticket"})
			return
		}

		oldStatus := ticket.TicketStatus
		// log.Printf("[EngineerChangeStatusHandler] Changing ticket %d status from %s to %s", ticket.ID, ticket.TicketStatus, input.Status)
		if err := db.Model(&ticket).Update("ticket_status", input.Status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
			return
		}
		// log.Printf("[EngineerChangeStatusHandler] Changed ticket %d status, now: %s", ticket.ID, input.Status)

		// Log status change activity with remarks
		activityService := services.NewActivityService(db)
		if input.Remarks != "" {
			activityService.LogStatusChangeWithRemarks(ticket.ID, &user.ID, oldStatus, input.Status, input.Remarks)
		} else {
			activityService.LogStatusChange(ticket.ID, &user.ID, oldStatus, input.Status)
		}

		// Audit log for engineer status change
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketStatusChanged,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Engineer changed ticket status: %s -> %s", oldStatus, input.Status),
			map[string]string{"status": oldStatus},
			map[string]string{"status": input.Status},
		)

		// Send notifications for engineer status change
		notificationService := services.NewNotificationService(db)
		if err := notificationService.NotifyStatusChanged(ticket, oldStatus, input.Status, &user.ID, "user"); err != nil {
			// Log error but don't fail the request
			fmt.Printf("❌ Failed to send engineer status change notifications: %v\n", err)
		} else {
			fmt.Printf("✅ Engineer status change notifications sent successfully\n")
		}

		c.JSON(http.StatusOK, gin.H{"ticket": ticket})
	}
}

// Engineer: List tasks assigned to the logged-in engineer
func EngineerListTasksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		userVal, _ := c.Get("user")
		user := userVal.(models.User)
		var tasks []models.Task
		if err := db.Where("assigned_to = ?", user.ID).Preload("AssignedUser").Preload("Creator").Order("created_at DESC").Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tasks"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

// Engineer: Get task by ID (only if assigned to engineer)
func EngineerGetTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		userVal, _ := c.Get("user")
		user := userVal.(models.User)

		id := c.Param("id")
		var task models.Task
		if err := db.Preload("AssignedUser").Preload("Creator").First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		// Check if task is assigned to this engineer
		if task.AssignedTo == nil || *task.AssignedTo != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not assigned to this task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

// Engineer: Change status of assigned task
func EngineerChangeTaskStatusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		userVal, _ := c.Get("user")
		user := userVal.(models.User)
		id := c.Param("id")
		var input struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		// Check if task is assigned to this engineer
		if task.AssignedTo == nil || *task.AssignedTo != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not assigned to this task"})
			return
		}

		oldStatus := task.TaskStatus

		// Update task status
		if err := db.Model(&task).Update("task_status", input.Status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
			return
		}

		// Audit log for engineer task status change
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTaskStatusChanged,
			models.EntityTypeTask,
			&task.ID,
			task.Subject,
			fmt.Sprintf("Engineer changed task status: %s -> %s", oldStatus, input.Status),
			map[string]string{"status": oldStatus},
			map[string]string{"status": input.Status},
		)

		// Send notifications for task status change
		fmt.Printf("🔥 DEBUG: EngineerChangeTaskStatusHandler - sending task status change notifications\n")
		fmt.Printf("🔥 DEBUG: Status change from '%s' to '%s' for task ID %s\n", oldStatus, input.Status, id)
		notificationService := services.NewNotificationService(db)

		// Get the task creator (manager) to notify them
		var taskCreator models.User
		if err := db.First(&taskCreator, task.CreatedBy).Error; err == nil {
			fmt.Printf("🔥 DEBUG: Found task creator (manager ID: %d) for task %d\n", taskCreator.ID, task.ID)

			// Only notify if the status change is not done by the creator themselves
			if user.ID != taskCreator.ID {
				taskData := map[string]string{
					"task_subject": task.Subject,
					"task_id":      fmt.Sprintf("%d", task.ID),
					"old_status":   oldStatus,
					"new_status":   input.Status,
				}

				fmt.Printf("🔥 DEBUG: Sending status change notification to manager (ID: %d)\n", taskCreator.ID)

				if err := notificationService.CreateNotification(services.NotificationData{
					RecipientID:      taskCreator.ID,
					RecipientType:    "user",
					NotificationType: models.NotificationTaskStatusUpdatedByEngineer,
					Variables:        taskData,
					RelatedID:        &task.ID,
					RelatedType:      "task",
					ActorID:          &user.ID,
					ActorType:        "user",
				}); err != nil {
					fmt.Printf("❌ Failed to send task status change notifications: %v\n", err)
				} else {
					fmt.Printf("✅ Task status change notifications sent successfully\n")
				}
			} else {
				fmt.Printf("🔥 DEBUG: Skipping notification - manager changed their own task status\n")
			}
		} else {
			fmt.Printf("❌ Failed to load task creator for notifications: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

// EngineerCreateTicketInput defines the input structure for engineer ticket creation
type EngineerCreateTicketInput struct {
	ContactID     uint   `json:"contact_id" binding:"required"`
	ProductID     uint   `json:"product_id" binding:"required"`
	Subject       string `json:"subject" binding:"required"`
	TicketDetails string `json:"ticket_details" binding:"required"`
	TicketStatus  string `json:"ticket_status" binding:"required"`
	Priority      string `json:"priority" binding:"required"`
	Channel       string `json:"channel" binding:"required"`
}

// Engineer: Create ticket on behalf of customer (unassigned - only managers can assign)
func EngineerCreateTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// Get logged-in engineer
		userVal, _ := c.Get("user")
		user := userVal.(models.User)

		// Bind input
		var input EngineerCreateTicketInput
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

		// Create ticket - leave unassigned (only managers can assign engineers)
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
			AssignedEngineer: nil, // Leave unassigned - only managers can assign
		}
		if err := db.Create(&ticket).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create ticket"})
			return
		}

		// Log ticket creation activity
		activityService := services.NewActivityService(db)
		activityService.LogTicketCreation(ticket.ID, &user.ID, ticket.TicketID)

		// Audit log for engineer ticket creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTicketCreated,
			models.EntityTypeTicket,
			&ticket.ID,
			ticket.TicketID,
			fmt.Sprintf("Ticket created by engineer: %s - %s", ticket.TicketID, ticket.Subject),
			nil,
			ticket,
		)

		// Send notifications for ticket creation
		notificationService := services.NewNotificationService(db)
		if err := notificationService.NotifyTicketCreated(ticket, &user.ID, "user"); err != nil {
			// Log error but don't fail the request
			fmt.Printf("❌ Failed to send engineer ticket creation notifications: %v\n", err)
		} else {
			fmt.Printf("✅ Engineer ticket creation notifications sent successfully\n")
		}

		c.JSON(http.StatusCreated, gin.H{"ticket": ticket})
	}
}
