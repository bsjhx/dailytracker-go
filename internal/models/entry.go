package models

// DailyEntry represents a daily tracking entry
type DailyEntry struct {
	ID            int    `json:"id"`
	EntryDate     string `json:"entry_date"`
	WorkScore     *int   `json:"work_score"`
	PersonalScore *int   `json:"personal_score"`
	Total         *int   `json:"total"`
	UserID        int    `json:"-"` // Don't expose in JSON
}

// CreateEntryRequest represents the request to create a new entry
type CreateEntryRequest struct {
	EntryDate     string `json:"entry_date"`
	WorkScore     *int   `json:"work_score"`
	PersonalScore *int   `json:"personal_score"`
}

// UpdateEntryRequest represents the request to update an entry
type UpdateEntryRequest struct {
	WorkScore     *int `json:"work_score"`
	PersonalScore *int `json:"personal_score"`
}
