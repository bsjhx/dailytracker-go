package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

		// Use PostgreSQL
		db, dbErr = initPostgres()

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
			log.Printf("Database initialized successfully")
		}
	})

	return db, dbErr
}

func initPostgres() (*sql.DB, error) {
	// Build connection string from environment variables
	dbURL := os.Getenv("DB_URL")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if dbURL == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required Postgres configuration: DB_URL, DB_USER, DB_PASSWORD, DB_NAME must all be set")
	}

	// Parse host and port from DB_URL (format: "host:port" or just "host")
	host := dbURL
	port := "5432" // default PostgreSQL port

	// Split host:port if port is provided
	if idx := strings.LastIndex(dbURL, ":"); idx != -1 {
		host = dbURL[:idx]
		port = dbURL[idx+1:]
	}

	// Build PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening Postgres database: %w", err)
	}

	log.Printf("Connected to PostgreSQL database: %s@%s:%s/%s", user, host, port, dbname)
	return db, nil
}

func runMigrations(db *sql.DB) error {
	var driver database.Driver
	var err error

	// Create Postgres driver instance
	driver, err = postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	// Determine migrations path
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}
	log.Printf("Using migrations: %s", migrationsPath)

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
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
