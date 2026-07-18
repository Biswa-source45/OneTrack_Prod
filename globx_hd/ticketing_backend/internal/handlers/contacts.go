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

type CreateContactInput struct {
	AccountID     *uint   `json:"account_id"` // Optional for Individual contacts
	DesignationID uint    `json:"designation_id" binding:"required"`
	ContactType   string  `json:"contact_type" binding:"required"`
	Department    string  `json:"department"`
	Location      string  `json:"location"`
	FirstName     string  `json:"first_name" binding:"required"`
	LastName      string  `json:"last_name"`
	Email         *string `json:"email" binding:"omitempty,email"`
	Mobile        string  `json:"mobile" binding:"required"`
	Password      *string `json:"password" binding:"omitempty,min=6"`
}

func CreateContact(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateContactInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate contact type
		validContactTypes := []string{"Govt.", "Private", "Individual"}
		isValidType := false
		for _, validType := range validContactTypes {
			if in.ContactType == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contact_type must be one of: Govt., Private, Individual"})
			return
		}

		// Validate account requirement for Govt. and Private
		if (in.ContactType == "Govt." || in.ContactType == "Private") && in.AccountID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "account_id is required for Govt. and Private contacts"})
			return
		}

		var customerCode string
		if in.AccountID != nil {
			// Validate account exists
			var acc models.Account
			if err := db.First(&acc, *in.AccountID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
				return
			}
			customerCode = acc.CustomerCode
		} else {
			// For Individual contacts, generate unique customer code
			code, err := utils.GenerateUniqueCustomerCode(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate customer code"})
				return
			}
			customerCode = code
		}

		var emailPtr *string
		if in.Email != nil && *in.Email != "" {
			emailPtr = in.Email
		}

		var hashPtr *string
		if in.Password != nil && *in.Password != "" {
			hash, err := utils.HashPassword(*in.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
				return
			}
			hashPtr = &hash
		}

		contact := models.Contact{
			AccountID:     in.AccountID,
			DesignationID: in.DesignationID,
			ContactType:   in.ContactType,
			Department:    in.Department,
			Location:      in.Location,
			FirstName:     in.FirstName,
			LastName:      in.LastName,
			Email:         emailPtr,
			Mobile:        in.Mobile,
			PasswordHash:  hashPtr,
			FirstLogin:    true,
			CustomerCode:  customerCode,
		}

		if err := db.Create(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log contact creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditContactCreated,
			models.EntityTypeContact,
			&contact.ID,
			fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
			fmt.Sprintf("Contact created: %s (%s) - Type: %s", derefString(contact.Email), contact.ContactType, contact.CustomerCode),
			nil,
			contact,
		)

		c.JSON(http.StatusCreated, contact)
	}
}

func GetContacts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var contacts []models.Contact
		if err := db.Preload("Designation").Preload("Account").Find(&contacts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, contacts)
	}
}

func GetContact(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var contact models.Contact
		if err := db.Preload("Designation").Preload("Account").First(&contact, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
			return
		}
		c.JSON(http.StatusOK, contact)
	}
}

type UpdateContactInput struct {
	AccountID     *uint   `json:"account_id"`
	DesignationID *uint   `json:"designation_id"`
	ContactType   *string `json:"contact_type"`
	Department    *string `json:"department"`
	Location      *string `json:"location"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	Email         *string `json:"email" binding:"omitempty,email"`
	Mobile        *string `json:"mobile"`
	Password      *string `json:"password"`
}

func UpdateContact(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var contact models.Contact
		if err := db.First(&contact, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
			return
		}

		// Store old values for audit
		oldContact := contact

		var in UpdateContactInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if in.AccountID != nil {
			contact.AccountID = in.AccountID
		}
		if in.ContactType != nil {
			// Validate contact type
			validContactTypes := []string{"Govt.", "Private", "Individual"}
			isValidType := false
			for _, validType := range validContactTypes {
				if *in.ContactType == validType {
					isValidType = true
					break
				}
			}
			if !isValidType {
				c.JSON(http.StatusBadRequest, gin.H{"error": "contact_type must be one of: Govt., Private, Individual"})
				return
			}
			contact.ContactType = *in.ContactType
		}
		if in.Department != nil {
			contact.Department = *in.Department
		}
		if in.Location != nil {
			contact.Location = *in.Location
		}
		if in.DesignationID != nil {
			contact.DesignationID = *in.DesignationID
		}
		if in.FirstName != nil {
			contact.FirstName = *in.FirstName
		}
		if in.LastName != nil {
			contact.LastName = *in.LastName
		}
		if in.Email != nil {
			if *in.Email == "" {
				contact.Email = nil
			} else {
				contact.Email = in.Email
			}
		}
		if in.Mobile != nil {
			contact.Mobile = *in.Mobile
		}
		if in.Password != nil {
			if *in.Password == "" {
				contact.PasswordHash = nil
			} else {
				if h, err := utils.HashPassword(*in.Password); err == nil {
					contact.PasswordHash = &h
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash error"})
					return
				}
			}
		}

		if err := db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log contact update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditContactUpdated,
			models.EntityTypeContact,
			&contact.ID,
			fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
			fmt.Sprintf("Contact updated: %s (%s)", derefString(contact.Email), contact.ContactType),
			oldContact,
			contact,
		)

		c.JSON(http.StatusOK, contact)
	}
}

func DeleteContact(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Get contact details before deletion for audit
		var contact models.Contact
		if err := db.First(&contact, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
			return
		}

		if err := db.Delete(&models.Contact{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log contact deletion
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditContactDeleted,
			models.EntityTypeContact,
			&contact.ID,
			fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
			fmt.Sprintf("Contact deleted: %s (%s)", derefString(contact.Email), contact.ContactType),
			contact,
			nil,
		)

		c.Status(http.StatusNoContent)
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
