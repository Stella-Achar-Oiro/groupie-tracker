package api

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/internal/models"
	"io"
	"net/http"
)

const (
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	locationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	datesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	relationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

func FetchArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if err := fetchJSON(artistsURL, &artists); err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}
	return artists, nil
}

func FetchLocations() (models.Locations, error) {
	var locations models.Locations
	if err := fetchJSON(locationsURL, &locations); err != nil {
		return models.Locations{}, fmt.Errorf("failed to fetch locations: %w", err)
	}
	return locations, nil
}

func FetchDates() (models.Dates, error) {
	var dates models.Dates
	if err := fetchJSON(datesURL, &dates); err != nil {
		return models.Dates{}, fmt.Errorf("failed to fetch dates: %w", err)
	}
	return dates, nil
}

func FetchRelations() (models.Relations, error) {
	var relations models.Relations
	if err := fetchJSON(relationURL, &relations); err != nil {
		return models.Relations{}, fmt.Errorf("failed to fetch relations: %w", err)
	}
	return relations, nil
}

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Print the raw JSON
	//fmt.Println("Raw JSON:", string(body))

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}
