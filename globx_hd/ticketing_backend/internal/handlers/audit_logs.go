package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAuditLogsHandler retrieves audit logs with filtering and pagination
func GetAuditLogsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		auditService := services.NewAuditService(db)

		// Parse query parameters
		filters := services.AuditLogFilters{
			Page:  1,
			Limit: 50,
		}

		// Parse page
		if pageStr := c.Query("page"); pageStr != "" {
			if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
				filters.Page = page
			}
		}

		// Parse limit
		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
				filters.Limit = limit
			}
		}

		// Parse actor_id
		if actorIDStr := c.Query("actor_id"); actorIDStr != "" {
			if actorID, err := strconv.ParseUint(actorIDStr, 10, 32); err == nil {
				actorIDUint := uint(actorID)
				filters.ActorID = &actorIDUint
			}
		}

		// Parse entity_id
		if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
			if entityID, err := strconv.ParseUint(entityIDStr, 10, 32); err == nil {
				entityIDUint := uint(entityID)
				filters.EntityID = &entityIDUint
			}
		}

		// Parse other filters
		filters.ActorType = c.Query("actor_type")
		filters.Action = c.Query("action")
		filters.EntityType = c.Query("entity_type")
		filters.Severity = c.Query("severity")
		filters.Status = c.Query("status")
		filters.Search = c.Query("search")

		// Parse date range
		if startDateStr := c.Query("start_date"); startDateStr != "" {
			if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
				filters.StartDate = startDate
			}
		}
		if endDateStr := c.Query("end_date"); endDateStr != "" {
			if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
				// Set to end of day
				filters.EndDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			}
		}

		// Get audit logs
		logs, total, err := auditService.GetAuditLogs(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
			return
		}

		// Calculate pagination metadata
		totalPages := (int(total) + filters.Limit - 1) / filters.Limit

		c.JSON(http.StatusOK, gin.H{
			"logs": logs,
			"pagination": gin.H{
				"current_page": filters.Page,
				"total_pages":  totalPages,
				"total_count":  total,
				"limit":        filters.Limit,
			},
		})
	}
}

// GetAuditLogByIDHandler retrieves a single audit log by ID
func GetAuditLogByIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audit log ID"})
			return
		}

		auditService := services.NewAuditService(db)
		log, err := auditService.GetAuditLogByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"log": log})
	}
}

// GetRecentAuditLogsHandler retrieves recent audit logs (last 30 days)
func GetRecentAuditLogsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		limit := 100
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
				limit = l
			}
		}

		auditService := services.NewAuditService(db)
		logs, err := auditService.GetRecentAuditLogs(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent audit logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// GetCriticalAuditLogsHandler retrieves critical severity audit logs
func GetCriticalAuditLogsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		limit := 100
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
				limit = l
			}
		}

		auditService := services.NewAuditService(db)
		logs, err := auditService.GetCriticalAuditLogs(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch critical audit logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// GetFailedAuditLogsHandler retrieves failed operation audit logs
func GetFailedAuditLogsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		limit := 100
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
				limit = l
			}
		}

		auditService := services.NewAuditService(db)
		logs, err := auditService.GetFailedAuditLogs(limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch failed audit logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// GetAuditLogsByEntityHandler retrieves audit logs for a specific entity
func GetAuditLogsByEntityHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		entityType := c.Param("entity_type")
		entityIDStr := c.Param("entity_id")
		entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
			return
		}

		limit := 100
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
				limit = l
			}
		}

		auditService := services.NewAuditService(db)
		logs, err := auditService.GetAuditLogsByEntity(entityType, uint(entityID), limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entity audit logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// GetAuditLogsByActorHandler retrieves audit logs for a specific actor
func GetAuditLogsByActorHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		actorType := c.Param("actor_type")
		actorIDStr := c.Param("actor_id")
		actorID, err := strconv.ParseUint(actorIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid actor ID"})
			return
		}

		limit := 100
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
				limit = l
			}
		}

		auditService := services.NewAuditService(db)
		logs, err := auditService.GetAuditLogsByActor(actorType, uint(actorID), limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch actor audit logs"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

// GetAuditLogStatsHandler retrieves audit log statistics
func GetAuditLogStatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only managers can view audit logs
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only managers can view audit logs"})
			return
		}

		// Get statistics for last 30 days
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

		var stats struct {
			TotalLogs         int64 `json:"total_logs"`
			CriticalLogs      int64 `json:"critical_logs"`
			FailedLogs        int64 `json:"failed_logs"`
			AuthEvents        int64 `json:"auth_events"`
			UserActions       int64 `json:"user_actions"`
			TicketActions     int64 `json:"ticket_actions"`
			MasterDataActions int64 `json:"master_data_actions"`
		}

		// Total logs in last 30 days
		db.Model(&models.AuditLog{}).Where("created_at >= ?", thirtyDaysAgo).Count(&stats.TotalLogs)

		// Critical logs
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND severity = ?", thirtyDaysAgo, models.SeverityCritical).Count(&stats.CriticalLogs)

		// Failed logs
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND status IN (?)", thirtyDaysAgo, []string{models.StatusFailure, models.StatusError}).Count(&stats.FailedLogs)

		// Auth events
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND entity_type = ?", thirtyDaysAgo, "authentication").Count(&stats.AuthEvents)

		// User actions
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND entity_type IN (?)", thirtyDaysAgo, []string{"user", "contact"}).Count(&stats.UserActions)

		// Ticket actions
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND entity_type = ?", thirtyDaysAgo, "ticket").Count(&stats.TicketActions)

		// Master data actions
		db.Model(&models.AuditLog{}).Where("created_at >= ? AND entity_type IN (?)", thirtyDaysAgo, []string{"product", "role", "designation", "account"}).Count(&stats.MasterDataActions)

		c.JSON(http.StatusOK, gin.H{"stats": stats})
	}
}
