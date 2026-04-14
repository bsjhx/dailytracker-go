package api

import (
	"context"
	"net/http"
)

// SessionAuthMiddleware validates the session and adds user context to the request
func SessionAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
	session, ok := r.Context().Value("session").(*Session)
	if !ok {
		return 0, false
	}
	return session.UserID, true
}
