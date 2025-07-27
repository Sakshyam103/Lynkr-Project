package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"lynkr/pkg/migrations"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var DB *sql.DB

// Config holds database configuration
type Config struct {
	DBPath        string
	MigrationsDir string
}

// Initialize sets up the database connection and runs migrations
func Initialize(config Config) error {
	// Ensure directory exists
	dir := filepath.Dir(config.DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open database connection
	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	log.Println("Database connection established")

	// Run migrations if migrations directory is provided
	if config.MigrationsDir != "" {
		log.Println("Running database migrations")
		migrator := migrations.NewMigrator(db, config.MigrationsDir)
		if err := migrator.Migrate(); err != nil {
			log.Printf("Migration error: %v", err)
			return err
		}
		log.Println("Database migrations completed")
	}

	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
