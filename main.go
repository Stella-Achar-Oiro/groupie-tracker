package main

import (
	"groupie-tracker/internal/routes"
	"log"
	"net/http"
)

func main() {
    mux := routes.SetupRoutes()
    
    // Serve static files
    fs := http.FileServer(http.Dir("./static"))
    mux.Handle("/", fs)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}