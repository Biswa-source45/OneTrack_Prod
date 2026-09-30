package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	maxFeedbackImageBytes = 5 << 20
	feedbackUploadDir     = "uploads/feedback"
)

// Sniffed from the file's bytes, never trusted from its name. SVG is excluded:
// it can carry <script> and would be a stored-XSS vector.
var feedbackImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var feedbackStatuses = map[string]bool{
	models.FeedbackOpen:       true,
	models.FeedbackInProgress: true,
	models.FeedbackCompleted:  true,
}

// IsSuperAdmin checks the role by name, since the superadmin role's ID
// depends on when migration 010 ran.
func IsSuperAdmin(c *gin.Context) bool {
	v, _ := c.Get("user")
	user, ok := v.(models.User)
	return ok && isSuperAdminUser(user)
}

// feedbackReporter returns who is making the request (staff user or customer contact).
func feedbackReporter(c *gin.Context) (id uint, kind, name, email string) {
	if v, ok := c.Get("user"); ok {
		u := v.(models.User)
		return u.ID, "user", strings.TrimSpace(u.FirstName + " " + u.LastName), u.Email
	}
	ct := c.MustGet("contact").(models.Contact)
	return ct.ID, "contact", strings.TrimSpace(ct.FirstName + " " + ct.LastName), derefString(ct.Email)
}

// validateFeedbackImage returns the extension to save under, given the declared
// size and the file's first bytes.
func validateFeedbackImage(size int64, head []byte) (string, error) {
	if size > maxFeedbackImageBytes {
		return "", fmt.Errorf("image must be 5 MB or smaller")
	}
	ext, ok := feedbackImageTypes[http.DetectContentType(head)]
	if !ok {
		return "", fmt.Errorf("only JPEG, PNG or WebP images are allowed")
	}
	return ext, nil
}

// saveFeedbackImage validates the upload and writes it under a random name.
func saveFeedbackImage(fh *multipart.FileHeader) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("could not read the uploaded image")
	}
	defer src.Close()

	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	ext, err := validateFeedbackImage(fh.Size, head[:n])
	if err != nil {
		return "", err
	}

	name := make([]byte, 16)
	if _, err := rand.Read(name); err != nil {
		return "", err
	}
	if err := os.MkdirAll(feedbackUploadDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.ToSlash(filepath.Join(feedbackUploadDir, hex.EncodeToString(name)+ext))
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(path)
		return "", err
	}
	return path, nil
}

// CreateFeedback: any logged-in user or contact. Multipart: title, description, page_url, image (optional).
func CreateFeedback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 5 MB image + form fields; anything bigger is rejected before it's buffered.
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFeedbackImageBytes+(1<<20))

		// Parse explicitly first: c.PostForm would swallow a too-large error.
		var maxBytesErr *http.MaxBytesError
		if err := c.Request.ParseMultipartForm(8 << 20); errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "image must be 5 MB or smaller"})
			return
		}
		fb := models.Feedback{
			// Single line: the title goes into the email Subject header, so CR/LF would allow header injection.
			Title:       strings.Join(strings.Fields(c.PostForm("title")), " "),
			Description: strings.TrimSpace(c.PostForm("description")),
			PageURL:     strings.TrimSpace(c.PostForm("page_url")),
			Status:      models.FeedbackOpen,
		}
		if fb.Title == "" || fb.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title and description are required"})
			return
		}
		if utf8.RuneCountInString(fb.Title) > 200 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title must be 200 characters or fewer"})
			return
		}
		if len(fb.PageURL) > 500 {
			fb.PageURL = "" // informational only; not worth failing the report over
		}
		fb.ReporterID, fb.ReporterType, fb.ReporterName, fb.ReporterEmail = feedbackReporter(c)

		fh, err := c.FormFile("image")
		switch {
		case err == nil:
			path, err := saveFeedbackImage(fh)
			if err != nil {
				log.Printf("[FEEDBACK] image rejected for %s #%d: %v", fb.ReporterType, fb.ReporterID, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fb.ImagePath = path
		case !errors.Is(err, http.ErrMissingFile):
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not read the uploaded image"})
			return
		}

		if err := db.Create(&fb).Error; err != nil {
			if fb.ImagePath != "" {
				os.Remove(fb.ImagePath)
			}
			log.Printf("[FEEDBACK] save failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save feedback"})
			return
		}

		services.NewAuditService(db).LogCRUD(c, models.AuditFeedbackCreated, models.EntityTypeFeedback, &fb.ID, fb.Title,
			fmt.Sprintf("Feedback submitted: %s", fb.Title), nil, fb)

		// Mail in the background: SMTP is slow and a mail failure must not lose the report.
		go func(fb models.Feedback) {
			if err := services.NewEmailNotificationServiceFromEnv(db).SendFeedbackEmail(&fb); err != nil {
				log.Printf("[FEEDBACK] email for #%d failed: %v", fb.ID, err)
			}
		}(fb)

		c.JSON(http.StatusCreated, fb)
	}
}

// ListMyFeedback: the caller's own reports, newest first.
func ListMyFeedback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, kind, _, _ := feedbackReporter(c)
		var list []models.Feedback
		if err := db.Where("reporter_type = ? AND reporter_id = ?", kind, id).
			Order("created_at DESC").Limit(200).Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load feedback"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

// ListAllFeedback: superadmin only, optional ?status= filter.
func ListAllFeedback(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsSuperAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		q := db.Order("created_at DESC")
		if s := c.Query("status"); s != "" {
			q = q.Where("status = ?", s)
		}
		var list []models.Feedback
		// ponytail: no paging, capped at 500 rows; add page/limit like dumped queries if the inbox grows past that.
		if err := q.Limit(500).Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load feedback"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

// UpdateFeedbackStatus: superadmin only. Body: {"status": "OPEN" | "IN PROGRESS" | "COMPLETED"}
func UpdateFeedbackStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsSuperAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var in struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&in); err != nil || !feedbackStatuses[in.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status must be OPEN, IN PROGRESS or COMPLETED"})
			return
		}
		fb, ok := findFeedback(c, db)
		if !ok {
			return
		}
		old := fb.Status
		if err := db.Model(&fb).Update("status", in.Status).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
			return
		}
		fb.Status = in.Status
		services.NewAuditService(db).LogCRUD(c, models.AuditFeedbackStatusChanged, models.EntityTypeFeedback, &fb.ID, fb.Title,
			fmt.Sprintf("Feedback status changed: %s -> %s", old, in.Status),
			gin.H{"status": old}, gin.H{"status": in.Status})
		c.JSON(http.StatusOK, fb)
	}
}

// findFeedback loads :id, writing a 404 if it's missing. The id is parsed to a
// number first: GORM treats a raw string argument to First as SQL.
func findFeedback(c *gin.Context, db *gorm.DB) (models.Feedback, bool) {
	var fb models.Feedback
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err == nil {
		err = db.First(&fb, id).Error
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "feedback not found"})
		return fb, false
	}
	return fb, true
}

// GetFeedbackImage: the reporter or a superadmin.
func GetFeedbackImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fb, ok := findFeedback(c, db)
		if !ok {
			return
		}
		if fb.ImagePath == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "this feedback has no image"})
			return
		}
		id, kind, _, _ := feedbackReporter(c)
		isOwner := fb.ReporterID == id && fb.ReporterType == kind
		if !isOwner && (kind != "user" || !IsSuperAdmin(c)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Header("X-Content-Type-Options", "nosniff")
		c.File(fb.ImagePath)
	}
}
