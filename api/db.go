package api

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbErr  error
	once   sync.Once
)

func GetDB() (*sql.DB, error) {
	once.Do(func() {
		log.Printf("GettingDB")
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./dailytracker.db"
		}
		db, dbErr = sql.Open("sqlite", dbPath)
		if dbErr != nil {
			log.Printf("Error creating DB")
			return
		}

		dbErr = db.Ping()
		if dbErr != nil {
			return
		}

		// Run migrations
		dbErr = runMigrations(db)
		if dbErr != nil {
			log.Printf("Migration error: %v", dbErr)
		} else {
			log.Printf("OK")
		}
	})

	return db, dbErr
}

func runMigrations(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Check if we need to create a default user (for fresh databases)
	var userCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&userCount)
	if err != nil {
		return err
	}

	if userCount == 0 {
		log.Printf("No users found, creating default admin user...")
		// We'll create user via the create-user script or API
		// For now, just log that manual user creation is needed
		log.Printf("IMPORTANT: No users exist. Please create a user manually:")
		log.Printf("  Option 1: Use ./scripts/create-user.sh username password")
		log.Printf("  Option 2: Use API: POST /api/users/create (see docs)")
	}

	// Create daily_entries table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS daily_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			entry_date DATE NOT NULL,
			work_score INTEGER CHECK (work_score BETWEEN 0 AND 5),
			personal_score INTEGER CHECK (personal_score BETWEEN 0 AND 5),
			total INTEGER,
			user_id INTEGER REFERENCES users(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Check if user_id column exists (for migrating old databases)
	var columnExists int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM pragma_table_info('daily_entries') WHERE name='user_id'
	`).Scan(&columnExists)
	if err != nil {
		return err
	}

	// If user_id doesn't exist, we need to migrate old database
	if columnExists == 0 {
		log.Printf("Migrating daily_entries table to add user_id column...")

		// Ensure at least one user exists for migration
		var userCountForMigration int
		err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&userCountForMigration)
		if err != nil {
			return err
		}

		if userCountForMigration == 0 {
			log.Printf("Creating default user for migration...")
			// Create a migration user - they should change this password
			_, err = db.Exec(`
				INSERT INTO users (id, username, password_hash)
				VALUES (1, 'admin', 'CHANGE_PASSWORD_IMMEDIATELY')
			`)
			if err != nil {
				return err
			}
		}

		// Since SQLite doesn't support adding columns with foreign keys directly,
		// we need to recreate the table
		_, err = db.Exec(`
			-- Create new table with user_id
			CREATE TABLE daily_entries_new (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				entry_date DATE NOT NULL,
				work_score INTEGER CHECK (work_score BETWEEN 0 AND 5),
				personal_score INTEGER CHECK (personal_score BETWEEN 0 AND 5),
				total INTEGER,
				user_id INTEGER REFERENCES users(id),
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);

			-- Copy existing data (assign to user_id = 1)
			INSERT INTO daily_entries_new (id, entry_date, work_score, personal_score, total, user_id, created_at, updated_at)
			SELECT id, entry_date, work_score, personal_score, total, 1, created_at, updated_at
			FROM daily_entries;

			-- Drop old table
			DROP TABLE daily_entries;

			-- Rename new table
			ALTER TABLE daily_entries_new RENAME TO daily_entries;
		`)
		if err != nil {
			return err
		}

		log.Printf("Migration complete: added user_id column and assigned existing entries to user_id=1")
	}

	// Create indexes
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_entry_date ON daily_entries(entry_date DESC);
		CREATE INDEX IF NOT EXISTS idx_user_entries ON daily_entries(user_id, entry_date DESC);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_date ON daily_entries(user_id, entry_date);
	`)
	if err != nil {
		return err
	}

	return nil
}
