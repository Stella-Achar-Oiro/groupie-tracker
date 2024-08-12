package main

import (
	"fmt"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// Serving templates files
	filesServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", filesServer))

	// Index handler
	http.HandleFunc("/", handlers.IndexHandler)
	// mainPage Handler
	http.HandleFunc("/home", handlers.HomeHandler)
	// artist handler
	http.HandleFunc("/artists/", handlers.ArtistHandler)
	// fliter handler
	http.HandleFunc("/filter", handlers.FilterHandler)

	fmt.Println("Server started at http://localhost:8010/")

	// Starting serveur
	http.ListenAndServe(":8010", nil)
}
