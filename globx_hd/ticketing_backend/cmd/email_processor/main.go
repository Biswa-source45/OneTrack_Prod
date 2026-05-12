package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Chinmay-Globx/ticketing-backend/internal/config"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services/email_service"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Email Processor Service...")

	// Load environment variables from .env.email
	if err := godotenv.Load(".env.email"); err != nil {
		// Try loading from parent directory
		if err := godotenv.Load("../../.env.email"); err != nil {
			log.Printf("Warning: Error loading .env.email file: %v", err)
			log.Println("Make sure .env.email exists with EMAIL_USERNAME and EMAIL_PASSWORD set")
			// Continue anyway, as env vars might be set directly
		}
	}

	// Parse command line flags
	runOnce := flag.Bool("once", false, "Run the processor once and exit")
	intervalMinutes := flag.Int("interval", 5, "Minutes between email processing cycles")
	flag.Parse()

	// Initialize configuration
	cfg := config.LoadConfigFromEnvOrDefault()

	// Connect to database
	err := config.ConnectDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Determine upload directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join(".", "uploads")
	}

	// Create email service
	service, err := email_service.NewEmailService(config.DB, uploadDir)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// If run-once flag is set, process emails once and exit
	if *runOnce {
		log.Println("Running single processing cycle...")
		if err := service.ProcessEmails(); err != nil {
			log.Printf("Error processing emails: %v", err)
			os.Exit(1)
		}
		log.Println("Processing complete")
		return
	}

	// Set up ticker for regular processing
	interval := time.Duration(*intervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Email processor running with %d minute interval", *intervalMinutes)

	// Process immediately on startup
	if err := service.ProcessEmails(); err != nil {
		log.Printf("Error processing emails: %v", err)
	}

	// Process on each tick
	for range ticker.C {
		log.Printf("Processing emails at %s", time.Now().Format(time.RFC3339))
		if err := service.ProcessEmails(); err != nil {
			log.Printf("Error processing emails: %v", err)
		}
	}
}
