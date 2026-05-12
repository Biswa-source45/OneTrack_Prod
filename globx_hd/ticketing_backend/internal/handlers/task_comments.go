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

// CreateTaskComment adds a comment to a task
func CreateTaskComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get task ID from URL
		taskIDStr := c.Param("id")
		taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
			return
		}

		// Verify task exists
		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
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
		var input struct {
			Content    string `json:"content" binding:"required"`
			IsInternal bool   `json:"is_internal"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create comment
		comment := models.TaskComment{
			TaskID:     uint(taskID),
			UserID:     userID,
			Content:    input.Content,
			IsInternal: input.IsInternal,
		}

		if err := db.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create comment"})
			return
		}

		// Log activity
		taskService := services.NewTaskService(db)
		taskService.LogTaskComment(uint(taskID), &userID)

		// Send notifications for task comment
		fmt.Printf("🔥 DEBUG: CreateTaskComment - sending task comment notifications\n")
		notificationService := services.NewNotificationService(db)

		// Send notifications to relevant parties
		var taskCreator models.User
		commentPreview := input.Content
		if len(commentPreview) > 100 {
			commentPreview = commentPreview[:100]
		}
		taskData := map[string]string{
			"task_subject": task.Subject,
			"task_id":      fmt.Sprintf("%d", task.ID),
			"comment_text": commentPreview,
		}

		// Get task creator (manager)
		if err := db.First(&taskCreator, task.CreatedBy).Error; err == nil {
			fmt.Printf("🔥 DEBUG: Task creator loaded - ID: %d, Comment author ID: %d\n", taskCreator.ID, userID)

			// If engineer comments, notify manager
			if userID != taskCreator.ID {
				fmt.Printf("🔥 DEBUG: Engineer commented on task - notifying manager (ID: %d)\n", taskCreator.ID)
				if err := notificationService.CreateNotification(services.NotificationData{
					RecipientID:      taskCreator.ID,
					RecipientType:    "user",
					NotificationType: models.NotificationTaskCommentAdded,
					Variables:        taskData,
					RelatedID:        &task.ID,
					RelatedType:      "task",
					RelatedSubID:     &comment.ID,
					ActorID:          &userID,
					ActorType:        "user",
				}); err != nil {
					fmt.Printf("❌ Failed to send task comment notifications to manager: %v\n", err)
				} else {
					fmt.Printf("✅ Task comment notifications sent to manager successfully\n")
				}
			} else {
				// Manager is commenting - check if task is assigned to notify engineer
				fmt.Printf("🔥 DEBUG: Manager commented - checking if task is assigned. task.AssignedTo: %v\n", task.AssignedTo)
				if task.AssignedTo != nil && *task.AssignedTo != 0 {
					fmt.Printf("🔥 DEBUG: Manager commented on task - notifying assigned engineer (ID: %d)\n", *task.AssignedTo)
					if err := notificationService.CreateNotification(services.NotificationData{
						RecipientID:      *task.AssignedTo,
						RecipientType:    "user",
						NotificationType: models.NotificationTaskManagerCommentAdded,
						Variables:        taskData,
						RelatedID:        &task.ID,
						RelatedType:      "task",
						RelatedSubID:     &comment.ID,
						ActorID:          &userID,
						ActorType:        "user",
					}); err != nil {
						fmt.Printf("❌ Failed to send task comment notifications to engineer: %v\n", err)
					} else {
						fmt.Printf("✅ Task comment notifications sent to engineer successfully\n")
					}
				} else {
					fmt.Printf("🔥 DEBUG: Task not assigned to anyone - no engineer to notify\n")
				}
			}
		} else {
			fmt.Printf("❌ Failed to load task creator for comment notifications: %v\n", err)
		}

		// Load user information for response
		db.Preload("User").First(&comment, comment.ID)

		c.JSON(http.StatusCreated, gin.H{"comment": comment})
	}
}

// GetTaskComments retrieves comments for a task
func GetTaskComments(db *gorm.DB) gin.HandlerFunc {
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

		// Get comments
		var comments []models.TaskComment
		err = db.Where("task_id = ?", taskID).
			Preload("User").
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&comments).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch comments"})
			return
		}

		// Get total count
		var total int64
		db.Model(&models.TaskComment{}).Where("task_id = ?", taskID).Count(&total)

		c.JSON(http.StatusOK, gin.H{
			"comments": comments,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		})
	}
}

// UpdateTaskComment updates a task comment
func UpdateTaskComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get comment ID from URL
		commentIDStr := c.Param("commentId")
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

		// Find comment
		var comment models.TaskComment
		if err := db.First(&comment, commentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		// Check if user owns the comment or is manager
		if comment.UserID != user.ID && !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "can only edit your own comments"})
			return
		}

		// Bind input
		var input struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update comment
		comment.Content = input.Content
		if err := db.Save(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update comment"})
			return
		}

		// Load user information for response
		db.Preload("User").First(&comment, comment.ID)

		c.JSON(http.StatusOK, gin.H{"comment": comment})
	}
}

// DeleteTaskComment deletes a task comment
func DeleteTaskComment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get comment ID from URL
		commentIDStr := c.Param("commentId")
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

		// Find comment
		var comment models.TaskComment
		if err := db.First(&comment, commentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		// Check if user owns the comment or is manager
		if comment.UserID != user.ID && !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "can only delete your own comments"})
			return
		}

		// Delete comment
		if err := db.Delete(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete comment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
	}
}
