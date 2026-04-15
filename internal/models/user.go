package models

import "time"

// User represents a user account
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID        string
	UserID    int
	Username  string
	ExpiresAt time.Time
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse represents user information returned to client
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// CreateUserRequest represents the user creation payload
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
