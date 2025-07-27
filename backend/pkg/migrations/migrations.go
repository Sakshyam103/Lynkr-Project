package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	// "os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Name string
	Path string
}

// Migrator handles database migrations
type Migrator struct {
	DB            *sql.DB
	MigrationsDir string
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		DB:            db,
		MigrationsDir: migrationsDir,
	}
}

// EnsureMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) EnsureMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.DB.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of already applied migrations
func (m *Migrator) GetAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := m.DB.Query("SELECT name FROM migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, nil
}

// FindMigrations returns a list of available migrations
func (m *Migrator) FindMigrations() ([]Migration, error) {
	files, err := ioutil.ReadDir(m.MigrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migrations = append(migrations, Migration{
			Name: strings.TrimSuffix(file.Name(), ".sql"),
			Path: filepath.Join(m.MigrationsDir, file.Name()),
		})
	}

	// Sort migrations by name
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

// ApplyMigration applies a single migration
func (m *Migrator) ApplyMigration(migration Migration) error {
	log.Printf("Applying migration: %s", migration.Name)

	// Read migration file
	content, err := ioutil.ReadFile(migration.Path)
	if err != nil {
		return err
	}

	// Extract up migration (everything before "-- Down Migration")
	upMigration := string(content)
	if idx := strings.Index(upMigration, "-- Down Migration"); idx != -1 {
		upMigration = upMigration[:idx]
	}

	// Begin transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// Execute migration
	_, err = tx.Exec(upMigration)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("migration failed: %v", err)
	}

	// Record migration
	_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Printf("Migration applied: %s", migration.Name)
	return nil
}

// Migrate applies all pending migrations
func (m *Migrator) Migrate() error {
	// Ensure migrations table exists
	if err := m.EnsureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Find available migrations
	migrations, err := m.FindMigrations()
	if err != nil {
		return fmt.Errorf("failed to find migrations: %v", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if !applied[migration.Name] {
			if err := m.ApplyMigration(migration); err != nil {
				return err
			}
		} else {
			log.Printf("Migration already applied: %s", migration.Name)
		}
	}

	return nil
}
