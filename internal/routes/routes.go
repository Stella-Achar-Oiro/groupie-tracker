package routes

import (
    "net/http"
    "groupie-tracker/internal/handlers"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

func SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", enableCORS(handlers.GetAllData))
    mux.HandleFunc("/api/search", enableCORS(handlers.SearchArtists))
    return mux
}