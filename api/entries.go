package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type DailyEntry struct {
	ID            int        `json:"id"`
	EntryDate     string     `json:"entry_date"`
	WorkScore     *int       `json:"work_score"`
	PersonalScore *int       `json:"personal_score"`
	Total         *int       `json:"total"`
}

type CreateEntryRequest struct {
	EntryDate     string `json:"entry_date"`
	WorkScore     *int   `json:"work_score"`
	PersonalScore *int   `json:"personal_score"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	db, err := GetDB()
	if err != nil {
		log.Printf("GetDB error: %v", err)
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		getEntries(w, db)
	case "POST":
		createEntry(w, r, db)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getEntries(w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		ORDER BY entry_date DESC
		LIMIT 30
	`)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, `{"error":"Failed to fetch entries"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []DailyEntry
	for rows.Next() {
		var entry DailyEntry
		var entryDate time.Time
		err := rows.Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)
		if err != nil {
			continue
		}
		entry.EntryDate = entryDate.Format("2006-01-02")
		entries = append(entries, entry)
	}

	if entries == nil {
		entries = []DailyEntry{}
	}

	json.NewEncoder(w).Encode(entries)
}

func createEntry(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Calculate total
	workScore := 0
	personalScore := 0
	if req.WorkScore != nil {
		workScore = *req.WorkScore
	}
	if req.PersonalScore != nil {
		personalScore = *req.PersonalScore
	}
	total := workScore + personalScore

	result, err := db.Exec(`
		INSERT INTO daily_entries (entry_date, work_score, personal_score, total)
		VALUES (?, ?, ?, ?)
	`, req.EntryDate, req.WorkScore, req.PersonalScore, total)

	if err != nil {
		http.Error(w, `{"error":"Failed to create entry"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Fetch the created entry
	var entry DailyEntry
	var entryDate time.Time
	err = db.QueryRow(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE id = ?
	`, id).Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)

	if err != nil {
		http.Error(w, `{"error":"Failed to fetch created entry"}`, http.StatusInternalServerError)
		return
	}

	entry.EntryDate = entryDate.Format("2006-01-02")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}
