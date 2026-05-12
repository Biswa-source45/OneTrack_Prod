package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============================================================================
// SMART TICKET RESOLVER - AI-POWERED EMAIL TO TICKET
// ============================================================================
// This is the main endpoint for n8n + Gemini integration.
// It receives AI-extracted hints and performs intelligent database lookups
// to resolve contacts, accounts, and products with high accuracy.
// ============================================================================

// SmartTicketRequest represents the AI-extracted data from an email
type SmartTicketRequest struct {
	// Raw email info
	SenderEmail string `json:"sender_email" binding:"required"`
	SenderName  string `json:"sender_name,omitempty"`
	Subject     string `json:"subject" binding:"required"`
	Body        string `json:"body" binding:"required"`

	// AI-extracted hints (all optional - resolver will try multiple strategies)
	PhoneNumbers  []string `json:"phone_numbers,omitempty"`  // Phones found in email
	PersonNames   []string `json:"person_names,omitempty"`   // Names from signature/body
	OrgNames      []string `json:"org_names,omitempty"`      // Organization names found
	ProductHints  []string `json:"product_hints,omitempty"`  // Product/service mentions
	PriorityHints []string `json:"priority_hints,omitempty"` // Urgency keywords

	// Attachments (base64 encoded)
	Attachments []N8nAttachment `json:"attachments,omitempty"`
}

// ResolutionResult contains details about how an entity was resolved
type ResolutionResult struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Method     string `json:"method"`     // How it was found (email_match, phone_match, etc.)
	Confidence int    `json:"confidence"` // 0-100
	Warning    string `json:"warning,omitempty"`
}

// SmartTicketResponse is the response from smart ticket creation
type SmartTicketResponse struct {
	Success    bool            `json:"success"`
	TicketID   string          `json:"ticket_id,omitempty"`
	Error      string          `json:"error,omitempty"`
	Resolution *ResolutionInfo `json:"resolution,omitempty"`
	Ticket     *TicketInfo     `json:"ticket,omitempty"`
	Warnings   []string        `json:"warnings,omitempty"`
}

// ResolutionInfo contains how each entity was resolved
type ResolutionInfo struct {
	Contact ResolutionResult `json:"contact"`
	Account ResolutionResult `json:"account"`
	Product ResolutionResult `json:"product"`
}

// TicketInfo contains basic ticket info for response
type TicketInfo struct {
	ID        uint      `json:"id"`
	TicketID  string    `json:"ticket_id"`
	Subject   string    `json:"subject"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

// ContactCandidate represents a potential contact match with scoring
type ContactCandidate struct {
	Contact     models.Contact
	Score       int
	Method      string
	AccountID   *uint
	AccountName string
}

// AccountCandidate represents a potential account match with scoring
type AccountCandidate struct {
	Account models.Account
	Score   int
	Method  string
}

// ============================================================================
// MAIN HANDLER
// ============================================================================

// SmartTicketHandler is the main endpoint for intelligent ticket creation
func SmartTicketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SmartTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, SmartTicketResponse{
				Success: false,
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		log.Printf("[SmartResolver] Processing email from: %s, Subject: %s", req.SenderEmail, req.Subject)

		resolver := NewSmartResolver(db)
		response := resolver.ResolveAndCreateTicket(&req)

		if response.Success {
			c.JSON(http.StatusCreated, response)
		} else {
			c.JSON(http.StatusBadRequest, response)
		}
	}
}

// ============================================================================
// SMART RESOLVER ENGINE
// ============================================================================

// SmartResolver handles the intelligent resolution logic
type SmartResolver struct {
	db       *gorm.DB
	warnings []string
}

// NewSmartResolver creates a new resolver instance
func NewSmartResolver(db *gorm.DB) *SmartResolver {
	return &SmartResolver{
		db:       db,
		warnings: make([]string, 0),
	}
}

// ResolveAndCreateTicket is the main entry point
func (r *SmartResolver) ResolveAndCreateTicket(req *SmartTicketRequest) SmartTicketResponse {
	// Step 1: Resolve Contact (most critical)
	contactResult := r.resolveContact(req)
	if contactResult.ID == 0 {
		return SmartTicketResponse{
			Success:  false,
			Error:    "Could not identify contact from email",
			Warnings: r.warnings,
		}
	}

	// Get full contact info
	var contact models.Contact
	if err := r.db.Preload("Account").First(&contact, contactResult.ID).Error; err != nil {
		return SmartTicketResponse{
			Success: false,
			Error:   "Failed to load contact details",
		}
	}

	// Step 2: Resolve Account (from contact or hints)
	accountResult := r.resolveAccount(req, &contact)

	// Step 3: Resolve Product
	productResult := r.resolveProduct(req)
	if productResult.ID == 0 {
		return SmartTicketResponse{
			Success:  false,
			Error:    "No products configured in system",
			Warnings: r.warnings,
		}
	}

	// Step 4: Determine Priority
	priority := r.determinePriority(req)

	// Step 5: Create Ticket
	ticket, err := r.createTicket(req, &contact, accountResult.ID, productResult.ID, priority)
	if err != nil {
		return SmartTicketResponse{
			Success:  false,
			Error:    "Failed to create ticket: " + err.Error(),
			Warnings: r.warnings,
		}
	}

	log.Printf("[SmartResolver] Successfully created ticket %s", ticket.TicketID)

	return SmartTicketResponse{
		Success:  true,
		TicketID: ticket.TicketID,
		Resolution: &ResolutionInfo{
			Contact: contactResult,
			Account: accountResult,
			Product: productResult,
		},
		Ticket: &TicketInfo{
			ID:        ticket.ID,
			TicketID:  ticket.TicketID,
			Subject:   ticket.Subject,
			Status:    ticket.TicketStatus,
			Priority:  ticket.Priority,
			CreatedAt: ticket.CreatedAt,
		},
		Warnings: r.warnings,
	}
}

// ============================================================================
// CONTACT RESOLUTION - WATERFALL STRATEGY
// ============================================================================

func (r *SmartResolver) resolveContact(req *SmartTicketRequest) ResolutionResult {
	candidates := make([]ContactCandidate, 0)

	// Strategy 1: Exact email match (highest confidence)
	if candidate := r.findContactByEmail(req.SenderEmail); candidate != nil {
		log.Printf("[SmartResolver] Contact found by email match: %s %s", candidate.Contact.FirstName, candidate.Contact.LastName)
		return ResolutionResult{
			ID:         candidate.Contact.ID,
			Name:       candidate.Contact.FirstName + " " + candidate.Contact.LastName,
			Method:     "email_exact_match",
			Confidence: 100,
		}
	}

	// Strategy 2: Phone number match
	for _, phone := range req.PhoneNumbers {
		if candidate := r.findContactByPhone(phone); candidate != nil {
			candidate.Method = "phone_match"
			candidate.Score = 95
			candidates = append(candidates, *candidate)
			log.Printf("[SmartResolver] Contact candidate from phone: %s %s (score: %d)",
				candidate.Contact.FirstName, candidate.Contact.LastName, candidate.Score)
		}
	}

	// Strategy 3: Parse email username for name hints
	emailNames := r.extractNamesFromEmail(req.SenderEmail)
	for _, name := range emailNames {
		matches := r.findContactsByName(name, nil)
		for _, match := range matches {
			match.Method = "email_username_parse"
			candidates = append(candidates, match)
			log.Printf("[SmartResolver] Contact candidate from email parse: %s %s (score: %d)",
				match.Contact.FirstName, match.Contact.LastName, match.Score)
		}
	}

	// Strategy 4: Name hints from AI
	for _, name := range req.PersonNames {
		matches := r.findContactsByName(name, nil)
		for _, match := range matches {
			match.Method = "ai_extracted_name"
			candidates = append(candidates, match)
			log.Printf("[SmartResolver] Contact candidate from AI name: %s %s (score: %d)",
				match.Contact.FirstName, match.Contact.LastName, match.Score)
		}
	}

	// Strategy 5: Find contacts via organization
	for _, orgName := range req.OrgNames {
		account := r.findAccountByName(orgName)
		if account != nil {
			contacts := r.findContactsInAccount(account.ID)
			for _, contact := range contacts {
				// Boost score if name matches any person names
				score := 50
				for _, personName := range req.PersonNames {
					nameScore := r.fuzzyMatchScore(
						strings.ToLower(contact.FirstName+" "+contact.LastName),
						strings.ToLower(personName),
					)
					if nameScore > score {
						score = nameScore
					}
				}
				candidates = append(candidates, ContactCandidate{
					Contact:     contact,
					Score:       score,
					Method:      "org_match",
					AccountID:   &account.ID,
					AccountName: account.AccountName,
				})
			}
		}
	}

	// Strategy 6: Email domain matching
	domain := r.extractEmailDomain(req.SenderEmail)
	if domain != "" && !r.isCommonEmailDomain(domain) {
		domainContacts := r.findContactsByEmailDomain(domain)
		for _, contact := range domainContacts {
			candidates = append(candidates, ContactCandidate{
				Contact: contact,
				Score:   40,
				Method:  "email_domain_match",
			})
		}
	}

	// Pick best candidate
	if len(candidates) == 0 {
		r.warnings = append(r.warnings, "No contact could be identified from the email")
		return ResolutionResult{}
	}

	// Sort by score descending
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	best := candidates[0]

	// Add warning if low confidence
	if best.Score < 70 {
		r.warnings = append(r.warnings, fmt.Sprintf(
			"Contact match confidence is low (%d%%). Please verify: %s %s",
			best.Score, best.Contact.FirstName, best.Contact.LastName,
		))
	}

	// Add warning if multiple high-scoring candidates
	if len(candidates) > 1 && candidates[1].Score >= best.Score-10 {
		r.warnings = append(r.warnings, fmt.Sprintf(
			"Multiple possible contacts found. Selected: %s %s. Other candidate: %s %s",
			best.Contact.FirstName, best.Contact.LastName,
			candidates[1].Contact.FirstName, candidates[1].Contact.LastName,
		))
	}

	return ResolutionResult{
		ID:         best.Contact.ID,
		Name:       best.Contact.FirstName + " " + best.Contact.LastName,
		Method:     best.Method,
		Confidence: best.Score,
	}
}

// ============================================================================
// CONTACT LOOKUP METHODS
// ============================================================================

func (r *SmartResolver) findContactByEmail(email string) *ContactCandidate {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil
	}

	var contact models.Contact
	err := r.db.Preload("Account").Where("LOWER(email) = ?", email).First(&contact).Error
	if err != nil {
		return nil
	}

	candidate := &ContactCandidate{
		Contact: contact,
		Score:   100,
		Method:  "email_exact_match",
	}

	if contact.AccountID != nil {
		candidate.AccountID = contact.AccountID
		candidate.AccountName = contact.Account.AccountName
	}

	return candidate
}

func (r *SmartResolver) findContactByPhone(phone string) *ContactCandidate {
	normalized := r.normalizePhone(phone)
	if len(normalized) < 10 {
		return nil
	}

	var contact models.Contact
	var err error

	// Try exact match
	err = r.db.Preload("Account").Where("mobile = ?", phone).First(&contact).Error
	if err != nil {
		// Try normalized
		err = r.db.Preload("Account").Where("mobile = ?", normalized).First(&contact).Error
	}
	if err != nil {
		// Try last 10 digits (handles country codes)
		last10 := normalized
		if len(normalized) > 10 {
			last10 = normalized[len(normalized)-10:]
		}
		err = r.db.Preload("Account").Where("mobile LIKE ?", "%"+last10).First(&contact).Error
	}

	if err != nil {
		return nil
	}

	candidate := &ContactCandidate{
		Contact: contact,
		Score:   95,
		Method:  "phone_match",
	}

	if contact.AccountID != nil {
		candidate.AccountID = contact.AccountID
		candidate.AccountName = contact.Account.AccountName
	}

	return candidate
}

func (r *SmartResolver) findContactsByName(name string, accountID *uint) []ContactCandidate {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	nameLower := strings.ToLower(name)
	parts := strings.Fields(nameLower)

	var contacts []models.Contact
	query := r.db.Preload("Account")

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}

	// Build search query
	if len(parts) >= 2 {
		// Try first + last name
		query.Where(
			"(LOWER(first_name) LIKE ? AND LOWER(last_name) LIKE ?) OR "+
				"(LOWER(first_name) LIKE ? AND LOWER(last_name) LIKE ?)",
			parts[0]+"%", "%"+strings.Join(parts[1:], " ")+"%",
			"%"+strings.Join(parts[1:], " ")+"%", parts[0]+"%", // Reversed
		).Limit(10).Find(&contacts)
	} else {
		query.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			"%"+nameLower+"%", "%"+nameLower+"%",
		).Limit(10).Find(&contacts)
	}

	candidates := make([]ContactCandidate, 0, len(contacts))
	for _, contact := range contacts {
		fullName := strings.ToLower(contact.FirstName + " " + contact.LastName)
		score := r.fuzzyMatchScore(fullName, nameLower)

		candidate := ContactCandidate{
			Contact: contact,
			Score:   score,
			Method:  "name_fuzzy_match",
		}

		if contact.AccountID != nil {
			candidate.AccountID = contact.AccountID
			candidate.AccountName = contact.Account.AccountName
		}

		candidates = append(candidates, candidate)
	}

	return candidates
}

func (r *SmartResolver) findContactsInAccount(accountID uint) []models.Contact {
	var contacts []models.Contact
	r.db.Where("account_id = ?", accountID).Find(&contacts)
	return contacts
}

func (r *SmartResolver) findContactsByEmailDomain(domain string) []models.Contact {
	var contacts []models.Contact
	r.db.Preload("Account").
		Where("LOWER(email) LIKE ?", "%@"+strings.ToLower(domain)).
		Limit(10).
		Find(&contacts)
	return contacts
}

// ============================================================================
// ACCOUNT RESOLUTION
// ============================================================================

func (r *SmartResolver) resolveAccount(req *SmartTicketRequest, contact *models.Contact) ResolutionResult {
	// If contact has an account, use it
	if contact.AccountID != nil {
		var account models.Account
		if err := r.db.First(&account, *contact.AccountID).Error; err == nil {
			return ResolutionResult{
				ID:         account.ID,
				Name:       account.AccountName,
				Method:     "from_contact",
				Confidence: 100,
			}
		}
	}

	// Try to find account from org names
	for _, orgName := range req.OrgNames {
		account := r.findAccountByName(orgName)
		if account != nil {
			return ResolutionResult{
				ID:         account.ID,
				Name:       account.AccountName,
				Method:     "org_name_match",
				Confidence: 80,
			}
		}
	}

	// Try email domain
	domain := r.extractEmailDomain(req.SenderEmail)
	if domain != "" && !r.isCommonEmailDomain(domain) {
		account := r.findAccountByEmailDomain(domain)
		if account != nil {
			return ResolutionResult{
				ID:         account.ID,
				Name:       account.AccountName,
				Method:     "email_domain_match",
				Confidence: 60,
			}
		}
	}

	// Individual contact - no account
	return ResolutionResult{
		ID:         0,
		Name:       "Individual",
		Method:     "no_account",
		Confidence: 100,
	}
}

func (r *SmartResolver) findAccountByName(name string) *models.Account {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}

	nameLower := strings.ToLower(name)
	var account models.Account

	// Try exact match first
	if err := r.db.Where("LOWER(account_name) = ?", nameLower).First(&account).Error; err == nil {
		return &account
	}

	// Try fuzzy match
	var accounts []models.Account
	r.db.Where("LOWER(account_name) LIKE ?", "%"+nameLower+"%").
		Limit(5).
		Find(&accounts)

	if len(accounts) == 0 {
		return nil
	}

	// Find best match
	var bestAccount *models.Account
	bestScore := 0
	for i := range accounts {
		score := r.fuzzyMatchScore(strings.ToLower(accounts[i].AccountName), nameLower)
		if score > bestScore {
			bestScore = score
			bestAccount = &accounts[i]
		}
	}

	if bestScore >= 50 {
		return bestAccount
	}

	return nil
}

func (r *SmartResolver) findAccountByEmailDomain(domain string) *models.Account {
	var contact models.Contact
	err := r.db.Preload("Account").
		Where("LOWER(email) LIKE ? AND account_id IS NOT NULL", "%@"+strings.ToLower(domain)).
		First(&contact).Error

	if err != nil {
		return nil
	}

	return &contact.Account
}

// ============================================================================
// PRODUCT RESOLUTION
// ============================================================================

func (r *SmartResolver) resolveProduct(req *SmartTicketRequest) ResolutionResult {
	// Try product hints from AI
	for _, hint := range req.ProductHints {
		product := r.findProductByKeyword(hint)
		if product != nil {
			return ResolutionResult{
				ID:         product.ID,
				Name:       product.ProductName,
				Method:     "keyword_match",
				Confidence: 80,
			}
		}
	}

	// Search in subject and body
	textToSearch := strings.ToLower(req.Subject + " " + req.Body)
	var products []models.MasterProduct
	r.db.Find(&products)

	for _, product := range products {
		if strings.Contains(textToSearch, strings.ToLower(product.ProductName)) {
			return ResolutionResult{
				ID:         product.ID,
				Name:       product.ProductName,
				Method:     "content_match",
				Confidence: 70,
			}
		}
	}

	// Use default product
	var defaultProduct models.MasterProduct
	if err := r.db.First(&defaultProduct).Error; err == nil {
		return ResolutionResult{
			ID:         defaultProduct.ID,
			Name:       defaultProduct.ProductName,
			Method:     "default",
			Confidence: 50,
		}
	}

	return ResolutionResult{}
}

func (r *SmartResolver) findProductByKeyword(keyword string) *models.MasterProduct {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	var product models.MasterProduct
	err := r.db.Where("LOWER(product_name) LIKE ?", "%"+strings.ToLower(keyword)+"%").
		First(&product).Error

	if err != nil {
		return nil
	}

	return &product
}

// ============================================================================
// PRIORITY DETERMINATION
// ============================================================================

func (r *SmartResolver) determinePriority(req *SmartTicketRequest) string {
	// Check AI-provided hints first
	for _, hint := range req.PriorityHints {
		hint = strings.ToLower(hint)
		if r.containsAny(hint, []string{"urgent", "critical", "emergency", "asap", "immediately", "high"}) {
			return "High"
		}
		if r.containsAny(hint, []string{"low", "minor", "when possible", "not urgent"}) {
			return "Low"
		}
	}

	// Check subject and body
	text := strings.ToLower(req.Subject + " " + req.Body)

	highPriorityKeywords := []string{
		"urgent", "critical", "emergency", "asap", "immediately",
		"down", "outage", "not working", "broken", "failed", "crash",
		"blocker", "production issue", "business critical",
	}

	lowPriorityKeywords := []string{
		"when you have time", "low priority", "not urgent",
		"question", "inquiry", "information", "clarification",
	}

	for _, kw := range highPriorityKeywords {
		if strings.Contains(text, kw) {
			return "High"
		}
	}

	for _, kw := range lowPriorityKeywords {
		if strings.Contains(text, kw) {
			return "Low"
		}
	}

	return "Medium"
}

// ============================================================================
// TICKET CREATION
// ============================================================================

func (r *SmartResolver) createTicket(req *SmartTicketRequest, contact *models.Contact, accountID uint, productID uint, priority string) (*models.Ticket, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Get customer code
	var customerCode string
	if accountID > 0 {
		var account models.Account
		if err := tx.First(&account, accountID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		customerCode = account.CustomerCode
	} else {
		customerCode = contact.CustomerCode
	}

	// Generate ticket ID
	now := time.Now()
	dateStr := utils.FormatDateForTicketID(now)

	// Get sequence based on whether we have an account or individual contact
	var seq int
	var err error
	if accountID > 0 {
		seq, err = utils.GetNextTicketSequence(tx, accountID, now)
	} else {
		// For individual contacts, use contact ID for sequence tracking
		seq, err = getNextContactTicketSequence(tx, contact.ID, now)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ticketID := utils.FormatTicketID(customerCode, dateStr, seq)

	// Prepare account ID pointer
	var accIDPtr *uint
	if accountID > 0 {
		accIDPtr = &accountID
	}

	// Create ticket
	ticket := &models.Ticket{
		TicketID:      ticketID,
		AccountID:     accIDPtr,
		ContactID:     contact.ID,
		ProductID:     productID,
		Subject:       truncateString(req.Subject, 255),
		TicketDetails: req.Body,
		TicketStatus:  "OPEN",
		Priority:      priority,
		Channel:       "Mail",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := tx.Create(ticket).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Process attachments
	r.processAttachments(tx, ticket, req.Attachments, contact.ID)

	// Commit
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Log activity (after commit)
	activityService := services.NewActivityService(r.db)
	activityService.LogTicketCreation(ticket.ID, nil, ticket.TicketID)

	description := fmt.Sprintf("Ticket created via AI-powered email processing from: %s", req.SenderEmail)
	activityService.LogActivity(ticket.ID, nil, models.ActivityTicketCreated, description)

	// Send notification asynchronously
	go sendN8nTicketNotification(r.db, ticket)

	return ticket, nil
}

func (r *SmartResolver) processAttachments(tx *gorm.DB, ticket *models.Ticket, attachments []N8nAttachment, contactID uint) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	for _, att := range attachments {
		if att.Filename == "" || att.Data == "" {
			continue
		}

		data, err := base64.StdEncoding.DecodeString(att.Data)
		if err != nil {
			log.Printf("[SmartResolver] Failed to decode attachment %s: %v", att.Filename, err)
			continue
		}

		uploadPath := filepath.Join(uploadDir, "ticket_attachments", ticket.TicketID)
		if err := os.MkdirAll(uploadPath, 0755); err != nil {
			log.Printf("[SmartResolver] Failed to create directory: %v", err)
			continue
		}

		storedFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), att.Filename)
		filePath := filepath.Join(uploadPath, storedFilename)

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			log.Printf("[SmartResolver] Failed to save attachment: %v", err)
			continue
		}

		attachment := models.TicketAttachment{
			TicketID:         ticket.TicketID,
			OriginalFilename: att.Filename,
			StoredFilename:   storedFilename,
			FilePath:         filePath,
			FileSize:         len(data),
			MimeType:         att.ContentType,
			UploadedBy:       contactID,
			UploadedAt:       time.Now(),
		}

		if err := tx.Create(&attachment).Error; err != nil {
			log.Printf("[SmartResolver] Failed to save attachment record: %v", err)
		}
	}
}

// ============================================================================
// NAME EXTRACTION FROM EMAIL
// ============================================================================

// extractNamesFromEmail parses email address to extract possible names
// Examples:
//   - doejohn8980@gmail.com → ["John Doe", "Doe John"]
//   - john.doe@company.com → ["John Doe"]
//   - johndoe@company.com → ["John Doe", "Johndoe"]
func (r *SmartResolver) extractNamesFromEmail(email string) []string {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil
	}

	// Get username part
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil
	}
	username := parts[0]

	names := make([]string, 0)

	// Strategy 1: Split by common separators
	for _, sep := range []string{".", "_", "-"} {
		if strings.Contains(username, sep) {
			parts := strings.Split(username, sep)
			if len(parts) >= 2 {
				// Remove numbers from parts
				cleanParts := make([]string, 0)
				for _, p := range parts {
					cleaned := r.removeNumbers(p)
					if len(cleaned) >= 2 {
						cleanParts = append(cleanParts, cleaned)
					}
				}
				if len(cleanParts) >= 2 {
					// First Last
					name := r.capitalizeWords(strings.Join(cleanParts, " "))
					names = append(names, name)
				}
			}
		}
	}

	// Strategy 2: Try to split concatenated names
	cleanUsername := r.removeNumbers(username)
	if len(cleanUsername) >= 4 {
		// Try to find common name patterns
		possibleSplits := r.findNameSplits(cleanUsername)
		names = append(names, possibleSplits...)
	}

	// Remove duplicates
	seen := make(map[string]bool)
	unique := make([]string, 0)
	for _, name := range names {
		nameLower := strings.ToLower(name)
		if !seen[nameLower] && len(name) >= 3 {
			seen[nameLower] = true
			unique = append(unique, name)
		}
	}

	return unique
}

// findNameSplits tries to intelligently split a concatenated string into names
func (r *SmartResolver) findNameSplits(s string) []string {
	s = strings.ToLower(s)
	results := make([]string, 0)

	// Common first name patterns to look for
	commonFirstNames := []string{
		"john", "jane", "james", "david", "michael", "robert", "william",
		"mary", "sarah", "jennifer", "lisa", "emily", "anna", "raj", "amit",
		"priya", "neha", "rahul", "vikram", "arun", "kumar", "singh", "sharma",
	}

	// Try to match common first names at the start
	for _, firstName := range commonFirstNames {
		if strings.HasPrefix(s, firstName) && len(s) > len(firstName)+1 {
			lastName := s[len(firstName):]
			if len(lastName) >= 2 {
				name := r.capitalizeWords(firstName + " " + lastName)
				results = append(results, name)
			}
		}
		// Also try at the end
		if strings.HasSuffix(s, firstName) && len(s) > len(firstName)+1 {
			firstPart := s[:len(s)-len(firstName)]
			if len(firstPart) >= 2 {
				name := r.capitalizeWords(firstPart + " " + firstName)
				results = append(results, name)
			}
		}
	}

	// If no common names found, try splitting at vowel-consonant boundaries
	if len(results) == 0 && len(s) >= 6 {
		// Try splitting in the middle third of the string
		start := len(s) / 3
		end := (len(s) * 2) / 3
		for i := start; i <= end; i++ {
			if i > 1 && i < len(s)-1 {
				// Look for consonant-vowel or vowel-consonant transitions
				if r.isVowel(rune(s[i-1])) != r.isVowel(rune(s[i])) {
					first := s[:i]
					second := s[i:]
					if len(first) >= 2 && len(second) >= 2 {
						name := r.capitalizeWords(first + " " + second)
						results = append(results, name)
					}
				}
			}
		}
	}

	return results
}

// ============================================================================
// FUZZY MATCHING ALGORITHM
// ============================================================================

// fuzzyMatchScore calculates a match score between 0-100
// Uses a combination of:
// - Exact match
// - Prefix match
// - Levenshtein distance
// - Word overlap
func (r *SmartResolver) fuzzyMatchScore(text, pattern string) int {
	text = strings.ToLower(strings.TrimSpace(text))
	pattern = strings.ToLower(strings.TrimSpace(pattern))

	if text == "" || pattern == "" {
		return 0
	}

	// Exact match
	if text == pattern {
		return 100
	}

	// Contains exact pattern
	if strings.Contains(text, pattern) {
		return 85
	}

	// Pattern contains text
	if strings.Contains(pattern, text) {
		return 80
	}

	// Prefix match
	if strings.HasPrefix(text, pattern) || strings.HasPrefix(pattern, text) {
		return 75
	}

	// Word-level matching
	textWords := strings.Fields(text)
	patternWords := strings.Fields(pattern)

	if len(textWords) == 0 || len(patternWords) == 0 {
		return 0
	}

	// Count matching words
	matchingWords := 0
	for _, pw := range patternWords {
		for _, tw := range textWords {
			if tw == pw {
				matchingWords++
				break
			}
			// Partial word match
			if strings.Contains(tw, pw) || strings.Contains(pw, tw) {
				matchingWords++
				break
			}
		}
	}

	wordScore := (matchingWords * 100) / len(patternWords)
	if wordScore >= 50 {
		return 60 + (wordScore * 30 / 100)
	}

	// Levenshtein-based similarity
	distance := r.levenshteinDistance(text, pattern)
	maxLen := max(len(text), len(pattern))
	if maxLen == 0 {
		return 0
	}

	similarity := ((maxLen - distance) * 100) / maxLen
	if similarity >= 70 {
		return 40 + (similarity * 40 / 100)
	}

	return similarity / 2
}

// levenshteinDistance calculates edit distance between two strings
func (r *SmartResolver) levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Optimization: swap if s2 is shorter (we iterate over s2)
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}

	prevRow := make([]int, len(s2)+1)
	currRow := make([]int, len(s2)+1)

	for j := 0; j <= len(s2); j++ {
		prevRow[j] = j
	}

	for i := 1; i <= len(s1); i++ {
		currRow[0] = i
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			currRow[j] = min(
				prevRow[j]+1,      // deletion
				currRow[j-1]+1,    // insertion
				prevRow[j-1]+cost, // substitution
			)
		}
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[len(s2)]
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

func (r *SmartResolver) normalizePhone(phone string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	return reg.ReplaceAllString(phone, "")
}

func (r *SmartResolver) extractEmailDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

func (r *SmartResolver) isCommonEmailDomain(domain string) bool {
	commonDomains := map[string]bool{
		"gmail.com":      true,
		"yahoo.com":      true,
		"hotmail.com":    true,
		"outlook.com":    true,
		"live.com":       true,
		"aol.com":        true,
		"icloud.com":     true,
		"mail.com":       true,
		"protonmail.com": true,
		"ymail.com":      true,
		"rediffmail.com": true,
	}
	return commonDomains[strings.ToLower(domain)]
}

func (r *SmartResolver) removeNumbers(s string) string {
	var result strings.Builder
	for _, r := range s {
		if !unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (r *SmartResolver) capitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func (r *SmartResolver) isVowel(c rune) bool {
	vowels := "aeiou"
	return strings.ContainsRune(vowels, c)
}

func (r *SmartResolver) containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getNextContactTicketSequence gets sequence for individual contacts
func getNextContactTicketSequence(db *gorm.DB, contactID uint, date time.Time) (int, error) {
	dateStr := utils.FormatDateForTicketID(date)
	pattern := fmt.Sprintf("%%-%s-%%", dateStr)

	var count int64
	err := db.Model(&models.Ticket{}).
		Where("contact_id = ? AND ticket_id LIKE ?", contactID, pattern).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count) + 1, nil
}

// Ensure math package is used
var _ = math.MaxInt

// ============================================================================
// N8N EMAIL PROCESSOR - NEW ARCHITECTURE
// ============================================================================
// This is the new endpoint for n8n integration where:
// 1. n8n sends only the prompt + email ID
// 2. Our backend calls Gemini API
// 3. Parses response and creates ticket
// 4. Returns response with same ID
// ============================================================================

// N8nEmailRequest represents the request from n8n with prompt and email ID
type N8nEmailRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	ID     string `json:"id" binding:"required"`
}

// N8nEmailResponse represents the response sent back to n8n
type N8nEmailResponse struct {
	Success       bool            `json:"success"`
	ID            string          `json:"id"`
	TicketID      string          `json:"ticket_id,omitempty"`
	DumpedQueryID uint            `json:"dumped_query_id,omitempty"` // New field
	Message       string          `json:"message,omitempty"`
	Error         string          `json:"error,omitempty"`
	Resolution    *ResolutionInfo `json:"resolution,omitempty"`
	Ticket        *TicketInfo     `json:"ticket,omitempty"`
	Warnings      []string        `json:"warnings,omitempty"`
	Debug         *N8nDebugInfo   `json:"debug,omitempty"`
}

// N8nDebugInfo contains debug information for troubleshooting
type N8nDebugInfo struct {
	ParsedEmail    *ParsedEmailDebug `json:"parsed_email,omitempty"`
	ExtractedData  interface{}       `json:"extracted_data,omitempty"`
	GeminiResponse string            `json:"gemini_response,omitempty"`
}

// ParsedEmailDebug contains parsed email info for debugging
type ParsedEmailDebug struct {
	SenderEmail string `json:"sender_email"`
	SenderName  string `json:"sender_name"`
	Subject     string `json:"subject"`
	BodyPreview string `json:"body_preview"`
}

// ProcessEmailHandler is the main endpoint for n8n email processing
// POST /n8n/process-email
// This endpoint:
// 1. Receives prompt + id from n8n
// 2. Calls Gemini 2.5 Flash API to extract structured data
// 3. Parses email content from the prompt
// 4. Uses SmartResolver to create ticket
// 5. Returns response with the same id
func ProcessEmailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req N8nEmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, N8nEmailResponse{
				Success: false,
				ID:      "",
				Error:   "Invalid request: " + err.Error(),
			})
			return
		}

		log.Printf("[ProcessEmail] Received email ID: %s", req.ID)

		// Initialize Gemini service
		geminiService := services.NewGeminiService()

		// Step 1: Parse email components from the prompt
		parsedEmail := geminiService.ParseEmailFromPrompt(req.Prompt)
		if parsedEmail.SenderEmail == "" {
			log.Printf("[ProcessEmail] Could not parse sender email. Dumping query for analysis.")

			// Create DumpedQuery for parsing failure
			dumpedQuery := models.DumpedQuery{
				N8nID:           req.ID,
				SenderEmail:     "unknown@parsing.failed", // Placeholder
				SenderName:      "Parsing Failed",
				Subject:         "Failed to parse email from prompt",
				Body:            req.Prompt, // Save full prompt as body since we couldn't parse it
				FailureReason:   "Could not parse sender email from prompt",
				AIExtractedData: "{}",
				Status:          "OPEN",
			}

			if err := db.Create(&dumpedQuery).Error; err != nil {
				log.Printf("[ProcessEmail] Failed to save dumped query: %v", err)
				c.JSON(http.StatusInternalServerError, N8nEmailResponse{
					Success: false,
					ID:      req.ID,
					Error:   "Parsing failed AND Dump creation failed: " + err.Error(),
				})
				return
			}

			// Return 200 OK with success=true (handled)
			c.JSON(http.StatusOK, N8nEmailResponse{
				Success:       true, // Marked as success so n8n continues
				ID:            req.ID,
				DumpedQueryID: dumpedQuery.ID,
				Message:       "Email parsing failed. Query dumped for manual review.",
				Error:         "Could not parse sender email from prompt",
				Debug: &N8nDebugInfo{
					ParsedEmail: &ParsedEmailDebug{
						SenderEmail: parsedEmail.SenderEmail,
						SenderName:  parsedEmail.SenderName,
						Subject:     parsedEmail.Subject,
						BodyPreview: truncateString(parsedEmail.Body, 200),
					},
				},
			})
			return
		}

		log.Printf("[ProcessEmail] Parsed email - From: %s <%s>, Subject: %s",
			parsedEmail.SenderName, parsedEmail.SenderEmail, parsedEmail.Subject)

		// Step 2: Call Gemini API to extract structured data
		extractedData, err := geminiService.ExtractEmailData(req.Prompt)
		if err != nil {
			log.Printf("[ProcessEmail] Gemini API error: %v", err)
			// Continue with empty extracted data - we can still try to resolve using email/name parsing
			extractedData = &services.ExtractedEmailData{
				PhoneNumbers:  []string{},
				PersonNames:   []string{},
				OrgNames:      []string{},
				ProductHints:  []string{},
				PriorityHints: []string{},
			}
		}

		log.Printf("[ProcessEmail] Extracted - Phones: %v, Names: %v, Orgs: %v, Products: %v, Priority: %v",
			extractedData.PhoneNumbers, extractedData.PersonNames, extractedData.OrgNames,
			extractedData.ProductHints, extractedData.PriorityHints)

		// Step 3: Build SmartTicketRequest from parsed data
		smartReq := &SmartTicketRequest{
			SenderEmail:   parsedEmail.SenderEmail,
			SenderName:    parsedEmail.SenderName,
			Subject:       parsedEmail.Subject,
			Body:          parsedEmail.Body,
			PhoneNumbers:  extractedData.PhoneNumbers,
			PersonNames:   extractedData.PersonNames,
			OrgNames:      extractedData.OrgNames,
			ProductHints:  extractedData.ProductHints,
			PriorityHints: extractedData.PriorityHints,
		}

		// Step 4: Use SmartResolver to create ticket
		resolver := NewSmartResolver(db)
		result := resolver.ResolveAndCreateTicket(smartReq)

		// Step 5: Build response with the same ID
		response := N8nEmailResponse{
			Success:    result.Success,
			ID:         req.ID,
			TicketID:   result.TicketID,
			Resolution: result.Resolution,
			Ticket:     result.Ticket,
			Warnings:   result.Warnings,
		}

		if result.Success {
			response.Message = fmt.Sprintf("Ticket %s created successfully", result.TicketID)
			log.Printf("[ProcessEmail] Successfully created ticket %s for email ID %s", result.TicketID, req.ID)
			c.JSON(http.StatusCreated, response)
		} else {
			// FAILED: Create DumpedQuery entry
			log.Printf("[ProcessEmail] Failed to create ticket. Dumping query for analysis.")

			extractedJson, _ := json.Marshal(extractedData)
			dumpedQuery := models.DumpedQuery{
				N8nID:           req.ID,
				SenderEmail:     parsedEmail.SenderEmail,
				SenderName:      parsedEmail.SenderName,
				Subject:         parsedEmail.Subject,
				Body:            parsedEmail.Body,
				FailureReason:   result.Error,
				AIExtractedData: string(extractedJson),
				Status:          "OPEN",
			}

			if err := db.Create(&dumpedQuery).Error; err != nil {
				log.Printf("[ProcessEmail] Failed to save dumped query: %v", err)
				// If dumping fails, we must return an error
				response.Error = "Ticket creation failed AND Dump creation failed: " + result.Error + " | " + err.Error()
				c.JSON(http.StatusInternalServerError, response)
				return
			}

			log.Printf("[ProcessEmail] Saved dumped query ID: %d", dumpedQuery.ID)

			// SUCCESSFUL DUMP: Return 200 OK so n8n continues
			response.Success = true
			response.DumpedQueryID = dumpedQuery.ID
			response.Message = "Ticket creation failed. Query dumped for manual review."
			response.Error = result.Error // Keep original error for context

			response.Debug = &N8nDebugInfo{
				ParsedEmail: &ParsedEmailDebug{
					SenderEmail: parsedEmail.SenderEmail,
					SenderName:  parsedEmail.SenderName,
					Subject:     parsedEmail.Subject,
					BodyPreview: truncateString(parsedEmail.Body, 200),
				},
				ExtractedData: extractedData,
			}
			log.Printf("[ProcessEmail] Handled failure by dumping query %d for email ID %s", dumpedQuery.ID, req.ID)
			c.JSON(http.StatusOK, response)
		}
	}
}
