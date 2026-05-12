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

type CreateTaskInput struct {
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"` // Format: YYYY-MM-DD
	AssignedTo  *uint  `json:"assigned_to"`
	Priority    string `json:"priority" binding:"required"`
}

// Manager: Create task
func ManagerCreateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		var createdBy uint
		if exists {
			if user, ok := userVal.(models.User); ok {
				createdBy = user.ID
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		var input CreateTaskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse due date if provided
		var dueDate *time.Time
		if input.DueDate != "" {
			if parsed, err := time.Parse("2006-01-02", input.DueDate); err == nil {
				dueDate = &parsed
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due date format, use YYYY-MM-DD"})
				return
			}
		}

		// Create task
		task := models.Task{
			Subject:     input.Subject,
			Description: input.Description,
			DueDate:     dueDate,
			TaskStatus:  "Not Started",
			Priority:    input.Priority,
			AssignedTo:  input.AssignedTo,
			CreatedBy:   createdBy,
		}

		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
			return
		}

		// Log task creation activity
		taskService := services.NewTaskService(db)
		taskService.LogTaskCreation(task.ID, &createdBy, task.Subject)

		// Audit log for task creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTaskCreated,
			models.EntityTypeTask,
			&task.ID,
			task.Subject,
			fmt.Sprintf("Task created: %s", task.Subject),
			nil,
			task,
		)

		// Send notifications for task creation and assignment
		fmt.Printf("🔥 DEBUG: ManagerCreateTaskHandler - sending task creation notifications\n")
		notificationService := services.NewNotificationService(db)

		if input.AssignedTo != nil && *input.AssignedTo != 0 {
			// Only notify if the assigned engineer is different from the creator (no self-notifications)
			if *input.AssignedTo != createdBy {
				// Notify engineer about task assignment
				taskData := map[string]string{
					"task_subject": task.Subject,
					"task_id":      fmt.Sprintf("%d", task.ID),
					"priority":     task.Priority,
				}

				fmt.Printf("🔥 DEBUG: Sending task assignment notification to engineer ID %d (creator ID: %d)\n", *input.AssignedTo, createdBy)

				if err := notificationService.CreateNotification(services.NotificationData{
					RecipientID:      *input.AssignedTo,
					RecipientType:    "user",
					NotificationType: models.NotificationTaskAssignedToYou,
					Variables:        taskData,
					RelatedID:        &task.ID,
					RelatedType:      "task",
					ActorID:          &createdBy,
					ActorType:        "user",
				}); err != nil {
					fmt.Printf("❌ Failed to send task assignment notifications: %v\n", err)
				} else {
					fmt.Printf("✅ Task assignment notifications sent successfully\n")
				}
			} else {
				fmt.Printf("🔥 DEBUG: Skipping self-notification - manager assigned task to themselves\n")
			}
		}

		c.JSON(http.StatusCreated, gin.H{"task": task})
	}
}

// Manager: Get all tasks
func ManagerGetTasksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		var tasks []models.Task
		err := db.Preload("AssignedUser").Preload("Creator").Order("created_at DESC").Find(&tasks).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tasks"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

// Manager: Get task by ID
func ManagerGetTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		id := c.Param("id")
		var task models.Task
		err := db.Preload("AssignedUser").Preload("Creator").First(&task, id).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

// Manager: Edit task
func ManagerEditTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		id := c.Param("id")
		var input struct {
			Subject     *string `json:"subject"`
			Description *string `json:"description"`
			DueDate     *string `json:"due_date"` // Format: YYYY-MM-DD
			AssignedTo  *uint   `json:"assigned_to"`
			TaskStatus  *string `json:"task_status"`
			Priority    *string `json:"priority"`
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

		// Get user from context for activity logging
		userVal, exists := c.Get("user")
		var userID *uint
		if exists {
			if user, ok := userVal.(models.User); ok {
				userID = &user.ID
			}
		}

		// Initialize task service
		taskService := services.NewTaskService(db)

		// Store original values for activity logging
		originalSubject := task.Subject
		originalDescription := task.Description
		originalDueDate := task.DueDate
		originalAssignedTo := task.AssignedTo
		originalTaskStatus := task.TaskStatus
		originalPriority := task.Priority

		// Update fields and log activities
		if input.Subject != nil && *input.Subject != originalSubject {
			task.Subject = *input.Subject
			taskService.LogTaskFieldChange(task.ID, userID, models.ActivityTaskUpdated, "Subject", originalSubject, *input.Subject)
		}

		if input.Description != nil && *input.Description != originalDescription {
			task.Description = *input.Description
			taskService.LogTaskFieldChange(task.ID, userID, models.ActivityTaskUpdated, "Description", originalDescription, *input.Description)
		}

		if input.DueDate != nil {
			var newDueDate *time.Time
			if *input.DueDate != "" {
				if parsed, err := time.Parse("2006-01-02", *input.DueDate); err == nil {
					newDueDate = &parsed
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due date format, use YYYY-MM-DD"})
					return
				}
			}

			// Compare due dates
			oldDateStr := "Not Set"
			if originalDueDate != nil {
				oldDateStr = originalDueDate.Format("2006-01-02")
			}
			newDateStr := "Not Set"
			if newDueDate != nil {
				newDateStr = newDueDate.Format("2006-01-02")
			}

			if oldDateStr != newDateStr {
				task.DueDate = newDueDate
				taskService.LogTaskFieldChange(task.ID, userID, models.ActivityTaskUpdated, "Due Date", oldDateStr, newDateStr)
			}
		}

		if input.AssignedTo != nil {
			var originalAssignedToID uint
			if originalAssignedTo != nil {
				originalAssignedToID = *originalAssignedTo
			}
			if *input.AssignedTo != originalAssignedToID {
				task.AssignedTo = input.AssignedTo

				// Get user names for activity logging
				oldVal := "Unassigned"
				if originalAssignedTo != nil {
					if oldUser, err := utils.GetUserName(db, *originalAssignedTo); err == nil {
						oldVal = oldUser
					}
				}
				newVal := "Unassigned"
				if input.AssignedTo != nil && *input.AssignedTo != 0 {
					if newUser, err := utils.GetUserName(db, *input.AssignedTo); err == nil {
						newVal = newUser
					}
				}
				taskService.LogTaskAssigneeChange(task.ID, userID, oldVal, newVal)
			}
		}

		if input.TaskStatus != nil && *input.TaskStatus != originalTaskStatus {
			task.TaskStatus = *input.TaskStatus
			taskService.LogTaskStatusChange(task.ID, userID, originalTaskStatus, *input.TaskStatus)
		}

		if input.Priority != nil && *input.Priority != originalPriority {
			task.Priority = *input.Priority
			taskService.LogTaskPriorityChange(task.ID, userID, originalPriority, *input.Priority)
		}

		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
			return
		}

		// Audit log for task update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditTaskUpdated,
			models.EntityTypeTask,
			&task.ID,
			task.Subject,
			fmt.Sprintf("Task updated: %s", task.Subject),
			nil,
			task,
		)

		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

// Manager: Delete task
func ManagerDeleteTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsManager(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		id := c.Param("id")

		// Start transaction for data integrity
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Get task details first (for logging and validation)
		var task models.Task
		if err := tx.First(&task, id).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find task"})
			}
			return
		}

		fmt.Printf("🗑️ DELETING TASK: \"%s\" (ID: %s)\n", task.Subject, id)

		// Delete related records in correct order

		// 1. Delete task comments
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskComment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task comments"})
			return
		}

		// 2. Delete task activities
		if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskActivity{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task activities"})
			return
		}

		// 3. Log deletion activity before deleting the task
		userVal, exists := c.Get("user")
		if exists {
			if user, ok := userVal.(models.User); ok {
				description := fmt.Sprintf("Task \"%s\" deleted by manager", task.Subject)
				activity := models.TaskActivity{
					TaskID:       task.ID,
					UserID:       &user.ID,
					ActivityType: models.ActivityTaskDeleted,
					Description:  description,
				}
				if err := tx.Create(&activity).Error; err != nil {
					fmt.Printf("⚠️ Warning: Could not log deletion activity: %v\n", err)
					// Don't fail the deletion for logging issues
				}

				// Audit log for task deletion
				auditService := services.NewAuditService(db)
				auditService.LogCRUD(
					c,
					models.AuditTaskDeleted,
					models.EntityTypeTask,
					&task.ID,
					task.Subject,
					fmt.Sprintf("Task deleted: %s", task.Subject),
					task,
					nil,
				)
			}
		}

		// 4. Finally, delete the main task
		if err := tx.Delete(&task).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete task deletion"})
			return
		}

		fmt.Printf("✅ TASK DELETED SUCCESSFULLY: \"%s\"\n", task.Subject)
		c.JSON(http.StatusOK, gin.H{
			"message":      "Task deleted successfully",
			"task_subject": task.Subject,
		})
	}
}
