package handlers

import (
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetNotificationsHandler retrieves notifications for the authenticated user
func GetNotificationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userVal, exists := c.Get("user")
		var userID uint
		var userType string
		
		if exists {
			// Staff user (manager/engineer)
			user, ok := userVal.(models.User)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
			userID = user.ID
			userType = "user"
		} else {
			// Customer contact
			contactIDVal, exists := c.Get("contact_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
				return
			}
			contactID, ok := contactIDVal.(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
			userID = contactID
			userType = "contact"
		}

		// Parse query parameters
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")
		category := c.DefaultQuery("category", "all")
		priority := c.DefaultQuery("priority", "all")
		isRead := c.DefaultQuery("is_read", "all")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			limit = 20
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		// Create notification service and get notifications
		notificationService := services.NewNotificationService(db)
		filters := services.NotificationFilters{
			Category: category,
			Priority: priority,
			IsRead:   isRead,
		}

		notifications, err := notificationService.GetUserNotifications(userID, userType, limit, offset, filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"notifications": notifications})
	}
}

// GetUnreadCountHandler returns the count of unread notifications
func GetUnreadCountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userVal, exists := c.Get("user")
		var userID uint
		var userType string
		
		if exists {
			// Staff user (manager/engineer)
			user, ok := userVal.(models.User)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
			userID = user.ID
			userType = "user"
		} else {
			// Customer contact
			contactIDVal, exists := c.Get("contact_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
				return
			}
			contactID, ok := contactIDVal.(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
			userID = contactID
			userType = "contact"
		}

		// Get unread count
		notificationService := services.NewNotificationService(db)
		count, err := notificationService.GetUnreadCount(userID, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch unread count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"unread_count": count})
	}
}

// MarkNotificationReadHandler marks a specific notification as read
func MarkNotificationReadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get notification ID from URL
		notificationIDStr := c.Param("id")
		notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		var userID uint
		var userType string
		
		if exists {
			// Staff user (manager/engineer)
			user, ok := userVal.(models.User)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
			userID = user.ID
			userType = "user"
		} else {
			// Customer contact
			contactIDVal, exists := c.Get("contact_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
				return
			}
			contactID, ok := contactIDVal.(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
			userID = contactID
			userType = "contact"
		}

		// Mark notification as read
		notificationService := services.NewNotificationService(db)
		err = notificationService.MarkAsRead(uint(notificationID), userID, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not mark notification as read"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
	}
}

// MarkAllNotificationsReadHandler marks all notifications as read for the user
func MarkAllNotificationsReadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userVal, exists := c.Get("user")
		var userID uint
		var userType string
		
		if exists {
			// Staff user (manager/engineer)
			user, ok := userVal.(models.User)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
			userID = user.ID
			userType = "user"
		} else {
			// Customer contact
			contactIDVal, exists := c.Get("contact_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
				return
			}
			contactID, ok := contactIDVal.(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
			userID = contactID
			userType = "contact"
		}

		// Mark all notifications as read
		notificationService := services.NewNotificationService(db)
		err := notificationService.MarkAllAsRead(userID, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not mark all notifications as read"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
	}
}

// DeleteNotificationHandler deletes a specific notification
func DeleteNotificationHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get notification ID from URL
		notificationIDStr := c.Param("id")
		notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
			return
		}

		// Get user from context
		userVal, exists := c.Get("user")
		var userID uint
		var userType string
		
		if exists {
			// Staff user (manager/engineer)
			user, ok := userVal.(models.User)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user in context"})
				return
			}
			userID = user.ID
			userType = "user"
		} else {
			// Customer contact
			contactIDVal, exists := c.Get("contact_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
				return
			}
			contactID, ok := contactIDVal.(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid contact in context"})
				return
			}
			userID = contactID
			userType = "contact"
		}

		// Delete notification
		notificationService := services.NewNotificationService(db)
		err = notificationService.DeleteNotification(uint(notificationID), userID, userType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete notification"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
	}
}
