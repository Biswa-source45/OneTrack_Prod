package handlers

import (
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTicketActivities retrieves activity history for a ticket
func GetTicketActivities(db *gorm.DB) gin.HandlerFunc {
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

		// Get activities using service
		activityService := services.NewActivityService(db)
		activities, err := activityService.GetTicketActivities(uint(ticketID), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch activities"})
			return
		}

		// Get total count
		total, err := activityService.GetTicketActivityCount(uint(ticketID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch activity count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"activities": activities,
			"total":      total,
			"limit":      limit,
			"offset":     offset,
		})
	}
}

// GetTicketTimeline retrieves a formatted timeline view of ticket activities
func GetTicketTimeline(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID from URL
		ticketIDStr := c.Param("id")
		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		// Verify ticket exists
		var ticket models.Ticket
		if err := db.First(&ticket, uint(ticketID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Get all activities for timeline (no pagination for timeline view)
		activityService := services.NewActivityService(db)
		activities, err := activityService.GetTicketActivities(uint(ticketID), 1000, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch timeline"})
			return
		}

		// Group activities by date for timeline display
		timeline := make(map[string][]models.TicketActivity)
		for _, activity := range activities {
			dateKey := activity.CreatedAt.Format("2006-01-02")
			timeline[dateKey] = append(timeline[dateKey], activity)
		}

		c.JSON(http.StatusOK, gin.H{
			"timeline": timeline,
			"ticket":   ticket,
		})
	}
}

// GetTicketFullDetails retrieves ticket with all related data (comments, calls, activities)
func GetTicketFullDetails(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ticket ID from URL
		ticketIDStr := c.Param("id")
		ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
			return
		}

		// Get ticket with basic relationships
		var ticket models.Ticket
		if err := db.Preload("Engineer").
			Preload("Product").
			Preload("Contact.Account").
			Preload("Account").
			First(&ticket, uint(ticketID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}

		// Get comments (limit recent ones)
		var comments []models.TicketComment
		commentQuery := db.Where("ticket_id = ?", ticketID).
			Preload("User").
			Order("created_at DESC").
			Limit(20)

		// Filter internal comments for non-managers
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok && user.RoleID != 2 { // Not a manager
				commentQuery = commentQuery.Where("is_internal = ?", false)
			}
		}
		commentQuery.Find(&comments)

		// Get recent calls
		var calls []models.TicketCall
		db.Where("ticket_id = ?", ticketID).
			Preload("User").
			Order("created_at DESC").
			Limit(10).
			Find(&calls)

		// Get recent activities
		activityService := services.NewActivityService(db)
		activities, _ := activityService.GetTicketActivities(uint(ticketID), 20, 0)

		// Get counts
		var commentCount, callCount, activityCount int64
		db.Model(&models.TicketComment{}).Where("ticket_id = ?", ticketID).Count(&commentCount)
		db.Model(&models.TicketCall{}).Where("ticket_id = ?", ticketID).Count(&callCount)
		activityCount, _ = activityService.GetTicketActivityCount(uint(ticketID))

		c.JSON(http.StatusOK, gin.H{
			"ticket":     ticket,
			"comments":   comments,
			"calls":      calls,
			"activities": activities,
			"counts": gin.H{
				"comments":   commentCount,
				"calls":      callCount,
				"activities": activityCount,
			},
		})
	}
}
