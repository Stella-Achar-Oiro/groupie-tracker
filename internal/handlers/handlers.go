package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/services"
	"net/http"
	"strings"
)

func GetAllData(w http.ResponseWriter, r *http.Request) {
	data, err := services.FetchData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	data, err := services.FetchData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filteredArtists []models.Artist
	for _, artist := range data.ArtistsData {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredArtists)
}
