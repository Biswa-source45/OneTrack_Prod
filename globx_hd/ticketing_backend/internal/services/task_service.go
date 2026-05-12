package services

import (
	"fmt"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// TaskService handles task-related operations
type TaskService struct {
	db *gorm.DB
}

// NewTaskService creates a new task service
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// LogTaskActivity logs a new activity for a task
func (s *TaskService) LogTaskActivity(taskID uint, userID *uint, activityType, description string) error {
	activity := models.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
	}
	return s.db.Create(&activity).Error
}

// LogTaskFieldChange logs a field change activity for a task
func (s *TaskService) LogTaskFieldChange(taskID uint, userID *uint, activityType, fieldName, oldValue, newValue string) error {
	description := fmt.Sprintf("%s changed from '%s' to '%s'", fieldName, oldValue, newValue)
	activity := models.TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		OldValue:     oldValue,
		NewValue:     newValue,
	}
	return s.db.Create(&activity).Error
}

// LogTaskCreation logs task creation activity
func (s *TaskService) LogTaskCreation(taskID uint, userID *uint, subject string) error {
	description := fmt.Sprintf("Task '%s' created", subject)
	return s.LogTaskActivity(taskID, userID, models.ActivityTaskCreated, description)
}

// LogTaskStatusChange logs task status change
func (s *TaskService) LogTaskStatusChange(taskID uint, userID *uint, oldStatus, newStatus string) error {
	return s.LogTaskFieldChange(taskID, userID, models.ActivityTaskStatusChanged, "Status", oldStatus, newStatus)
}

// LogTaskAssigneeChange logs task assignee change
func (s *TaskService) LogTaskAssigneeChange(taskID uint, userID *uint, oldAssignee, newAssignee string) error {
	return s.LogTaskFieldChange(taskID, userID, models.ActivityTaskAssigneeChanged, "Assigned To", oldAssignee, newAssignee)
}

// LogTaskPriorityChange logs task priority change
func (s *TaskService) LogTaskPriorityChange(taskID uint, userID *uint, oldPriority, newPriority string) error {
	return s.LogTaskFieldChange(taskID, userID, models.ActivityTaskPriorityChanged, "Priority", oldPriority, newPriority)
}

// LogTaskComment logs comment addition
func (s *TaskService) LogTaskComment(taskID uint, userID *uint) error {
	description := "Comment added"
	return s.LogTaskActivity(taskID, userID, models.ActivityTaskCommentAdded, description)
}

// GetTaskActivities retrieves activities for a task with pagination
func (s *TaskService) GetTaskActivities(taskID uint, limit, offset int) ([]models.TaskActivity, error) {
	var activities []models.TaskActivity
	err := s.db.Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	return activities, err
}

// GetTaskActivityCount gets total count of activities for a task
func (s *TaskService) GetTaskActivityCount(taskID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.TaskActivity{}).Where("task_id = ?", taskID).Count(&count).Error
	return count, err
}
