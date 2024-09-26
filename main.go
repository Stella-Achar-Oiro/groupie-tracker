package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/handlers"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/service"
)

var (
	indexTpl *template.Template
	logger   *log.Logger
)

const cacheDuration = 1 * time.Hour

func main() {
	// Initialize logger
	logger = log.New(os.Stdout, "GROUPIE-TRACKER: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize models package with required constants
	models.InitConstants(service.GetMapboxAccessToken(), service.GetMapboxGeocodingAPI())
	logger.Println("Models initialized with Mapbox constants")

	// Initialize cache
	cache.Init(cacheDuration)
	logger.Println("Cache initialized with duration:", cacheDuration)

	// Initial data fetch
	if err := cache.RefreshCache(); err != nil {
		logger.Fatalf("Failed to fetch initial data: %v", err)
	}
	logger.Println("Initial data fetched successfully")

	// Parse HTML template
	var err error
	indexTpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Fatalf("Failed to parse template: %v", err)
	}
	logger.Println("HTML template parsed successfully")

	// Set up routes
	http.HandleFunc("/", handlers.HandleIndex(indexTpl))
	http.HandleFunc("/api/search", handlers.HandleSearch)
	http.HandleFunc("/api/artist/", handlers.HandleArtist)
	http.HandleFunc("/api/suggestions", handlers.HandleSuggestions)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	logger.Println("Routes and static file server set up")

	// Start server
	port := ":8080"
	logger.Printf("Server starting on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
