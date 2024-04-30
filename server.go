package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("faq.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db.AutoMigrate(&FAQ{}, &Answer{})
}

// LoadTranslations loads the translation files
func LoadTranslations() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("./locales/en-US.json")
	bundle.MustLoadMessageFile("./locales/tr-TR.json")
}


// Middleware for logging requests
func logRequests(c *fiber.Ctx) error {
	log.Printf("Request: %s, IP: %s, Method: %s ", c.OriginalURL(), c.IP(), c.Method())
	return c.Next()
}

// Function to create a backup of the database
func createBackup() error {
	// Generate a filename with the current timestamp
	backupFileName := "backup.db"

	// Open the original database file
	originalDB, err := os.Open("faq.db")
	if err != nil {
		return err
	}
	defer originalDB.Close()

	// Create the backup file
	backupDB, err := os.Create(backupFileName)
	if err != nil {
		return err
	}
	defer backupDB.Close()

	// Copy the contents of the original database to the backup file
	_, err = io.Copy(backupDB, originalDB)
	if err != nil {
		return err
	}

	fmt.Printf("Backup updated: %s\n", backupFileName)
	return nil
}
