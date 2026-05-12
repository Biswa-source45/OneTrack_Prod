package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
}

// DB is the package-level DB handle used by handlers
var DB *gorm.DB

// LoadConfigFromEnvOrDefault reads env vars or returns defaults
func LoadConfigFromEnvOrDefault() *Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// default to your provided creds with IST timezone
		dsn = "host=localhost user=postgres password=postgres dbname=globx_hd port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	}
	sa := os.Getenv("SERVER_ADDR")
	if sa == "" {
		sa = ":8080"
	}
	return &Config{
		DatabaseURL: dsn,
		ServerAddr:  sa,
	}
}

// ConnectDatabase opens the DB connection and sets package DB
func ConnectDatabase(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// optional: tune connection pool if you need the *sql.DB
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(25)
		// sqlDB.SetConnMaxLifetime(time.Minute * 5)
	}

	DB = db
	fmt.Println("✅ Connected to DB")
	return nil
}
