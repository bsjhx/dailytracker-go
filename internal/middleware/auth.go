package middleware

import (
	"context"
	"dailytracker/internal/models"
	"net/http"
)

// SessionStore interface for session management
type SessionStore interface {
	GetSession(sessionID string) (*models.Session, bool)
}

var sessionStore SessionStore

// SetSessionStore sets the session store to be used by the middleware
func SetSessionStore(store SessionStore) {
	sessionStore = store
}

// SessionAuth validates the session and adds user context to the request
func SessionAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, `{"error":"Not authenticated"}`, http.StatusUnauthorized)
			return
		}

		// Validate session
		session, exists := sessionStore.GetSession(cookie.Value)
		if !exists {
			// Session invalid or expired, clear the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
			http.Error(w, `{"error":"Session expired"}`, http.StatusUnauthorized)
			return
		}

		// Add session to request context
		ctx := context.WithValue(r.Context(), "session", session)
		next(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
	session, ok := r.Context().Value("session").(*models.Session)
	if !ok {
		return 0, false
	}
	return session.UserID, true
}
