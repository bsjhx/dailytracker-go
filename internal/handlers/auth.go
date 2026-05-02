package handlers

import (
	"crypto/rand"
	"dailytracker/internal/models"
	"dailytracker/internal/repository"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SessionStore manages active sessions
type SessionStore struct {
	sessions map[string]*models.Session
	mu       sync.RWMutex
}

var sessionStore = &SessionStore{
	sessions: make(map[string]*models.Session),
}

// GetSessionStore returns the global session store (for middleware)
func GetSessionStore() *SessionStore {
	return sessionStore
}

// generateSessionID creates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AddSession adds a new session to the store
func (s *SessionStore) AddSession(userID int, username string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		Username:  username,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
	}
	s.sessions[sessionID] = session

	return sessionID, nil
}

// GetSession retrieves a session by ID
func (s *SessionStore) GetSession(sessionID string) (*models.Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, false
	}

	return session, true
}

// RemoveSession removes a session from the store
func (s *SessionStore) RemoveSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
}

// CleanExpiredSessions removes expired sessions (should be called periodically)
func (s *SessionStore) CleanExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	// Fetch user from database
	var userID int
	var username string
	var passwordHash string
	err = db.QueryRow(`
		SELECT id, username, password_hash
		FROM users
		WHERE username = $1
	`, req.Username).Scan(&userID, &username, &passwordHash)

	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// Create session
	sessionID, err := sessionStore.AddSession(userID, username)
	if err != nil {
		http.Error(w, `{"error":"Failed to create session"}`, http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	// Return user info
	response := models.UserResponse{
		ID:       userID,
		Username: username,
	}
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Remove session from store
		sessionStore.RemoveSession(cookie.Value)
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Delete cookie
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}

// CurrentUser returns the current user's information
func CurrentUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Get session from context (set by middleware)
	session, ok := r.Context().Value("session").(*models.Session)
	if !ok {
		http.Error(w, `{"error":"Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	response := models.UserResponse{
		ID:       session.UserID,
		Username: session.Username,
	}
	json.NewEncoder(w).Encode(response)
}

// CreateUser creates a new user (for manual user creation)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		http.Error(w, `{"error":"Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	db, err := repository.GetDB()
	if err != nil {
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}

	// Insert user
	log.Println("Insertign user:", string(passwordHash))
	result, err := db.Exec(`
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
	`, req.Username, string(passwordHash))

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			http.Error(w, `{"error":"Username already exists"}`, http.StatusConflict)
			return
		}
		log.Println("Error inserting user:", err)
		http.Error(w, `{"error":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()

	response := models.UserResponse{
		ID:       int(userID),
		Username: req.Username,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// StartSessionCleaner starts a background goroutine to clean expired sessions
func StartSessionCleaner() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			sessionStore.CleanExpiredSessions()
		}
	}()
}
