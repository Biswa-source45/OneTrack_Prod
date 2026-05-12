package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// N8nTicketRequest represents the JSON payload from n8n for ticket creation
type N8nTicketRequest struct {
	// Option 1: Direct IDs (if n8n did the lookup)
	AccountID uint `json:"account_id,omitempty"`
	ContactID uint `json:"contact_id,omitempty"`
	ProductID uint `json:"product_id,omitempty"`

	// Option 2: Names for lookup (if n8n sends raw data)
	AccountName  string `json:"account_name,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ProductName  string `json:"product_name,omitempty"`

	// Ticket data
	Subject       string `json:"subject" binding:"required"`
	Details       string `json:"details" binding:"required"`
	Priority      string `json:"priority,omitempty"`      // High, Medium, Low - defaults to Medium
	SenderEmail   string `json:"sender_email,omitempty"`  // Original email sender
	SenderUserID  uint   `json:"sender_user_id,omitempty"` // User ID for activity logging

	// Attachments (base64 encoded)
	Attachments []N8nAttachment `json:"attachments,omitempty"`

	// Raw email data (for debugging/audit)
	RawEmailSubject string `json:"raw_email_subject,omitempty"`
	RawEmailFrom    string `json:"raw_email_from,omitempty"`
	RawEmailDate    string `json:"raw_email_date,omitempty"`
}

// N8nAttachment represents a base64 encoded attachment from n8n
type N8nAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"` // Base64 encoded
}

// N8nLookupAccountResponse represents account lookup response
type N8nLookupAccountResponse struct {
	ID           uint   `json:"id"`
	AccountName  string `json:"account_name"`
	CustomerCode string `json:"customer_code"`
	AccountOwner string `json:"account_owner"`
	Score        int    `json:"score"` // Match score for fuzzy matching
}

// N8nLookupContactResponse represents contact lookup response
type N8nLookupContactResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	AccountID   *uint  `json:"account_id"`
	AccountName string `json:"account_name,omitempty"`
	Score       int    `json:"score"` // Match score for fuzzy matching
}

// N8nAPIKeyMiddleware validates the API key for n8n webhook requests
func N8nAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("N8N_API_KEY")
		if apiKey == "" {
			log.Println("Warning: N8N_API_KEY not configured. Using default key (insecure).")
			apiKey = "default-n8n-key-change-me"
		}

		// Check for API key in header or query param
		providedKey := c.GetHeader("X-N8N-API-Key")
		if providedKey == "" {
			providedKey = c.Query("api_key")
		}

		if providedKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid API key",
				"message": "Provide valid API key in X-N8N-API-Key header or api_key query parameter",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// N8nCreateTicketHandler handles ticket creation from n8n webhook
func N8nCreateTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req N8nTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		log.Printf("[n8n] Received ticket creation request: Subject=%s, AccountName=%s, ContactEmail=%s",
			req.Subject, req.AccountName, req.ContactEmail)

		// Start transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Step 1: Resolve Account
		var account models.Account
		var accountID uint

		if req.AccountID > 0 {
			// Direct ID provided
			if err := tx.First(&account, req.AccountID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Account not found",
					"details": fmt.Sprintf("No account with ID %d", req.AccountID),
				})
				return
			}
			accountID = account.ID
			log.Printf("[n8n] Using provided account ID: %d (%s)", account.ID, account.AccountName)
		} else if req.AccountName != "" {
			// Lookup by name
			account, err := lookupAccountByName(tx, req.AccountName)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Account not found",
					"details": fmt.Sprintf("Could not find account: %s", req.AccountName),
					"hint":    "Use /n8n/lookup/accounts endpoint to find the correct account",
				})
				return
			}
			accountID = account.ID
			log.Printf("[n8n] Found account by name: %d (%s)", account.ID, account.AccountName)
		} else {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Account required",
				"details": "Provide either account_id or account_name",
			})
			return
		}

		// Reload account to get customer_code
		if err := tx.First(&account, accountID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load account"})
			return
		}

		// Step 2: Resolve Contact
		var contact models.Contact
		var contactID uint

		if req.ContactID > 0 {
			// Direct ID provided
			if err := tx.First(&contact, req.ContactID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Contact not found",
					"details": fmt.Sprintf("No contact with ID %d", req.ContactID),
				})
				return
			}
			contactID = contact.ID
			log.Printf("[n8n] Using provided contact ID: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
		} else if req.ContactEmail != "" {
			// Lookup by email
			contact, err := lookupContactByEmail(tx, req.ContactEmail, accountID)
			if err != nil {
				// Try to find any contact in the account
				contact, err = lookupFirstContactInAccount(tx, accountID)
				if err != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Contact not found",
						"details": fmt.Sprintf("Could not find contact with email: %s", req.ContactEmail),
						"hint":    "Use /n8n/lookup/contacts endpoint to find the correct contact",
					})
					return
				}
				log.Printf("[n8n] Using fallback contact from account: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
			} else {
				log.Printf("[n8n] Found contact by email: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
			}
			contactID = contact.ID
		} else if req.ContactName != "" {
			// Lookup by name
			contact, err := lookupContactByName(tx, req.ContactName, accountID)
			if err != nil {
				// Try to find any contact in the account
				contact, err = lookupFirstContactInAccount(tx, accountID)
				if err != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Contact not found",
						"details": fmt.Sprintf("Could not find contact: %s", req.ContactName),
					})
					return
				}
				log.Printf("[n8n] Using fallback contact from account: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
			} else {
				log.Printf("[n8n] Found contact by name: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
			}
			contactID = contact.ID
		} else {
			// No contact specified, use first contact from account
			contact, err := lookupFirstContactInAccount(tx, accountID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "No contact found",
					"details": "Account has no contacts and no contact was specified",
				})
				return
			}
			contactID = contact.ID
			log.Printf("[n8n] Using first contact from account: %d (%s %s)", contact.ID, contact.FirstName, contact.LastName)
		}

		// Step 3: Resolve Product
		var product models.MasterProduct
		var productID uint

		if req.ProductID > 0 {
			if err := tx.First(&product, req.ProductID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Product not found",
					"details": fmt.Sprintf("No product with ID %d", req.ProductID),
				})
				return
			}
			productID = product.ID
		} else if req.ProductName != "" {
			// Lookup by name
			if err := tx.Where("LOWER(product_name) LIKE LOWER(?)", "%"+req.ProductName+"%").First(&product).Error; err != nil {
				// Use default product
				if err := tx.First(&product).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No products configured"})
					return
				}
			}
			productID = product.ID
		} else {
			// Use default product (first one)
			if err := tx.First(&product).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No products configured"})
				return
			}
			productID = product.ID
		}
		log.Printf("[n8n] Using product: %d (%s)", productID, product.ProductName)

		// Step 4: Generate Ticket ID
		now := time.Now()
		dateStr := utils.FormatDateForTicketID(now)
		seq, err := utils.GetNextTicketSequence(tx, accountID, now)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ticket sequence"})
			return
		}
		ticketID := utils.FormatTicketID(account.CustomerCode, dateStr, seq)

		// Step 5: Normalize priority
		priority := normalizePriority(req.Priority)

		// Step 6: Create Ticket
		ticket := &models.Ticket{
			TicketID:      ticketID,
			AccountID:     &accountID,
			ContactID:     contactID,
			ProductID:     productID,
			Subject:       req.Subject,
			TicketDetails: req.Details,
			TicketStatus:  "OPEN",
			Priority:      priority,
			Channel:       "Mail",
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := tx.Create(ticket).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create ticket",
				"details": err.Error(),
			})
			return
		}
		log.Printf("[n8n] Created ticket: %s", ticketID)

		// Step 7: Process Attachments
		uploadDir := os.Getenv("UPLOAD_DIR")
		if uploadDir == "" {
			uploadDir = "./uploads"
		}

		for _, att := range req.Attachments {
			if att.Filename == "" || att.Data == "" {
				continue
			}

			// Decode base64 data
			data, err := base64.StdEncoding.DecodeString(att.Data)
			if err != nil {
				log.Printf("[n8n] Warning: Failed to decode attachment %s: %v", att.Filename, err)
				continue
			}

			// Create directory
			uploadPath := filepath.Join(uploadDir, "ticket_attachments", ticketID)
			if err := os.MkdirAll(uploadPath, 0755); err != nil {
				log.Printf("[n8n] Warning: Failed to create attachment directory: %v", err)
				continue
			}

			// Save file
			storedFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), att.Filename)
			filePath := filepath.Join(uploadPath, storedFilename)

			if err := os.WriteFile(filePath, data, 0644); err != nil {
				log.Printf("[n8n] Warning: Failed to save attachment %s: %v", att.Filename, err)
				continue
			}

			// Create attachment record
			attachment := models.TicketAttachment{
				TicketID:         ticketID,
				OriginalFilename: att.Filename,
				StoredFilename:   storedFilename,
				FilePath:         filePath,
				FileSize:         len(data),
				MimeType:         att.ContentType,
				UploadedBy:       contactID,
				UploadedAt:       now,
				CreatedAt:        now,
				UpdatedAt:        now,
			}

			if err := tx.Create(&attachment).Error; err != nil {
				log.Printf("[n8n] Warning: Failed to create attachment record: %v", err)
				continue
			}
			log.Printf("[n8n] Saved attachment: %s (%d bytes)", att.Filename, len(data))
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		// Step 8: Log Activity (after commit)
		activityService := services.NewActivityService(db)
		var userID *uint
		if req.SenderUserID > 0 {
			userID = &req.SenderUserID
		}
		activityService.LogTicketCreation(ticket.ID, userID, ticket.TicketID)

		// Log additional context
		description := fmt.Sprintf("Ticket created via n8n webhook for account: %s, contact: %s %s",
			account.AccountName, contact.FirstName, contact.LastName)
		if req.RawEmailFrom != "" {
			description += fmt.Sprintf(" (from: %s)", req.RawEmailFrom)
		}
		activityService.LogActivity(ticket.ID, userID, models.ActivityTicketCreated, description)

		// Step 9: Send notification (optional)
		go sendN8nTicketNotification(db, ticket)

		log.Printf("[n8n] Successfully created ticket %s", ticketID)

		c.JSON(http.StatusCreated, gin.H{
			"success":   true,
			"ticket_id": ticketID,
			"ticket": gin.H{
				"id":         ticket.ID,
				"ticket_id":  ticket.TicketID,
				"subject":    ticket.Subject,
				"status":     ticket.TicketStatus,
				"priority":   ticket.Priority,
				"account":    account.AccountName,
				"contact":    fmt.Sprintf("%s %s", contact.FirstName, contact.LastName),
				"product":    product.ProductName,
				"created_at": ticket.CreatedAt,
			},
		})
	}
}

// N8nLookupAccountsHandler searches for accounts by name
func N8nLookupAccountsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		if search == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "search parameter required"})
			return
		}

		var accounts []models.Account
		searchLower := strings.ToLower(search)

		// Find accounts with fuzzy matching
		err := db.Where("LOWER(account_name) LIKE ?", "%"+searchLower+"%").
			Order("account_name").
			Limit(10).
			Find(&accounts).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Calculate match scores
		results := make([]N8nLookupAccountResponse, 0, len(accounts))
		for _, acc := range accounts {
			score := calculateMatchScore(strings.ToLower(acc.AccountName), searchLower)
			results = append(results, N8nLookupAccountResponse{
				ID:           acc.ID,
				AccountName:  acc.AccountName,
				CustomerCode: acc.CustomerCode,
				AccountOwner: acc.AccountOwner,
				Score:        score,
			})
		}

		// Sort by score (highest first) - simple bubble sort for small arrays
		for i := 0; i < len(results)-1; i++ {
			for j := i + 1; j < len(results); j++ {
				if results[j].Score > results[i].Score {
					results[i], results[j] = results[j], results[i]
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    len(results),
			"accounts": results,
		})
	}
}

// N8nLookupContactsHandler searches for contacts by name or email
func N8nLookupContactsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		accountIDStr := c.Query("account_id")

		if search == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "search parameter required"})
			return
		}

		var contacts []models.Contact
		searchLower := strings.ToLower(search)

		query := db.Preload("Account").Where(
			"LOWER(email) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			"%"+searchLower+"%", "%"+searchLower+"%", "%"+searchLower+"%",
		)

		// Filter by account if provided
		if accountIDStr != "" {
			query = query.Where("account_id = ?", accountIDStr)
		}

		err := query.Order("first_name, last_name").
			Limit(10).
			Find(&contacts).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Calculate match scores and build response
		results := make([]N8nLookupContactResponse, 0, len(contacts))
		for _, cont := range contacts {
			fullName := strings.ToLower(cont.FirstName + " " + cont.LastName)
			emailLower := strings.ToLower(cont.Email)

			// Calculate best score from name or email
			nameScore := calculateMatchScore(fullName, searchLower)
			emailScore := calculateMatchScore(emailLower, searchLower)
			score := nameScore
			if emailScore > score {
				score = emailScore
			}

			resp := N8nLookupContactResponse{
				ID:        cont.ID,
				FirstName: cont.FirstName,
				LastName:  cont.LastName,
				Email:     cont.Email,
				Mobile:    cont.Mobile,
				AccountID: cont.AccountID,
				Score:     score,
			}

			if cont.Account.ID > 0 {
				resp.AccountName = cont.Account.AccountName
			}

			results = append(results, resp)
		}

		// Sort by score (highest first)
		for i := 0; i < len(results)-1; i++ {
			for j := i + 1; j < len(results); j++ {
				if results[j].Score > results[i].Score {
					results[i], results[j] = results[j], results[i]
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    len(results),
			"contacts": results,
		})
	}
}

// N8nLookupProductsHandler searches for products by name
func N8nLookupProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		var products []models.MasterProduct

		if search != "" {
			searchLower := strings.ToLower(search)
			db.Where("LOWER(product_name) LIKE ?", "%"+searchLower+"%").
				Order("product_name").
				Limit(10).
				Find(&products)
		} else {
			// Return all products if no search
			db.Order("product_name").Find(&products)
		}

		c.JSON(http.StatusOK, gin.H{
			"count":    len(products),
			"products": products,
		})
	}
}

// N8nHealthCheckHandler returns the status of the n8n webhook endpoint
func N8nHealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "n8n-webhook",
			"time":    time.Now().Format(time.RFC3339),
		})
	}
}

// Helper functions

func lookupAccountByName(db *gorm.DB, name string) (models.Account, error) {
	var account models.Account
	nameLower := strings.ToLower(strings.TrimSpace(name))

	// Try exact match first
	err := db.Where("LOWER(account_name) = ?", nameLower).First(&account).Error
	if err == nil {
		return account, nil
	}

	// Try LIKE match
	err = db.Where("LOWER(account_name) LIKE ?", "%"+nameLower+"%").First(&account).Error
	return account, err
}

func lookupContactByEmail(db *gorm.DB, email string, accountID uint) (models.Contact, error) {
	var contact models.Contact
	emailLower := strings.ToLower(strings.TrimSpace(email))

	query := db.Where("LOWER(email) = ?", emailLower)
	if accountID > 0 {
		query = query.Where("account_id = ?", accountID)
	}

	err := query.First(&contact).Error
	return contact, err
}

func lookupContactByName(db *gorm.DB, name string, accountID uint) (models.Contact, error) {
	var contact models.Contact
	nameLower := strings.ToLower(strings.TrimSpace(name))

	// Try to split into first/last name
	parts := strings.Split(nameLower, " ")

	query := db.Where("account_id = ?", accountID)

	if len(parts) >= 2 {
		firstName := parts[0]
		lastName := strings.Join(parts[1:], " ")
		query = query.Where("LOWER(first_name) = ? AND LOWER(last_name) = ?", firstName, lastName)
	} else {
		// Search by first name only or LIKE
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", "%"+nameLower+"%", "%"+nameLower+"%")
	}

	err := query.First(&contact).Error
	return contact, err
}

func lookupFirstContactInAccount(db *gorm.DB, accountID uint) (models.Contact, error) {
	var contact models.Contact
	err := db.Where("account_id = ?", accountID).First(&contact).Error
	return contact, err
}

func normalizePriority(priority string) string {
	switch strings.ToLower(strings.TrimSpace(priority)) {
	case "high", "urgent", "critical", "h", "1":
		return "High"
	case "low", "minor", "l", "3":
		return "Low"
	default:
		return "Medium"
	}
}

func calculateMatchScore(text, search string) int {
	// Exact match
	if text == search {
		return 100
	}

	// Contains exact search
	if strings.Contains(text, search) {
		return 80
	}

	// Starts with search
	if strings.HasPrefix(text, search) {
		return 70
	}

	// Word match
	textWords := strings.Fields(text)
	searchWords := strings.Fields(search)
	matches := 0
	for _, sw := range searchWords {
		for _, tw := range textWords {
			if strings.Contains(tw, sw) || strings.Contains(sw, tw) {
				matches++
				break
			}
		}
	}

	if len(searchWords) > 0 {
		return (matches * 50) / len(searchWords)
	}

	return 0
}

func sendN8nTicketNotification(db *gorm.DB, ticket *models.Ticket) {
	// Create email notification service and send notification
	smtpServer := os.Getenv("EMAIL_SMTP_SERVER")
	smtpPort := os.Getenv("EMAIL_SMTP_PORT")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	notificationAddr := os.Getenv("EMAIL_NOTIFICATION_ADDRESS")

	if smtpServer == "" || username == "" || password == "" || notificationAddr == "" {
		log.Println("[n8n] Email notification skipped - not configured")
		return
	}

	emailService := services.NewEmailNotificationService(db, smtpServer, smtpPort, username, password, notificationAddr)
	if err := emailService.SendTicketCreationEmail(ticket); err != nil {
		log.Printf("[n8n] Failed to send notification email: %v", err)
	} else {
		log.Printf("[n8n] Notification email sent for ticket %s", ticket.TicketID)
	}
}
