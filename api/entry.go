package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type UpdateEntryRequest struct {
	WorkScore     *int `json:"work_score"`
	PersonalScore *int `json:"personal_score"`
}

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

	db, err := GetDB()
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
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	var entry DailyEntry
	var entryDate time.Time
	err := db.QueryRow(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date = ? AND user_id = ?
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
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	var req UpdateEntryRequest
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
		SET work_score = ?, personal_score = ?, total = ?, updated_at = CURRENT_TIMESTAMP
		WHERE entry_date = ? AND user_id = ?
	`, req.WorkScore, req.PersonalScore, total, date, userID)

	if err != nil {
		http.Error(w, `{"error":"Failed to update entry"}`, http.StatusInternalServerError)
		return
	}

	// Fetch the updated entry
	var entry DailyEntry
	var entryDate time.Time
	err = db.QueryRow(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date = ? AND user_id = ?
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
