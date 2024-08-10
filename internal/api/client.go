package api

import (
	"encoding/json"
	"groupie-tracker/internal/models"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}

	artistCache     []models.Artist
	locationCache   models.Location
	dateCache       models.Date
	relationCache   models.Relation
	artistCacheTime time.Time
	cacheDuration   = 5 * time.Minute
)

func FetchArtists() ([]models.Artist, error) {
	if time.Since(artistCacheTime) < cacheDuration {
		return artistCache, nil
	}

	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	artistCache = artists
	artistCacheTime = time.Now()
	return artists, nil
}

func FetchLocations() (models.Location, error) {
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return models.Location{}, err
	}
	defer resp.Body.Close()

	var location models.Location
	err = json.NewDecoder(resp.Body).Decode(&location)
	return location, err
}

func FetchDates() (models.Date, error) {
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		return models.Date{}, err
	}
	defer resp.Body.Close()

	var date models.Date
	err = json.NewDecoder(resp.Body).Decode(&date)
	return date, err
}

func FetchRelation() (models.Relation, error) {
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return models.Relation{}, err
	}
	defer resp.Body.Close()

	var relation models.Relation
	err = json.NewDecoder(resp.Body).Decode(&relation)
	return relation, err
}