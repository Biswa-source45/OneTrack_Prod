package services

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// EmailNotificationService handles sending email notifications
type EmailNotificationService struct {
	db               *gorm.DB
	smtpServer       string
	smtpPort         string
	emailUsername    string
	emailPassword    string
	notificationAddr string
}

// NewEmailNotificationService creates a new email notification service
func NewEmailNotificationService(db *gorm.DB, smtpServer, smtpPort, username, password, notificationAddr string) *EmailNotificationService {
	return &EmailNotificationService{
		db:               db,
		smtpServer:       smtpServer,
		smtpPort:         smtpPort,
		emailUsername:    username,
		emailPassword:    password,
		notificationAddr: notificationAddr,
	}
}

// SendTicketCreationEmail sends an email notification when a ticket is created
func (s *EmailNotificationService) SendTicketCreationEmail(ticket *models.Ticket) error {
	log.Printf("Sending notification about ticket %s to %s", ticket.TicketID, s.notificationAddr)

	// Skip if notification address is empty
	if s.notificationAddr == "" {
		return fmt.Errorf("notification email address not configured")
	}

	// Get contact and account info for the notification
	var contact models.Contact
	var account models.Account
	var product models.MasterProduct

	// Fetch the contact details
	if err := s.db.First(&contact, ticket.ContactID).Error; err != nil {
		return fmt.Errorf("failed to get contact info for notification: %w", err)
	}

	// Fetch the account details if available
	if ticket.AccountID != nil && *ticket.AccountID > 0 {
		if err := s.db.First(&account, ticket.AccountID).Error; err != nil {
			log.Printf("Warning: could not get account info for notification: %v", err)
		}
	}

	// Fetch the product details
	if err := s.db.First(&product, ticket.ProductID).Error; err != nil {
		log.Printf("Warning: could not get product info for notification: %v", err)
	}

	// Set up authentication information
	auth := smtp.PlainAuth("", s.emailUsername, s.emailPassword, s.smtpServer)

	// Compose email content
	subject := fmt.Sprintf("New Ticket Created: %s - %s", ticket.TicketID, ticket.Subject)

	// Build a simple HTML email body
	htmlBody := fmt.Sprintf(`
	<html>
	<body>
	<h2>New Ticket Created</h2>
	<p>A new support ticket has been created in the ticketing system.</p>
	<table border="0" cellpadding="5">
		<tr><td><strong>Ticket ID:</strong></td><td>%s</td></tr>
		<tr><td><strong>Subject:</strong></td><td>%s</td></tr>
		<tr><td><strong>Status:</strong></td><td>%s</td></tr>
		<tr><td><strong>Priority:</strong></td><td>%s</td></tr>
		<tr><td><strong>Product:</strong></td><td>%s</td></tr>
		<tr><td><strong>Account:</strong></td><td>%s</td></tr>
		<tr><td><strong>Contact:</strong></td><td>%s %s</td></tr>
		<tr><td><strong>Created:</strong></td><td>%s</td></tr>
	</table>
	
	<h3>Ticket Details:</h3>
	<p>%s</p>
	
	<p>Please login to the ticketing system to view and respond to this ticket.</p>
	</body>
	</html>`,
		ticket.TicketID,
		ticket.Subject,
		ticket.TicketStatus,
		ticket.Priority,
		product.ProductName,
		account.AccountName,
		contact.FirstName, contact.LastName,
		ticket.CreatedAt.Format("Jan 2, 2006 15:04:05"),
		strings.ReplaceAll(ticket.TicketDetails, "\n", "<br>"),
	)

	// Compose the full email message with headers
	headers := make(map[string]string)
	headers["From"] = s.emailUsername
	headers["To"] = s.notificationAddr
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build the message
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + htmlBody

	// Connect to the SMTP server
	smtpAddr := fmt.Sprintf("%s:%s", s.smtpServer, s.smtpPort)

	// Send the email
	err := smtp.SendMail(smtpAddr, auth, s.emailUsername, []string{s.notificationAddr}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send notification email: %w", err)
	}

	log.Printf("Successfully sent notification email about ticket %s to %s", ticket.TicketID, s.notificationAddr)
	return nil
}

// ListenForTicketNotifications adds a hook to the notification service
// This method subscribes to ticket creation events from the notification service
func (s *EmailNotificationService) ListenForTicketNotifications(notification *models.Notification) {
	// Only send emails for ticket creation notifications
	if notification.NotificationType == models.NotificationTicketCreatedByCustomer ||
		notification.NotificationType == models.NotificationTicketCreatedConfirmation {

		// Get ticket details from related ID
		if notification.RelatedID != nil && notification.RelatedType == "ticket" {
			var ticket models.Ticket
			if err := s.db.First(&ticket, *notification.RelatedID).Error; err == nil {
				// Send email notification
				if err := s.SendTicketCreationEmail(&ticket); err != nil {
					log.Printf("Failed to send ticket creation email: %v", err)
				}
			}
		}
	}
}
