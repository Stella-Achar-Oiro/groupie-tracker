// internal/handlers/search.go

package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
	"net/http"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	filteredArtists := filterArtists(artists, query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredArtists)
}

func filterArtists(artists []models.Artist, query string) []models.Artist {
	var filtered []models.Artist
	lowercaseQuery := strings.ToLower(query)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) {
			filtered = append(filtered, artist)
		} else {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), lowercaseQuery) {
					filtered = append(filtered, artist)
					break
				}
			}
		}
	}

	return filtered
}
