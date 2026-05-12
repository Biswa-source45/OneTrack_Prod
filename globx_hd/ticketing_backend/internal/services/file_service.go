package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"gorm.io/gorm"
)

type FileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{db: db}
}

// AllowedFileTypes defines the allowed file extensions and their MIME types
var AllowedFileTypes = map[string][]string{
	"image": {"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
	"document": {"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", 
				"text/plain", "application/rtf"},
	"spreadsheet": {"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "text/csv"},
}

// AllowedExtensions defines the allowed file extensions
var AllowedExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".doc", ".docx", ".txt", ".rtf", ".xls", ".xlsx", ".csv"}

const MaxFileSize = 3 * 1024 * 1024 // 3MB in bytes

// ValidateFile checks if the uploaded file meets our requirements
func (fs *FileService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of 3MB")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	isValidExt := false
	for _, allowedExt := range AllowedExtensions {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return fmt.Errorf("file type not allowed. Allowed types: %v", AllowedExtensions)
	}

	return nil
}

// SanitizeFilename removes or replaces invalid characters in filename
func (fs *FileService) SanitizeFilename(filename string) string {
	// Replace spaces with underscores and remove special characters
	sanitized := strings.ReplaceAll(filename, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "(", "")
	sanitized = strings.ReplaceAll(sanitized, ")", "")
	sanitized = strings.ReplaceAll(sanitized, "[", "")
	sanitized = strings.ReplaceAll(sanitized, "]", "")
	sanitized = strings.ReplaceAll(sanitized, "{", "")
	sanitized = strings.ReplaceAll(sanitized, "}", "")
	sanitized = strings.ReplaceAll(sanitized, "&", "and")
	return sanitized
}

// GenerateStoredFilename creates the stored filename with ticket_id prefix
func (fs *FileService) GenerateStoredFilename(ticketID, originalFilename string) string {
	sanitizedFilename := fs.SanitizeFilename(originalFilename)
	return fmt.Sprintf("%s_%s", ticketID, sanitizedFilename)
}

// CreateUploadDirectory ensures the upload directory exists
func (fs *FileService) CreateUploadDirectory() (string, error) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	
	uploadDir := filepath.Join("uploads", "tickets", fmt.Sprintf("%d", year), fmt.Sprintf("%02d", int(month)))
	
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}
	
	return uploadDir, nil
}

// SaveFile saves the uploaded file to the filesystem
func (fs *FileService) SaveFile(fileHeader *multipart.FileHeader, ticketID string) (*models.TicketAttachment, error) {
	// Validate file
	if err := fs.ValidateFile(fileHeader); err != nil {
		return nil, err
	}

	// Create upload directory
	uploadDir, err := fs.CreateUploadDirectory()
	if err != nil {
		return nil, err
	}

	// Generate stored filename
	storedFilename := fs.GenerateStoredFilename(ticketID, fileHeader.Filename)
	filePath := filepath.Join(uploadDir, storedFilename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Detect MIME type
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Create attachment record
	attachment := &models.TicketAttachment{
		TicketID:         ticketID,
		OriginalFilename: fileHeader.Filename,
		StoredFilename:   storedFilename,
		FilePath:         filePath,
		FileSize:         int(fileHeader.Size),
		MimeType:         mimeType,
	}

	return attachment, nil
}

// DeleteFile removes the file from filesystem and database
func (fs *FileService) DeleteFile(attachmentID uint) error {
	var attachment models.TicketAttachment
	if err := fs.db.First(&attachment, attachmentID).Error; err != nil {
		return fmt.Errorf("attachment not found: %v", err)
	}

	// Delete file from filesystem
	if err := os.Remove(attachment.FilePath); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: Failed to delete file from filesystem: %v\n", err)
	}

	// Delete from database
	if err := fs.db.Delete(&attachment).Error; err != nil {
		return fmt.Errorf("failed to delete attachment from database: %v", err)
	}

	return nil
}

// GetTicketAttachments retrieves all attachments for a ticket
func (fs *FileService) GetTicketAttachments(ticketID string) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	if err := fs.db.Where("ticket_id = ?", ticketID).Preload("Contact").Find(&attachments).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve attachments: %v", err)
	}
	return attachments, nil
}
