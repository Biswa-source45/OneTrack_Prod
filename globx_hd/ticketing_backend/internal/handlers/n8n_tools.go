package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========================================
// AI AGENT TOOL ENDPOINTS
// ========================================
// These endpoints are designed to be called by AI agents (like OpenAI/Claude)
// through n8n's AI Agent node with Tools/Function Calling.
//
// Each endpoint returns simple, structured data that the AI can use
// to make decisions and resolve entities to their database IDs.
// ========================================

// ----------------------------------------
// CONTACT LOOKUP TOOLS
// ----------------------------------------

// ToolSearchContactByEmail finds a contact by exact email match
// This is the highest confidence lookup method
func ToolSearchContactByEmail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := strings.TrimSpace(strings.ToLower(c.Query("email")))
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "email parameter is required",
			})
			return
		}

		var contact models.Contact
		err := db.Preload("Account").Where("LOWER(email) = ?", email).First(&contact).Error

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No contact found with this email",
			})
			return
		}

		response := gin.H{
			"found":      true,
			"contact_id": contact.ID,
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"email":      contact.Email,
			"mobile":     contact.Mobile,
		}

		if contact.AccountID != nil {
			response["account_id"] = *contact.AccountID
			response["account_name"] = contact.Account.AccountName
			response["customer_code"] = contact.Account.CustomerCode
		} else {
			response["account_id"] = nil
			response["account_name"] = "Individual Contact"
			response["customer_code"] = contact.CustomerCode
		}

		c.JSON(http.StatusOK, response)
	}
}

// ToolSearchContactByPhone finds a contact by phone number
// Handles various phone formats
func ToolSearchContactByPhone(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		phone := c.Query("phone")
		if phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "phone parameter is required",
			})
			return
		}

		// Normalize phone number - remove all non-digits
		normalizedPhone := normalizePhoneNumber(phone)
		if len(normalizedPhone) < 10 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "Phone number too short (need at least 10 digits)",
			})
			return
		}

		// Search strategies:
		// 1. Exact match on normalized number
		// 2. Last 10 digits match (handles country code variations)
		var contact models.Contact
		var err error

		// Try exact match first
		err = db.Preload("Account").Where("mobile = ?", phone).First(&contact).Error

		// If not found, try with normalized number
		if err != nil {
			err = db.Preload("Account").Where("mobile = ?", normalizedPhone).First(&contact).Error
		}

		// If still not found, try last 10 digits
		if err != nil && len(normalizedPhone) >= 10 {
			last10 := normalizedPhone[len(normalizedPhone)-10:]
			err = db.Preload("Account").Where("mobile LIKE ?", "%"+last10).First(&contact).Error
		}

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No contact found with this phone number",
			})
			return
		}

		response := gin.H{
			"found":      true,
			"contact_id": contact.ID,
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"email":      contact.Email,
			"mobile":     contact.Mobile,
		}

		if contact.AccountID != nil {
			response["account_id"] = *contact.AccountID
			response["account_name"] = contact.Account.AccountName
			response["customer_code"] = contact.Account.CustomerCode
		} else {
			response["account_id"] = nil
			response["account_name"] = "Individual Contact"
			response["customer_code"] = contact.CustomerCode
		}

		c.JSON(http.StatusOK, response)
	}
}

// ToolSearchContactByName finds contacts by name (fuzzy search)
// Optionally filter by account_id for more precise results
func ToolSearchContactByName(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimSpace(c.Query("name"))
		accountID := c.Query("account_id")

		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "name parameter is required",
			})
			return
		}

		nameLower := strings.ToLower(name)
		parts := strings.Fields(nameLower)

		query := db.Preload("Account")

		// If account_id is provided, filter by it
		if accountID != "" {
			query = query.Where("account_id = ?", accountID)
		}

		var contacts []models.Contact

		if len(parts) >= 2 {
			// Try first + last name match
			firstName := parts[0]
			lastName := strings.Join(parts[1:], " ")
			query.Where(
				"(LOWER(first_name) LIKE ? AND LOWER(last_name) LIKE ?) OR "+
					"(LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?)",
				"%"+firstName+"%", "%"+lastName+"%",
				"%"+nameLower+"%", "%"+nameLower+"%",
			).Limit(5).Find(&contacts)
		} else {
			// Single word - search in both fields
			query.Where(
				"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
				"%"+nameLower+"%", "%"+nameLower+"%",
			).Limit(5).Find(&contacts)
		}

		if len(contacts) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No contacts found matching this name",
			})
			return
		}

		// Return all matches so AI can choose
		results := make([]gin.H, 0, len(contacts))
		for _, contact := range contacts {
			result := gin.H{
				"contact_id": contact.ID,
				"first_name": contact.FirstName,
				"last_name":  contact.LastName,
				"email":      contact.Email,
				"mobile":     contact.Mobile,
			}

			if contact.AccountID != nil {
				result["account_id"] = *contact.AccountID
				result["account_name"] = contact.Account.AccountName
			} else {
				result["account_id"] = nil
				result["account_name"] = "Individual Contact"
			}

			results = append(results, result)
		}

		c.JSON(http.StatusOK, gin.H{
			"found":    true,
			"count":    len(results),
			"contacts": results,
			"message":  "Multiple contacts found. Use additional information to identify the correct one.",
		})
	}
}

// ----------------------------------------
// ACCOUNT LOOKUP TOOLS
// ----------------------------------------

// ToolSearchAccountByName finds an account by name (fuzzy search)
func ToolSearchAccountByName(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := strings.TrimSpace(c.Query("name"))
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "name parameter is required",
			})
			return
		}

		nameLower := strings.ToLower(name)

		var accounts []models.Account

		// Search by name (fuzzy)
		db.Where("LOWER(account_name) LIKE ?", "%"+nameLower+"%").
			Limit(5).
			Find(&accounts)

		if len(accounts) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No accounts found matching this name",
			})
			return
		}

		results := make([]gin.H, 0, len(accounts))
		for _, account := range accounts {
			results = append(results, gin.H{
				"account_id":    account.ID,
				"account_name":  account.AccountName,
				"customer_code": account.CustomerCode,
				"account_owner": account.AccountOwner,
			})
		}

		if len(results) == 1 {
			c.JSON(http.StatusOK, gin.H{
				"found":         true,
				"account_id":    results[0]["account_id"],
				"account_name":  results[0]["account_name"],
				"customer_code": results[0]["customer_code"],
				"account_owner": results[0]["account_owner"],
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"found":    true,
			"count":    len(results),
			"accounts": results,
			"message":  "Multiple accounts found. Use the most relevant one.",
		})
	}
}

// ToolSearchAccountByDomain finds accounts by email domain
// Looks for contacts with matching email domain and returns their account
func ToolSearchAccountByDomain(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := strings.TrimSpace(strings.ToLower(c.Query("domain")))
		if domain == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "domain parameter is required",
			})
			return
		}

		// Remove @ if included
		domain = strings.TrimPrefix(domain, "@")

		// Find contacts with this email domain
		var contacts []models.Contact
		db.Preload("Account").
			Where("LOWER(email) LIKE ?", "%@"+domain).
			Limit(10).
			Find(&contacts)

		if len(contacts) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No contacts found with this email domain",
			})
			return
		}

		// Group by account
		accountMap := make(map[uint]gin.H)
		for _, contact := range contacts {
			if contact.AccountID != nil {
				if _, exists := accountMap[*contact.AccountID]; !exists {
					accountMap[*contact.AccountID] = gin.H{
						"account_id":    *contact.AccountID,
						"account_name":  contact.Account.AccountName,
						"customer_code": contact.Account.CustomerCode,
						"contact_count": 0,
					}
				}
				accountMap[*contact.AccountID]["contact_count"] = accountMap[*contact.AccountID]["contact_count"].(int) + 1
			}
		}

		if len(accountMap) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "Contacts found but they are Individual contacts without an account",
			})
			return
		}

		results := make([]gin.H, 0, len(accountMap))
		for _, acc := range accountMap {
			results = append(results, acc)
		}

		if len(results) == 1 {
			c.JSON(http.StatusOK, gin.H{
				"found":         true,
				"account_id":    results[0]["account_id"],
				"account_name":  results[0]["account_name"],
				"customer_code": results[0]["customer_code"],
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"found":    true,
			"count":    len(results),
			"accounts": results,
			"message":  "Multiple accounts found with contacts using this domain.",
		})
	}
}

// ToolGetAccountContacts returns all contacts for a given account
// Useful after finding an account to pick the right contact
func ToolGetAccountContacts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.Query("account_id")
		if accountID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"found": false,
				"error": "account_id parameter is required",
			})
			return
		}

		var contacts []models.Contact
		db.Where("account_id = ?", accountID).Find(&contacts)

		if len(contacts) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No contacts found for this account",
			})
			return
		}

		results := make([]gin.H, 0, len(contacts))
		for _, contact := range contacts {
			results = append(results, gin.H{
				"contact_id": contact.ID,
				"first_name": contact.FirstName,
				"last_name":  contact.LastName,
				"email":      contact.Email,
				"mobile":     contact.Mobile,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"found":    true,
			"count":    len(results),
			"contacts": results,
		})
	}
}

// ----------------------------------------
// PRODUCT LOOKUP TOOLS
// ----------------------------------------

// ToolSearchProduct finds products by keyword
func ToolSearchProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := strings.TrimSpace(c.Query("query"))

		var products []models.MasterProduct

		if query == "" {
			// Return all products if no query
			db.Order("product_name").Find(&products)
		} else {
			// Search by name
			queryLower := strings.ToLower(query)
			db.Where("LOWER(product_name) LIKE ?", "%"+queryLower+"%").
				Order("product_name").
				Find(&products)
		}

		if len(products) == 0 {
			// Return default product if nothing found
			var defaultProduct models.MasterProduct
			if err := db.First(&defaultProduct).Error; err == nil {
				c.JSON(http.StatusOK, gin.H{
					"found":        true,
					"product_id":   defaultProduct.ID,
					"product_name": defaultProduct.ProductName,
					"message":      "No matching product found. Using default product.",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"found":   false,
				"message": "No products found",
			})
			return
		}

		if len(products) == 1 {
			c.JSON(http.StatusOK, gin.H{
				"found":        true,
				"product_id":   products[0].ID,
				"product_name": products[0].ProductName,
			})
			return
		}

		results := make([]gin.H, 0, len(products))
		for _, product := range products {
			results = append(results, gin.H{
				"product_id":   product.ID,
				"product_name": product.ProductName,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"found":    true,
			"count":    len(results),
			"products": results,
			"message":  "Multiple products found. Choose the most relevant one.",
		})
	}
}

// ToolListAllProducts returns all products (for AI to show options)
func ToolListAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.MasterProduct
		db.Order("product_name").Find(&products)

		results := make([]gin.H, 0, len(products))
		for _, product := range products {
			results = append(results, gin.H{
				"product_id":   product.ID,
				"product_name": product.ProductName,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    len(results),
			"products": results,
		})
	}
}

// ----------------------------------------
// HELPER FUNCTIONS
// ----------------------------------------

// normalizePhoneNumber removes all non-digit characters from a phone number
func normalizePhoneNumber(phone string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(phone, "")
}
