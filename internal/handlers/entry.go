package handlers

import (
	"dailytracker/internal/middleware"
	"dailytracker/internal/models"
	"dailytracker/internal/repository"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Entry handles GET /api/entries/:date and PUT /api/entries/:date
func Entry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract date from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"Date parameter required"}`, http.StatusBadRequest)
		return
	}
	date := parts[len(parts)-1]

	db, err := repository.GetDB()
	if err != nil {
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		getEntry(w, db, date, r)
	case "PUT":
		updateEntry(w, r, db, date)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func getEntry(w http.ResponseWriter, db *sql.DB, date string, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	var entry models.DailyEntry
	var entryDate time.Time
	err := db.QueryRow(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date = $1 AND user_id = $2
	`, date, userID).Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Entry not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch entry"}`, http.StatusInternalServerError)
		return
	}

	entry.EntryDate = entryDate.Format("2006-01-02")
	json.NewEncoder(w).Encode(entry)
}

func updateEntry(w http.ResponseWriter, r *http.Request, db *sql.DB, date string) {
	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	var req models.UpdateEntryRequest
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

	_, err := db.Exec(`
		UPDATE daily_entries
		SET work_score = $1, personal_score = $2, total = $3, updated_at = CURRENT_TIMESTAMP
		WHERE entry_date = $4 AND user_id = $5
	`, req.WorkScore, req.PersonalScore, total, date, userID)

	if err != nil {
		http.Error(w, `{"error":"Failed to update entry"}`, http.StatusInternalServerError)
		return
	}

	// Fetch the updated entry
	var entry models.DailyEntry
	var entryDate time.Time
	err = db.QueryRow(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date = $1 AND user_id = $2
	`, date, userID).Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Entry not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch updated entry"}`, http.StatusInternalServerError)
		return
	}

	entry.EntryDate = entryDate.Format("2006-01-02")
	json.NewEncoder(w).Encode(entry)
}
