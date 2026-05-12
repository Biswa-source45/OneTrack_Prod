package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadTicketAttachments handles file uploads for tickets
func UploadTicketAttachments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("🚨 DEBUG: UploadTicketAttachments ENDPOINT HIT! Method: %s, URL: %s\n", c.Request.Method, c.Request.URL.Path)
		ticketID := c.Param("id")
		fmt.Printf("🔥 DEBUG: UploadTicketAttachments called with ticketID: %s\n", ticketID)

		// Get user context - support both contact and user authentication
		var uploaderID uint
		var isContact bool

		// Try contact authentication first (customer side)
		if contactIDVal, exists := c.Get("contact_id"); exists {
			if contactID, ok := contactIDVal.(uint); ok {
				uploaderID = contactID
				isContact = true
				fmt.Printf("🔥 DEBUG: Contact authentication - Contact ID: %d\n", contactID)
			}
		}

		// Try user authentication (manager side)
		if uploaderID == 0 {
			if userVal, exists := c.Get("user"); exists {
				if user, ok := userVal.(models.User); ok {
					uploaderID = user.ID
					isContact = false
					fmt.Printf("🔥 DEBUG: User authentication - User ID: %d\n", user.ID)
				}
			}
		}

		if uploaderID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required - no valid contact_id or user_id found"})
			return
		}

		// Verify ticket exists
		var ticket models.Ticket
		var err error
		if isContact {
			// For contacts, verify they own the ticket
			err = db.Where("ticket_id = ? AND contact_id = ?", ticketID, uploaderID).First(&ticket).Error
		} else {
			// For users/managers, just verify ticket exists (they can upload to any ticket)
			err = db.Where("ticket_id = ?", ticketID).First(&ticket).Error
		}

		if err != nil {
			fmt.Printf("❌ DEBUG: Ticket not found - ticketID: %s, uploaderID: %d, isContact: %v, error: %v\n", ticketID, uploaderID, isContact, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found or access denied"})
			return
		}
		fmt.Printf("✅ DEBUG: Ticket found - ID: %s, Contact: %d\n", ticket.TicketID, ticket.ContactID)

		// Parse multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
			return
		}

		files := form.File["files"]
		fmt.Printf("🔥 DEBUG: Found %d files in form\n", len(files))
		if len(files) == 0 {
			fmt.Printf("❌ DEBUG: No files found in multipart form\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
			return
		}

		fileService := services.NewFileService(db)
		var savedAttachments []models.TicketAttachment
		var errors []string

		// Process each file
		for _, fileHeader := range files {
			fmt.Printf("🔥 DEBUG: Processing file: %s, Size: %d bytes\n", fileHeader.Filename, fileHeader.Size)

			// Save file
			attachment, err := fileService.SaveFile(fileHeader, ticketID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to save %s: %v", fileHeader.Filename, err))
				continue
			}

			// Set uploaded_by - always use the ticket's contact_id for foreign key consistency
			attachment.UploadedBy = ticket.ContactID

			// Save to database
			if err := db.Create(attachment).Error; err != nil {
				errors = append(errors, fmt.Sprintf("Failed to save %s to database: %v", fileHeader.Filename, err))
				continue
			}

			// Log attachment upload
			auditService := services.NewAuditService(db)
			auditService.LogCRUD(c, models.AuditAttachmentUploaded, models.EntityTypeAttachment, &attachment.ID, attachment.OriginalFilename, fmt.Sprintf("Attachment uploaded: %s (Ticket: %s)", attachment.OriginalFilename, ticketID), nil, attachment)

			savedAttachments = append(savedAttachments, *attachment)
			fmt.Printf("✅ Successfully saved file: %s as %s\n", fileHeader.Filename, attachment.StoredFilename)
		}

		// Prepare response
		response := gin.H{
			"message":     fmt.Sprintf("Processed %d files", len(files)),
			"saved_count": len(savedAttachments),
			"attachments": savedAttachments,
		}

		if len(errors) > 0 {
			response["errors"] = errors
			response["error_count"] = len(errors)
		}

		statusCode := http.StatusCreated
		if len(savedAttachments) == 0 {
			statusCode = http.StatusBadRequest
		} else if len(errors) > 0 {
			statusCode = http.StatusPartialContent
		}

		c.JSON(statusCode, response)
	}
}

// GetTicketAttachments retrieves all attachments for a ticket
func GetTicketAttachments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID := c.Param("id")

		// Get user context - support both contact and user authentication
		var userID uint
		var isContact bool

		// Try contact authentication first (customer side)
		if contactIDVal, exists := c.Get("contact_id"); exists {
			if contactID, ok := contactIDVal.(uint); ok {
				userID = contactID
				isContact = true
			}
		}

		// Try user authentication (manager side)
		if userID == 0 {
			if userVal, exists := c.Get("user"); exists {
				if user, ok := userVal.(models.User); ok {
					userID = user.ID
					isContact = false
				}
			}
		}

		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required - no valid contact_id or user_id found"})
			return
		}

		// Verify ticket exists
		var ticket models.Ticket
		var err error
		if isContact {
			// For contacts, verify they own the ticket
			err = db.Where("ticket_id = ? AND contact_id = ?", ticketID, userID).First(&ticket).Error
		} else {
			// For users/managers, just verify ticket exists (they can access any ticket)
			err = db.Where("ticket_id = ?", ticketID).First(&ticket).Error
		}

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found or access denied"})
			return
		}

		// Get attachments
		fileService := services.NewFileService(db)
		attachments, err := fileService.GetTicketAttachments(ticketID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attachments"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ticket_id":   ticketID,
			"attachments": attachments,
		})
	}
}

// DownloadAttachment serves the attachment file for download
func DownloadAttachment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attachmentIDStr := c.Param("id")
		attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment ID"})
			return
		}

		// Get attachment first
		var attachment models.TicketAttachment
		if err := db.First(&attachment, uint(attachmentID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "attachment not found"})
			return
		}

		// Check authentication - support both contact and user (manager)
		var hasAccess bool

		// Try contact authentication first (customer side)
		if contactIDVal, exists := c.Get("contact_id"); exists {
			if contactID, ok := contactIDVal.(uint); ok {
				// Verify ticket belongs to this contact
				var ticket models.Ticket
				if err := db.Where("ticket_id = ? AND contact_id = ?", attachment.TicketID, contactID).First(&ticket).Error; err == nil {
					hasAccess = true
				}
			}
		}

		// Try user authentication (manager side)
		if !hasAccess {
			if userVal, exists := c.Get("user"); exists {
				if _, ok := userVal.(models.User); ok {
					// Managers can access any attachment (they have admin privileges)
					hasAccess = true
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		// Log attachment download
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditAttachmentDownloaded, models.EntityTypeAttachment, &attachment.ID, attachment.OriginalFilename, fmt.Sprintf("Attachment downloaded: %s", attachment.OriginalFilename), nil, nil)

		// Serve the file
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachment.OriginalFilename))
		c.Header("Content-Type", attachment.MimeType)
		c.File(attachment.FilePath)
	}
}

// DeleteAttachment removes an attachment
func DeleteAttachment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		attachmentIDStr := c.Param("id")
		attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment ID"})
			return
		}

		// Get attachment first
		var attachment models.TicketAttachment
		if err := db.First(&attachment, uint(attachmentID)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "attachment not found"})
			return
		}

		// Check authentication - support both contact and user (manager)
		var hasAccess bool

		// Try contact authentication first (customer side)
		if contactIDVal, exists := c.Get("contact_id"); exists {
			if contactID, ok := contactIDVal.(uint); ok {
				// Verify ticket belongs to this contact
				var ticket models.Ticket
				if err := db.Where("ticket_id = ? AND contact_id = ?", attachment.TicketID, contactID).First(&ticket).Error; err == nil {
					hasAccess = true
				}
			}
		}

		// Try user authentication (manager side)
		if !hasAccess {
			if userVal, exists := c.Get("user"); exists {
				if _, ok := userVal.(models.User); ok {
					// Managers can delete any attachment (they have admin privileges)
					hasAccess = true
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		// Delete attachment
		fileService := services.NewFileService(db)
		if err := fileService.DeleteFile(uint(attachmentID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attachment"})
			return
		}

		// Log attachment deletion
		auditService := services.NewAuditService(db)
		auditService.LogCRUD(c, models.AuditAttachmentDeleted, models.EntityTypeAttachment, &attachment.ID, attachment.OriginalFilename, fmt.Sprintf("Attachment deleted: %s", attachment.OriginalFilename), attachment, nil)

		c.JSON(http.StatusOK, gin.H{"message": "attachment deleted successfully"})
	}
}
