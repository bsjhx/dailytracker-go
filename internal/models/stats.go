package models

// WeeklyStats represents weekly statistics
type WeeklyStats struct {
	Average float64      `json:"average"`
	Entries []DailyEntry `json:"entries"`
}
