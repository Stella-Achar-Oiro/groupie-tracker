package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var artists []models.Artist
	var locations models.Locations
	var dates models.Dates
	var relations models.Relations
	var errs []error

	wg.Add(4)
	go func() {
		defer wg.Done()
		var err error
		artists, err = api.FetchArtists()
		if err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		locations, err = api.FetchLocations()
		if err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		dates, err = api.FetchDates()
		if err != nil {
			errs = append(errs, err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		relations, err = api.FetchRelations()
		if err != nil {
			errs = append(errs, err)
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("Error fetching data: %v", err)
		}
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	fullData := models.FullData{
		Artists:   artists,
		Locations: locations,
		Dates:     dates,
		Relations: relations,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fullData); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("artist")
	if artistID == "" {
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(artistID)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	relations, err := api.FetchRelations()
	if err != nil {
		log.Printf("Error fetching relations: %v", err)
		http.Error(w, "Failed to fetch event data", http.StatusInternalServerError)
		return
	}

	var artistEvents models.Relation
	for _, relation := range relations.Index {
		if relation.ID == uint64(id) {
			artistEvents = relation
			break
		}
	}

	if artistEvents.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artistEvents); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
