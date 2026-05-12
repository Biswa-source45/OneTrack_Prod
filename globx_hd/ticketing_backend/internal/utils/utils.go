package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateUnique3Digit returns a unique 3-digit string (zero-padded).
// It uses a DB check to avoid collisions.
func GenerateUnique3Digit(db *gorm.DB) (string, error) {
	// try many times (space: 000..999)
	for i := 0; i < 2000; i++ {
		code := fmt.Sprintf("%03d", rand.Intn(1000))
		var a models.Account
		// lock the row scan to prevent race? using plain check is OK for now.
		// We check existence by searching for a matching code
		err := db.Where("customer_code = ?", code).Take(&a).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return code, nil
			}
			return "", err
		}
		// found -> try again
	}
	return "", fmt.Errorf("failed to generate unique 3-digit code")
}

// GenerateUniqueCustomerCode returns a unique 3-digit string (zero-padded).
// It checks for uniqueness across both accounts and contacts to avoid collisions.
func GenerateUniqueCustomerCode(db *gorm.DB) (string, error) {
	// try many times (space: 000..999)
	for i := 0; i < 2000; i++ {
		code := fmt.Sprintf("%03d", rand.Intn(1000))
		
		// Check if code exists in accounts table
		var account models.Account
		err := db.Where("customer_code = ?", code).Take(&account).Error
		if err == nil {
			// found in accounts -> try again
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return "", err
		}
		
		// Check if code exists in contacts table
		var contact models.Contact
		err = db.Where("customer_code = ?", code).Take(&contact).Error
		if err == nil {
			// found in contacts -> try again
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return "", err
		}
		
		// Not found in either table -> code is unique
		return code, nil
	}
	return "", fmt.Errorf("failed to generate unique customer code")
}

// GetNextTicketSequence returns the next sequence number for tickets for an account in a given year
func GetNextTicketSequence(db *gorm.DB, accountID uint, date time.Time) (int, error) {
	var count int64
	var err error
	
	// Extract year from the date for annual sequence counting
	year := date.Year()
	
	if accountID == 0 {
		// For individual contacts (account_id is NULL)
		err = db.Model(&models.Ticket{}).
			Where("account_id IS NULL AND EXTRACT(YEAR FROM created_at) = ?", year).
			Count(&count).Error
	} else {
		// For contacts with accounts
		err = db.Model(&models.Ticket{}).
			Where("account_id = ? AND EXTRACT(YEAR FROM created_at) = ?", accountID, year).
			Count(&count).Error
	}
	
	if err != nil {
		return 0, err
	}
	return int(count) + 1, nil
}

// FormatDateForTicketID formats a date to DDMMYY format for ticket IDs
func FormatDateForTicketID(date time.Time) string {
	return date.Format("020106") // DDMMYY format
}

// FormatTicketID creates the ticket_id string as per requirements
func FormatTicketID(customerCode string, dateStr string, seq int) string {
	return fmt.Sprintf("%s-%s-%04d", customerCode, dateStr, seq)
}

// UintToString converts uint to string for activity logging
func UintToString(value uint) string {
	return fmt.Sprintf("%d", value)
}

// GetUserName returns the full name of a user by ID for activity logging
func GetUserName(db *gorm.DB, userID uint) (string, error) {
	var user models.User
	err := db.First(&user, userID).Error
	if err != nil {
		return "", err
	}
	return user.FirstName + " " + user.LastName, nil
}
