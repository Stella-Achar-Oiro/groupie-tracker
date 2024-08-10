package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	var artist models.Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	relation, err := api.FetchRelation()
	if err != nil {
		http.Error(w, "Failed to fetch relation data", http.StatusInternalServerError)
		return
	}

	var artistRelation models.Relation
	for _, r := range relation.Index {
		if r.ID == id {
			artistRelation.Index = []struct {
				ID             int                    `json:"id"`
				DatesLocations map[string]interface{} `json:"datesLocations"`
			}{r}
			break
		}
	}

	data := struct {
		Artist   models.Artist
		Relation models.Relation
	}{
		Artist:   artist,
		Relation: artistRelation,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}