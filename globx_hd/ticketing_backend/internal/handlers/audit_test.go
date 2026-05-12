package handlers

import (
	"fmt"
	"net/http"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestAuditLog is a test endpoint to verify audit logging is working
// This should be removed in production or protected with admin-only access
func TestAuditLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auditService := services.NewAuditService(db)

		// Test 1: Simple authentication log
		testID := uint(999)
		err1 := auditService.LogAuthentication(
			models.ActorTypeSystem,
			&testID,
			"Test User",
			"test@example.com",
			"TEST_ACTION",
			true,
			c.ClientIP(),
			c.Request.UserAgent(),
			"",
		)

		// Test 2: Check if table exists
		var tableExists bool
		err2 := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'audit_logs')").Scan(&tableExists).Error

		// Test 3: Count audit logs
		var count int64
		err3 := db.Model(&models.AuditLog{}).Count(&count).Error

		// Test 4: Get recent logs
		var recentLogs []models.AuditLog
		err4 := db.Model(&models.AuditLog{}).Order("created_at DESC").Limit(5).Find(&recentLogs).Error

		response := gin.H{
			"test_1_create_log": map[string]interface{}{
				"success": err1 == nil,
				"error":   fmt.Sprintf("%v", err1),
			},
			"test_2_table_exists": map[string]interface{}{
				"exists": tableExists,
				"error":  fmt.Sprintf("%v", err2),
			},
			"test_3_count_logs": map[string]interface{}{
				"count": count,
				"error": fmt.Sprintf("%v", err3),
			},
			"test_4_recent_logs": map[string]interface{}{
				"count": len(recentLogs),
				"logs":  recentLogs,
				"error": fmt.Sprintf("%v", err4),
			},
		}

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
