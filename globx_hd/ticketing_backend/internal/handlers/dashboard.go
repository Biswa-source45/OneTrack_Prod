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
		if err := openQuery.Where("ticket_status IN ?", []string{"OPEN", "Open"}).Count(&openCount).Error; err != nil {
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
		if err := inProgressQuery.Where("ticket_status IN ?", []string{"IN PROGRESS", "IN_PROGRESS", "In Progress"}).Count(&inProgressCount).Error; err != nil {
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

		// --- TASK STATISTICS ---
		// Total tasks count
		var totalTasksCount int64
		baseTaskQuery := db.Model(&models.Task{})
		if hasDateFilter {
			baseTaskQuery = baseTaskQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := baseTaskQuery.Count(&totalTasksCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch total tasks count"})
			return
		}

		// Not completed tasks count (In Progress / Not Started / TODO / ON_HOLD / Deferred / Waiting)
		var notCompletedTasksCount int64
		notCompletedTaskQuery := db.Model(&models.Task{})
		if hasDateFilter {
			notCompletedTaskQuery = notCompletedTaskQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := notCompletedTaskQuery.Where("task_status IN ?", []string{
			"Not Started", "TODO", "In Progress", "IN PROGRESS", "IN_PROGRESS", "ON_HOLD", "Deferred", "Waiting on someone else",
		}).Count(&notCompletedTasksCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch not completed tasks count"})
			return
		}

		// Review Pending / Waiting Verification tasks count
		var reviewPendingTasksCount int64
		reviewPendingTaskQuery := db.Model(&models.Task{})
		if hasDateFilter {
			reviewPendingTaskQuery = reviewPendingTaskQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := reviewPendingTaskQuery.Where("task_status IN ?", []string{
			"Review Pending", "REVIEW_PENDING", "Under Review", "Manager Review", "Pending Approval", "Waiting Verification",
		}).Count(&reviewPendingTasksCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch review pending tasks count"})
			return
		}

		// Completed tasks count
		var completedTasksCount int64
		completedTaskQuery := db.Model(&models.Task{})
		if hasDateFilter {
			completedTaskQuery = completedTaskQuery.Where("created_at >= ? AND created_at < ?", startDate, endDate)
		}
		if err := completedTaskQuery.Where("task_status IN ?", []string{
			"Completed", "COMPLETED", "Accepted",
		}).Count(&completedTasksCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch completed tasks count"})
			return
		}

		// Return statistics
		c.JSON(http.StatusOK, gin.H{
			"total_tickets":        totalCount,
			"open_tickets":         openCount,
			"closed_tickets":       closedCount,
			"in_progress_tickets":  inProgressCount,
			"resolved_tickets":     resolvedCount,
			"total_tasks":          totalTasksCount,
			"not_completed_tasks":  notCompletedTasksCount,
			"review_pending_tasks": reviewPendingTasksCount,
			"completed_tasks":      completedTasksCount,
			"month":                monthParam,
			"year":                 yearParam,
		})
	}
}

// EngineerDashboardStatsHandler returns ticket and task statistics for engineer dashboard
func EngineerDashboardStatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user, ok := userVal.(models.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		// Ticket statistics assigned to this engineer
		var totalTickets, openTickets, inProgressTickets, resolvedTickets, closedTickets int64

		db.Model(&models.Ticket{}).Where("assigned_engineer_id = ?", user.ID).Count(&totalTickets)
		db.Model(&models.Ticket{}).Where("assigned_engineer_id = ? AND ticket_status IN ?", user.ID, []string{"OPEN", "Open"}).Count(&openTickets)
		db.Model(&models.Ticket{}).Where("assigned_engineer_id = ? AND ticket_status IN ?", user.ID, []string{"IN PROGRESS", "IN_PROGRESS", "In Progress"}).Count(&inProgressTickets)
		db.Model(&models.Ticket{}).Where("assigned_engineer_id = ? AND ticket_status IN ?", user.ID, []string{"RESOLVED", "Resolved"}).Count(&resolvedTickets)
		db.Model(&models.Ticket{}).Where("assigned_engineer_id = ? AND ticket_status IN ?", user.ID, []string{"CLOSED", "Closed"}).Count(&closedTickets)

		// Task statistics assigned to this engineer
		var totalTasks, inProgressTasks, reviewPendingTasks, completedTasks int64

		db.Model(&models.Task{}).Where("assigned_user_id = ?", user.ID).Count(&totalTasks)
		db.Model(&models.Task{}).Where("assigned_user_id = ? AND task_status IN ?", user.ID, []string{"Not Started", "TODO", "In Progress", "IN PROGRESS", "IN_PROGRESS", "ON_HOLD", "Deferred"}).Count(&inProgressTasks)
		db.Model(&models.Task{}).Where("assigned_user_id = ? AND task_status IN ?", user.ID, []string{"Review Pending", "REVIEW_PENDING", "Under Review", "Manager Review", "Pending Approval", "Waiting Verification"}).Count(&reviewPendingTasks)
		db.Model(&models.Task{}).Where("assigned_user_id = ? AND task_status IN ?", user.ID, []string{"Completed", "COMPLETED", "ACCEPTED"}).Count(&completedTasks)

		c.JSON(http.StatusOK, gin.H{
			"total_tickets":        totalTickets,
			"open_tickets":         openTickets,
			"in_progress_tickets":  inProgressTickets,
			"resolved_tickets":     resolvedTickets,
			"closed_tickets":       closedTickets,
			"total_tasks":          totalTasks,
			"in_progress_tasks":    inProgressTasks,
			"review_pending_tasks": reviewPendingTasks,
			"completed_tasks":      completedTasks,
		})
	}
}
