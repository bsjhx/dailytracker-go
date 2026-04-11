package main

import (
	"dailytracker/api"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// API routes
	http.HandleFunc("/api/entries", api.Handler)
	http.HandleFunc("/api/entries/", func(w http.ResponseWriter, r *http.Request) {
		// Route entries/:date to the Entry handler
		if strings.HasPrefix(r.URL.Path, "/api/entries/") {
			api.Entry(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/stats/weekly", api.Stats)

	// Serve static files from public directory
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
