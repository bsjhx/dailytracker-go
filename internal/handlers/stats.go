package handlers

import (
	"dailytracker/internal/middleware"
	"dailytracker/internal/models"
	"dailytracker/internal/repository"
	"encoding/json"
	"net/http"
	"time"
)

// WeeklyStats handles GET /api/stats/weekly
func WeeklyStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	// Get user ID from context
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date >= CURRENT_DATE - INTERVAL '7 days' AND user_id = $1
		ORDER BY entry_date DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []models.DailyEntry
	var sum int
	for rows.Next() {
		var entry models.DailyEntry
		var entryDate time.Time
		err := rows.Scan(&entry.ID, &entryDate, &entry.WorkScore, &entry.PersonalScore, &entry.Total)
		if err != nil {
			continue
		}
		entry.EntryDate = entryDate.Format("2006-01-02")
		entries = append(entries, entry)
		if entry.Total != nil {
			sum += *entry.Total
		}
	}

	if entries == nil {
		entries = []models.DailyEntry{}
	}

	average := 0.0
	if len(entries) > 0 {
		average = float64(sum) / float64(len(entries))
	}

	stats := models.WeeklyStats{
		Average: average,
		Entries: entries,
	}

	json.NewEncoder(w).Encode(stats)
}
