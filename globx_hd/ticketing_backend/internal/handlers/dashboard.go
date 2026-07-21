package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ManagerDashboardStatsHandler returns ticket statistics for manager dashboard
func ManagerDashboardStatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// Get month and year parameters (optional)
		monthParam := c.Query("month")
		yearParam := c.Query("year")

		var startDate, endDate time.Time
		var hasDateFilter bool

		// If month and year are provided, filter by that month
		if monthParam != "" && yearParam != "" {
			month, err1 := strconv.Atoi(monthParam)
			year, err2 := strconv.Atoi(yearParam)

			if err1 != nil || err2 != nil || month < 1 || month > 12 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month or year parameter"})
				return
			}

			// Create start and end dates for the month
			startDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, 0) // First day of next month
			hasDateFilter = true
		}

		// Build base query for total count
		baseQuery := db.Model(&models.Ticket{})
		if hasDateFilter {
			baseQuery = baseQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}

		// Get total tickets count
		var totalCount int64
		if err := baseQuery.Count(&totalCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch total tickets count"})
			return
		}

		// Get open tickets count (OPEN status only)
		var openCount int64
		openQuery := db.Model(&models.Ticket{})
		if hasDateFilter {
			openQuery = openQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := openQuery.Where("ticket_status = ?", "OPEN").Count(&openCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch open tickets count"})
			return
		}

		// Get closed tickets count
		var closedCount int64
		closedQuery := db.Model(&models.Ticket{})
		if hasDateFilter {
			closedQuery = closedQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := closedQuery.Where("ticket_status = ?", "CLOSED").Count(&closedCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch closed tickets count"})
			return
		}

		// Get in progress tickets count
		var inProgressCount int64
		inProgressQuery := db.Model(&models.Ticket{})
		if hasDateFilter {
			inProgressQuery = inProgressQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := inProgressQuery.Where("ticket_status = ?", "IN PROGRESS").Count(&inProgressCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch in progress tickets count"})
			return
		}

		// Get resolved tickets count
		var resolvedCount int64
		resolvedQuery := db.Model(&models.Ticket{})
		if hasDateFilter {
			resolvedQuery = resolvedQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := resolvedQuery.Where("ticket_status = ?", "RESOLVED").Count(&resolvedCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch resolved tickets count"})
			return
		}

		// Return statistics
		c.JSON(http.StatusOK, gin.H{
			"total_tickets":       totalCount,
			"open_tickets":        openCount,
			"closed_tickets":      closedCount,
			"in_progress_tickets": inProgressCount,
			"resolved_tickets":    resolvedCount,
			"month":               monthParam,
			"year":                yearParam,
		})
	}
}
