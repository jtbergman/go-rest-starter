package mocks

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

// ============================================================================
// Create Test Database
// ============================================================================

// Opens a new test database with all migrations applied
func createTestDB(t *testing.T, dsn string) *sql.DB {
	db := connectToDB(t, dsn)

	// Cleanup the database
	migrations := getMigrations(t, db, "down.sql")
	applyMigrations(t, db, migrations)

	// Create the database
	migrations = getMigrations(t, db, "up.sql")
	applyMigrations(t, db, migrations)

	return db
}

// ============================================================================
// Create Database
// ============================================================================

// Opens a connection to the database and pings it to create a connection
func connectToDB(t *testing.T, dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		t.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

// Gets the migration files in sorted alphabetical order
func getMigrations(t *testing.T, db *sql.DB, suffix string) []string {
	dir := findMigrations(t, db)

	files, err := os.ReadDir(dir)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to read migration files: %v", err)
	}

	var migrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			migrations = append(migrations, filepath.Join(dir, file.Name()))
		}
	}

	sort.Strings(migrations)
	return migrations
}

// Apply the migrations in order
func applyMigrations(t *testing.T, db *sql.DB, migrations []string) {
	for _, file := range migrations {
		script, err := os.ReadFile(file)
		if err != nil {
			db.Close()
			t.Fatalf("Failed to read migration (%s): %v", file, err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			db.Close()
			t.Fatalf("Failed to exec migration (%s): %v", file, err)
		}
	}
}

// Finds the migrations directory
func findMigrations(t *testing.T, db *sql.DB) string {
	absPath, err := os.Getwd()
	if err != nil {
		db.Close()
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	rootDir := string(filepath.Separator)

	for {
		potentialPath := filepath.Join(absPath, "migrations")
		if fileInfo, err := os.Stat(potentialPath); err == nil && fileInfo.IsDir() {
			return potentialPath
		}

		parentPath := filepath.Dir(absPath)
		if parentPath == absPath || parentPath == rootDir {
			db.Close()
			t.Fatalf("Failed to find migrations directory")
		}
		absPath = parentPath
	}
}
