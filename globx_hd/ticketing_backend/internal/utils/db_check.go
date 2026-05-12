package utils

import (
	"fmt"
	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// CheckNotificationTables checks if notification tables exist and have data
func CheckNotificationTables(db *gorm.DB) {
	fmt.Println("🔍 Checking notification system...")
	
	// Check if notifications table exists
	if !db.Migrator().HasTable(&models.Notification{}) {
		fmt.Println("❌ notifications table does not exist")
		return
	}
	fmt.Println("✅ notifications table exists")
	
	// Check if notification_templates table exists
	if !db.Migrator().HasTable(&models.NotificationTemplate{}) {
		fmt.Println("❌ notification_templates table does not exist")
		return
	}
	fmt.Println("✅ notification_templates table exists")
	
	// Check template count
	var templateCount int64
	db.Model(&models.NotificationTemplate{}).Count(&templateCount)
	fmt.Printf("📊 Found %d notification templates\n", templateCount)
	
	if templateCount == 0 {
		fmt.Println("❌ No notification templates found - this is likely the issue!")
		return
	}
	
	// Check for specific templates that are failing
	requiredTemplates := []string{
		models.NotificationTicketAssignedToYou,
		models.NotificationTicketEngineerAssigned,
		models.NotificationTicketStatusChanged,
		models.NotificationTicketCallLogged,
	}
	
	for _, templateType := range requiredTemplates {
		var template models.NotificationTemplate
		err := db.Where("notification_type = ? AND is_active = ?", templateType, true).First(&template).Error
		if err != nil {
			fmt.Printf("❌ Missing template: %s\n", templateType)
		} else {
			fmt.Printf("✅ Found template: %s\n", templateType)
		}
	}
	
	// Check notification count
	var notificationCount int64
	db.Model(&models.Notification{}).Count(&notificationCount)
	fmt.Printf("📊 Found %d notifications in database\n", notificationCount)
}
