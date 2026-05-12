package handlers

import (
	"fmt"
	"net/http"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateAccountInput struct {
	AccountName  string `json:"account_name" binding:"required"`
	AccountOwner string `json:"account_owner"`
	Address      string `json:"address"`
}

func CreateAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateAccountInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code, err := utils.GenerateUnique3Digit(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate customer code"})
			return
		}

		acc := models.Account{
			AccountName:  in.AccountName,
			AccountOwner: in.AccountOwner,
			Address:      in.Address,
			CustomerCode: code,
		}

		if err := db.Create(&acc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log account creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditAccountCreated,
			models.EntityTypeAccount,
			&acc.ID,
			acc.AccountName,
			fmt.Sprintf("Account created: %s (Code: %s)", acc.AccountName, acc.CustomerCode),
			nil,
			acc,
		)

		c.JSON(http.StatusCreated, acc)
	}
}

func GetAccounts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []models.Account
		if err := db.Preload("Contacts").Find(&accounts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, accounts)
	}
}

func GetAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var acc models.Account
		if err := db.Preload("Contacts").First(&acc, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusOK, acc)
	}
}

type UpdateAccountInput struct {
	AccountName  *string `json:"account_name"`
	AccountOwner *string `json:"account_owner"`
	Address      *string `json:"address"`
}

func UpdateAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var acc models.Account
		if err := db.First(&acc, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}

		// Store old values for audit
		oldAcc := acc

		var in UpdateAccountInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Only update non-nil fields
		if in.AccountName != nil {
			acc.AccountName = *in.AccountName
		}
		if in.AccountOwner != nil {
			acc.AccountOwner = *in.AccountOwner
		}
		if in.Address != nil {
			acc.Address = *in.Address
		}

		if err := db.Save(&acc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log account update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditAccountUpdated,
			models.EntityTypeAccount,
			&acc.ID,
			acc.AccountName,
			fmt.Sprintf("Account updated: %s (Code: %s)", acc.AccountName, acc.CustomerCode),
			oldAcc,
			acc,
		)

		c.JSON(http.StatusOK, acc)
	}
}

func DeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Get account details before deletion for audit
		var acc models.Account
		if err := db.First(&acc, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}

		if err := db.Delete(&models.Account{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log account deletion
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditAccountDeleted,
			models.EntityTypeAccount,
			&acc.ID,
			acc.AccountName,
			fmt.Sprintf("Account deleted: %s (Code: %s)", acc.AccountName, acc.CustomerCode),
			acc,
			nil,
		)

		c.Status(http.StatusNoContent)
	}
}
