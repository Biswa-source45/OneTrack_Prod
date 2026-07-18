package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

// EmailNotificationService handles sending email notifications
type EmailNotificationService struct {
	db            *gorm.DB
	smtpServer    string
	smtpPort      string
	emailUsername string
	emailPassword string
}

// NewEmailNotificationService creates a new email notification service
func NewEmailNotificationService(db *gorm.DB, smtpServer, smtpPort, username, password, _ string) *EmailNotificationService {
	return &EmailNotificationService{
		db:            db,
		smtpServer:    smtpServer,
		smtpPort:      smtpPort,
		emailUsername: username,
		emailPassword: password,
	}
}

// NewEmailNotificationServiceFromEnv creates service using env vars directly
func NewEmailNotificationServiceFromEnv(db *gorm.DB) *EmailNotificationService {
	return &EmailNotificationService{
		db:            db,
		smtpServer:    getenvOrDefault("EMAIL_SMTP_SERVER", "smtp.gmail.com"),
		smtpPort:      getenvOrDefault("EMAIL_SMTP_PORT", "587"),
		emailUsername: os.Getenv("EMAIL_USERNAME"),
		emailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// sendMail sends a single HTML email
func (s *EmailNotificationService) sendMail(to []string, subject, htmlBody string) error {
	if s.emailUsername == "" || s.emailPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	// Build recipients (To + self-CC for tracking)
	allTo := append([]string{}, to...)
	allTo = append(allTo, s.emailUsername) // self-CC always

	// Deduplicate
	seen := map[string]bool{}
	unique := allTo[:0]
	for _, addr := range allTo {
		a := strings.ToLower(strings.TrimSpace(addr))
		if a != "" && !seen[a] {
			seen[a] = true
			unique = append(unique, addr)
		}
	}

	toHeader := strings.Join(to, ", ")
	msg := "From: GlobX Support <" + s.emailUsername + ">\r\n" +
		"To: " + toHeader + "\r\n" +
		"Cc: " + s.emailUsername + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + htmlBody

	smtpAddr := fmt.Sprintf("%s:%s", s.smtpServer, s.smtpPort)
	auth := smtp.PlainAuth("", s.emailUsername, s.emailPassword, s.smtpServer)

	// Try STARTTLS (port 587)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:465", s.smtpServer), &tls.Config{ServerName: s.smtpServer})
	if err == nil {
		// SSL on port 465 worked
		defer conn.Close()
		client, err := smtp.NewClient(conn, s.smtpServer)
		if err != nil {
			return err
		}
		if err = client.Auth(auth); err != nil {
			return err
		}
		if err = client.Mail(s.emailUsername); err != nil {
			return err
		}
		for _, addr := range unique {
			if err = client.Rcpt(addr); err != nil {
				log.Printf("Warning: could not add recipient %s: %v", addr, err)
			}
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, msg)
		if err != nil {
			return err
		}
		w.Close()
		return client.Quit()
	}

	// Fall back to STARTTLS on port 587
	return smtp.SendMail(smtpAddr, auth, s.emailUsername, unique, []byte(msg))
}

// SendTicketCreationEmail sends ticket creation email to contact + self-CC
func (s *EmailNotificationService) SendTicketCreationEmail(ticket *models.Ticket) error {
	var contact models.Contact
	var account models.Account
	var product models.MasterProduct

	if err := s.db.First(&contact, ticket.ContactID).Error; err != nil {
		return fmt.Errorf("contact not found: %w", err)
	}
	if ticket.AccountID != nil && *ticket.AccountID > 0 {
		s.db.First(&account, ticket.AccountID)
	}
	s.db.First(&product, ticket.ProductID)

	var engineer models.User
	engineerName := "Unassigned"
	if ticket.AssignedEngineer != nil && *ticket.AssignedEngineer > 0 {
		if err := s.db.First(&engineer, ticket.AssignedEngineer).Error; err == nil {
			engineerName = engineer.FirstName + " " + engineer.LastName
		}
	}

	subject := fmt.Sprintf("[GlobX Support] Ticket Raised: %s – %s", ticket.TicketID, ticket.Subject)
	body := buildTicketCreationHTML(ticket, &contact, &account, &product, engineerName)

	recipients := []string{}
	if contact.Email != nil && *contact.Email != "" {
		recipients = append(recipients, *contact.Email)
	}

	if len(recipients) == 0 {
		// No contact email – still send to self for internal tracking
		recipients = append(recipients, s.emailUsername)
	}

	if err := s.sendMail(recipients, subject, body); err != nil {
		return fmt.Errorf("failed to send ticket creation email: %w", err)
	}
	log.Printf("[EMAIL] Ticket creation email sent for %s", ticket.TicketID)
	return nil
}

// SendTicketUpdateEmail sends ticket update email to contact + self-CC
func (s *EmailNotificationService) SendTicketUpdateEmail(ticket *models.Ticket, changedBy string, changes map[string]string) error {
	var contact models.Contact
	var account models.Account
	var product models.MasterProduct

	if err := s.db.First(&contact, ticket.ContactID).Error; err != nil {
		return fmt.Errorf("contact not found: %w", err)
	}
	if ticket.AccountID != nil && *ticket.AccountID > 0 {
		s.db.First(&account, ticket.AccountID)
	}
	s.db.First(&product, ticket.ProductID)

	subject := fmt.Sprintf("[GlobX Support] Ticket Updated: %s – %s", ticket.TicketID, ticket.Subject)
	body := buildTicketUpdateHTML(ticket, &contact, &account, &product, changedBy, changes)

	recipients := []string{}
	if contact.Email != nil && *contact.Email != "" {
		recipients = append(recipients, *contact.Email)
	}
	if len(recipients) == 0 {
		recipients = append(recipients, s.emailUsername)
	}

	if err := s.sendMail(recipients, subject, body); err != nil {
		return fmt.Errorf("failed to send ticket update email: %w", err)
	}
	log.Printf("[EMAIL] Ticket update email sent for %s", ticket.TicketID)
	return nil
}

// ListenForTicketNotifications legacy hook (no-op now)
func (s *EmailNotificationService) ListenForTicketNotifications(notification *models.Notification) {}

// ─── HTML Templates ───────────────────────────────────────────────────────────

func buildTicketCreationHTML(ticket *models.Ticket, contact *models.Contact, account *models.Account, product *models.MasterProduct, engineer string) string {
	contactEmail := ""
	if contact.Email != nil {
		contactEmail = *contact.Email
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Ticket Created</title></head>
<body style="margin:0;padding:0;background:#f4f6f8;font-family:Arial,Helvetica,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f8;padding:40px 0;">
<tr><td align="center">
<table width="620" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:6px;overflow:hidden;border:1px solid #dde2e8;">

  <!-- Header -->
  <tr><td style="background:#1a3c6e;padding:28px 36px;">
    <table width="100%%" cellpadding="0" cellspacing="0">
    <tr>
      <td><span style="color:#ffffff;font-size:20px;font-weight:700;letter-spacing:0.5px;">GlobX Enterprise Support</span></td>
      <td align="right"><span style="color:#93c5fd;font-size:12px;font-weight:600;text-transform:uppercase;letter-spacing:1px;">New Ticket</span></td>
    </tr>
    </table>
  </td></tr>

  <!-- Ticket ID Banner -->
  <tr><td style="background:#eef2ff;padding:16px 36px;border-bottom:1px solid #c7d2fe;">
    <span style="font-size:13px;color:#4338ca;font-weight:700;text-transform:uppercase;letter-spacing:0.8px;">Ticket Reference: </span>
    <span style="font-size:15px;color:#1e1b4b;font-weight:700;">%s</span>
  </td></tr>

  <!-- Body -->
  <tr><td style="padding:32px 36px;">
    <p style="margin:0 0 6px;font-size:14px;color:#64748b;text-transform:uppercase;letter-spacing:0.7px;font-weight:600;">Subject</p>
    <p style="margin:0 0 24px;font-size:17px;color:#1e293b;font-weight:700;">%s</p>

    <!-- Details Grid -->
    <table width="100%%" cellpadding="0" cellspacing="0" style="border:1px solid #e2e8f0;border-radius:4px;overflow:hidden;margin-bottom:24px;">
      <tr style="background:#f8fafc;">
        <td style="padding:10px 16px;font-size:12px;color:#64748b;font-weight:700;text-transform:uppercase;letter-spacing:0.6px;width:38%%%%;border-bottom:1px solid #e2e8f0;">Field</td>
        <td style="padding:10px 16px;font-size:12px;color:#64748b;font-weight:700;text-transform:uppercase;letter-spacing:0.6px;border-bottom:1px solid #e2e8f0;">Value</td>
      </tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Status</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Priority</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Product</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Account</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Contact</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s %s &lt;%s&gt;</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Assigned To</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;">Created</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;">%s</td></tr>
    </table>

    <!-- Ticket Details -->
    <p style="margin:0 0 8px;font-size:13px;color:#64748b;text-transform:uppercase;letter-spacing:0.7px;font-weight:600;">Ticket Description</p>
    <div style="background:#f8fafc;border-left:4px solid #1a3c6e;padding:14px 18px;border-radius:0 4px 4px 0;font-size:13px;color:#334155;line-height:1.7;">%s</div>
  </td></tr>

  <!-- Footer -->
  <tr><td style="background:#f8fafc;padding:20px 36px;border-top:1px solid #e2e8f0;text-align:center;">
    <p style="margin:0;font-size:11px;color:#94a3b8;">This is an automated notification from GlobX Enterprise Support System.<br>Please do not reply directly to this email.</p>
    <p style="margin:8px 0 0;font-size:11px;color:#94a3b8;">&copy; %d GlobX. All rights reserved.</p>
  </td></tr>

</table>
</td></tr>
</table>
</body>
</html>`,
		ticket.TicketID,
		ticket.Subject,
		ticket.TicketStatus,
		ticket.Priority,
		product.ProductName,
		account.AccountName,
		contact.FirstName, contact.LastName, contactEmail,
		engineer,
		ticket.CreatedAt.In(time.FixedZone("IST", 5*3600+30*60)).Format("02 Jan 2006, 03:04 PM MST"),
		strings.ReplaceAll(ticket.TicketDetails, "\n", "<br>"),
		time.Now().Year(),
	)
}

func buildTicketUpdateHTML(ticket *models.Ticket, contact *models.Contact, account *models.Account, product *models.MasterProduct, changedBy string, changes map[string]string) string {
	changesHTML := ""
	for field, val := range changes {
		changesHTML += fmt.Sprintf(`
      <tr>
        <td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">%s</td>
        <td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td>
      </tr>`, field, val)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Ticket Updated</title></head>
<body style="margin:0;padding:0;background:#f4f6f8;font-family:Arial,Helvetica,sans-serif;">
<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f8;padding:40px 0;">
<tr><td align="center">
<table width="620" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:6px;overflow:hidden;border:1px solid #dde2e8;">

  <!-- Header -->
  <tr><td style="background:#1a3c6e;padding:28px 36px;">
    <table width="100%%" cellpadding="0" cellspacing="0">
    <tr>
      <td><span style="color:#ffffff;font-size:20px;font-weight:700;letter-spacing:0.5px;">GlobX Enterprise Support</span></td>
      <td align="right"><span style="color:#fbbf24;font-size:12px;font-weight:600;text-transform:uppercase;letter-spacing:1px;">Ticket Updated</span></td>
    </tr>
    </table>
  </td></tr>

  <!-- Ticket ID Banner -->
  <tr><td style="background:#fffbeb;padding:16px 36px;border-bottom:1px solid #fde68a;">
    <span style="font-size:13px;color:#92400e;font-weight:700;text-transform:uppercase;letter-spacing:0.8px;">Ticket Reference: </span>
    <span style="font-size:15px;color:#78350f;font-weight:700;">%s</span>
  </td></tr>

  <!-- Body -->
  <tr><td style="padding:32px 36px;">
    <p style="margin:0 0 6px;font-size:14px;color:#64748b;text-transform:uppercase;letter-spacing:0.7px;font-weight:600;">Subject</p>
    <p style="margin:0 0 8px;font-size:17px;color:#1e293b;font-weight:700;">%s</p>
    <p style="margin:0 0 24px;font-size:13px;color:#64748b;">Updated by: <strong>%s</strong> on %s</p>

    <!-- Changes -->
    <p style="margin:0 0 8px;font-size:13px;color:#64748b;text-transform:uppercase;letter-spacing:0.7px;font-weight:600;">What Changed</p>
    <table width="100%%" cellpadding="0" cellspacing="0" style="border:1px solid #e2e8f0;border-radius:4px;overflow:hidden;margin-bottom:24px;">
      <tr style="background:#f8fafc;">
        <td style="padding:10px 16px;font-size:12px;color:#64748b;font-weight:700;text-transform:uppercase;letter-spacing:0.6px;width:38%%%%;border-bottom:1px solid #e2e8f0;">Field</td>
        <td style="padding:10px 16px;font-size:12px;color:#64748b;font-weight:700;text-transform:uppercase;letter-spacing:0.6px;border-bottom:1px solid #e2e8f0;">New Value</td>
      </tr>
      %s
    </table>

    <!-- Current State -->
    <p style="margin:0 0 8px;font-size:13px;color:#64748b;text-transform:uppercase;letter-spacing:0.7px;font-weight:600;">Current Ticket State</p>
    <table width="100%%" cellpadding="0" cellspacing="0" style="border:1px solid #e2e8f0;border-radius:4px;overflow:hidden;margin-bottom:24px;">
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;width:38%%%%;">Status</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;border-bottom:1px solid #f1f5f9;">Priority</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;border-bottom:1px solid #f1f5f9;">%s</td></tr>
      <tr><td style="padding:10px 16px;font-size:13px;color:#475569;">Product</td><td style="padding:10px 16px;font-size:13px;color:#1e293b;font-weight:600;">%s</td></tr>
    </table>
  </td></tr>

  <!-- Footer -->
  <tr><td style="background:#f8fafc;padding:20px 36px;border-top:1px solid #e2e8f0;text-align:center;">
    <p style="margin:0;font-size:11px;color:#94a3b8;">This is an automated notification from GlobX Enterprise Support System.<br>Please do not reply directly to this email.</p>
    <p style="margin:8px 0 0;font-size:11px;color:#94a3b8;">&copy; %d GlobX. All rights reserved.</p>
  </td></tr>

</table>
</td></tr>
</table>
</body>
</html>`,
		ticket.TicketID,
		ticket.Subject,
		changedBy,
		time.Now().In(time.FixedZone("IST", 5*3600+30*60)).Format("02 Jan 2006, 03:04 PM MST"),
		changesHTML,
		ticket.TicketStatus,
		ticket.Priority,
		product.ProductName,
		time.Now().Year(),
	)
}
