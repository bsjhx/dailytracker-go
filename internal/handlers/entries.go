package handlers

import (
	"dailytracker/internal/middleware"
	"dailytracker/internal/models"
	"dailytracker/internal/repository"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ListEntries handles GET /api/entries
func ListEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		log.Printf("GetDB error: %v", err)
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		getEntries(w, db, r)
	case "POST":
		createEntry(w, r, db)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getEntries(w http.ResponseWriter, db *sql.DB, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE user_id = ?
		ORDER BY entry_date DESC
		LIMIT 30
	`, userID)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, `{"error":"Failed to fetch entries"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []models.DailyEntry
	for rows.Next() {
		var entry models.DailyEntry
		var entryDate time.Time
		err := rows.Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)
		if err != nil {
			continue
		}
		entry.EntryDate = entryDate.Format("2006-01-02")
		entries = append(entries, entry)
	}

	if entries == nil {
		entries = []models.DailyEntry{}
	}

	json.NewEncoder(w).Encode(entries)
}

func createEntry(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	var req models.CreateEntryRequest
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
		INSERT INTO daily_entries (entry_date, work_score, personal_score, total, user_id)
		VALUES (?, ?, ?, ?, ?)
	`, req.EntryDate, req.WorkScore, req.PersonalScore, total, userID)

	if err != nil {
		http.Error(w, `{"error":"Failed to create entry"}`, http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Fetch the created entry
	var entry models.DailyEntry
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
