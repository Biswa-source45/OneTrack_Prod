// Remove duplicate 'package handlers' at the end of the file.
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Auth middleware to protect routes
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token format"})
			return
		}
		tokenString := authHeader[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		idFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid subject"})
			return
		}
		userType, ok := claims["type"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid type"})
			return
		}
		id := uint(idFloat)
		if userType == "user" {
			var user models.User
			if err := db.First(&user, id).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			c.Set("user", user)
		} else if userType == "contact" {
			var contact models.Contact
			if err := db.First(&contact, id).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "contact not found"})
				return
			}
			c.Set("contact", contact)
			c.Set("contact_id", contact.ID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user type"})
			return
		}
		c.Next()
	}
}

// Logout endpoint (stateless)
func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auditService := services.NewAuditService(db)

		// Log logout event
		if userVal, exists := c.Get("user"); exists {
			if user, ok := userVal.(models.User); ok {
				auditService.LogAuthentication(
					models.ActorTypeUser,
					&user.ID,
					fmt.Sprintf("%s %s", user.FirstName, user.LastName),
					user.Email,
					models.AuditLogout,
					true,
					c.ClientIP(),
					c.Request.UserAgent(),
					"",
				)
			}
		} else if contactVal, exists := c.Get("contact"); exists {
			if contact, ok := contactVal.(models.Contact); ok {
				auditService.LogAuthentication(
					models.ActorTypeContact,
					&contact.ID,
					fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
					derefString(contact.Email),
					models.AuditLogout,
					true,
					c.ClientIP(),
					c.Request.UserAgent(),
					"",
				)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

type ResetPasswordInput struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// Password reset for users
func ResetUserPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in ResetPasswordInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("username = ? OR email = ?", in.UsernameOrEmail, in.UsernameOrEmail).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		hash, err := utils.HashPassword(in.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}
		user.PasswordHash = hash
		user.FirstLogin = false
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log password reset
		auditService := services.NewAuditService(db)
		auditService.LogAuthentication(
			models.ActorTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			user.Email,
			models.AuditPasswordReset,
			true,
			c.ClientIP(),
			c.Request.UserAgent(),
			"",
		)

		c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
	}
}

// Password reset for contacts
func ResetContactPassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in ResetPasswordInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var contact models.Contact
		if err := db.Where("email = ?", in.UsernameOrEmail).First(&contact).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "contact not found"})
			return
		}
		hash, err := utils.HashPassword(in.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}
		contact.PasswordHash = &hash
		contact.FirstLogin = false
		if err := db.Save(&contact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Log password reset
		auditService := services.NewAuditService(db)
		auditService.LogAuthentication(
			models.ActorTypeContact,
			&contact.ID,
			fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
			derefString(contact.Email),
			models.AuditPasswordReset,
			true,
			c.ClientIP(),
			c.Request.UserAgent(),
			"",
		)

		c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
	}
}

var jwtSecret = []byte("your_secret_key") // Change to env/config in production

const (
	AccessTokenValidity  = 2 * time.Hour
	RefreshTokenValidity = 30 * 24 * time.Hour
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// User login handler
func UserLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in LoginInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		auditService := services.NewAuditService(db)

		// Preload Role from MasterRole
		if err := db.Preload("Role").Where("username = ?", in.Username).First(&user).Error; err != nil {
			// Log failed login attempt
			auditService.LogAuthentication(
				models.ActorTypeUser,
				nil,
				in.Username,
				in.Username,
				models.AuditUserLoginFailure,
				false,
				c.ClientIP(),
				c.Request.UserAgent(),
				"User not found",
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if !utils.CheckPasswordHash(in.Password, user.PasswordHash) {
			// Log failed login attempt
			auditService.LogAuthentication(
				models.ActorTypeUser,
				&user.ID,
				fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				user.Email,
				models.AuditUserLoginFailure,
				false,
				c.ClientIP(),
				c.Request.UserAgent(),
				"Invalid password",
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		accessToken, err := generateJWT(user.ID, "user", AccessTokenValidity, user.Role.RoleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		refreshToken, err := generateJWT(user.ID, "user", RefreshTokenValidity, user.Role.RoleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}

		// Log successful login
		auditService.LogAuthentication(
			models.ActorTypeUser,
			&user.ID,
			fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			user.Email,
			models.AuditUserLoginSuccess,
			true,
			c.ClientIP(),
			c.Request.UserAgent(),
			"",
		)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    int64(AccessTokenValidity.Seconds()),
			"role":          user.Role.RoleName,
			"first_login":   user.FirstLogin,
			"user":          user,
		})
	}
}

// Contact login handler
func ContactLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in LoginInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var contact models.Contact
		auditService := services.NewAuditService(db)

		if err := db.Where("email = ?", in.Username).First(&contact).Error; err != nil {
			// Log failed login attempt
			auditService.LogAuthentication(
				models.ActorTypeContact,
				nil,
				in.Username,
				in.Username,
				models.AuditContactLoginFailure,
				false,
				c.ClientIP(),
				c.Request.UserAgent(),
				"Contact not found",
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if contact.PasswordHash == nil || !utils.CheckPasswordHash(in.Password, *contact.PasswordHash) {
			// Log failed login attempt
			auditService.LogAuthentication(
				models.ActorTypeContact,
				&contact.ID,
				fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
				derefString(contact.Email),
				models.AuditContactLoginFailure,
				false,
				c.ClientIP(),
				c.Request.UserAgent(),
				"Invalid password",
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		accessToken, err := generateJWT(contact.ID, "contact", AccessTokenValidity, "contact")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}
		refreshToken, err := generateJWT(contact.ID, "contact", RefreshTokenValidity, "contact")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
			return
		}

		// Log successful login
		if err := auditService.LogAuthentication(
			models.ActorTypeContact,
			&contact.ID,
			fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
			derefString(contact.Email),
			models.AuditContactLoginSuccess,
			true,
			c.ClientIP(),
			c.Request.UserAgent(),
			"",
		); err != nil {
			fmt.Printf("WARNING: Failed to log contact login: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    int64(AccessTokenValidity.Seconds()),
			"first_login":   contact.FirstLogin,
			"contact":       contact,
		})
	}
}

// JWT generation helper
func generateJWT(id uint, userType string, validity time.Duration, roleName string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"type": userType,
		"role": roleName,
		"exp":  time.Now().Add(validity).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func RefreshToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in RefreshInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := jwt.Parse(in.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		idFloat, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid subject"})
			return
		}
		userType, ok := claims["type"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid type"})
			return
		}
		id := uint(idFloat)
		var accessToken string
		if userType == "user" {
			var user models.User
			if err := db.Preload("Role").First(&user, id).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				return
			}
			accessToken, err = generateJWT(user.ID, "user", AccessTokenValidity, user.Role.RoleName)
		} else if userType == "contact" {
			var contact models.Contact
			if err := db.First(&contact, id).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "contact not found"})
				return
			}
			accessToken, err = generateJWT(contact.ID, "contact", AccessTokenValidity, "contact")
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user type"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
			return
		}
		c.JSON(http.StatusOK, RefreshResponse{
			AccessToken: accessToken,
			ExpiresIn:   int64(AccessTokenValidity.Seconds()),
		})
	}
}
