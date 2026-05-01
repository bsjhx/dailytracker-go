package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbErr  error
	once   sync.Once
	dbType string // "sqlite" or "postgres"
)

// GetDB returns the database instance, initializing it if necessary
func GetDB() (*sql.DB, error) {
	once.Do(func() {
		log.Printf("Initializing database connection...")

		// Determine environment
		env := os.Getenv("ENV")
		if env == "" {
			env = "dev"
		}
		log.Printf("Environment: %s", env)

		if env == "prod" {
			// Production: Use Postgres
			dbType = "postgres"
			db, dbErr = initPostgres()
		} else {
			// Development: Use SQLite
			dbType = "sqlite"
			db, dbErr = initSQLite()
		}

		if dbErr != nil {
			log.Printf("Error initializing database: %v", dbErr)
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
			log.Printf("Database initialized successfully (%s)", dbType)
		}
	})

	return db, dbErr
}

func initSQLite() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dailytracker.db"
	}

	// Create database file if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Ensure parent directory exists
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("error creating database directory: %w", err)
		}

		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("error creating database file: %w", err)
		}
		file.Close()
		log.Printf("Created database file: %s", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %w", err)
	}

	log.Printf("Connected to SQLite database: %s", dbPath)
	return db, nil
}

func initPostgres() (*sql.DB, error) {
	// Build connection string from environment variables
	host := os.Getenv("DB_URL")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required Postgres configuration: DB_URL, DB_USER, DB_PASSWORD, DB_NAME must all be set in production mode")
	}

	// Build PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening Postgres database: %w", err)
	}

	log.Printf("Connected to PostgreSQL database: %s@%s/%s", user, host, dbname)
	return db, nil
}

func runMigrations(db *sql.DB) error {
	var driver database.Driver
	var databaseName string
	var err error

	// Create appropriate driver instance
	if dbType == "postgres" {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
		databaseName = "postgres"
	} else {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
		databaseName = "sqlite3"
	}

	if err != nil {
		return err
	}

	// Determine migrations path - try database-specific folder first
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		// Try database-specific folder first
		dbSpecificPath := fmt.Sprintf("file://migrations/%s", dbType)
		if _, err := os.Stat(fmt.Sprintf("migrations/%s", dbType)); err == nil {
			migrationsPath = dbSpecificPath
			log.Printf("Using database-specific migrations: %s", migrationsPath)
		} else {
			// Fall back to generic migrations folder
			migrationsPath = "file://migrations"
			log.Printf("Using generic migrations: %s", migrationsPath)
		}
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		databaseName,
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
