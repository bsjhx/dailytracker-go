package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type WeeklyStats struct {
	Average float64      `json:"average"`
	Entries []DailyEntry `json:"entries"`
}

func Stats(w http.ResponseWriter, r *http.Request) {
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

	db, err := GetDB()
	if err != nil {
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	// Get user ID from context
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		http.Error(w, `{"error":"User not found in context"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT id, entry_date, work_score, personal_score, total
		FROM daily_entries
		WHERE entry_date >= date('now', '-7 days') AND user_id = ?
		ORDER BY entry_date DESC
	`, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch stats"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []DailyEntry
	var sum int
	for rows.Next() {
		var entry DailyEntry
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
		entries = []DailyEntry{}
	}

	average := 0.0
	if len(entries) > 0 {
		average = float64(sum) / float64(len(entries))
	}

	stats := WeeklyStats{
		Average: average,
		Entries: entries,
	}

	json.NewEncoder(w).Encode(stats)
}
