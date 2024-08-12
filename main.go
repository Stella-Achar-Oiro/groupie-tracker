package main

import (
	"groupie-tracker/internal/handlers"
	"groupie-tracker/internal/middleware"
	"log"
	"net/http"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes
	http.HandleFunc("/", middleware.LogRequest(handlers.HomeHandler))
	http.HandleFunc("/api/artists", middleware.RateLimit(handlers.ArtistsHandler))
	http.HandleFunc("/api/search", middleware.RateLimit(handlers.SearchHandler))
	http.HandleFunc("/artist/", middleware.LogRequest(handlers.ArtistDetailHandler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
