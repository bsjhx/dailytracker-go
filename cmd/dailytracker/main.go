package main

import (
	"dailytracker/internal/handlers"
	"dailytracker/internal/middleware"
	"dailytracker/internal/repository"
	"dailytracker/internal/version"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	log.Printf("🎯 DailyTracker v%s", version.Version)

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database (runs migrations)
	_, err := repository.GetDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize session store for middleware
	middleware.SetSessionStore(handlers.GetSessionStore())

	// Start session cleaner
	handlers.StartSessionCleaner()

	// Authentication routes (no auth required)
	http.HandleFunc("/api/auth/login", handlers.Login)
	http.HandleFunc("/api/auth/logout", handlers.Logout)
	http.HandleFunc("/api/auth/me", middleware.SessionAuth(handlers.CurrentUser))

	// User management routes - unprotected for initial user creation
	http.HandleFunc("/api/users/create", handlers.CreateUser)

	// API routes - protected by session auth
	http.HandleFunc("/api/entries", middleware.SessionAuth(handlers.ListEntries))
	http.HandleFunc("/api/entries/", middleware.SessionAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/entries/") {
			handlers.Entry(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	http.HandleFunc("/api/stats/weekly", middleware.SessionAuth(handlers.WeeklyStats))

	// Serve static files - no auth required for login page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index page for root
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web/templates/index.html")
			return
		}

		// Serve HTML templates
		if r.URL.Path == "/index.html" {
			http.ServeFile(w, r, "./web/templates/index.html")
			return
		}
		if r.URL.Path == "/login.html" {
			http.ServeFile(w, r, "./web/templates/login.html")
			return
		}

		// Serve static files from web/static
		fs := http.FileServer(http.Dir("./web/static"))
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	log.Printf("🚀 Server starting on port %s...", port)
	log.Printf("📍 http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
