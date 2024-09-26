package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"

	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/models"
)

// HandleSuggestions handles the suggestions API endpoint
func HandleSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	cachedData, err := cache.GetCachedData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	suggestions := getSuggestions(query, cachedData.ArtistsData)
	json.NewEncoder(w).Encode(suggestions)
}

func getSuggestions(query string, artists []models.Artist) []models.Suggestion {
	var suggestions []models.Suggestion
	lowercaseQuery := strings.ToLower(query)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) {
			suggestions = append(suggestions, models.Suggestion{Text: artist.Name, Type: "artist"})
		}
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), lowercaseQuery) {
				suggestions = append(suggestions, models.Suggestion{Text: member, Type: "member"})
			}
		}
		if strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) {
			suggestions = append(suggestions, models.Suggestion{Text: artist.FirstAlbum, Type: "album"})
		}
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(suggestions, models.Suggestion{Text: strconv.Itoa(artist.CreationDate), Type: "year"})
		}
	}

	return suggestions
}