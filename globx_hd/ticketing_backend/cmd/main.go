package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	// "path/filepath" // DISABLED: Legacy IMAP email processor
	// "time" // DISABLED: Legacy IMAP email processor

	"github.com/Chinmay-Globx/ticketing-backend/internal/config"
	"github.com/Chinmay-Globx/ticketing-backend/internal/routes"
	"github.com/Chinmay-Globx/ticketing-backend/internal/services/email_service"
	"github.com/Chinmay-Globx/ticketing-backend/internal/websocket"

	// "github.com/Chinmay-Globx/ticketing-backend/internal/services/email_service" // DISABLED: Using n8n instead
	"github.com/Chinmay-Globx/ticketing-backend/internal/utils"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env.email FIRST so DATABASE_URL is available
	if err := godotenv.Load(".env.email"); err != nil {
		log.Printf("⚠️  Warning: Could not load .env.email file: %v", err)
		log.Println("Using default configuration or environment variables")
	} else {
		log.Println("✅ Configuration loaded from .env.email")
	}

	cfg := config.LoadConfigFromEnvOrDefault()

	if err := config.ConnectDatabase(cfg.DatabaseURL); err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	// Check notification system
	utils.CheckNotificationTables(config.DB)

	// Initialize WebSocket hub for real-time notifications
	log.Println("🔌 Initializing WebSocket Hub...")
	hub := websocket.NewHub()
	go hub.Run()
	log.Println("✅ WebSocket Hub started successfully")

	// DISABLED: Legacy IMAP email processor - now using n8n + Gemini integration
	// The n8n workflow handles email fetching and sends to POST /n8n/process-email
	// go startEmailProcessor()

	r := routes.SetupRouter(config.DB, hub) // pass global DB and WebSocket hub
	addr := cfg.ServerAddr
	if v := os.Getenv("SERVER_ADDR"); v != "" {
		addr = v
	}

	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

// startEmailProcessor initializes and runs the email processor service
func startEmailProcessor() {
	log.Println("📧 Initializing Email Processor Service...")

	// Load environment variables from .env.email
	if err := godotenv.Load(".env.email"); err != nil {
		log.Printf("⚠️  Warning: Could not load .env.email file: %v", err)
		log.Println("Email processor will use environment variables if set")
	}

	// Check if email service is configured
	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	if emailUsername == "" || emailPassword == "" {
		log.Println("⚠️  Email processor not configured (EMAIL_USERNAME or EMAIL_PASSWORD missing)")
		log.Println("Email processor will not start. To enable, configure .env.email file")
		return
	}

	// Determine upload directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join(".", "uploads")
	}

	// Create email service
	service, err := email_service.NewEmailService(config.DB, uploadDir)
	if err != nil {
		log.Printf("❌ Failed to initialize email service: %v", err)
		log.Println("Email processor will not start")
		return
	}

	// Get processing interval from environment or use default (5 minutes)
	intervalMinutes := 5
	if envInterval := os.Getenv("EMAIL_POLL_INTERVAL_MINUTES"); envInterval != "" {
		if parsed, err := time.ParseDuration(envInterval + "m"); err == nil {
			intervalMinutes = int(parsed.Minutes())
		}
	}

	interval := time.Duration(intervalMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("✅ Email processor started successfully (polling every %d minutes)", intervalMinutes)

	// Process immediately on startup
	log.Println("📬 Processing emails on startup...")
	if err := service.ProcessEmails(); err != nil {
		log.Printf("❌ Error processing emails on startup: %v", err)
	}

	// Process on each tick
	for range ticker.C {
		log.Printf("📬 Processing emails at %s", time.Now().Format("2006-01-02 15:04:05"))
		if err := service.ProcessEmails(); err != nil {
			log.Printf("❌ Error processing emails: %v", err)
		}
	}
}
