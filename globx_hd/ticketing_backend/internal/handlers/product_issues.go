package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create Product Issue
func CreateProductIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {
			ProductID uint   `json:"product_id" binding:"required"`
			IssueName string `json:"issue_name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		issue := models.MasterProductIssue{
			ProductID: in.ProductID,
			IssueName: in.IssueName,
			CreatedAt: time.Now(),
		}
		if err := db.Create(&issue).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Audit log
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditProductIssueCreated,
			models.EntityTypeProductIssue,
			&issue.ID,
			issue.IssueName,
			fmt.Sprintf("Product issue created: %s", issue.IssueName),
			nil,
			issue,
		)

		c.JSON(http.StatusOK, issue)
	}
}

// Get All Product Issues
func GetProductIssues(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var issues []models.MasterProductIssue
		if err := db.Preload("Product").Find(&issues).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, issues)
	}
}

// Get Single Product Issue
func GetProductIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var issue models.MasterProductIssue
		id := c.Param("id")
		if err := db.Preload("Product").First(&issue, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}
		c.JSON(http.StatusOK, issue)
	}
}

// Update Product Issue
func UpdateProductIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {
			IssueName string `json:"issue_name"`
		}
		id := c.Param("id")
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var issue models.MasterProductIssue
		if err := db.First(&issue, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}
		oldIssueName := issue.IssueName
		if in.IssueName != "" {
			issue.IssueName = in.IssueName
		}
		if err := db.Save(&issue).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Audit log
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditProductIssueUpdated,
			models.EntityTypeProductIssue,
			&issue.ID,
			issue.IssueName,
			fmt.Sprintf("Product issue updated: %s", issue.IssueName),
			map[string]string{"issue_name": oldIssueName},
			map[string]string{"issue_name": issue.IssueName},
		)

		c.JSON(http.StatusOK, issue)
	}
}

// Delete Product Issue
func DeleteProductIssue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Get issue first for audit logging
		var issue models.MasterProductIssue
		if err := db.First(&issue, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "issue not found"})
			return
		}

		if err := db.Delete(&models.MasterProductIssue{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Audit log
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditProductIssueDeleted,
			models.EntityTypeProductIssue,
			&issue.ID,
			issue.IssueName,
			fmt.Sprintf("Product issue deleted: %s", issue.IssueName),
			issue,
			nil,
		)

		c.JSON(http.StatusOK, gin.H{"message": "issue deleted"})
	}
}
