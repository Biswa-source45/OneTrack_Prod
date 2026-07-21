package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint
	Username  string
	Email     string
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	RoleID    uint   `gorm:"column:role_id"`
}

type Contact struct {
	ID        uint
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     *string
}

func main() {
	// Try loading from .env.email
	_ = godotenv.Load(".env.email")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Println("--- USERS ---")
	var users []User
	if err := db.Table("users").Find(&users).Error; err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	for _, u := range users {
		fmt.Printf("ID: %d | Username: %s | Email: %s | Name: %s %s | RoleID: %d\n", u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.RoleID)
	}

	fmt.Println("--- CONTACTS ---")
	var contacts []Contact
	if err := db.Table("contacts").Find(&contacts).Error; err != nil {
		log.Fatalf("Failed to query contacts: %v", err)
	}
	for _, c := range contacts {
		email := "N/A"
		if c.Email != nil {
			email = *c.Email
		}
		fmt.Printf("ID: %d | Name: %s %s | Email: %s\n", c.ID, c.FirstName, c.LastName, email)
	}
}
