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
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./dailytracker.db"
		}
		db, dbErr = sql.Open("sqlite", dbPath)
		if dbErr != nil {
			return
		}

		dbErr = db.Ping()
		if dbErr != nil {
			return
		}

		// Run migrations
		_, dbErr = db.Exec(`
			CREATE TABLE IF NOT EXISTS daily_entries (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				entry_date DATE NOT NULL UNIQUE,
				work_score INTEGER CHECK (work_score BETWEEN 0 AND 5),
				personal_score INTEGER CHECK (personal_score BETWEEN 0 AND 5),
				total INTEGER,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX IF NOT EXISTS idx_entry_date ON daily_entries(entry_date DESC);
		`)
		if dbErr != nil {
			log.Printf("Migration error: %v", dbErr)
		}
	})

	return db, dbErr
}
