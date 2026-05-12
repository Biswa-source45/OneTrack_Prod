package handlers

import (
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTaskActivities retrieves activity history for a task
func GetTaskActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get task ID from URL
		taskIDStr := c.Param("id")
		taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
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
			limit = 100
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		// Get activities using service
		taskService := services.NewTaskService(db)
		activities, err := taskService.GetTaskActivities(uint(taskID), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch activities"})
			return
		}

		// Get total count
		total, err := taskService.GetTaskActivityCount(uint(taskID))
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
