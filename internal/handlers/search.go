package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/internal/cache"
	"groupie-tracker/internal/models"
)

// HandleSearch handles the search API endpoint
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	
	var filters models.FilterParams
	err := json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		http.Error(w, "Invalid filter parameters", http.StatusBadRequest)
		return
	}

	cachedData, err := cache.GetCachedData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	results := searchArtists(query, cachedData.ArtistsData, filters)
	json.NewEncoder(w).Encode(results)
}

func searchArtists(query string, artists []models.Artist, filters models.FilterParams) models.SearchResult {
	var results models.SearchResult
	lowercaseQuery := strings.ToLower(query)

	for _, artist := range artists {
		if matchesFilters(artist, filters) &&
			(query == "" || // If query is empty, include all artists that match filters
				strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) ||
				containsAny(artist.Members, lowercaseQuery) ||
				strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) ||
				strconv.Itoa(artist.CreationDate) == query) {
			results.Artists = append(results.Artists, artist)
		}
	}

	return results
}

func containsAny(slice []string, substr string) bool {
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), substr) {
			return true
		}
	}
	return false
}

func matchesFilters(artist models.Artist, filters models.FilterParams) bool {
	// Check creation year
	if artist.CreationDate < filters.CreationYearMin || artist.CreationDate > filters.CreationYearMax {
		return false
	}

	// Check first album year
	firstAlbumYear, _ := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
	if firstAlbumYear < filters.FirstAlbumYearMin || firstAlbumYear > filters.FirstAlbumYearMax {
		return false
	}

	// Check number of members
	if len(filters.Members) > 0 {
		memberCount := len(artist.Members)
		if !contains(filters.Members, memberCount) {
			return false
		}
	}

	// Check locations
	if len(filters.Locations) > 0 {
		artistLocations := strings.Split(artist.Locations, ",")
		matched := false
		for _, loc := range artistLocations {
			for _, filterLoc := range filters.Locations {
				if strings.Contains(strings.ToLower(loc), strings.ToLower(filterLoc)) {
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// contains checks if a slice contains a value
func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}