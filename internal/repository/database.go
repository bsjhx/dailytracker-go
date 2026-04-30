package repository

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

var (
	db    *sql.DB
	dbErr error
	once  sync.Once
)

// GetDB returns the database instance, initializing it if necessary
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		log.Printf("Initializing database connection...")
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./dailytracker.db"
		}

		// Create database file if it doesn't exist
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			// Ensure parent directory exists
			dir := filepath.Dir(dbPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Printf("Error creating database directory: %v", err)
				dbErr = err
				return
			}

			file, err := os.Create(dbPath)
			if err != nil {
				log.Printf("Error creating database file: %v", err)
				dbErr = err
				return
			}
			file.Close()
			log.Printf("Created database file: %s", dbPath)
		}

		db, dbErr = sql.Open("sqlite", dbPath)
		if dbErr != nil {
			log.Printf("Error opening database: %v", dbErr)
			return
		}

		dbErr = db.Ping()
		if dbErr != nil {
			log.Printf("Error pinging database: %v", dbErr)
			return
		}

		// Run migrations
		dbErr = runMigrations(db)
		if dbErr != nil {
			log.Printf("Migration error: %v", dbErr)
		} else {
			log.Printf("Database initialized successfully")
		}
	})

	return db, dbErr
}

func runMigrations(db *sql.DB) error {
	// Create driver instance for golang-migrate
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	// Run migrations - this is idempotent, safe to run multiple times
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Printf("✅ Database schema is up to date")
	} else {
		log.Printf("✅ Migrations applied successfully")
	}

	return nil
}
