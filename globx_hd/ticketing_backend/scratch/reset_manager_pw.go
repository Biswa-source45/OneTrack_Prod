package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint
	Username     string
	PasswordHash string `gorm:"column:password_hash"`
	FirstLogin   bool   `gorm:"column:first_login"`
}

func main() {
	_ = godotenv.Load(".env.email")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	password := "Password@123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var user User
	// Let's reset chinmay@manager
	if err := db.Table("users").Where("username = ?", "chinmay@manager").First(&user).Error; err != nil {
		log.Fatalf("User chinmay@manager not found: %v", err)
	}

	user.PasswordHash = string(hash)
	user.FirstLogin = false // bypass any first login reset requirement for ease of test
	if err := db.Table("users").Save(&user).Error; err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	fmt.Printf("Successfully updated password for %s to %s\n", user.Username, password)
}
