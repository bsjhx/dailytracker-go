package main

import (
	"dailytracker/api"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start session cleaner
	api.StartSessionCleaner()

	// Authentication routes (no auth required)
	http.HandleFunc("/api/auth/login", api.LoginHandler)
	http.HandleFunc("/api/auth/logout", api.LogoutHandler)
	http.HandleFunc("/api/auth/me", api.SessionAuthMiddleware(api.CurrentUserHandler))

	// User management routes - unprotected for initial user creation
	http.HandleFunc("/api/users/create", api.CreateUserHandler)

	// API routes - protected by session auth
	http.HandleFunc("/api/entries", api.SessionAuthMiddleware(api.Handler))
	http.HandleFunc("/api/entries/", api.SessionAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/entries/") {
			api.Entry(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	http.HandleFunc("/api/stats/weekly", api.SessionAuthMiddleware(api.Stats))

	// Serve static files - no auth required for login page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Allow access to login page and its assets
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/login.html" {
			fs := http.FileServer(http.Dir("./public"))
			fs.ServeHTTP(w, r)
			return
		}

		// For other static files, serve them without auth
		// (CSS, JS, etc. should be accessible)
		fs := http.FileServer(http.Dir("./public"))
		fs.ServeHTTP(w, r)
	})

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
