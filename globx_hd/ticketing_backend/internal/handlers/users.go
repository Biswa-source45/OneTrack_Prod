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

type CreateUserInput struct {
	EmployeeID    string `json:"employee_id" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required,min=6"`
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
	DesignationID uint   `json:"designation_id" binding:"required"`
	RoleID        uint   `json:"role_id" binding:"required"`
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in CreateUserInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hash, err := utils.HashPassword(in.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}
		user := models.User{
			EmployeeID:    in.EmployeeID,
			Username:      in.Username,
			PasswordHash:  hash,
			FirstName:     in.FirstName,
			LastName:      in.LastName,
			Email:         in.Email,
			Phone:         in.Phone,
			DesignationID: in.DesignationID,
			RoleID:        in.RoleID,
			FirstLogin:    true,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log user creation
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditUserCreated,
			models.EntityTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			fmt.Sprintf("User created: %s (%s)", user.Username, user.Email),
			nil,
			user,
		)

		c.JSON(http.StatusCreated, user)
	}
}

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Preload("Designation").Preload("Role").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var u models.User
		if err := db.Preload("Designation").Preload("Role").First(&u, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	}
}

type UpdateUserInput struct {
	EmployeeID    *string `json:"employee_id"`
	Username      *string `json:"username"`
	Password      *string `json:"password"`
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
	DesignationID *uint   `json:"designation_id"`
	RoleID        *uint   `json:"role_id"`
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Store old values for audit
		oldUser := user

		var in UpdateUserInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if in.EmployeeID != nil {
			user.EmployeeID = *in.EmployeeID
		}
		if in.Username != nil {
			user.Username = *in.Username
		}
		if in.Password != nil {
			if h, err := utils.HashPassword(*in.Password); err == nil {
				user.PasswordHash = h
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash error"})
				return
			}
		}
		if in.FirstName != nil {
			user.FirstName = *in.FirstName
		}
		if in.LastName != nil {
			user.LastName = *in.LastName
		}
		if in.Email != nil {
			user.Email = *in.Email
		}
		if in.Phone != nil {
			user.Phone = *in.Phone
		}
		if in.DesignationID != nil {
			user.DesignationID = *in.DesignationID
		}
		if in.RoleID != nil {
			user.RoleID = *in.RoleID
		}

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log user update
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditUserUpdated,
			models.EntityTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			fmt.Sprintf("User updated: %s (%s)", user.Username, user.Email),
			oldUser,
			user,
		)

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Get user details before deletion for audit
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err := db.Delete(&models.User{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log user deletion
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditUserDeleted,
			models.EntityTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			fmt.Sprintf("User deleted: %s (%s)", user.Username, user.Email),
			user,
			nil,
		)

		c.Status(http.StatusNoContent)
	}
}

type ResetManagedUserPasswordInput struct {
	Password string `json:"password" binding:"required,min=6"`
}

func ResetManagedUserPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		oldUser := user

		var in ResetManagedUserPasswordInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := utils.HashPassword(in.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}

		if err := db.Model(&user).Update("password_hash", hash).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log administrative password reset
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(
			c,
			models.AuditUserUpdated,
			models.EntityTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			fmt.Sprintf("Administrative password reset for user: %s", user.Username),
			oldUser,
			user,
		)

		c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
	}
}

