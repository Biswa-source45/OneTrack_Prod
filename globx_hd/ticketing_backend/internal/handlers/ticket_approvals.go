package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateApprovalRequest creates a new approval request
func CreateApprovalRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is manager or engineer
		if !IsManager(c) && !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		ticketID := c.Param("id")
		var input struct {
			ApproverID uint   `json:"approver_id" binding:"required"`
			Subject    string `json:"subject" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get current user
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		user := userVal.(models.User)

		// Verify ticket exists
		var ticket models.Ticket
		if err := db.First(&ticket, ticketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Verify approver exists and is a manager
		var approver models.User
		if err := db.Where("id = ? AND role_id = ?", input.ApproverID, 2).First(&approver).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "approver must be a manager"})
			return
		}

		// Create approval request
		approval := models.TicketApproval{
			TicketID:    ticket.ID,
			RequesterID: user.ID,
			ApproverID:  input.ApproverID,
			Subject:     input.Subject,
			Status:      "PENDING",
		}

		if err := db.Create(&approval).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create approval request"})
			return
		}

		// Load relationships for response
		db.Preload("Requester").Preload("Approver").First(&approval, approval.ID)

		// Log activity
		activityService := services.NewActivityService(db)
		approverName := fmt.Sprintf("%s %s", approver.FirstName, approver.LastName)
		activityService.LogApprovalRequested(ticket.ID, &user.ID, approverName, input.Subject)

		// Send notification to approver
		notificationService := services.NewNotificationService(db)
		variables := map[string]string{
			"ticket_id":      ticket.TicketID,
			"requester_name": fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			"subject":        input.Subject,
		}

		if err := notificationService.CreateNotification(services.NotificationData{
			RecipientID:      input.ApproverID,
			RecipientType:    "user",
			NotificationType: models.NotificationTicketApprovalRequested,
			Variables:        variables,
			RelatedID:        &ticket.ID,
			RelatedType:      "ticket",
			RelatedSubID:     &approval.ID,
			ActorID:          &user.ID,
			ActorType:        "user",
		}); err != nil {
			fmt.Printf("❌ Failed to send approval request notification: %v\n", err)
		}

		c.JSON(http.StatusCreated, gin.H{"approval": approval})
	}
}

// ListApprovalRequests lists all approval requests for a ticket
func ListApprovalRequestsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is manager or engineer
		if !IsManager(c) && !IsEngineer(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		ticketID := c.Param("id")

		// Verify ticket exists
		var ticket models.Ticket
		if err := db.First(&ticket, ticketID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Get pagination parameters
		limit := 50
		offset := 0
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}
		if offsetStr := c.Query("offset"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = o
			}
		}

		// Fetch approvals
		var approvals []models.TicketApproval
		query := db.Where("ticket_id = ?", ticket.ID).
			Preload("Requester").
			Preload("Approver").
			Order("created_at DESC").
			Limit(limit).
			Offset(offset)

		if err := query.Find(&approvals).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch approvals"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"approvals": approvals})
	}
}

// ApproveRequest approves an approval request
func ApproveRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is manager
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only managers can approve requests"})
			return
		}

		ticketID := c.Param("id")
		approvalID := c.Param("approvalId")

		var input struct {
			Remarks string `json:"remarks" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get current user
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		user := userVal.(models.User)

		// Find approval request
		var approval models.TicketApproval
		if err := db.Where("id = ? AND ticket_id = ? AND approver_id = ?", approvalID, ticketID, user.ID).
			Preload("Requester").Preload("Approver").First(&approval).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "approval request not found or not authorized"})
			return
		}

		// Check if already processed
		if approval.Status != "PENDING" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "approval request already processed"})
			return
		}

		// Update approval
		approval.Status = "APPROVED"
		approval.Remarks = input.Remarks

		if err := db.Save(&approval).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not approve request"})
			return
		}

		// Log activity
		activityService := services.NewActivityService(db)
		ticketIDUint, _ := strconv.ParseUint(ticketID, 10, 32)
		activityService.LogApprovalApproved(uint(ticketIDUint), &user.ID, approval.Subject, input.Remarks)

		// Send notification to requester
		notificationService := services.NewNotificationService(db)
		variables := map[string]string{
			"approver_name": fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			"subject":       approval.Subject,
			"remarks":       input.Remarks,
		}

		if err := notificationService.CreateNotification(services.NotificationData{
			RecipientID:      approval.RequesterID,
			RecipientType:    "user",
			NotificationType: models.NotificationTicketApprovalApproved,
			Variables:        variables,
			RelatedID:        &approval.TicketID,
			RelatedType:      "ticket",
			RelatedSubID:     &approval.ID,
			ActorID:          &user.ID,
			ActorType:        "user",
		}); err != nil {
			fmt.Printf("❌ Failed to send approval approved notification: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"approval": approval})
	}
}

// RejectRequest rejects an approval request
func RejectRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is manager
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only managers can reject requests"})
			return
		}

		ticketID := c.Param("id")
		approvalID := c.Param("approvalId")

		var input struct {
			Remarks string `json:"remarks" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get current user
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		user := userVal.(models.User)

		// Find approval request
		var approval models.TicketApproval
		if err := db.Where("id = ? AND ticket_id = ? AND approver_id = ?", approvalID, ticketID, user.ID).
			Preload("Requester").Preload("Approver").First(&approval).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "approval request not found or not authorized"})
			return
		}

		// Check if already processed
		if approval.Status != "PENDING" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "approval request already processed"})
			return
		}

		// Update approval
		approval.Status = "REJECTED"
		approval.Remarks = input.Remarks

		if err := db.Save(&approval).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reject request"})
			return
		}

		// Log activity
		activityService := services.NewActivityService(db)
		ticketIDUint, _ := strconv.ParseUint(ticketID, 10, 32)
		activityService.LogApprovalRejected(uint(ticketIDUint), &user.ID, approval.Subject, input.Remarks)

		// Send notification to requester
		notificationService := services.NewNotificationService(db)
		variables := map[string]string{
			"approver_name": fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			"subject":       approval.Subject,
			"remarks":       input.Remarks,
		}

		if err := notificationService.CreateNotification(services.NotificationData{
			RecipientID:      approval.RequesterID,
			RecipientType:    "user",
			NotificationType: models.NotificationTicketApprovalRejected,
			Variables:        variables,
			RelatedID:        &approval.TicketID,
			RelatedType:      "ticket",
			RelatedSubID:     &approval.ID,
			ActorID:          &user.ID,
			ActorType:        "user",
		}); err != nil {
			fmt.Printf("❌ Failed to send approval rejected notification: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"approval": approval})
	}
}
