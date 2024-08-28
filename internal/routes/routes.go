package routes

import (
	"net/http"

	"groupie-tracker/internal/handlers"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		done := make(chan bool)

		// Start a goroutine to set CORS headers
		go func() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			done <- true
		}()

		// Wait for the goroutine to finish
		<-done

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next(w, r)
	}
}

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", enableCORS(handlers.GetAllData))
	mux.HandleFunc("/api/search", enableCORS(handlers.SearchArtists))
	return mux
}
