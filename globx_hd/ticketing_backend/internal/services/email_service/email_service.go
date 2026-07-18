package email_service

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
)

// EmailService handles email processing for ticket creation
type EmailService struct {
	db                *gorm.DB
	imapServer        string
	imapPort          string
	smtpServer        string
	smtpPort          string
	emailUsername     string
	emailPassword     string
	allowedSenders    []string
	notificationAddr  string
	uploadDir         string
	activityService   *services.ActivityService
	emailNotification *services.EmailNotificationService
	senderUserID      uint // Store the sender's user ID for activity logging
}

// EmailTicketData represents structured data extracted from an email
type EmailTicketData struct {
	Account          string // Account name
	Contact          string // Contact name or email
	Product          string
	Priority         string
	Details          string
	Subject          string
	AdditionalFields map[string]string
}

// AttachmentData represents an email attachment
type AttachmentData struct {
	Filename    string
	ContentType string
	Data        []byte
}

// NewEmailService creates a new email service instance
func NewEmailService(db *gorm.DB, uploadDir string) (*EmailService, error) {
	// Try to load from .env.email file first, then fall back to .env
	err := godotenv.Load(".env.email")
	if err != nil {
		// Try loading from root project directory
		err = godotenv.Load("../../../.env.email")
	}
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Continue anyway, as env vars might be set directly in the environment
	}

	// Load configuration from environment variables
	imapServer := getEnvOrDefault("EMAIL_IMAP_SERVER", "imap.gmail.com")
	imapPort := getEnvOrDefault("EMAIL_IMAP_PORT", "993")
	smtpServer := getEnvOrDefault("EMAIL_SMTP_SERVER", "smtp.gmail.com") // Add SMTP server
	smtpPort := getEnvOrDefault("EMAIL_SMTP_PORT", "587")                // Common TLS port
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	allowedSendersStr := os.Getenv("EMAIL_ALLOWED_SENDERS")
	notificationAddr := os.Getenv("EMAIL_NOTIFICATION_ADDRESS")

	// Validate required configuration
	if username == "" || password == "" {
		return nil, fmt.Errorf("email username and password must be set in environment variables")
	}

	if allowedSendersStr == "" {
		return nil, fmt.Errorf("allowed senders must be set in environment variables")
	}

	if notificationAddr == "" {
		return nil, fmt.Errorf("notification email address must be set in environment variables")
	}

	// Split allowed senders string into slice
	allowedSenders := strings.Split(allowedSendersStr, ",")
	for i, sender := range allowedSenders {
		allowedSenders[i] = strings.TrimSpace(sender)
	}

	activityService := services.NewActivityService(db)

	// Create email notification service
	emailNotification := services.NewEmailNotificationService(
		db,
		smtpServer,
		smtpPort,
		username,
		password,
		notificationAddr,
	)

	return &EmailService{
		db:                db,
		imapServer:        imapServer,
		imapPort:          imapPort,
		smtpServer:        smtpServer,
		smtpPort:          smtpPort,
		emailUsername:     username,
		emailPassword:     password,
		allowedSenders:    allowedSenders,
		notificationAddr:  notificationAddr,
		uploadDir:         uploadDir,
		activityService:   activityService,
		emailNotification: emailNotification,
		senderUserID:      0, // Initialize to 0 (no user yet)
	}, nil
}

// ProcessEmails fetches and processes unread emails to create tickets
func (s *EmailService) ProcessEmails() error {
	log.Println("Starting email processing...")

	// Connect to the IMAP server
	imapAddr := fmt.Sprintf("%s:%s", s.imapServer, s.imapPort)
	c, err := client.DialTLS(imapAddr, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	log.Println("Connected to IMAP server")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(s.emailUsername, s.emailPassword); err != nil {
		return fmt.Errorf("failed to login to email account: %w", err)
	}
	log.Println("Logged in to email account")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("failed to select inbox: %w", err)
	}
	log.Printf("Mailbox status: %d messages, %d recent, %d unseen\n", mbox.Messages, mbox.Recent, mbox.Unseen)

	// Set search criteria to find unread messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	// Search for unread messages
	ids, err := c.Search(criteria)
	if err != nil {
		return fmt.Errorf("failed to search for emails: %w", err)
	}
	log.Printf("Found %d unread messages\n", len(ids))

	if len(ids) == 0 {
		log.Println("No unread messages found to process")
		return nil
	}

	// Process emails in batches of 10
	batchSize := 10
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batchIds := ids[i:end]
		seqset := new(imap.SeqSet)
		seqset.AddNum(batchIds...)

		// Get the whole message body
		section := &imap.BodySectionName{
			BodyPartName: imap.BodyPartName{},
			Peek:         true, // Don't mark emails as read yet
		}

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)

		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, section.FetchItem()}, messages)
		}()

		// Process each message
		for msg := range messages {
			err = s.processEmail(msg, section, c)
			if err != nil {
				log.Printf("Error processing email %d: %v", msg.SeqNum, err)
				// Continue processing other emails even if one fails
			}
		}

		if err := <-done; err != nil {
			log.Printf("Error during fetch: %v", err)
		}
	}

	log.Println("Email processing completed")
	return nil
}

// processEmail processes a single email message
func (s *EmailService) processEmail(msg *imap.Message, section *imap.BodySectionName, imapClient *client.Client) error {
	log.Printf("Processing email: Subject=%s, From=%v\n", msg.Envelope.Subject, msg.Envelope.From)

	// Check if sender is allowed
	senderEmail := msg.Envelope.From[0].MailboxName + "@" + msg.Envelope.From[0].HostName
	if !s.isSenderAllowed(senderEmail) {
		log.Printf("Email from unauthorized sender: %s - marking as read and skipping", senderEmail)
		// Mark as read since we don't want to process it again
		seqSet := new(imap.SeqSet)
		seqSet.AddNum(msg.SeqNum)
		err := imapClient.Store(seqSet, imap.AddFlags, []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Printf("Error marking unauthorized email as read: %v", err)
		}
		return fmt.Errorf("sender not in allowed list: %s", senderEmail)
	}

	// Find the user ID for this sender email (to use in activity logs)
	var user models.User
	err := s.db.Where("LOWER(email) = LOWER(?)", senderEmail).First(&user).Error
	if err != nil {
		log.Printf("Warning: No user found with email %s. Activity logging may fail.", senderEmail)
	}

	// Store the user ID in a context that will be accessible during ticket creation
	s.senderUserID = user.ID
	log.Printf("Found user ID %d for sender email %s", user.ID, senderEmail)

	// Get message body
	r := msg.GetBody(section)
	if r == nil {
		return fmt.Errorf("no message body found")
	}

	// Parse the message
	mr, err := mail.CreateReader(r)
	if err != nil {
		return fmt.Errorf("failed to parse email: %w", err)
	}

	// Extract email data
	var emailData EmailTicketData
	emailData.AdditionalFields = make(map[string]string)

	var attachments []AttachmentData

	// Process email headers
	header := mr.Header
	if date, err := header.Date(); err == nil {
		emailData.AdditionalFields["Date"] = date.Format(time.RFC1123)
	}
	if from, err := header.AddressList("From"); err == nil && len(from) > 0 {
		emailData.AdditionalFields["Sender"] = from[0].String()
	}
	if to, err := header.AddressList("To"); err == nil && len(to) > 0 {
		emailData.AdditionalFields["Recipient"] = to[0].String()
	}

	// Set subject from email subject
	emailData.Subject = msg.Envelope.Subject

	// Process each email part
	var bodyText string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error getting next part: %v", err)
			continue
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message body
			contentType, _, _ := h.ContentType()
			if contentType == "text/plain" {
				bodyBytes, _ := ioutil.ReadAll(p.Body)
				bodyText += string(bodyBytes) + "\n"
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			contentType, _, _ := h.ContentType()

			attData, err := ioutil.ReadAll(p.Body)
			if err != nil {
				log.Printf("Error reading attachment %s: %v", filename, err)
				continue
			}

			attachments = append(attachments, AttachmentData{
				Filename:    filename,
				ContentType: contentType,
				Data:        attData,
			})
			log.Printf("Found attachment: %s (%s, %d bytes)", filename, contentType, len(attData))
		}
	}

	// Parse the email body to extract ticket data
	err = s.parseEmailBody(bodyText, &emailData)
	if err != nil {
		return fmt.Errorf("failed to parse email body: %w", err)
	}

	// Create the ticket
	ticket, err := s.createTicketFromEmail(emailData, attachments)
	if err != nil {
		return fmt.Errorf("failed to create ticket from email: %w", err)
	}

	// Mark email as read after successful processing
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(msg.SeqNum)
	markErr := imapClient.Store(seqSet, imap.AddFlags, []interface{}{imap.SeenFlag}, nil)
	if markErr != nil {
		log.Printf("Warning: Failed to mark email as read: %v", markErr)
		// Continue anyway since we've successfully processed the email
	} else {
		log.Printf("Marked email as read: %d", msg.SeqNum)
	}

	// Send notification about the new ticket
	err = s.sendTicketCreationNotification(ticket)
	if err != nil {
		log.Printf("Error sending ticket creation notification: %v", err)
		// Don't return error here, as ticket is already created
	}

	log.Printf("Successfully created ticket %s from email\n", ticket.TicketID)
	return nil
}

// parseEmailBody extracts structured ticket data from email body text
func (s *EmailService) parseEmailBody(body string, data *EmailTicketData) error {
	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")

	// Define patterns to match required fields
	accountPattern := regexp.MustCompile(`(?i)^Account:\s*(.+)$`)
	contactPattern := regexp.MustCompile(`(?i)^Contact:\s*(.+)$`)
	detailsStartPattern := regexp.MustCompile(`(?i)^Details:\s*$`)
	// Also match details that are on the same line
	detailsInlinePattern := regexp.MustCompile(`(?i)^Details:\s+(.+)$`)

	inDetails := false
	var detailsLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// If we're in the details section, collect lines until we hit another field
		if inDetails {
			if accountPattern.MatchString(line) || contactPattern.MatchString(line) ||
				detailsStartPattern.MatchString(line) || detailsInlinePattern.MatchString(line) {
				inDetails = false
			} else {
				detailsLines = append(detailsLines, line)
				continue
			}
		}

		// Check for each field pattern
		if match := accountPattern.FindStringSubmatch(line); match != nil {
			// Store the account name in data.Account
			data.Account = strings.TrimSpace(match[1])
		} else if match := contactPattern.FindStringSubmatch(line); match != nil {
			// Store the contact info in data.Contact
			data.Contact = strings.TrimSpace(match[1])
		} else if match := detailsInlinePattern.FindStringSubmatch(line); match != nil {
			// Handle "Details: content" on same line
			detailsLines = append(detailsLines, match[1])
			inDetails = true
		} else if detailsStartPattern.MatchString(line) {
			inDetails = true
		} else if !inDetails {
			// If this line doesn't match any pattern and isn't part of details,
			// include it in details anyway
			detailsLines = append(detailsLines, line)
		}
	}

	// Join collected details lines
	data.Details = strings.TrimSpace(strings.Join(detailsLines, "\n"))

	// Set default values for non-required fields
	data.Product = "Other"   // Default product
	data.Priority = "Medium" // Default priority

	// Validate required fields
	if data.Account == "" {
		return fmt.Errorf("account name not found in email - required format: Account: [account name]")
	}
	if data.Details == "" {
		return fmt.Errorf("ticket details not found in email - required format: Details: [description]")
	}
	// Contact is optional but highly recommended
	if data.Contact == "" {
		log.Printf("Warning: No contact specified in email. Will attempt to find an appropriate contact from the account.")
	}

	return nil
}

// createTicketFromEmail creates a ticket from parsed email data
func (s *EmailService) createTicketFromEmail(data EmailTicketData, attachments []AttachmentData) (*models.Ticket, error) {
	// Start database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Step 1: Find the account by name
	var account models.Account
	accountName := strings.TrimSpace(data.Account)
	log.Printf("Looking for account with name: %s", accountName)

	// Try exact match on account name first
	err := tx.Where("LOWER(account_name) = LOWER(?)", accountName).First(&account).Error

	// If exact match fails, try LIKE search
	if err != nil {
		err = tx.Where("LOWER(account_name) LIKE LOWER(?)", "%"+accountName+"%").First(&account).Error
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("account not found: %s: %w", accountName, err)
		}
		log.Printf("Found account using LIKE search: %s (ID: %d)", account.AccountName, account.ID)
	} else {
		log.Printf("Found account with exact match: %s (ID: %d)", account.AccountName, account.ID)
	}

	// Step 2: Find the contact for this account
	var contact models.Contact
	var contactFound bool

	// If a contact was specified, try to find it
	if data.Contact != "" {
		contactInfo := strings.TrimSpace(data.Contact)
		log.Printf("Looking for contact with name/email: %s in account: %s", contactInfo, account.AccountName)

		// Try to find contact by email within this account
		err = tx.Where("account_id = ? AND LOWER(email) = LOWER(?)", account.ID, contactInfo).First(&contact).Error

		// If not found by email, try to find by first name + last name (if provided)
		if err != nil {
			names := strings.Split(contactInfo, " ")
			if len(names) > 1 {
				firstName := names[0]
				lastName := strings.Join(names[1:], " ")
				log.Printf("Trying to find contact by first_name='%s' and last_name='%s'", firstName, lastName)

				err = tx.Where("account_id = ? AND LOWER(first_name) = LOWER(?) AND LOWER(last_name) = LOWER(?)",
					account.ID, firstName, lastName).First(&contact).Error
			} else {
				// Try by first name only
				log.Printf("Trying to find contact by first_name='%s'", contactInfo)
				err = tx.Where("account_id = ? AND LOWER(first_name) = LOWER(?)", account.ID, contactInfo).First(&contact).Error
			}

			// If still not found, try LIKE searches
			if err != nil {
				log.Printf("Trying fuzzy search for contacts in account %d with name containing '%s'", account.ID, contactInfo)

				// Try LIKE on first_name
				err = tx.Where("account_id = ? AND LOWER(first_name) LIKE LOWER(?)", account.ID, "%"+contactInfo+"%").First(&contact).Error

				// Try LIKE on last_name if first_name search failed
				if err != nil {
					err = tx.Where("account_id = ? AND LOWER(last_name) LIKE LOWER(?)", account.ID, "%"+contactInfo+"%").First(&contact).Error
				}
			}
		}

		if err == nil {
			contactFound = true
			log.Printf("Found specified contact: %s %s (ID: %d) in account %s",
				contact.FirstName, contact.LastName, contact.ID, account.AccountName)
		} else {
			log.Printf("Could not find specified contact '%s' in account '%s'. Will try to find default contact.",
				contactInfo, account.AccountName)
		}
	}

	// If contact was not found or not specified, find the first contact from this account
	if !contactFound {
		err = tx.Where("account_id = ?", account.ID).First(&contact).Error
		if err != nil {
			log.Printf("No contacts found for account '%s'. Will attempt to use admin contact.", account.AccountName)

			// No contacts found for this account, try to find an admin contact (typically ID 1)
			err = tx.First(&contact).Error
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("no contacts found for account and no default contact exists: %w", err)
			}
			log.Printf("Using admin contact as fallback: %s %s (ID: %d)",
				contact.FirstName, contact.LastName, contact.ID)
		} else {
			log.Printf("Using first contact from account: %s %s (ID: %d)",
				contact.FirstName, contact.LastName, contact.ID)
		}
	}

	// Get default product (first one in database or create a generic one if none exists)
	var product models.MasterProduct
	err = tx.First(&product).Error
	if err != nil {
		// If no products exist, create a default one
		product = models.MasterProduct{
			ProductName:        "Other",
			ProductDescription: "Default product for email-created tickets",
			CreatedAt:          time.Now(),
		}

		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create default product: %w", err)
		}
	}

	// Create new ticket
	now := time.Now()

	// Generate ticket ID using the account's customer code
	dateStr := utils.FormatDateForTicketID(now)
	customerCode := account.CustomerCode

	// Get next sequence number
	seq, err := utils.GetNextTicketSequence(tx, account.ID, now)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get next ticket sequence: %w", err)
	}

	ticketID := utils.FormatTicketID(customerCode, dateStr, seq)

	// Create the ticket record - now including the contact ID
	ticket := &models.Ticket{
		TicketID:      ticketID,
		AccountID:     &account.ID,
		ContactID:     contact.ID, // Use the found contact
		ProductID:     product.ID,
		Subject:       data.Subject,
		TicketDetails: data.Details,
		TicketStatus:  "OPEN",
		Priority:      "Medium", // Default priority
		Channel:       "Mail",   // Email-created tickets use "Mail" channel
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save the ticket
	if err := tx.Create(ticket).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Process attachments
	for _, att := range attachments {
		// Create directory for attachments if it doesn't exist
		uploadPath := filepath.Join(s.uploadDir, "ticket_attachments", ticketID)
		if err := os.MkdirAll(uploadPath, 0755); err != nil {
			log.Printf("Error creating directory for attachments: %v", err)
			continue
		}

		// Generate unique filename
		storedFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), att.Filename)
		filePath := filepath.Join(uploadPath, storedFilename)

		// Save attachment to disk
		err := ioutil.WriteFile(filePath, att.Data, 0644)
		if err != nil {
			log.Printf("Error saving attachment %s: %v", att.Filename, err)
			continue
		}

		// Create attachment record - use the contact ID we found
		attachment := models.TicketAttachment{
			TicketID:         ticketID,
			OriginalFilename: att.Filename,
			StoredFilename:   storedFilename,
			FilePath:         filePath,
			FileSize:         len(att.Data),
			MimeType:         att.ContentType,
			UploadedBy:       contact.ID, // Use the found contact ID
			UploadedAt:       now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := tx.Create(&attachment).Error; err != nil {
			log.Printf("Error creating attachment record %s: %v", att.Filename, err)
			continue
		}
	}

	// Commit transaction first
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Now log activities AFTER the commit so the ticket ID exists in the database
	// Use the sender's user ID for activity logging
	var userID *uint
	// Only set userID if we found a valid user
	if s.senderUserID > 0 {
		userID = &s.senderUserID
	}

	activityErr := s.activityService.LogTicketCreation(ticket.ID, userID, ticket.TicketID)
	if activityErr != nil {
		log.Printf("Error logging ticket creation activity: %v", activityErr)
		// Continue anyway since ticket is created
	}

	// Log additional information about email-based creation
	description := fmt.Sprintf("Ticket created via email for account: %s by contact: %s %s",
		account.AccountName, contact.FirstName, contact.LastName)
	s.activityService.LogActivity(ticket.ID, userID, models.ActivityTicketCreated, description)

	return ticket, nil
}

// calculateSimilarity calculates a similarity score between two strings
func calculateSimilarity(s1, s2 string) int {
	// Exact match gets highest score
	if s1 == s2 {
		return 100
	}

	// Check if one string contains the other
	if strings.Contains(s1, s2) || strings.Contains(s2, s1) {
		return 80
	}

	// Split into words and check for word matches
	words1 := strings.Split(s1, " ")
	words2 := strings.Split(s2, " ")

	matchedWords := 0
	for _, w1 := range words1 {
		for _, w2 := range words2 {
			if w1 == w2 && len(w1) > 2 { // Only count matches for words longer than 2 chars
				matchedWords++
			}
		}
	}

	// Calculate score based on matched word ratio
	totalWords := len(words1) + len(words2)
	if totalWords == 0 {
		return 0
	}

	return (matchedWords * 100) / totalWords
}

// mapPriority standardizes priority string values
func mapPriority(priority string) string {
	priority = strings.ToLower(strings.TrimSpace(priority))

	switch priority {
	case "high", "urgent", "critical", "h":
		return "High"
	case "low", "minor", "l":
		return "Low"
	default:
		return "Medium" // Default to medium priority
	}
}

// sendTicketCreationNotification sends an email notification about new ticket
func (s *EmailService) sendTicketCreationNotification(ticket *models.Ticket) error {
	// Use the centralized email notification service instead
	return s.emailNotification.SendTicketCreationEmail(ticket)
}

// isSenderAllowed checks if the sender email is in the allowlist
func (s *EmailService) isSenderAllowed(email string) bool {
	email = strings.ToLower(email)
	for _, allowed := range s.allowedSenders {
		if strings.ToLower(allowed) == email {
			return true
		}
	}
	return false
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
