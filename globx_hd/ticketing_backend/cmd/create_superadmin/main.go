// One-off: creates the Super Admin who triages feedback. Run once on the
// server from ticketing_backend/ after migration 010:
//
//	SA_USERNAME=superadmin SA_EMAIL=you@company.com SA_PASSWORD='temp-pass' go run ./cmd/create_superadmin
//
// Optional: SA_FIRST_NAME (default "Super"), SA_EMPLOYEE_ID (default "SUPERADMIN").
// The account must change its password on first login.
package main

import (
	"log"
	"os"

	"github.com/Chinmay-Globx/ticketing-backend/internal/config"
	"github.com/Chinmay-Globx/ticketing-backend/internal/models"
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/joho/godotenv"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	_ = godotenv.Load(".env.email")
	username, email, password := os.Getenv("SA_USERNAME"), os.Getenv("SA_EMAIL"), os.Getenv("SA_PASSWORD")
	if username == "" || email == "" || len(password) < 8 {
		log.Fatal("set SA_USERNAME, SA_EMAIL and SA_PASSWORD (min 8 chars)")
	}

	if err := config.ConnectDatabase(config.LoadConfigFromEnvOrDefault().DatabaseURL); err != nil {
		log.Fatalf("db connect: %v", err)
	}
	db := config.DB

	var role models.MasterRole
	if err := db.Where("LOWER(TRIM(role_name)) = ?", "superadmin").First(&role).Error; err != nil {
		log.Fatal("superadmin role not found — run migrations/010_create_feedback.sql first")
	}
	var n int64
	db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&n)
	if n > 0 {
		log.Fatalf("a user with username %q or email %q already exists — not touching it", username, email)
	}
	// designation_id is NOT NULL; any existing designation will do for this account.
	var desig models.MasterUserDesignation
	if err := db.Order("id").First(&desig).Error; err != nil {
		log.Fatal("no user designation exists — create one in Master Data first")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	u := models.User{
		EmployeeID:    env("SA_EMPLOYEE_ID", "SUPERADMIN"),
		Username:      username,
		PasswordHash:  hash,
		FirstName:     env("SA_FIRST_NAME", "Super"),
		LastName:      "Admin",
		Email:         email,
		DesignationID: desig.ID,
		RoleID:        role.ID,
		FirstLogin:    true,
	}
	if err := db.Create(&u).Error; err != nil {
		log.Fatalf("create user: %v", err)
	}
	log.Printf("✅ Super Admin created: id=%d username=%s — log in at /login/user and set a new password", u.ID, u.Username)
}
