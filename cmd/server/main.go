package main

import (
	"groupie-tracker/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artists", handlers.ArtistsHandler)
	http.HandleFunc("/events", handlers.EventsHandler)

	fs := http.FileServer(http.Dir("web"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
