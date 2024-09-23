package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"groupie-tracker/internal/handlers"
	"groupie-tracker/internal/routes"
)
func main() {
	if len(os.Args) != 1 {
		fmt.Println("Usage: go run main.go")
		return
	}
	mux := routes.SetupRoutes()
	// Create a new file server handler for static files
	fs := http.FileServer(http.Dir("./static"))
	// Serve static files from the /static/ URL path
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	// Serve index.html from the templates folder
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/static/") && r.URL.Path != "/" {
			handlers.ErrorHandler(w, r, http.StatusNotFound, "Page Not Found")
			return
		}
		// Check if index.html exists before proceeding
	tmplPath := "templates/index.html"
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
		return
	}
		http.ServeFile(w, r, "./templates/index.html")
		log.Println("Template executed successfully")
	})
	// Create a custom server with timeouts
	server := &http.Server{
		Addr:         ":8020",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// Create a channel to listen for errors
	errChan := make(chan error)
	// Start the server in a goroutine
	go func() {
		log.Println("Server starting on http://localhost:8020")
		errChan <- server.ListenAndServe()
	}()
	// Wait for an error from the server
	err := <-errChan
	log.Fatal(err)
}