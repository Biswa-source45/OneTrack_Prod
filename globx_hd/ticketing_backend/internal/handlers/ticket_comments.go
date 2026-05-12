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

// Input structures for ticket comments
type CreateCommentInput struct {
	Type       string `json:"type" binding:"required,oneof=comment resolution"`
	Content    string `json:"content" binding:"required"`
	IsInternal bool   `json:"is_internal"`
}

type UpdateCommentInput struct {
	Content    string `json:"content" binding:"required"`
	IsInternal *bool  `json:"is_internal"`
}

// CreateTicketComment creates a new comment or resolution for a ticket
func CreateTicketComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID from URL
		ticketIDStr := c.Param("id")
		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		user, ok := userVal.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
			return
		}
		userID := user.ID

		// Bind input
		var input CreateCommentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify ticket exists
		var ticket models.Ticket
		if err := db.First(&ticket, uint(ticketID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Create comment
		comment := models.TicketComment{
			TicketID:   uint(ticketID),
			UserID:     userID,
			Type:       input.Type,
			Content:    input.Content,
			IsInternal: input.IsInternal,
		}

		if err := db.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create comment"})
			return
		}

		// Log activity
		activityService := services.NewActivityService(db)
		activityService.LogComment(uint(ticketID), &userID, input.Type)

		// Send notifications for comment
		notificationService := services.NewNotificationService(db)
		if err := notificationService.NotifyCommentAdded(ticket, comment, &userID, "user"); err != nil {
			// Log error but don't fail the request
			// log.Printf("Failed to send comment notifications: %v", err)
		}

		// Load user information for response
		db.Preload("User").First(&comment, comment.ID)

		c.JSON(http.StatusCreated, gin.H{"comment": comment})
	}
}

// GetTicketComments retrieves all comments for a ticket
func GetTicketComments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID from URL
		ticketIDStr := c.Param("id")
		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		// Get pagination parameters
		limitStr := c.DefaultQuery("limit", "50")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 50
		}
		if limit > 100 {
			limit = 100 // Max limit
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		// Verify ticket exists
		var ticket models.Ticket
		if err := db.First(&ticket, uint(ticketID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Get comments
		var comments []models.TicketComment
		query := db.Where("ticket_id = ?", ticketID).
			Preload("User").
			Order("created_at ASC").
			Limit(limit).
			Offset(offset)

		// Filter internal comments for non-managers
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok && user.RoleID != 2 { // Not a manager
				query = query.Where("is_internal = ?", false)
			}
		}

		if err := query.Find(&comments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch comments"})
			return
		}

		// Get total count
		var total int64
		countQuery := db.Model(&models.TicketComment{}).Where("ticket_id = ?", ticketID)
		if exists {
			if user, ok := userVal.(models.User); ok && user.RoleID != 2 {
				countQuery = countQuery.Where("is_internal = ?", false)
			}
		}
		countQuery.Count(&total)

		c.JSON(http.StatusOK, gin.H{
			"comments": comments,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		})
	}
}

// UpdateTicketComment updates an existing comment
func UpdateTicketComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID and comment ID from URL
		ticketIDStr := c.Param("id")
		commentIDStr := c.Param("comment_id")

		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		user, ok := userVal.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
			return
		}
		userID := user.ID

		// Bind input
		var input UpdateCommentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find comment
		var comment models.TicketComment
		if err := db.Where("id = ? AND ticket_id = ?", commentID, ticketID).First(&comment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		// Check if user owns the comment or is a manager
		if exists {
			if user, ok := userVal.(models.User); ok {
				if comment.UserID != userID && user.RoleID != 2 { // Not owner and not manager
					c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to edit this comment"})
					return
				}
			}
		}

		// Update comment
		comment.Content = input.Content
		if input.IsInternal != nil {
			comment.IsInternal = *input.IsInternal
		}

		oldComment := comment

		if err := db.Save(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update comment"})
			return
		}

		// Log comment update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditCommentUpdated, models.EntityTypeComment, &comment.ID, fmt.Sprintf("Comment #%d", comment.ID), fmt.Sprintf("Comment updated on ticket %d", ticketID), oldComment, comment)

		// Load user information for response
		db.Preload("User").First(&comment, comment.ID)

		c.JSON(http.StatusOK, gin.H{"comment": comment})
	}
}

// DeleteTicketComment deletes a comment
func DeleteTicketComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID and comment ID from URL
		ticketIDStr := c.Param("id")
		commentIDStr := c.Param("comment_id")

		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		user, ok := userVal.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
			return
		}
		userID := user.ID

		// Find comment
		var comment models.TicketComment
		if err := db.Where("id = ? AND ticket_id = ?", commentID, ticketID).First(&comment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		// Check if user owns the comment or is a manager
		if exists {
			if user, ok := userVal.(models.User); ok {
				if comment.UserID != userID && user.RoleID != 2 { // Not owner and not manager
					c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to delete this comment"})
					return
				}
			}
		}

		// Delete comment
		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete comment"})
			return
		}

		// Log comment deletion
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditCommentDeleted, models.EntityTypeComment, &comment.ID, fmt.Sprintf("Comment #%d", comment.ID), fmt.Sprintf("Comment deleted from ticket %d", ticketID), comment, nil)

		c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
	}
}
